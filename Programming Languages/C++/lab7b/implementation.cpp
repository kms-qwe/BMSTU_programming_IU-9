#include "declaration.h"

void AsciiSequence::reserve(int new_capacity) {
    char* new_data = new char[new_capacity];
    std::memcpy(new_data, data, length * sizeof(char));
    delete[] data;
    data = new_data;
    capacity = new_capacity;
}

AsciiSequence::AsciiSequence() : data(nullptr), length(0), capacity(0) {}

AsciiSequence::AsciiSequence(const char* str) {
    length = std::strlen(str);
    capacity = length + 1;
    data = new char[capacity];
    std::strcpy(data, str);
}

AsciiSequence::AsciiSequence(const AsciiSequence& other) {
    length = other.length;
    capacity = other.capacity;
    data = new char[capacity];
    std::memcpy(data, other.data, length * sizeof(char));
    data[length] = '\0';
}

AsciiSequence::~AsciiSequence() {
    delete[] data;
}

AsciiSequence& AsciiSequence::operator=(const AsciiSequence& other) {
    if (this != &other) {
        delete[] data;
        length = other.length;
        capacity = other.capacity;
        data = new char[capacity];
        std::memcpy(data, other.data, length * sizeof(char));
        data[length] = '\0';
    }
    return *this;
}

int AsciiSequence::getLength() const {
    return length;
}

char& AsciiSequence::operator[](int index) {
    return data[index];
}

void AsciiSequence::insert(int index, char c) {
    if (length == capacity - 1) {
        reserve(2 * capacity);
    }

    for (int i = length; i > index; i--) {
        data[i] = data[i - 1];
    }

    data[index] = c;
    length++;
    data[length] = '\0';
}

bool AsciiSequence::isPalindrome() const {
    for (int i = 0; i < length / 2; i++) {
        if (data[i] != data[length - i - 1]) {
            return false;
        }
    }
    return true;
}