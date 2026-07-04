# Operators

Go provides a standard set of operators to manipulate variables and values. 

## 1. Arithmetic Operators
Used for standard mathematical operations.

| Operator | Description | Example |
| :---: | :--- | :--- |
| `+` | Addition | `a + b` |
| `-` | Subtraction | `a - b` |
| `*` | Multiplication | `a * b` |
| `/` | Division | `a / b` |
| `%` | Modulus (Remainder) | `a % b` |

*Note: Go does not support prefix increment/decrement (`++i` or `--i`). It only supports postfix increment/decrement as independent statements (`i++` and `i--`), not as expressions.*

## 2. Comparison Operators
Used to compare two values. They always return a `bool` (`true` or `false`).

| Operator | Description | Example |
| :---: | :--- | :--- |
| `==` | Equal to | `a == b` |
| `!=` | Not equal to | `a != b` |
| `>` | Greater than | `a > b` |
| `<` | Less than | `a < b` |
| `>=` | Greater than or equal to | `a >= b` |
| `<=` | Less than or equal to | `a <= b` |

## 3. Logical Operators
Used to combine conditional statements.

| Operator | Description | Example |
| :---: | :--- | :--- |
| `&&` | Logical AND (true if both are true) | `a < 5 && a < 10` |
| <code>\|\|</code> | Logical OR (true if one is true) | `a < 5 \|\| a < 4` |
| `!` | Logical NOT (reverses the boolean) | `!(a < 5 && a < 10)` |

## 4. Bitwise Operators
Used to manipulate data at the binary level. Common in cryptography, networking, and performance-critical systems.

| Operator | Description | Example |
| :---: | :--- | :--- |
| `&` | Bitwise AND | `a & b` |
| <code>\|</code> | Bitwise OR | <code>a \| b</code> |
| `^` | Bitwise XOR | `a ^ b` |
| `<<` | Left shift | `a << 2` |
| `>>` | Right shift | `a >> 2` |

## 5. Assignment Operators
Shorthands for assigning values based on an operation.

* `=` (Assign)
* `+=` (Add and assign: `x += 5` is `x = x + 5`)
* `-=`, `*=`, `/=`, `%=`, `&=`, `|=`, `^=`, `<<=`, `>>=`
