#include <stdio.h>

int lab_func(int *p, int n) {
    int arr[8];
    int i = 0;
    int sum = 0;

    while (i < 8) {
        arr[i] = i * 3 + n;
        i = i + 1;
    }

    i = 0;

    while (i < n) {
        if ((i & 1) == 0) {
            sum = sum + arr[i & 7];
        } else {
            sum = sum + p[i];
        }
        i = i + 1;
    }

    if (sum > 100) {
        sum = sum - *(p + 1);
    } else {
        sum = sum + arr[2];
    }

    return sum;
}

int main() {
    int data[10] = {5, 7, 2, 8, 1, 4, 6, 3, 9, 10};
    int r = lab_func(data, 6);
    printf("%d\n", r);
    return 0;
}
