

#ifndef SEGMENT_TREE_H
#define SEGMENT_TREE_H

#include <vector>
#include <functional>
#include <string>

template<typename T, size_t N>
class SegmentTree {
private:
    std::vector<T> tree;
    std::function<T(T, T)> operation;
    size_t size;

    void build(const std::vector<T>& arr, size_t v, size_t tl, size_t tr) {
        if (tl == tr) {
            tree[v] = arr[tl];
        } else {
            size_t tm = (tl + tr) / 2;
            build(arr, v * 2, tl, tm);
            build(arr, v * 2 + 1, tm + 1, tr);
            tree[v] = operation(tree[v * 2], tree[v * 2 + 1]);
        }
    }

    void update(size_t v, size_t tl, size_t tr, size_t pos, const T& new_val) {
        if (tl == tr) {
            tree[v] = new_val;
        } else {
            size_t tm = (tl + tr) / 2;
            if (pos <= tm)
                update(v * 2, tl, tm, pos, new_val);
            else
                update(v * 2 + 1, tm + 1, tr, pos, new_val);
            tree[v] = operation(tree[v * 2], tree[v * 2 + 1]);
        }
    }

    T query(size_t v, size_t tl, size_t tr, size_t l, size_t r) {
        if (l > r)
            return T(); 

        if (l == tl && r == tr)
            return tree[v];

        size_t tm = (tl + tr) / 2;
        return operation(query(v * 2, tl, tm, l, std::min(r, tm)),
                         query(v * 2 + 1, tm + 1, tr, std::max(l, tm + 1), r));
    }

public:
    SegmentTree(const std::vector<T>& arr, const std::function<T(T, T)>& op) {
        size = N;
        tree.resize(4 * N);
        operation = op;
        build(arr, 1, 0, size - 1);
    }

    void update(size_t pos, const T& new_val) {
        update(1, 0, size - 1, pos, new_val);
    }

    T query(size_t l, size_t r) {
        return query(1, 0, size - 1, l, r);
    }
};


template<size_t N>
class SegmentTree<std::string, N> {
private:
    std::string tree[N];
    size_t size;

public:
    SegmentTree(const std::string arr[N]) {
        size = N;
        for (size_t i = 0; i < N; ++i) {
            tree[i] = arr[i];
        }
    }

    void update(size_t pos, const std::string& new_val) {
        if (pos < N)
            tree[pos] = new_val;
        else
            std::cerr << "Error: Index out of bounds" << std::endl;
    }
    
    std::string query(size_t l, size_t r) {
        if (l > r || r >= N) {
            std::cerr << "Error: Invalid range" << std::endl;
            return "";
        }
        
        std::string result;
        for (size_t i = l; i <= r; ++i) {
            result += tree[i];
        }
        return result;
    }
};
#endif 
