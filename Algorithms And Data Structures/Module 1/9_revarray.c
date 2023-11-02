#include <stdio.h>
#include <string.h>
#include <stdlib.h>
void revarray(void *base, size_t nel, size_t width);
int main(int argc, char const *argv[])
{
	int n;
	scanf("%d", &n);
	int arr[n];
	for (int i = 0; i < n; i++)
		scanf("%d", arr + i);
	printf("\n");
	revarray(arr, n, sizeof(int));
	for (int i = 0; i < n; i++)
		printf("%d\n", *(arr + i));
	return 0;
}

void revarray(void *base, size_t nel, size_t width)
{
	for(int i = 0; i < nel / 2;++i)
	{
		char *tmp = (char *)malloc(width);
		char *left = (char *)base + i*width;
		char *right = (char *)base + (nel - i - 1)*width;
		for(int j = 0; j < width; j++)
		{
			memcpy(tmp + j,right + j, 1);
			memcpy(right + j, left + j, 1);
			memcpy( left+j, tmp+j, 1);
		}
		free(tmp);
	}


}