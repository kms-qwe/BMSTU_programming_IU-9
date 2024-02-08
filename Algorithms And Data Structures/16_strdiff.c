#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <iso646.h>
#include <math.h>
int strdiff(char *a, char *b)
{
	int ind = 0, len_a = strlen(a), len_b = strlen(b), pow2;
	// printf("%d %d %d\n", len_a, len_b, (len_a > len_b) ? len_b : len_a);
	for(int i = 0; i <= ((len_a > len_b) ? len_b : len_a); i++)
	{
		for(int j = 0; j < 8; j++)
		{

			pow2 = (int)pow(2,j);
			// printf("%d %d %d %c %c\n", pow2 & a[i], pow2 & b[i], pow2, a[i], b[i]);
			if ((pow2 & a[i]) != (pow2 & b[i]))
				return ind;
			ind += 1;
		}

	}

	return (len_a == len_b) ? -1 : ind + 1;
}
int main()
{
	char s1[1000], s2[1000];
	fgets(s1, 1000, stdin);
	fgets(s2, 1000, stdin);
	if (s1[strlen(s1) - 1] == '\n')
		s1[strlen(s1) - 1] = '\0';
	if (s2[strlen(s2) - 1] == '\n')
		s2[strlen(s2) - 1] = '\0';
	// printf("|%s|%s|\n", s1, s2);
	printf("%d\n", strdiff(s1, s2));
	return 0;
}