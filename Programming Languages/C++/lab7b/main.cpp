#include "declaration.h"
#include <iostream>

void printSequence(AsciiSequence& seq) {
    for (int i = 0; i < seq.getLength(); i++) {
        std::cout << seq[i];
    }
    std::cout << std::endl;
}

int main() {
    AsciiSequence seq1("hello");
    std::cout << "Original sequence: ";
    printSequence(seq1);

    seq1.insert(2, 'x');
    std::cout << "After insertion: ";
    printSequence(seq1);

    std::cout << "Is palindrome? " << std::boolalpha << seq1.isPalindrome() << std::endl;

    AsciiSequence seq2 = seq1;
    std::cout << "Copied sequence: ";
    printSequence(seq2);

    seq2[0] = 'a';
    std::cout << "Modified copy: ";
    printSequence(seq2);

    AsciiSequence seq3("racecar");
    std::cout << "Palindrome sequence: ";
    printSequence(seq3);
    std::cout << "Is palindrome? " << std::boolalpha << seq3.isPalindrome() << std::endl;

    return 0;
}