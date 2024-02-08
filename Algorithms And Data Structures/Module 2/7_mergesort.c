#include <stdio.h>
#include <iso646.h>
#include <stdlib.h>
#include <math.h>

void merge(int l, int m, int r, int arr[])
{
	int *tmp = (int *)malloc(sizeof(int)*(r - l + 1));
	int i = l, j = m + 1, t = 0;
	while (i <= m and j <= r)
	{
		if (abs(arr[i]) <= abs(arr[j]))
		{
			tmp[t] = arr[i];
			t++;
			i++;
		}
		else
		{
			tmp[t] = arr[j];
			t++;
			j++;
		}

	}
	while (i<=m)
	{
		tmp[t] = arr[i];
		t++;
		i++;
	}
	while (j <= r)
	{
		tmp[t] = arr[j];
		t++;
		j++;
	}
	for(int k = l; k <= r; k++)	
		arr[k] = tmp[k-l];
	free(tmp);
}

void insertion_sort(int l, int r, int arr[])
{
	for(int i = l + 1; i <= r; i++)
	{
		int j = i;
		while (j > l and abs(arr[j]) < abs(arr[j-1]))
		{
			int tmp = arr[j];
			arr[j] = arr[j-1];
			arr[j-1] = tmp; 
			j--;
		}
	}
}
void merge_sort(int l, int r, int arr[])
{
	if (r - l + 1 >= 5)
	{
		int m = (l + r)/2;
		merge_sort(l, m, arr);
		merge_sort(m + 1, r, arr);
		merge(l, m, r, arr);
	}
	else
	{	
		insertion_sort(l, r, arr);
	}
}

int main()
{
	int n;
	scanf("%d", &n);
	int arr[n];
	for(int i = 0; i < n; i++)
	{
		scanf("%d", &arr[i]);
	}
	merge_sort(0, n - 1, arr);
	// insertion_sort(0, n - 1, arr);
	// printf("Sorted array:\n");
	for(int i = 0; i < n; i++)
		printf("   %d", arr[i]);
	printf("\n");
	return 0;
}