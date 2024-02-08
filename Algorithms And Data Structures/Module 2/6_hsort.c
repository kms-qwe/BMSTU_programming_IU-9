#include <stdio.h>
#include <stdlib.h>
#include <iso646.h>
#include <string.h>

void swap(void *a, void* b, size_t width)
{	
	void *tmp = (void*)malloc(width);
	memcpy(tmp, a, width);
	memcpy(a, b, width);
	memcpy(b, tmp, width);
	free(tmp);
}

void heapify(void *base, size_t width, int i, int n,
	int(*compare)(const void *a, const void *b))
{
	while(1)
	{	

		int l = 2*i + 1;
		int r = l + 1;
		int j = i;
		if (l < n and compare((char *)base + i*width, (char *)base + l*width) < 0)
			i = l;
		if (r < n and compare((char *)base + i*width, (char *)base + r*width) < 0)
			i = r;
		if (i == j)
			break;
		swap((char *)base + i*width, (char *)base + j*width, width);
	}
}

void buildheap(void *base, size_t nel, size_t width,
	int(*compare)(const void *a, const void *b))
{
	int i = nel/2 - 1;
	while(i >= 0)
	{
		heapify(base, width, i, nel, compare);
		i--;
	}
}

void hsort(void *base, size_t nel, size_t width,
	int (*compare)(const void *a, const void *b))
{
	buildheap(base, nel, width, compare);
	for(size_t i = nel - 1; i > 0; i--)
	{
		swap(base, (char *)base + i*width, width);
		heapify(base, width, 0, i, compare);
	}

}
int compare(const void* a, const void* b)
{
	char *start_a = *(char **)a;
	char *start_b = *(char **)b;
	int cntA = 0, cntB = 0;
	int i = 0;
	while (*(char*)(start_a +  i) != '\0')
	{
		if (*(char*)(start_a + i) == 'a')
			cntA++;
		i++;
	}
	i = 0;
	while (*(char*)(start_b + i) != '\0')
	{
		if(*(char*)(start_b + i) == 'a')
			cntB++;
		i++;
	}
	return cntA - cntB;

}
int main()
{
	size_t n;
	scanf("%lu\n", &n);
	char *s[n];
	for(size_t i = 0; i < n; i++)
	{
		s[i] = (char *)malloc(1001*sizeof(char));
		fgets(s[i], 1001, stdin);
	}
	hsort(s, n, sizeof(char*), compare);
	// printf("Sorted array:\n");
	for(size_t i = 0; i < n; i++)
	{
		printf("%s", s[i]);
		free(s[i]);
	}

}