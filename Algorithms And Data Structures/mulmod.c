#include <stdio.h>
#include <math.h>

long long gorner_mod (long long a, long long b, long long m)
{
    long long ans;
    ans = (b / (long long)pow(2, 63)) * a % m;
    b = b % (long long)pow(2,63);
    for (int i = 62; i >= 0; i--)
    {
        ans = (ans % m * 2 % m) + (b / (long long)pow(2,i)) * a % m;
        b = b % (long long)pow(2, i);
    } 
    return ans;
}

int main()
{
    long long a, b, m;
    scanf("%lld%lld%lld", &a, &b, &m);
    printf("%lld", gorner_mod(a,b,m));
    return 0;
}