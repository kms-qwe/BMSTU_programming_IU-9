#include <stdio.h>
#include <stdlib.h>
#include <string.h>

char *concat(char **s, int n)
{
    int N = 0;
    for(int i = 0; i < n; i++)
        N += strlen(s[i]);
    char *res = (char*)malloc(N+1);
    int top = 0;
    for(int i = 0; i < n; i++)
    {
        for(int j =0; (j < N) && (s[i][j] != '\0'); j++)
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
    scanf("%d ", &n);
    char *s[n];

    if (n == 0)
        return 0;
    for(int i = 0; i < n; i++)
    {
        s[i] = (char*)malloc(1100);
    }

    for(int i = 0; i < n; i++)
    {
        fgets(s[i], 1100, stdin);
        if (s[i][strlen(s[i]) - 1] == '\n')
            s[i][strlen(s[i]) - 1] = '\0';
    }
    char *res = concat(s,n);
    printf("\n%s\n", res);
    free(res);
    for(int i = 0; i< n; i++)
    {
        free(s[i]);
    }
    return 0;
}