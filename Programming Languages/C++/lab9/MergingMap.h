#include <map>
#include <string>

template<typename K, typename V>
class MergingMap {
private:
    std::map<K, V> data;

public:
    MergingMap() {}
    MergingMap(int ) {}

    V& operator[](const K& key) {
        return data[key];
    }

    MergingMap<K, V> operator+(const MergingMap<K, V>& other) const {
        MergingMap<K, V> result(0);
        for (const auto& pair : data) {
            result[pair.first] = pair.second;
        }
        for (const auto& pair : other.data) {
            result[pair.first] = result[pair.first] + pair.second;
        }
        return result;
    }

    bool operator==(const MergingMap<K, V>& other) const {
        return data == other.data;
    }

    bool operator!=(const MergingMap<K, V>& other) const {
        return !(*this == other);
    }
    typename std::map<K, V>::iterator begin() {
        return data.begin();
    }
    typename std::map<K, V>::iterator end() {
        return data.end();
    }
};
