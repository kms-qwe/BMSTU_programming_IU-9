#ifndef MATRIX_H
#define MATRIX_H

#include <iostream>

class Matrix {
private:
    int** data;
    int size;

public:
    Matrix(int n);
    ~Matrix();
    int getSize() const;

    class Row {
    private:
        Matrix& matrix;
        int row_idx;

    public:
        Row(Matrix& mat, int idx);
        int& operator[](int col_idx);
        friend std::ostream& operator<<(std::ostream& os, const Matrix& matrix);
    };

    Row operator[](int row_idx);
    bool isMinMaxElement(int row_idx, int col_idx) const;
    friend std::ostream& operator<<(std::ostream& os, const Matrix& matrix);
};

std::ostream& operator<<(std::ostream& os, const Matrix& matrix);

#endif