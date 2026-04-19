from __future__ import annotations
import re
from dataclasses import dataclass
from typing import Any, Optional

_TOKEN_RE = re.compile(r"""
    \s*(
        -?\d+            |
        [^\s]+
    )
""", re.VERBOSE)

CONTROL = {"define", "end", "if", "else", "endif"}

def tokenize(src: str) -> list[str]:
    return [m.group(1) for m in _TOKEN_RE.finditer(src)]

def atom(tok: str) -> Any:
    if re.fullmatch(r"-?\d+", tok):
        return int(tok)
    return tok

@dataclass
class Parser:
    toks: list[str]
    i: int = 0

    def at_end(self) -> bool:
        return self.i >= len(self.toks)

    def peek(self) -> Optional[str]:
        return None if self.at_end() else self.toks[self.i]

    def pop(self) -> Optional[str]:
        if self.at_end():
            return None
        t = self.toks[self.i]
        self.i += 1
        return t

    def expect(self, value: str) -> bool:
        if self.peek() == value:
            self.i += 1
            return True
        return False

    def parse_program(self) -> Optional[list]:
        articles = {}
        while self.peek() == "define":
            art = self.parse_article()
            if art is None:
                return None
            name, body = art
            articles[name] = body

        body = self.parse_body(stoppers=set())
        if body is None:
            return None

        if not self.at_end():
            return None

        return [articles, body]

    def parse_article(self) -> Optional[tuple[str, list]]:
        if not self.expect("define"):
            return None

        name = self.pop()
        if name is None:
            return None
        if name in CONTROL or re.fullmatch(r"-?\d+", name):
            return None

        body = self.parse_body(stoppers={"end"})
        if body is None:
            return None

        if not self.expect("end"):
            return None

        return (name, body)

    def parse_body(self, stoppers: set[str]) -> Optional[list]:
        out: list[Any] = []

        while True:
            t = self.peek()
            if t is None:
                break
            if t in stoppers:
                break

            if t == "endif" or t == "end":
                return None

            if t == "else":
                return None

            if t == "define":
                return None

            if t == "if":
                node = self.parse_if()
                if node is None:
                    return None
                out.append(node)
                continue

            out.append(atom(self.pop()))

        return out

    def parse_if(self) -> Optional[list]:
        if not self.expect("if"):
            return None

        then_body = self.parse_body(stoppers={"else", "endif"})
        if then_body is None:
            return None

        t = self.peek()
        if t == "else":
            self.pop()
            else_body = self.parse_body(stoppers={"endif"})
            if else_body is None:
                return None
            if not self.expect("endif"):
                return None
            return ["if", then_body, else_body]

        if t == "endif":
            self.pop()
            return ["if", then_body]

        return None


def parse(src: str) -> Optional[list]:
    p = Parser(tokenize(src))
    return p.parse_program()


if __name__ == "__main__":
    print(parse("1 2 +"))
    print(parse("x dup 0 swap if drop -1 endif"))
    print(parse("x dup 0 swap if drop -1 else swap 1 + endif"))
    print(parse("define abs dup 0 < if -1 * endif end 10 abs -10 abs"))
    print(parse("define word w1 w2 w3"))
