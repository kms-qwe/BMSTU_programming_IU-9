#ifndef ASCII_SEQUENCE_H
#define ASCII_SEQUENCE_H

#include <cstring>

class AsciiSequence {
private:
    char* data;
    int length;
    int capacity;

    void reserve(int new_capacity);

public:
    AsciiSequence();
    AsciiSequence(const char* str);
    AsciiSequence(const AsciiSequence& other);
    ~AsciiSequence();

    AsciiSequence& operator=(const AsciiSequence& other);

    int getLength() const;
    char& operator[](int index);
    void insert(int index, char c);
    bool isPalindrome() const;
};

#endif