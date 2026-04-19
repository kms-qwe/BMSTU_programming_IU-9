# Грамматика языка (struct / union / enum)

## Программа

```
Program → Declarations

Declarations → ε
             | Declarations Declaration

Declaration → BuiltinDecl
            | StructDecl
            | UnionDecl
            | EnumDecl
```

---

## Встроенные типы

```
BuiltinDecl → BuiltinType DeclaratorsOpt ;

BuiltinType → int | char | double | float
```

---

## Struct

```
StructDecl → struct ;
           | struct IDENT ;
           | struct { FieldDeclarations } ;
           | struct IDENT { FieldDeclarations } ;
           | struct { FieldDeclarations } DeclaratorList ;
           | struct IDENT { FieldDeclarations } DeclaratorList ;
           | struct DeclaratorNoIdent ;
           | struct IDENT DeclaratorList ;
```

---

## Union

```
UnionDecl → union ;
          | union IDENT ;
          | union { FieldDeclarations } ;
          | union IDENT { FieldDeclarations } ;
          | union { FieldDeclarations } DeclaratorList ;
          | union IDENT { FieldDeclarations } DeclaratorList ;
          | union DeclaratorNoIdent ;
          | union IDENT DeclaratorList ;
```

---

## Enum

```
EnumDecl → enum ;
         | enum IDENT ;
         | enum { } ;
         | enum IDENT { } ;
         | enum { Enumerators } ;
         | enum IDENT { Enumerators } ;
         | enum { Enumerators } DeclaratorList ;
         | enum IDENT { Enumerators } DeclaratorList ;
         | enum DeclaratorNoIdent ;
         | enum IDENT DeclaratorList ;
```

---

## Поля

```
FieldDeclarations → ε
                  | FieldDeclarations Declaration
```

---

## Перечисления

```
Enumerators → Enumerator
            | Enumerators , Enumerator

Enumerator → IDENT
           | IDENT = Expr

TrailingCommaOpt → ε | ,
```

---

## Деклараторы

```
DeclaratorsOpt → ε
               | DeclaratorList

DeclaratorList → Declarator
               | DeclaratorList , Declarator

Declarator → PointerOpt OptName ArraySuffixes

DeclaratorNoIdent → PointerPlus OptName ArraySuffixes
                  | PointerOpt OptName ArraySuffixesPlus
```

---

## Указатели и массивы

```
PointerOpt → ε | PointerPlus
PointerPlus → * | PointerPlus *

OptName → ε | IDENT

ArraySuffixes → ε | ArraySuffixesPlus

ArraySuffixesPlus → ArraySuffix
                  | ArraySuffixesPlus ArraySuffix

ArraySuffix → [ ]
            | [ Expr ]
```

---

## Выражения

```
Expr → AddExpr

AddExpr → MulExpr
        | AddExpr + MulExpr
        | AddExpr - MulExpr

MulExpr → UnaryExpr
        | MulExpr * UnaryExpr
        | MulExpr / UnaryExpr

UnaryExpr → PrimaryExpr
          | + UnaryExpr
          | - UnaryExpr
          | sizeof UnaryExpr
          | sizeof ( TypeSpec )

PrimaryExpr → IDENT
            | INT_CONST
            | ( Expr )
```
