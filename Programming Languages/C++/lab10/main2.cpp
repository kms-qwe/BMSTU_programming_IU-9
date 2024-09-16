#include <iostream>
#include <vector>
#include <iterator>

class IntSequence {
private:
    std::vector<uint32_t> sequence;

public:
    IntSequence() = default;
    IntSequence(std::initializer_list<uint32_t> init) : sequence(init) {}

    void add(uint32_t num) {
        sequence.push_back(num);
    }

    class ZeroRunIterator : public std::iterator<std::input_iterator_tag, uint32_t> {
    private:
        const std::vector<uint32_t>& seq;
        std::vector<uint32_t>::const_iterator sequenceIter;
        uint32_t currentNum;
        uint32_t mask;
        uint32_t zeroRunLength;

    public:
        ZeroRunIterator(const std::vector<uint32_t>& sequence)
            : seq(sequence), sequenceIter(sequence.begin()), currentNum(sequence.empty() ? 0 : *sequenceIter), mask(0x80000000), zeroRunLength(0) {
            skipNonZeroRuns();
        }

        ZeroRunIterator(const std::vector<uint32_t>& sequence, std::vector<uint32_t>::const_iterator iter)
            : seq(sequence), sequenceIter(iter), currentNum(iter == sequence.end() ? 0 : *iter), mask(0), zeroRunLength(0) {}

        ZeroRunIterator& operator++() {
            if (mask == 0) {
                ++sequenceIter;
                if (sequenceIter != seq.end()) {
                    currentNum = *sequenceIter;
                    mask = 0x80000000;
                    skipNonZeroRuns();
                } else {
                    zeroRunLength = 0;
                }
            } else {
                if ((currentNum & mask) == 0) {
                    ++zeroRunLength;
                    mask >>= 1;
                } else {
                    mask >>= 1;
                    skipNonZeroRuns();
                }
            }
            return *this;
        }

        ZeroRunIterator operator++(int) {
            ZeroRunIterator tmp(*this);
            ++(*this);
            return tmp;
        }

        uint32_t operator*() const {
            return zeroRunLength;
        }

        bool operator==(const ZeroRunIterator& other) const {
            return sequenceIter == other.sequenceIter && mask == other.mask;
        }

        bool operator!=(const ZeroRunIterator& other) const {
            return !(*this == other);
        }

    private:
        void skipNonZeroRuns() {
            while (mask != 0 && (currentNum & mask) != 0) {
                mask >>= 1;
            }
            if (mask == 0) {
                zeroRunLength = 0;
            } else {
                zeroRunLength = 1;
            }
        }
    };

    ZeroRunIterator zero_run_begin() const {
        return ZeroRunIterator(sequence);
    }

    ZeroRunIterator zero_run_end() const {
        return ZeroRunIterator(sequence, sequence.end());
    }
};

int main() {
    IntSequence seq{0x80000000, 0x00000000, 0x80000000, 0xFFFFFFFF};

    std::cout << "Zero run lengths: ";
    for (auto it = seq.zero_run_begin(); it != seq.zero_run_end(); ++it) {
        std::cout << *it << " ";
    }
    std::cout << std::endl;

    return 0;
}