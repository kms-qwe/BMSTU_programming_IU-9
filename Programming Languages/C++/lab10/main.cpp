#include <iostream>
#include <vector>
#include <iterator>

class IntSequence {
private:
    std::vector<int> sequence;

public:
    IntSequence() = default;

    IntSequence(std::initializer_list<int> init) : sequence(init) {}

    void add(int num) {
        sequence.push_back(num);
    }

    class DigitIterator : public std::iterator<std::bidirectional_iterator_tag, int> {
    private:
        const std::vector<int>& seq;
        std::vector<int>::const_iterator sequenceIter;
        int currentNum;
        int digit;

    public:
        DigitIterator(const std::vector<int>& sequence)
            : seq(sequence), sequenceIter(sequence.begin()), currentNum(sequence.empty() ? 0 : *sequenceIter), digit(digitCount(currentNum)) {}

        DigitIterator(const std::vector<int>& sequence, std::vector<int>::const_iterator iter)
            : seq(sequence), sequenceIter(iter), currentNum(*iter), digit(digitCount(currentNum)) {}

        DigitIterator& operator++() {
            if (digit == 0) {
                ++sequenceIter;
                if (sequenceIter != seq.end()) {
                    currentNum = *sequenceIter;
                    digit = digitCount(currentNum);
                }
            } else {
                --digit;
            }
            return *this;
        }

        DigitIterator operator++(int) {
            DigitIterator tmp(*this);
            ++(*this);
            return tmp;
        }

        DigitIterator& operator--() {
            if (digit == digitCount(currentNum)) {
                if (sequenceIter == seq.begin()) {
                    currentNum = 0;
                    digit = 0;
                } else {
                    --sequenceIter;
                    currentNum = *sequenceIter;
                    digit = digitCount(currentNum);
                }
            } else {
                ++digit;
            }
            return *this;
        }

        DigitIterator operator--(int) {
            DigitIterator tmp(*this);
            --(*this);
            return tmp;
        }

        int operator*() const {
            int divisor = 1;
            for (int i = 0; i < digit; ++i) {
                divisor *= 10;
            }
            return (currentNum / divisor) % 10;
        }

        bool operator==(const DigitIterator& other) const {
            return sequenceIter == other.sequenceIter && digit == other.digit;
        }

        bool operator!=(const DigitIterator& other) const {
            return !(*this == other);
        }

    private:
        static int digitCount(int num) {
            int count = 0;
            while (num != 0) {
                num /= 10;
                ++count;
            }
            return count;
        }
    };

    DigitIterator digit_begin() const {
        return DigitIterator(sequence);
    }

    DigitIterator digit_end() const {
        return DigitIterator(sequence, sequence.end());
    }
};

int main() {
    IntSequence seq{123, 45, 6789};

    for (auto it = seq.digit_begin(); it != seq.digit_end(); ++it) {
        std::cout << *it;
    }
    std::cout << std::endl;

    return 0;
}