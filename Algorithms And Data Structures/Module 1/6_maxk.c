#include <stdio.h>

int main()
{
    int n, k;
    scanf("%d", &n);
    long long A[n], sum = 0, mx;
    for(int i = 0; i < n; i ++)
    {
        scanf("%lld", &A[i]);
    }
    scanf("%d", &k);
    for(int i = 0; i < k; i ++)
    {
        sum += A[i];
    }
    mx = sum;
    for(int i = k; i < n; i++)
    {
        sum = sum - A[i-k] + A[i];
        mx = (sum > mx) ? sum : mx;
    }
    printf("%lld", mx);
    return 0;


}