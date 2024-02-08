#include <stdio.h>
#include <stdlib.h>
#include <iso646.h>


void bubblesort(unsigned long nel, 
	int (* compare)(unsigned long i, unsigned long j),
	void (*swap)(unsigned long i, unsigned long j))
{
	int i = 0, left_bound = 0, right_bound = nel - 1;
	while (i < nel)
	{
		if(i % 2)
		{
			int j = right_bound;
			while (j > left_bound)
			{
				if(compare(j,j - 1) < 0)
					swap(j, j - 1);
				j--;
			}
			left_bound++;
		}
		else
		{
			int j = left_bound;
			while (j < right_bound)
			{
				if(compare(j,j + 1) > 0)
					swap(j, j + 1);
				j++;

			}
			right_bound--;
		}
		i++;
	}
}