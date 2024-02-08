#include <stdio.h>
#include <stdlib.h>
#include <iso646.h>
#include <string.h>

int *prefix(char *s)
{
	int len_s = strlen(s), t = 0;
	int *pi = calloc(len_s, sizeof(int));
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

int overlap(char *s, char *t)
{
	int len_s = strlen(s), len_t = strlen(t), len_q = len_s + len_t + 1;
	char *q = calloc(len_q + 1, sizeof(char));
	strcat(q, s);
	strcat(q, "#");
	strcat(q, t);
	int *pi = prefix(q);
	int res = pi[len_q - 1];
	free(pi);
	free(q);
	return res;
}


int supstr(char **s, int len_s)
{
	int lengths[len_s];
	for(int i = 0; i < len_s; i++)
		lengths[i] = strlen(s[i]);

	int **overlap_matrix = (int **)calloc(len_s, sizeof(int *));
	for(int i = 0; i < len_s; i++)
		overlap_matrix[i] = (int *)calloc(len_s, sizeof(int));

	for(int i = 0; i < len_s; i++)
		for(int j = 0; j < len_s; j++)
			overlap_matrix[i][j] = overlap(s[j], s[i]);

	// for(int i = 0; i < len_s; i++)
	// {
	// 	for(int j = 0; j < len_s; j++)
	// 		printf("m[%d][%d] = %d\t", i, j, overlap_matrix[i][j]);
	// 	printf("\n");
	// }

	char used_str[len_s];
	for(int i = 0; i < len_s; i++)
		used_str[i] = 0;
	int res = 1 << 30;
	void generation_supstr(int n, int cnt, int len, int prev)
	{
		// printf("%d %d %d\n", cnt, len, prev);
		if (cnt == n)
		{
			res = (res > len ? len : res);
			return;
		}

		for(int i = 0; i < len_s; i++)
			if (used_str[i] == 0)
			{
				used_str[i] = 1;
				if (prev != -1)
				{
					// printf("%d %d %d\n", len, lengths[i], overlap_matrix[prev][i]);
					generation_supstr(n, cnt + 1, len + lengths[i] - overlap_matrix[prev][i], i);
				}
				else
				{
					// printf("%d %d\n", len, lengths[i]);
					generation_supstr(n, cnt + 1, len + lengths[i], i);
				}
				used_str[i] = 0;
			}
	}
	generation_supstr(len_s, 0, 0, -1);
	for(int i = 0; i < len_s; i++)
		free(overlap_matrix[i]);
	free(overlap_matrix);
	return res;
}
int main()
{	
	int n;
	scanf("%d\n", &n);
	char *s[n];
	for(int i = 0; i < n; i++)
	{
		s[i] = calloc(1000, sizeof(char));
		scanf("%s", s[i]);
	}
	int res = supstr(s, n);
	// printf("RESULT: %d\n", res);
	printf("%d\n", res);

	for(int i = 0; i < n; i++)
		free(s[i]);
	return 0;
}