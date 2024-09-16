
#include <iostream>
#include <vector>
#include <string>
#include "segment_tree.h"

using namespace std;
int sum(int a, int b) {
    return a + b;
}
int main() {
    vector<int> arr = {1, 3, 5, 7, 9, 11};
    SegmentTree<int, 6> seg_tree(arr, sum);

    cout << "Sum from index 1 to 3: " << seg_tree.query(1, 3) << endl; 

    seg_tree.update(2, 10);
    cout << "Sum from index 1 to 3 after updating index 2: " << seg_tree.query(1, 3) << endl; 

    std::string arr_str[5] = {"Hello", ", ", "world", "!", "\n"};

    SegmentTree<std::string, 5> seg_tree2(arr_str);

    std::cout << "Query result: " << seg_tree2.query(0, 4) << std::endl; 

    seg_tree2.update(3, "?");
    std::cout << "After update: " << seg_tree2.query(0, 4) << std::endl;
    return 0;
}
 