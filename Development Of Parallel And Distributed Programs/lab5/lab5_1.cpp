#include <iostream>
#include <vector>
#include <thread>
#include <random>
#include <chrono>
#include <mutex>
#include <condition_variable>

using namespace std;

using Matrix = vector<vector<int>>;

struct Mailbox {
    mutex m;
    condition_variable cv;
    bool has_request = false;
    bool has_response = false;
    int requested_row = -1;
    vector<int> row_data;
};

mutex barrier_mutex;
condition_variable barrier_cv;
int barrier_count = 0;
int barrier_phase = 0;

vector<int> start_row;
vector<int> end_row;
vector<Mailbox> mailboxes;

void barrier(int total_threads) {
    unique_lock<mutex> lock(barrier_mutex);
    int phase = barrier_phase;
    ++barrier_count;
    if (barrier_count == total_threads) {
        barrier_count = 0;
        ++barrier_phase;
        barrier_cv.notify_all();
    } else {
        barrier_cv.wait(lock, [&] { return barrier_phase != phase; });
    }
}

int neighborsAliveSingle(const Matrix& grid, int x, int y, int rows, int cols) {
    static const int dirs[8][2] = {
        {-1, -1}, {-1, 0}, {-1, 1},
        { 0, -1},          { 0, 1},
        { 1, -1}, { 1, 0}, { 1, 1}
    };
    int alive = 0;
    for (auto& d : dirs) {
        int nx = (x + d[0] + rows) % rows;
        int ny = (y + d[1] + cols) % cols;
        alive += grid[nx][ny];
    }
    return alive;
}

void stepCellSingle(const Matrix& cur, Matrix& next, int x, int y, int rows, int cols) {
    int alive = neighborsAliveSingle(cur, x, y, rows, cols);
    if (cur[x][y] == 1) {
        next[x][y] = (alive == 2 || alive == 3) ? 1 : 0;
    } else {
        next[x][y] = (alive == 3) ? 1 : 0;
    }
}

int ownerForRow(int row, int num_threads) {
    for (int id = 0; id < num_threads; ++id) {
        if (row >= start_row[id] && row < end_row[id]) return id;
    }
    return 0;
}

void serveRequests(int id, const Matrix& cur, int cols) {
    Mailbox& mb = mailboxes[id];
    unique_lock<mutex> lock(mb.m);
    if (mb.has_request && !mb.has_response) {
        int r = mb.requested_row;
        mb.row_data.assign(cols, 0);
        for (int j = 0; j < cols; ++j) {
            mb.row_data[j] = cur[r][j];
        }
        mb.has_response = true;
        mb.cv.notify_all();
    }
}

vector<int> requestRow(int requester_id, int row, int cols, int num_threads) {
    int owner = ownerForRow(row, num_threads);
    Mailbox& mb = mailboxes[owner];

    unique_lock<mutex> lock(mb.m);
    mb.cv.wait(lock, [&] { return !mb.has_request && !mb.has_response; });

    mb.requested_row = row;
    mb.has_request = true;
    mb.has_response = false;
    mb.cv.notify_all();

    mb.cv.wait(lock, [&] { return mb.has_response; });

    vector<int> result = mb.row_data;
    mb.has_request = false;
    mb.has_response = false;
    mb.cv.notify_all();

    return result;
}

int neighborsAliveThread(const Matrix& cur,
                         int x, int y,
                         int rows, int cols,
                         int thread_id,
                         int num_threads) {
    static const int dirs[8][2] = {
        {-1, -1}, {-1, 0}, {-1, 1},
        { 0, -1},          { 0, 1},
        { 1, -1}, { 1, 0}, { 1, 1}
    };

    int alive = 0;
    for (auto& d : dirs) {
        int nx0 = x + d[0];
        int ny = (y + d[1] + cols) % cols;

        int nx;
        if (nx0 < 0) {
            nx = rows - 1;
        } else if (nx0 >= rows) {
            nx = 0;
        } else {
            nx = nx0;
        }

        if (nx >= start_row[thread_id] && nx < end_row[thread_id]) {
            alive += cur[nx][ny];
        } else {
            vector<int> row_data = requestRow(thread_id, nx, cols, num_threads);
            alive += row_data[ny];
        }
    }

    return alive;
}

void stepCellThread(const Matrix& cur, Matrix& next,
                    int x, int y,
                    int rows, int cols,
                    int thread_id,
                    int num_threads) {
    int alive = neighborsAliveThread(cur, x, y, rows, cols, thread_id, num_threads);
    if (cur[x][y] == 1) {
        next[x][y] = (alive == 2 || alive == 3) ? 1 : 0;
    } else {
        next[x][y] = (alive == 3) ? 1 : 0;
    }
}

void worker(int id, const Matrix& cur, Matrix& next, int rows, int cols, int num_threads) {
    int start = start_row[id];
    int finish = end_row[id];

    for (int i = start; i < finish; ++i) {
        serveRequests(id, cur, cols);
        for (int j = 0; j < cols; ++j) {
            stepCellThread(cur, next, i, j, rows, cols, id, num_threads);
        }
    }

    serveRequests(id, cur, cols);
    barrier(num_threads);
}

int main() {
    const int rows = 200;
    const int cols = 200;
    const int num_threads = 4;
    const int num_steps = 5;

    Matrix current(rows, vector<int>(cols));
    Matrix next(rows, vector<int>(cols));

    start_row.assign(num_threads, 0);
    end_row.assign(num_threads, 0);
    mailboxes.assign(num_threads, Mailbox{});

    int base = rows / num_threads;
    int rem = rows % num_threads;
    int cur_start = 0;
    for (int i = 0; i < num_threads; ++i) {
        int block = base + (i < rem ? 1 : 0);
        start_row[i] = cur_start;
        end_row[i] = cur_start + block;
        cur_start += block;
    }

    random_device rd;
    mt19937 gen(rd());
    uniform_int_distribution<int> dist(0, 1);

    for (int i = 0; i < rows; ++i) {
        for (int j = 0; j < cols; ++j) {
            current[i][j] = dist(gen);
        }
    }

    auto run_with_threads = [&](Matrix& a, Matrix& b) {
        auto t1 = chrono::high_resolution_clock::now();
        for (int step = 0; step < num_steps; ++step) {
            barrier_count = 0;
            barrier_phase = 0;

            vector<thread> pool;
            pool.reserve(num_threads);
            for (int i = 0; i < num_threads; ++i) {
                pool.emplace_back(worker, i, cref(a), ref(b), rows, cols, num_threads);
            }
            for (auto& th : pool) {
                th.join();
            }
            swap(a, b);
        }
        auto t2 = chrono::high_resolution_clock::now();
        chrono::duration<double> d = t2 - t1;
        return d.count() / num_steps;
    };

    auto run_single = [&](Matrix& a, Matrix& b) {
        auto t1 = chrono::high_resolution_clock::now();
        for (int step = 0; step < num_steps; ++step) {
            for (int i = 0; i < rows; ++i) {
                for (int j = 0; j < cols; ++j) {
                    stepCellSingle(a, b, i, j, rows, cols);
                }
            }
            swap(a, b);
        }
        auto t2 = chrono::high_resolution_clock::now();
        chrono::duration<double> d = t2 - t1;
        return d.count() / num_steps;
    };

    double avg_parallel = run_with_threads(current, next);
    cout << "Среднее время одного шага с потоками: " << avg_parallel << " секунд" << endl;

    double avg_single = run_single(current, next);
    cout << "Среднее время одного шага без потоков: " << avg_single << " секунд" << endl;

    return 0;
}
