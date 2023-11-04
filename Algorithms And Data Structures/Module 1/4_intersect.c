#include <stdio.h>
#include <math.h>
int main()
{
    int a = 0, b = 0, len_x, xi, ans;
    scanf("%d", &len_x);
    for (int i = 0; i < len_x; i++)
    {
        scanf("%d", &xi);
        a = a | (int)pow(2, xi);
    }
    scanf("%d", &len_x);
    for (int i = 0; i < len_x; i ++)
    {
        scanf("%d", &xi);
        b = b | (int)pow(2,xi);
    }
    ans = a & b;
    for (int i = 0; i < 32; i++)
    {   
        if ((int)pow(2,i) == (ans & (int)pow(2,i)))
        {
            printf("%d ", i); 
        }
    }
}
