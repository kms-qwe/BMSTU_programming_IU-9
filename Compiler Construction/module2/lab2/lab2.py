import abc
import sys
from dataclasses import dataclass
from pprint import pprint

import parser_edsl as pe


# =========================
# Абстрактный синтаксис
# =========================


class TypeSpec(abc.ABC):
    pass


class Expr(abc.ABC):
    pass


# Program → Declarations
@dataclass
class Program:
    declarations: list['Declaration']


# Declaration → TypeSpec Declarators? ;
@dataclass
class Declaration:
    type_spec: TypeSpec
    declarators: list['Declarator']


# BuiltinType → int | char | double | float
@dataclass
class BuiltinType(TypeSpec):
    name: str


# StructType → struct Tag? Fields?
@dataclass
class StructType(TypeSpec):
    tag: str | None
    fields: list['Declaration'] | None


# UnionType → union Tag? Fields?
@dataclass
class UnionType(TypeSpec):
    tag: str | None
    fields: list['Declaration'] | None


# EnumType → enum Tag? Enumerators?
@dataclass
class EnumType(TypeSpec):
    tag: str | None
    enumerators: list['Enumerator'] | None


# Enumerator → IDENT | IDENT = Expr
@dataclass
class Enumerator:
    name: str
    value: Expr | None

# ClassType → class Tag : IDENT { FieldDeclarations }
@dataclass
class ClassType(TypeSpec):
    tag: str
    parent: str | None
    fields: list['Declaration']

# Declarator → Pointer* Name ArraySuffix*
@dataclass
class Declarator:
    pointer_level: int
    name: str | None
    arrays: list[Expr | None]

# Expr → IDENT
@dataclass
class IdentifierExpr(Expr):
    name: str


# Expr → INT_CONST
@dataclass
class IntExpr(Expr):
    value: int


# Expr → Expr BinOp Expr
@dataclass
class BinaryExpr(Expr):
    left: Expr
    op: str
    right: Expr


# Expr → UnOp Expr
@dataclass
class UnaryExpr(Expr):
    op: str
    expr: Expr


# Expr → sizeof ( TypeSpec )
@dataclass
class SizeofTypeExpr(Expr):
    type_spec: TypeSpec


# =========================
# Служебные значения
# =========================

UNSIZED_ARRAY = object()


# =========================
# Лексическая структура
# =========================


IDENT = pe.Terminal('IDENT', r'[A-Za-z_][A-Za-z0-9_]*', str)
INT_CONST = pe.Terminal('INT_CONST', r'[0-9]+', int, priority=7)


# =========================
# Нетерминальные символы
# =========================


NProgram = pe.NonTerminal('Program')
NDeclarations = pe.NonTerminal('Declarations')
NDeclaration = pe.NonTerminal('Declaration')

NBuiltinDecl = pe.NonTerminal('BuiltinDecl')
NStructDecl = pe.NonTerminal('StructDecl')
NUnionDecl = pe.NonTerminal('UnionDecl')
NEnumDecl = pe.NonTerminal('EnumDecl')
NClassDecl = pe.NonTerminal('ClassDecl')

NTypeSpec = pe.NonTerminal('TypeSpec')
NBuiltinType = pe.NonTerminal('BuiltinType')

NDeclaratorsOpt = pe.NonTerminal('DeclaratorsOpt')
NDeclaratorList = pe.NonTerminal('DeclaratorList')
NDeclarator = pe.NonTerminal('Declarator')
NDeclaratorNoIdent = pe.NonTerminal('DeclaratorNoIdent')
NPointerOpt = pe.NonTerminal('PointerOpt')
NPointerPlus = pe.NonTerminal('PointerPlus')
NArraySuffixes = pe.NonTerminal('ArraySuffixes')
NArraySuffixesPlus = pe.NonTerminal('ArraySuffixesPlus')
NArraySuffix = pe.NonTerminal('ArraySuffix')

NFieldDeclarations = pe.NonTerminal('FieldDeclarations')

NEnumerators = pe.NonTerminal('Enumerators')
NEnumerator = pe.NonTerminal('Enumerator')
NTrailingCommaOpt = pe.NonTerminal('TrailingCommaOpt')

NExpr = pe.NonTerminal('Expr')
NAddExpr = pe.NonTerminal('AddExpr')
NMulExpr = pe.NonTerminal('MulExpr')
NUnaryExpr = pe.NonTerminal('UnaryExpr')
NPrimaryExpr = pe.NonTerminal('PrimaryExpr')


# =========================
# Вспомогательные функции
# =========================


def make_binary(op):
    return lambda left, right: BinaryExpr(left, op, right)


def build_declarator(pointer_level, name, arrays):
    real_arrays = [None if dim is UNSIZED_ARRAY else dim for dim in arrays]
    return Declarator(pointer_level, name, real_arrays)


# =========================
# Грамматика: программа и объявления
# =========================


NProgram |= NDeclarations, Program

NDeclarations |= lambda: []
NDeclarations |= NDeclarations, NDeclaration, lambda ds, d: ds + [d]

NDeclaration |= NBuiltinDecl
NDeclaration |= NStructDecl
NDeclaration |= NUnionDecl
NDeclaration |= NEnumDecl
NDeclaration |= NClassDecl

# =========================
# Грамматика: объявления встроенных типов
# =========================


# Declaration → BuiltinType DeclaratorsOpt ;
NBuiltinDecl |= NBuiltinType, NDeclaratorsOpt, ';', Declaration


# =========================
# Грамматика: объявления struct
# =========================


# struct;
NStructDecl |= 'struct', ';', lambda: Declaration(StructType(None, None), [])

# struct Tag;
NStructDecl |= 'struct', IDENT, ';', lambda tag: Declaration(StructType(tag, None), [])

# struct { Fields };
NStructDecl |= 'struct', '{', NFieldDeclarations, '}', ';', \
    lambda fields: Declaration(StructType(None, fields), [])

# struct Tag { Fields };
NStructDecl |= 'struct', IDENT, '{', NFieldDeclarations, '}', ';', \
    lambda tag, fields: Declaration(StructType(tag, fields), [])

# struct { Fields } Declarators ;
NStructDecl |= 'struct', '{', NFieldDeclarations, '}', NDeclaratorList, ';', \
    lambda fields, ds: Declaration(StructType(None, fields), ds)

# struct Tag { Fields } Declarators ;
NStructDecl |= 'struct', IDENT, '{', NFieldDeclarations, '}', NDeclaratorList, ';', \
    lambda tag, fields, ds: Declaration(StructType(tag, fields), ds)

# struct *p;
NStructDecl |= 'struct', NDeclaratorNoIdent, ';', \
    lambda d: Declaration(StructType(None, None), [d])

# struct Type x;
NStructDecl |= 'struct', IDENT, NDeclaratorList, ';', \
    lambda tag, ds: Declaration(StructType(tag, None), ds)


# =========================
# Грамматика: объявления union
# =========================


# union;
NUnionDecl |= 'union', ';', lambda: Declaration(UnionType(None, None), [])

# union Tag;
NUnionDecl |= 'union', IDENT, ';', lambda tag: Declaration(UnionType(tag, None), [])

# union { Fields };
NUnionDecl |= 'union', '{', NFieldDeclarations, '}', ';', \
    lambda fields: Declaration(UnionType(None, fields), [])

# union Tag { Fields };
NUnionDecl |= 'union', IDENT, '{', NFieldDeclarations, '}', ';', \
    lambda tag, fields: Declaration(UnionType(tag, fields), [])

# union { Fields } Declarators ;
NUnionDecl |= 'union', '{', NFieldDeclarations, '}', NDeclaratorList, ';', \
    lambda fields, ds: Declaration(UnionType(None, fields), ds)

# union Tag { Fields } Declarators ;
NUnionDecl |= 'union', IDENT, '{', NFieldDeclarations, '}', NDeclaratorList, ';', \
    lambda tag, fields, ds: Declaration(UnionType(tag, fields), ds)

# union *p;
NUnionDecl |= 'union', NDeclaratorNoIdent, ';', \
    lambda d: Declaration(UnionType(None, None), [d])

# union Type x;
NUnionDecl |= 'union', IDENT, NDeclaratorList, ';', \
    lambda tag, ds: Declaration(UnionType(tag, None), ds)


# =========================
# Грамматика: объявления enum
# =========================


# enum;
NEnumDecl |= 'enum', ';', lambda: Declaration(EnumType(None, None), [])

# enum Tag;
NEnumDecl |= 'enum', IDENT, ';', lambda tag: Declaration(EnumType(tag, None), [])

# enum { };
NEnumDecl |= 'enum', '{', '}', ';', lambda: Declaration(EnumType(None, []), [])

# enum Tag { };
NEnumDecl |= 'enum', IDENT, '{', '}', ';', \
    lambda tag: Declaration(EnumType(tag, []), [])

# enum { Enumerators };
NEnumDecl |= 'enum', '{', NEnumerators, NTrailingCommaOpt, '}', ';', \
    lambda enums, _comma: Declaration(EnumType(None, enums), [])

# enum Tag { Enumerators };
NEnumDecl |= 'enum', IDENT, '{', NEnumerators, NTrailingCommaOpt, '}', ';', \
    lambda tag, enums, _comma: Declaration(EnumType(tag, enums), [])

# enum { Enumerators } Declarators ;
NEnumDecl |= 'enum', '{', NEnumerators, NTrailingCommaOpt, '}', NDeclaratorList, ';', \
    lambda enums, _comma, ds: Declaration(EnumType(None, enums), ds)

# enum Tag { Enumerators } Declarators ;
NEnumDecl |= 'enum', IDENT, '{', NEnumerators, NTrailingCommaOpt, '}', NDeclaratorList, ';', \
    lambda tag, enums, _comma, ds: Declaration(EnumType(tag, enums), ds)

# enum [];
NEnumDecl |= 'enum', NDeclaratorNoIdent, ';', \
    lambda d: Declaration(EnumType(None, None), [d])

# enum Color x;
NEnumDecl |= 'enum', IDENT, NDeclaratorList, ';', \
    lambda tag, ds: Declaration(EnumType(tag, None), ds)


# =========================
# Грамматика: поля struct/union
# =========================


NFieldDeclarations |= lambda: []
NFieldDeclarations |= NFieldDeclarations, NDeclaration, lambda ds, d: ds + [d]


# =========================
# Грамматика: классы
# =========================

NClassDecl |= 'class', IDENT, ':', IDENT, '{', NFieldDeclarations, '}', lambda tag, parent, fields: Declaration(ClassType(tag, parent, fields), [])

NClassDecl |= 'class', IDENT, ':', 'public', IDENT, '{', NFieldDeclarations, '}', lambda tag, parent, fields: Declaration(ClassType(tag, parent, fields), [])

NClassDecl |= 'class', IDENT, ':', '{', NFieldDeclarations, '}', lambda tag, fields: Declaration(ClassType(tag, None, fields), [])


# =========================
# Грамматика: деклараторы
# =========================


NDeclaratorsOpt |= lambda: []
NDeclaratorsOpt |= NDeclaratorList

NDeclaratorList |= NDeclarator, lambda d: [d]
NDeclaratorList |= NDeclaratorList, ',', NDeclarator, lambda ds, d: ds + [d]

# Declarator → PointerOpt IDENT ArraySuffixes
NDeclarator |= NPointerOpt, IDENT, NArraySuffixes, build_declarator

# DeclaratorNoIdent → PointerPlus IDENT ArraySuffixes
#                  | PointerOpt IDENT ArraySuffixesPlus
NDeclaratorNoIdent |= NPointerPlus, IDENT, NArraySuffixes, build_declarator
NDeclaratorNoIdent |= NPointerOpt, IDENT, NArraySuffixesPlus, build_declarator

NPointerOpt |= lambda: 0
NPointerOpt |= NPointerPlus

NPointerPlus |= '*', lambda: 1
NPointerPlus |= NPointerPlus, '*', lambda n: n + 1

NArraySuffixes |= lambda: []
NArraySuffixes |= NArraySuffixesPlus

NArraySuffixesPlus |= NArraySuffix, lambda s: [s]
NArraySuffixesPlus |= NArraySuffixesPlus, NArraySuffix, lambda ss, s: ss + [s]

NArraySuffix |= '[', ']', lambda: UNSIZED_ARRAY
NArraySuffix |= '[', NExpr, ']', lambda expr: expr


# =========================
# Грамматика: встроенные типы
# =========================


NTypeSpec |= NBuiltinType
NTypeSpec |= 'struct', lambda: StructType(None, None)
NTypeSpec |= 'struct', IDENT, lambda tag: StructType(tag, None)
NTypeSpec |= 'struct', '{', NFieldDeclarations, '}', lambda fields: StructType(None, fields)
NTypeSpec |= 'struct', IDENT, '{', NFieldDeclarations, '}', lambda tag, fields: StructType(tag, fields)

NTypeSpec |= 'union', lambda: UnionType(None, None)
NTypeSpec |= 'union', IDENT, lambda tag: UnionType(tag, None)
NTypeSpec |= 'union', '{', NFieldDeclarations, '}', lambda fields: UnionType(None, fields)
NTypeSpec |= 'union', IDENT, '{', NFieldDeclarations, '}', lambda tag, fields: UnionType(tag, fields)

NTypeSpec |= 'enum', lambda: EnumType(None, None)
NTypeSpec |= 'enum', IDENT, lambda tag: EnumType(tag, None)
NTypeSpec |= 'enum', '{', '}', lambda: EnumType(None, [])
NTypeSpec |= 'enum', '{', NEnumerators, NTrailingCommaOpt, '}', \
    lambda enums, _comma: EnumType(None, enums)
NTypeSpec |= 'enum', IDENT, '{', '}', lambda tag: EnumType(tag, [])
NTypeSpec |= 'enum', IDENT, '{', NEnumerators, NTrailingCommaOpt, '}', \
    lambda tag, enums, _comma: EnumType(tag, enums)

NBuiltinType |= 'int', lambda: BuiltinType('int')
NBuiltinType |= 'char', lambda: BuiltinType('char')
NBuiltinType |= 'double', lambda: BuiltinType('double')
NBuiltinType |= 'float', lambda: BuiltinType('float')


# =========================
# Грамматика: перечислители
# =========================


NTrailingCommaOpt |= lambda: []
NTrailingCommaOpt |= ',', lambda: []

NEnumerators |= NEnumerator, lambda e: [e]
NEnumerators |= NEnumerators, ',', NEnumerator, lambda es, e: es + [e]

NEnumerator |= IDENT, lambda name: Enumerator(name, None)
NEnumerator |= IDENT, '=', NExpr, lambda name, value: Enumerator(name, value)


# =========================
# Грамматика: выражения для enum
# =========================


NExpr |= NAddExpr

NAddExpr |= NMulExpr
NAddExpr |= NAddExpr, '+', NMulExpr, make_binary('+')
NAddExpr |= NAddExpr, '-', NMulExpr, make_binary('-')

NMulExpr |= NUnaryExpr
NMulExpr |= NMulExpr, '*', NUnaryExpr, make_binary('*')
NMulExpr |= NMulExpr, '/', NUnaryExpr, make_binary('/')

NUnaryExpr |= NPrimaryExpr
NUnaryExpr |= '+', NUnaryExpr, lambda expr: UnaryExpr('+', expr)
NUnaryExpr |= '-', NUnaryExpr, lambda expr: UnaryExpr('-', expr)
NUnaryExpr |= 'sizeof', NUnaryExpr, lambda expr: UnaryExpr('sizeof', expr)
NUnaryExpr |= 'sizeof', '(', NTypeSpec, ')', SizeofTypeExpr

NPrimaryExpr |= IDENT, IdentifierExpr
NPrimaryExpr |= INT_CONST, IntExpr
NPrimaryExpr |= '(', NExpr, ')'


# =========================
# Парсер
# =========================


parser = pe.Parser(NProgram, method=pe.EARLEY)
parser.add_skipped_domain(r'\s+')
parser.add_skipped_domain(r'//[^\n]*')
parser.add_skipped_domain(r'/\*.*?\*/')


def parse_text(text: str) -> Program:
    return parser.parse(text)


def main() -> None:
    if len(sys.argv) < 2:
        print('Использование: python lab2.py <file1> [file2 ...]')
        return

    for filename in sys.argv[1:]:
        print(f'\n===== {filename} =====')
        try:
            with open(filename, 'r', encoding='utf-8') as f:
                text = f.read()
            tree = parse_text(text)
            pprint(tree, sort_dicts=False)
        except pe.Error as e:
            print(f'Ошибка {e.pos}: {e.message}')
        except Exception as e:
            print(f'Непредвиденная ошибка: {e}')


if __name__ == '__main__':
    main()

# добавим ООП с классами и public, тело считаем обязательным и тег тоже, а предок может отсутсвовать.
