#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int maxarray(void *base, size_t nel, size_t width, int (*compare)(void *a, void *b))
{
	void *max = (void *)malloc(width);
	memcpy(max, base, width);
	int maxind = 0;
	for (int i = 1; i < nel; ++i)
	{
		void *ai = (char *)base + i*width;
		if (compare(ai, max) > 0)
		{
			memcpy(max,ai, width);
			maxind = i;
		}
	}
	free(max);
	return maxind;
}