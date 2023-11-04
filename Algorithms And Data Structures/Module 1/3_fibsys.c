#include <stdio.h>

void fib_sys_print (unsigned long x)
{
	if (x == 0)
	{
		printf("%d", 0);
	}
    int ans[100] = {0};
    int mx_ind = 0;
    while (x > 0)
    {
    unsigned long long fib_0[2] = {0, 1}, fib_1[2] = {1, 1}, fib_2[2] = {2, 2};
        while (fib_2[1] <= x)
        {
            fib_0[1] = fib_1[1]; fib_0[0] += 1;
            fib_1[1] = fib_2[1]; fib_1[0] += 1;
            fib_2[1] = fib_0[1] + fib_1[1]; fib_2[0] += 1;
        }
        ans[fib_1[0]] = 1;
        mx_ind = (fib_1[0] > mx_ind) ? fib_1[0] : mx_ind;
        x -= fib_1[1];

    }
    for (int i = mx_ind; i > 0; i--)
    {
        printf("%d", ans[i]);
    }

}

int main()
{
    unsigned long long x;
    scanf("%llu", &x);
    fib_sys_print(x);
    return 0;
}
