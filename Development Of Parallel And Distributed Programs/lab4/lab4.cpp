#include <atomic>
#include <chrono>
#include <iomanip>
#include <iostream>
#include <mutex>
#include <random>
#include <thread>
#include <vector>

using std::atomic;
using std::cout;
using std::endl;
using std::lock_guard;
using std::mutex;
using std::ref;
using std::string;
using std::thread;
using std::vector;
using namespace std::chrono;

mutex                   log_mutex;
steady_clock::time_point start_time;

void log_state(int id, const string& state) {
    lock_guard<mutex> lock(log_mutex);

    const auto now = steady_clock::now();
    const auto ms  = duration_cast<milliseconds>(now - start_time).count();

    cout << "[" << std::setw(6) << ms << " мс] Философ " << id << " " << state << endl;
}

void philosopher(int id,
                 int N,
                 vector<mutex>& forks,
                 atomic<bool>& running)
{
    const int left  = id;
    const int right = (id + 1) % N;

    std::random_device rd;
    std::mt19937       gen(rd() ^ (id * 7919));

    std::uniform_int_distribution<int> think_dist(300, 1000);
    std::uniform_int_distribution<int> eat_dist(300, 800);
    std::uniform_int_distribution<int> retry_dist(50, 150);

    while (running.load(std::memory_order_relaxed)) {
        log_state(id, "размышляет");
        std::this_thread::sleep_for(milliseconds(think_dist(gen)));

        if (!running.load(std::memory_order_relaxed)) {
            break;
        }

        log_state(id, "проголодался и хочет есть");

        while (running.load(std::memory_order_relaxed)) {
            log_state(id, "пытается взять ЛЕВУЮ вилку");
            forks[left].lock();
            log_state(id, "взял ЛЕВУЮ вилку");

            log_state(id, "пытается взять ПРАВУЮ вилку");
            if (forks[right].try_lock()) {
                log_state(id, "взял ПРАВУЮ вилку");

                log_state(id, "ЕСТ");
                std::this_thread::sleep_for(milliseconds(eat_dist(gen)));

                log_state(id, "кладёт обе вилки на стол");
                forks[right].unlock();
                forks[left].unlock();
                break;
            } else {
                log_state(id, "не смог взять ПРАВУЮ вилку, кладёт ЛЕВУЮ вилку обратно");
                forks[left].unlock();

                std::this_thread::sleep_for(milliseconds(retry_dist(gen)));
            }
        }
    }

    log_state(id, "останавливается");
}

int main(int argc, char* argv[]) {
    int N                  = 5;
    int simulation_seconds = 1;

    if (argc >= 2) {
        N = std::stoi(argv[1]);
    }
    if (argc >= 3) {
        simulation_seconds = std::stoi(argv[2]);
    }

    start_time = steady_clock::now();

    vector<mutex> forks(N);
    atomic<bool>  running(true);

    vector<thread> threads;
    threads.reserve(N);

    for (int i = 0; i < N; ++i) {
        threads.emplace_back(philosopher, i, N, ref(forks), ref(running));
    }

    std::this_thread::sleep_for(seconds(simulation_seconds));

    running.store(false, std::memory_order_relaxed);

    for (auto& t : threads) {
        if (t.joinable()) {
            t.join();
        }
    }

    return 0;
}
