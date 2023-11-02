#include <stdio.h>
#include <stdlib.h>
#include <string.h>

unsigned long binsearch(unsigned long nel, int(*compare)(unsigned long i))
{
    unsigned long left = 0, right = nel - 1;
    while (left <= right)
    {
        unsigned long m = (left+right)/2;
        if (0 == compare(m))
            return m;
        if (1 == compare(m))
            right = m - 1;
        if (-1 == compare(m))
            left = m + 1;
        
    }
    return nel;
}
int comp(unsigned long i)
{   
    int se = 3;
    int arr[] = {0,2,4,6,8,10,12,14,16,18};
    if (arr[i] == se)
        return 0;
    if (arr[i] < se)
        return -1;
    return 1;
}

int main()
{
    printf("%lu\n", binsearch(10, comp));
    return 0;
}
