#include <stdio.h>
#include <iso646.h>
#include <math.h>

int fibel(int x)
{
	int f1 = 0, f2 = 1;
	while (f2 < x)
	{
		f2 = f2 + f1;
		f1 = f2 - f1;
	}
	return f1;
}

void shellsort(unsigned long nel,
	int (*compare)(unsigned long i, unsigned long j),
	void (*swap)(unsigned long i, unsigned long j))
{
	for(int step = fibel(nel); step > 0; step = fibel(step))
		for(int i = 0; i < nel - step; i++)
			for(int j = i; j >= 0 and compare(j, j + step) > 0; j -= step)
				swap(j, j + step);
}

int main()
{
	return 0;
}