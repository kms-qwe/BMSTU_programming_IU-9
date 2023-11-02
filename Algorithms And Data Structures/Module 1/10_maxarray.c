#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int maxarray(void *base, size_t nel, size_t width, int (*compare)(void *a, void *b))
{
	void *max = (void *)malloc(width);
	int maxind = 0;
	for (int i = 1; i < nel; ++i)
	{
		void *ai = base + i*width;
		if (compare(ai, max) > 0)
		{
			memcpy(max,ai, width);
			maxind = i;
		}
	}
	free(max);
	return maxind;
}

int comp(void *a, void *b)
{
	return *((int *)a) - *((int *)b);
}
int main(int argc, char const *argv[])
{
	int arr[10] = {0, 1, 2, 3,5,43,123, 1232,4,1};
	printf("%d\n", maxarray(arr, 10, sizeof(int), comp));
	return 0;
}