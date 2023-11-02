#include <stdio.h>
#include <stdlib.h>
#include <string.h>

char *concat(char **s, int n)
{
    char *res = (char*)malloc(n*2000);
    int top = 0;
    for(int i = 0; i < n; i++)
    {
        for(int j =0; (j < 10000) && (s[i][j] != '\0'); j++)
        {
            memcpy((void *)(res + top), (const void*)&s[i][j], 1);
            top += 1;
        }
    }
    res[top] = '\0';
    return res;
}
int main()
{
    int n;
    scanf("%d", &n);
    char *s[n];

    for(int i = 0; i < n; i++)
    {
        s[i] = (char*)malloc(2000);
    }

    fflush(stdin);
    for(int i = 0; i < n; i++)
    {
        gets(s[i]);
    }

 
    printf("\n%s\n", concat(s, n));
    for(int i = 0; i< n; i++)
    {
        free(s[i]);
    }
    return 0;
}