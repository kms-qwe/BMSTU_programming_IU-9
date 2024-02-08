#include <iso646.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

struct date 
{
	int day, month, year;
};

int keys(struct date d, int ind_sort_dmy)
{
	if (ind_sort_dmy == 0)
		return d.day - 1;
	if (ind_sort_dmy == 1)
		return d.month - 1;
	return d.year - 1970;
}

void distribution_sort(int(*keys)(struct date d, int ind_sort_dmy),
 int ind_sort_dmy, struct date* arr, int len, int len_count)
{
	int count[len_count];
	for(int i = 0; i < len_count; i++)
		count[i] = 0;

	for(int i = 0; i < len; i++)
	{
		count[keys(arr[i], ind_sort_dmy)]++;
	}

	for(int i = 1; i < len_count; i++)
	{
		count[i] += count[i - 1];
	}

	struct date arr_new[len];

	for(int j = len - 1; j >= 0; j--)
	{
		int k = keys(arr[j], ind_sort_dmy);
		int i = count[k] - 1;
		count[k] = i;
		arr_new[i] = arr[j];
	}
	for(int i = 0; i < len; i++)
		arr[i] = arr_new[i];
}
void radix_sort(struct date* arr, 
	int(*keys)(struct date d, int ind_sort_dmy), int len) 
{
	int len_count[3] = {31, 12, 61};
	for(int i = 0; i < 3; i++)
	{
		distribution_sort(keys, i, arr, len, len_count[i]);
	}
}

int main()
{
	int n;
	scanf("%d", &n);
	struct date arr[n];
	for(int i = 0; i < n; i++)
		scanf("%d %d %d", &arr[i].year, &arr[i].month, &arr[i].day);

	radix_sort(arr, keys, n);

	// printf("Sorted array:\n");
	for(int i = 0; i < n; i++)
		printf("%04d %02d %02d\n", arr[i].year, arr[i].month, arr[i].day);
}