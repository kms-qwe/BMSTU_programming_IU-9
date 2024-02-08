#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <iso646.h>
#include <stdbool.h>

int *suffix(char *s)
{
	int len_s = strlen(s), t = len_s - 1;
	int *sigma = calloc(len_s, sizeof(int));
	sigma[len_s - 1] = t;

	for(int i = len_s - 2; i >=0; i--)
	{
		while (t < len_s - 1 and s[t] !=  s[i])
			t = sigma[t + 1];
		if (s[t] == s[i])
			t--;
		sigma[i] = t;
	}
	// for(int i = 0; i < len_s; i++)
	// 	printf("%d %d\n", i, sigma[i]);
	return sigma;
}

int *delta1(char *s, int size)
{
	int len_s = strlen(s);
	int *delta = calloc(size, sizeof(int));
	for(int i = 0; i < size; i++)
		delta[i] = len_s;
	for(int i = 0; i < len_s; i++)
		delta[ (int)s[i] - 33 ] = len_s - i - 1;
	return delta;
}
int *delta2(char *s)
{
	int* sigma = suffix(s);
	int len_s = strlen(s), t = sigma[0];

	int *delta_2 = calloc(len_s, sizeof(int));
	for(int i = 0; i < len_s; i++)
	{
		while (t < i)
			t = sigma[t + 1];
		delta_2[i] = -i + t + len_s;
	}
	for(int i = 0; i < len_s - 1; i++)
	{
		t = i;
		while (t < len_s - 1)
		{
			t = sigma[t + 1];
			if (s[i] != s[t])
				delta_2[t] = -(i + 1) + len_s;
		}
	}
	free(sigma);
	return delta_2;
}

void simple_BMS(char *s, int size, char *t)
{
	int *delta_1 = delta1(s, size);
	int len_s = strlen(s), k = len_s - 1, len_t = strlen(t);
	while (k < len_t)
	{
		int i = len_s - 1;
		while (t[k] == s[i])
		{
			if (i == 0)
			{
				printf("%d\n", k);
				// free(delta_1);
				// return;
				break;
			}
			i--;
			k--;

		}
		k += (delta_1[ (int)t[k] - 33] >= len_s - i ? delta_1[ (int)t[k] - 33] : len_s - i);
	}
	free(delta_1);
}

void BMS(char *s, int size, char *t)
{
	int *delta_1 = delta1(s, size);
	int *delta_2 = delta2(s);
	int len_s = strlen(s), len_t = strlen(t), k = len_s - 1;
	while (k < len_t)
	{
		int i = len_s - 1;
		while (t[k] == s[i])
		{
			if (i == 0)
			{
				// free(delta_1);
				// free(delta_2);
				printf("%d\n", k);
				// return;
				break;
			}
			i--;
			k--;
		}
		k += (delta_1[ (int)t[k] -33] >= delta_2[i] ? 
			delta_1[ (int)t[k] -33] : delta_2[i]);
	}
	free(delta_1);
	free(delta_2);
}

int main(int argc, char *argv[])
{
    if (argc == 1)
    {
        printf("Missing argument\n");
        return 0;
    }
    int size = 126 - 33 + 1;
    // simple_BMS(*(argv + 1), size, *(argv + 2));
    // printf("_____________\n");
    BMS(*(argv + 1), size, *(argv + 2));
	return 0;
}