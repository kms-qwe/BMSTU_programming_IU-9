#include <stdio.h>
#include <stdlib.h>
#include <iso646.h>

void select_sort(int l, int r, int *arr)
{
	int j = r;
	while (j > l)
	{
		int k = j;
		int i = j - 1;
		while (i >= l)
		{
			if (arr[k] < arr[i])
				k = i;
			i--;;
		}
		int tmp = arr[j];
		arr[j]  = arr[k];
		arr[k] = tmp;
		j--;
	} 
}

int partition(int low, int high, int *arr)
{
	int i = low, j = low;
	while (j < high)
	{
		if (arr[j] < arr[high])
		{
			int tmp = arr[j];
			arr[j] = arr[i];
			arr[i] = tmp;
			i++;
		}
		j++;
	}
	int tmp = arr[i];
	arr[i] = arr[high];
	arr[high] = tmp;
	return i;
}
void quick_sort(int low, int high, int m, int *arr)
{
	if ((high - low + 1 >= m) and (low < high))
	{
		int q = partition(low, high, arr);
		quick_sort(low, q - 1, m, arr);
		quick_sort(q + 1, high, m, arr);
	}
	if ((high - low + 1 < m) and (low < high))
	{
		select_sort(low, high, arr);

	}
}
int main()
{
	int n,m;
	scanf("%d %d", &n, &m);
	int arr[n];
	for(int i = 0; i < n; i++) 
		scanf("%d", &arr[i]);
	quick_sort(0, n - 1, m, arr);
	for(int i = 0; i < n; i++)
		printf("%d ", arr[i]);
	printf("\n");
	return 0;
}