#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <iso646.h>
#include <stdbool.h>
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

void pword(char *s, char *t)
{
	int len_s = strlen(s), len_t = strlen(t), len_q = len_s + len_t + 1;
	char *q = calloc(len_q + 1, sizeof(char));
	strcat(q, s);
	strcat(q, "#");
	strcat(q, t);
	int *pi = prefix(q);

	for(int i = len_s + 1; i < len_q; i++)
		if (!pi[i])
		{
			printf("no\n");
			free(pi);
			free(q);
			return;
		}
	printf("yes\n");
	free(pi);
	free(q);
}


int main(int argc, char *argv[])
{	
    if (argc == 1)
    {
        printf("Missing argument\n");
        return 0;
    }

    pword(*(argv + 1), *(argv + 2));
	return 0;
}