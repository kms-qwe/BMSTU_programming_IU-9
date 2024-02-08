#include <stdio.h>
#include <stdlib.h>
#include <iso646.h>
#include <string.h>
#include <math.h>

long long fib(int n)
{
	float fi1 = (1 + sqrt(5))/2;
	float fi2 = (1 - sqrt(5))/2;

	return (long long)round( (pow(fi1, n) - pow(fi2, n)) / sqrt(5));
}

char *fibstr(int n)
{
	long long len = fib(n) + 1;
	char *s1 = calloc(len, sizeof(char));
	char *s2 = calloc(len, sizeof(char));
	s1[0] = 'a';
	s2[0] = 'b';

	for(int i = 2; i < n; i += 2)
	{
		strcat(s1, s2);
		if (n - i == 1)
			break;
		strcat(s2, s1);
	}

	if (n % 2 == 1)
	{
		free(s2);
		return s1;
	}
	else
	{
		free(s1);
		return s2;
	}

}
int main()
{
	int n;
	scanf("%d", &n);
	char* res = fibstr(n);
	printf("%s\n", res);
	free(res);
	return 0;
}