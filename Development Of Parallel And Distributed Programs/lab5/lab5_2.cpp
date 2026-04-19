#include <iostream>
#include <thread>
#include <random>
#include <vector>
#include <shared_mutex>

using namespace std;

class ListNode {
public:
    int value;
    ListNode* next;

    explicit ListNode(int v) : value(v), next(nullptr) {}
};

class IntList {
private:
    ListNode* head;
    mutable shared_mutex rw_mutex;

    bool has_unlocked(int v) const {
        ListNode* cur = head;
        while (cur) {
            if (cur->value == v) {
                return true;
            }
            cur = cur->next;
        }
        return false;
    }

public:
    IntList() : head(nullptr) {}

    bool has(int v) const {
        shared_lock<shared_mutex> lock(rw_mutex);
        return has_unlocked(v);
    }

    void push_back(int v) {
        {
            shared_lock<shared_mutex> rlock(rw_mutex);
            if (has_unlocked(v)) {
                return;
            }
        }

        {
            unique_lock<shared_mutex> wlock(rw_mutex);
            if (has_unlocked(v)) {
                return;
            }

            auto* node = new ListNode(v);
            if (!head) {
                head = node;
            } else {
                ListNode* cur = head;
                while (cur->next) {
                    cur = cur->next;
                }
                cur->next = node;
            }
        }
    }

    void print() const {
        shared_lock<shared_mutex> lock(rw_mutex);
        ListNode* cur = head;
        while (cur) {
            cout << cur->value << " ";
            cur = cur->next;
        }
        cout << endl;
    }

    ~IntList() {
        ListNode* cur = head;
        while (cur) {
            ListNode* nxt = cur->next;
            delete cur;
            cur = nxt;
        }
    }
};

void producer(int count, IntList& list) {
    random_device rd;
    mt19937 gen(rd());
    uniform_int_distribution<int> dist(0, 1000);

    for (int i = 0; i < count; ++i) {
        int v = dist(gen);
        list.push_back(v);
    }
}

int main() {
    int threads_count = 4;
    int values_per_thread = 4;

    IntList list;
    vector<thread> workers;
    workers.reserve(threads_count);

    for (int i = 0; i < threads_count; ++i) {
        workers.emplace_back(producer, values_per_thread, ref(list));
    }

    for (auto& t : workers) {
        t.join();
    }

    cout << "Список: ";
    list.print();

    return 0;
}
