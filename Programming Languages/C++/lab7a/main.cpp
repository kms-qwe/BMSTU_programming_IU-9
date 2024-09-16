#include "declaration.h"
#include <iostream>

int main() {
    Matrix matrix(3);

    matrix[0][0] = 1;
    matrix[0][1] = 2;
    matrix[0][2] = 3;
    matrix[1][0] = 4;
    matrix[1][1] = 5;
    matrix[1][2] = 6;
    matrix[2][0] = 7;
    matrix[2][1] = 8;
    matrix[2][2] = 9;

    std::cout << "Matrix:\n" << matrix << "\n";
    std::cout << "Size: " << matrix.getSize() << "\n";

    std::cout << "Element (2, 0) is min in row and max in column? "
               << std::boolalpha << matrix.isMinMaxElement(2, 0) << "\n";

    return 0;
}