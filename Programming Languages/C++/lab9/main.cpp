#include <iostream>
#include "MergingMap.h"

int main() {
    MergingMap<std::string, int> map1(0);
    map1["a"] = 1;
    map1["b"] = 2;

    MergingMap<std::string, int> map2(0);
    map2["b"] = 3;
    map2["c"] = 4;

    MergingMap<std::string, int> result = map1 + map2;

    std::cout << "Result of merging map1 and map2:\n";
    for (const auto& pair : result) {
        std::cout << pair.first << ": " << pair.second << std::endl;
    }
    MergingMap<std::string, MergingMap<std::string, int>> nested_map(0);
    nested_map["outer1"]["inner1"] = 10;
    nested_map["outer1"]["inner2"] = 20;
    nested_map["outer2"]["inner1"] = 100;
    nested_map["outer2"]["intter1"] = 200;

    MergingMap<std::string, MergingMap<std::string, int>> nested_map2(0);
    nested_map2["outer1"]["inner2"] = 30;
    nested_map2["outer1"]["inner3"] = 40;

    MergingMap<std::string, MergingMap<std::string, int>> nested_result = nested_map + nested_map2;

    std::cout << "\nResult of merging nested_map and nested_map2:\n";
    for (auto& outer_pair : nested_result) {
        std::cout << outer_pair.first << ":\n";
        for (const auto& inner_pair : outer_pair.second) {
            std::cout << "  " << inner_pair.first << ": " << inner_pair.second << std::endl;
        }
    }

    return 0;
}
