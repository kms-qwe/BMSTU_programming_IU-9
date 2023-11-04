#include <stdio.h>
#include <math.h>

unsigned long gorner_mod (unsigned long a, unsigned long b, unsigned long m)
{
    unsigned long ans;

    ans = (b / (unsigned long)pow(2, 63)) * a % m;
    b = b % (unsigned long)pow(2,63);
    for (int i = 62; i >= 0; i--)
    {
        ans = ((ans * 2 % m) + (b / (unsigned long)pow(2,i)) * a % m) % m;
        b = b % (unsigned long)pow(2, i);
    } 
    return ans;
}

int main()
{
    unsigned long a, b, m;
    scanf("%lu%lu%lu", &a, &b, &m);
    printf("%lu\n", gorner_mod(a,b,m));
    return 0;
}
