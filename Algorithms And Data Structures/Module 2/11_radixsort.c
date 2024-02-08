#include <stdio.h>
#include <stdlib.h>
#include <iso646.h>

union Int32 {
	int x;
	unsigned char bytes[4];
};

int keys(union Int32 var, int sort_index)
{
	if (sort_index == 0)
		return var.bytes[0];
	if (sort_index == 1)
		return var.bytes[1];
	if (sort_index == 2)
		return var.bytes[2];
	return var.bytes[3] ^ 128;
}

void distribution_sort(int(*keys)(union Int32 var, int sort_index),
	union Int32 *arr, int len, int sort_index)
{
	int count[256];
	for(int i = 0; i < 256; i++) 
		count[i] = 0;

	for(int i = 0; i < len; i++)
		count[keys(arr[i], sort_index)]++;

	for(int i = 1; i < 256; i++)
		count[i] += count[i - 1];

	union Int32 sort_arr[len];
	
	for(int j = len - 1; j>= 0; j--)
	{
		int k = keys(arr[j], sort_index);
		int i = count[k] - 1;
		count[k] = i;
		sort_arr[i] = arr[j];
	}

	for(int i = 0; i < len; i++)
		arr[i] = sort_arr[i];
}

void radix_sort(int(*keys)(union Int32 var, int sort_index),
	union Int32 *arr, int len)
{
	for(int i = 0; i <= 3; i++)
		distribution_sort(keys, arr, len, i);
}

int main()
{
	int n;
	scanf("%d", &n);
	union Int32 *arr = calloc(n, sizeof(int));
	for(int i = 0; i < n; i++)
	{
		scanf("%d", &arr[i].x);
	}

	radix_sort(keys, arr, n);
	// printf("Sorted array:\n");
	for(int i = 0; i < n; i ++)
	{
		printf("%d\n", arr[i].x);
	}
	free(arr);
	return 0;
}