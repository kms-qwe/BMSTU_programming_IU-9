#include "declaration.h"
#include <algorithm>
#include <limits>

Matrix::Matrix(int n) : size(n) {
    data = new int*[size];
    for (int i = 0; i < size; ++i) {
        data[i] = new int[size];
    }
}

Matrix::~Matrix() {
    for (int i = 0; i < size; ++i) {
        delete[] data[i];
    }
    delete[] data;
}

int Matrix::getSize() const {
    return size;
}

Matrix::Row::Row(Matrix& mat, int idx) : matrix(mat), row_idx(idx) {}

int& Matrix::Row::operator[](int col_idx) {
    return matrix.data[row_idx][col_idx];
}

Matrix::Row Matrix::operator[](int row_idx) {
    return Row(*this, row_idx);
}

bool Matrix::isMinMaxElement(int row_idx, int col_idx) const {
    int val = data[row_idx][col_idx];


    int* row = data[row_idx];
    int row_min = *std::min_element(row, row + size);
    if (row_min != val) {
        return false;
    }

    int col_max = data[0][col_idx];
    for (int i = 1; i < size; ++i) {
        col_max = std::max(col_max, data[i][col_idx]);
    }
    return col_max == val;
}

std::ostream& operator<<(std::ostream& os, const Matrix& matrix) {
    for (int i = 0; i < matrix.size; ++i) {
        for (int j = 0; j < matrix.size; ++j) {
            os << matrix.data[i][j] << " ";
        }
        os << "\n";
    }
    return os;
}