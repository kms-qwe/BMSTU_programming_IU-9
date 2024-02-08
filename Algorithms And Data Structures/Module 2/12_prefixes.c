#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <iso646.h>
int* prefix(char * s)
{
    size_t len_s = strlen(s);
    int *pi = calloc(len_s, sizeof(int));
    int t = 0;
    pi[0] = 0;
    for(int i = 1; i < len_s; i++)
    {
        while (t > 0 and s[t] != s[i])
            t = pi[t - 1];
        if (s[t] == s[i])
            t++;
        pi[i] = t;
    }
    return pi;
}

int main(int argc, char *argv[])
{

    if (argc == 1)
    {
        printf("Missing argument\n");
        return 0;
    }
    char *s = *(argv + 1);
    int len_s = strlen(s);

    int *pi = prefix(s);
    int pref_len = 0;

    for(int i = 0; i < len_s; i++)
    {
        pref_len = i + 1;
        if(pi[i] and !(pref_len % (pref_len - pi[i])))
            printf("%d %d\n", pref_len, pref_len / (pref_len - pi[i]));
    }
    free(pi);
    return 0;
}