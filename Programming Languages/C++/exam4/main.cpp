#include <iostream>
#include <vector>
#include <utility>

class BinaryRelation {
public:
    using Pair = std::pair<int, int>;

    BinaryRelation(int n) : n_(n), matrix_(n, std::vector<bool>(n, false)) {}

    void addRelation(int i, int j) {
        if (i >= 0 && i < n_ && j >= 0 && j < n_) {
            matrix_[i][j] = true;
        }
    }

    bool hasRelation(int i, int j) const {
        if (i >= 0 && i < n_ && j >= 0 && j < n_) {
            return matrix_[i][j];
        }
        return false;
    }

    class ConstIterator {
    public:
        ConstIterator(const BinaryRelation& relation, bool end = false)
            : relation_(relation), i_(0), j_(0) {
            if (end) {
                i_ = relation_.n_;
                j_ = 0;
            } else {
                advanceToNextValid();
            }
        }

        ConstIterator& operator++() {
            ++j_;
            advanceToNextValid();
            return *this;
        }

        Pair operator*() const {
            return {i_, j_};
        }

        bool operator==(const ConstIterator& other) const {
            return i_ == other.i_ && j_ == other.j_ && &relation_ == &other.relation_;
        }

        bool operator!=(const ConstIterator& other) const {
            return !(*this == other);
        }

    private:
        void advanceToNextValid() {
            while (i_ < relation_.n_) {
                while (j_ < relation_.n_) {
                    if (relation_.matrix_[i_][j_]) {
                        return;
                    }
                    ++j_;
                }
                j_ = 0;
                ++i_;
            }
        }

        const BinaryRelation& relation_;
        int i_;
        int j_;
    };

    ConstIterator begin() const {
        return ConstIterator(*this);
    }

    ConstIterator end() const {
        return ConstIterator(*this, true);
    }

private:
    int n_;
    std::vector<std::vector<bool>> matrix_;
};

int main() {
    BinaryRelation relation(5);

    relation.addRelation(0, 1);
    relation.addRelation(2, 3);
    relation.addRelation(4, 0);

    for (auto it = relation.begin(); it != relation.end(); ++it) {
        auto [i, j] = *it;
        std::cout << "(" << i << ", " << j << ")\n";
    }

    return 0;
}
