#include <stdio.h>
#include <stdlib.h>
#include <iso646.h>
#include <string.h>
int key(const char c)
{
	return (int)(c - 'a');
}

void distribution_sort(int(*key)(const char c), char* s)
{
	int len_s = strlen(s);
	int count[26] = {0};

	for(int i = 0; i < len_s; i++)
		count[key(s[i])]++;

	for(int i = 1; i < 26; i++)
		count[i] += count[i-1];

	char* snew = (char*)calloc(len_s + 1, 1);
	for(int j = len_s - 1; j >= 0; j--)
	{
		int k = key(s[j]);
		int i = count[k] - 1;
		count[k] = i;
		snew[i] = s[j];
	}
	// printf("|%s|\n|%s|\n%lu\t%lu\n", s, snew, strlen(s), strlen(snew));
	strcpy(s, snew);
	free(snew);
}

int main()
{
	char* s = (char*)calloc(1000100, 1);
	fgets(s, 1000100, stdin);
	s[strlen(s) - 1] = '\0';
	distribution_sort(key, s);
	printf("%s\n", s);
	free(s);
}