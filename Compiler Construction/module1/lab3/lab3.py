from __future__ import annotations

import sys
from dataclasses import dataclass
from enum import Enum, auto
from typing import Dict, List, Optional


EOF_CHAR = "\0"
MAX_INT64 = 2**63 - 1


@dataclass(order=True)
class Position:
    text: str
    index: int = 0
    line: int = 1
    col: int = 1

    @property
    def cp(self) -> str:
        if self.index >= len(self.text):
            return EOF_CHAR
        return self.text[self.index]

    @property
    def is_white_space(self) -> bool:
        return self.cp != EOF_CHAR and self.cp.isspace()

    @property
    def is_newline(self) -> bool:
        return self.cp in ("\n", "\r", EOF_CHAR)

    def copy(self) -> "Position":
        return Position(self.text, self.index, self.line, self.col)

    def advance(self) -> "Position":
        if self.index >= len(self.text):
            return self

        ch = self.text[self.index]
        if ch == "\r":
            if self.index + 1 < len(self.text) and self.text[self.index + 1] == "\n":
                self.index += 2
            else:
                self.index += 1
            self.line += 1
            self.col = 1
            return self

        if ch == "\n":
            self.index += 1
            self.line += 1
            self.col = 1
            return self

        self.index += 1
        self.col += 1
        return self

    def __str__(self) -> str:
        return f"({self.line}, {self.col})"


@dataclass
class Fragment:
    starting: Position
    following: Position

    def __str__(self) -> str:
        return f"{self.starting}-{self.following}"


@dataclass
class Message:
    is_error: bool
    position: Position
    text: str


class DomainTag(Enum):
    IDENT = auto()
    DECIMAL = auto()
    BINARY = auto()
    OCT = auto()
    STRING = auto()
    END_OF_PROGRAM = auto()


@dataclass
class Token:
    tag: DomainTag
    coords: Fragment
    attr: Optional[object] = None
    lexeme: Optional[str] = None

    def format_for_output(self) -> str:
        if self.attr is None:
            return f"{self.tag.name} {self.coords}"
        return f"{self.tag.name} {self.coords}: {self.attr}"


class Compiler:
    def __init__(self) -> None:
        self.messages: List[Message] = []
        self.name_codes: Dict[str, int] = {}
        self.names: List[str] = []

    def add_name(self, name: str) -> int:
        if name in self.name_codes:
            return self.name_codes[name]
        code = len(self.names)
        self.names.append(name)
        self.name_codes[name] = code
        return code

    def get_name(self, code: int) -> str:
        return self.names[code]

    def add_message(self, is_error: bool, position: Position, text: str) -> None:
        self.messages.append(Message(is_error, position.copy(), text))

    def get_scanner(self, program: str) -> "Scanner":
        return Scanner(program, self)


class Scanner:
    def __init__(self, program: str, compiler: Compiler) -> None:
        self.program = program
        self.compiler = compiler
        self.cur = Position(program)
        self.comments: List[Fragment] = []

    def _peek(self, offset: int = 1) -> str:
        idx = self.cur.index + offset
        if idx >= len(self.program):
            return EOF_CHAR
        return self.program[idx]

    def _skip_whitespace(self) -> None:
        while self.cur.is_white_space:
            self.cur.advance()

    def _is_ident_start(self, ch: str) -> bool:
        return ch in "?*|"

    def _is_ident_part(self, ch: str) -> bool:
        return ch.isdigit() or ch in "?*|"

    def _read_identifier(self) -> Token:
        start = self.cur.copy()
        chars: List[str] = []
        while self._is_ident_part(self.cur.cp):
            chars.append(self.cur.cp)
            self.cur.advance()
        name = "".join(chars)
        code = self.compiler.add_name(name)
        return Token(DomainTag.IDENT, Fragment(start, self.cur.copy()), code, name)

   # Восьмеричные литералы, заканчиваются на t
    def _read_number(self) -> Token:
        start = self.cur.copy()
        digits: List[str] = []
        while self.cur.cp.isdigit():
            digits.append(self.cur.cp)
            self.cur.advance()

        lexeme = "".join(digits)

        if self.cur.cp == "t":
            if all(ch in '01234567' for ch in lexeme):
                self.cur.advance()
                full_lexeme = lexeme + "t"
                value = int(lexeme, 8) if lexeme else 0
                if value > MAX_INT64:
                    self.compiler.add_message(True, start, "octal literal is too large for int64")
                return Token(DomainTag.OCT, Fragment(start, self.cur.copy()), value, full_lexeme)

        if self.cur.cp == "b":
            if all(ch in "01" for ch in lexeme):
                self.cur.advance()
                full_lexeme = lexeme + "b"
                value = int(lexeme, 2) if lexeme else 0
                if value > MAX_INT64:
                    self.compiler.add_message(True, start, "binary literal is too large for int64")
                return Token(DomainTag.BINARY, Fragment(start, self.cur.copy()), value, full_lexeme)

            self.compiler.add_message(True, self.cur.copy(), "invalid binary literal: only 0 and 1 are allowed before 'b'")
            self.cur.advance()
            return Token(DomainTag.DECIMAL, Fragment(start, self.cur.copy()), int(lexeme), lexeme + "b")

        value = int(lexeme)
        if value > MAX_INT64:
            self.compiler.add_message(True, start, "decimal literal is too large for int64")
        return Token(DomainTag.DECIMAL, Fragment(start, self.cur.copy()), value, lexeme)

    def _read_string(self) -> Token:
        start = self.cur.copy()
        self.cur.advance()  # opening backtick
        value_chars: List[str] = []

        while self.cur.cp != EOF_CHAR:
            if self.cur.cp == "`":
                if self._peek() == "`":
                    value_chars.append("`")
                    self.cur.advance()
                    self.cur.advance()
                    continue
                self.cur.advance()
                return Token(
                    DomainTag.STRING,
                    Fragment(start, self.cur.copy()),
                    "".join(value_chars),
                )

            value_chars.append(self.cur.cp)
            self.cur.advance()

        self.compiler.add_message(True, start, "unterminated string literal: closing backtick expected")
        return Token(
            DomainTag.STRING,
            Fragment(start, self.cur.copy()),
            "".join(value_chars),
        )

    def next_token(self) -> Token:
        while self.cur.cp != EOF_CHAR:
            self._skip_whitespace()
            start = self.cur.copy()

            if self.cur.cp == EOF_CHAR:
                break

            if self._is_ident_start(self.cur.cp):
                return self._read_identifier()

            if self.cur.cp.isdigit():
                return self._read_number()

            if self.cur.cp == "`":
                return self._read_string()

            self.compiler.add_message(True, self.cur.copy(), f"unexpected character: {self.cur.cp!r}")
            self.cur.advance()

        pos = self.cur.copy()
        return Token(DomainTag.END_OF_PROGRAM, Fragment(pos, pos))

    def scan_all(self):
        tokens: List[Token] = []
        while True:
            token = self.next_token()
            tokens.append(token)
            if token.tag == DomainTag.END_OF_PROGRAM:
                break
        return tokens, self.compiler.names[:], self.comments[:], self.compiler.messages[:]


def main() -> None:
    if len(sys.argv) != 2:
        print("Usage: python lexer_lab.py <input-file>")
        return

    with open(sys.argv[1], "r", encoding="utf-8") as f:
        program = f.read()

    compiler = Compiler()
    scanner = compiler.get_scanner(program)
    tokens, identifiers, comments, messages = scanner.scan_all()

    print("TOKENS:")
    for token in tokens:
        print(token.format_for_output())

    print("\nIDENTIFIERS:")
    for i, name in enumerate(identifiers):
        print(f"{i}: {name}")

    print("\nCOMMENTS:")
    if not comments:
        print("<none>")
    else:
        for comment in comments:
            print(comment)

    print("\nMESSAGES:")
    if not messages:
        print("<none>")
    else:
        for msg in messages:
            kind = "Error" if msg.is_error else "Warning"
            print(f"{kind} {msg.position}: {msg.text}")


if __name__ == "__main__":
    main()
