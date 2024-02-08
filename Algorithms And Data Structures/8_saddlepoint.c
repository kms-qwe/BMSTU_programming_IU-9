#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <limits.h>

bool find(int A, int* B, int len_B);
int any_in(int *A, int *B, int len_A, int len_B);

int main(int argc, char const *argv[])
{
	int N, M;
	scanf("%d%d", &N, &M);
	int arr[10][10], max_in_row[N], min_in_colomn[M];

	for(int i = 0; i < N; i++)
		for(int j = 0; j < M; j++)
			scanf("%d", &arr[i][j]);

	for(int i = 0; i < N; ++i)
		max_in_row[i] = INT_MIN;

	for(int i = 0; i < M; i++)
		min_in_colomn[i] = INT_MAX;

	for(int i = 0; i < N; ++i)
		for(int j = 0; j < M; ++j)
		{
			max_in_row[i] = (max_in_row[i] < arr[i][j]) ? arr[i][j] : max_in_row[i];
			// printf("%d %d %d %d %d\n", i, j, arr[i][j], min_in_colomn[j], min_in_colomn[j] > arr[i][j]);
			min_in_colomn[j] = (min_in_colomn[j] > arr[i][j]) ? arr[i][j] : min_in_colomn[j];
		}

	// for(int i = 0; i < N; i++)
	// 	printf("%d\t", max_in_row[i]);
	// printf("\n");
	// for(int i = 0; i < M; i++)
	// 	printf("%d\t", min_in_colomn[i]);
	// printf("\n");
	any_in(min_in_colomn, max_in_row, M, N);

	return 0;
}

bool find(int A, int* B, int len_B)
{
	for(int i = 0; i < len_B; i++)
		if (A == *(B + i))
			{
				printf("%d ", i);
				return true;
			}
	return false;
}

int any_in(int *A, int *B, int len_A, int len_B)
{
	for(int i = 0; i < len_A; ++i)
		if (find(*(A + i), B, len_B))
			{
				printf("%d\n", i);
				return 1;
			}
	printf("none\n");
	return 0; 
}