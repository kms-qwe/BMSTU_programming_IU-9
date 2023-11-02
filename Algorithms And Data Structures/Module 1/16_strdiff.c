#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <iso646.h>
#include <math.h>
int strdiff(char *a, char *b)
{
	int ind = 0, len_a = strlen(a), len_b = strlen(b), pow2, ls;
	for(int i = 0; (i < (len_a > len_b) ? len_b : len_a) && a[i] != '\0' && b[i] != '\0' && a[i] != '\n' && b[i] != '\n'; i++)
	{
		for(int j = 0; j < 8; j++)
		{
			pow2 = (int)pow(2,j);
			if ((pow2 & a[i]) != (pow2 & b[i]))
				return ind;
			ind += 1;
		}
		// printf("%c\t%c\n", a[i], b[i]);
		ls = 0;
	}
	// printf("%d\n", ls);
	// printf("%d %d\n", len_a, len_b);
	return ((a[ls] == b[ls]) && (len_a == len_b)) ? -1 : ind;
}
int main()
{
	char s1[100], s2[10];
	fflush(stdin);
	fgets(s1, 100, stdin);
	fflush(stdin);
	fgets(s2, 10, stdin);
	// printf("%s", s1);
	// printf("%s", s2);
	// printf("123\n123\n");
	printf("%d\n", strdiff(s1, s2));
	return 0;
}