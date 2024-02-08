#include <stdio.h>
#include <stdbool.h>
#include <iso646.h>
#include <stdlib.h>
#include <math.h>

int main()
{
	long x;
	scanf("%ld", &x);
	if (x < 0)
		x = - x;
	long long k = 2;
	while (k * k <= (long long)x)
	{
		if(x % k == 0)
			break;
		k++;
	}
	if(k * k > (long long)x)
	{
		printf("%ld\n", x);
		return 0;
	}
	// printf("TIME\n");
	int divisors[100];
	int top = 0;
	int d = 2;
	while (x != 1)
	{
		while (x % d == 0)
		{
			divisors[top] = d;
			x /= d;
			top += 1;
		}
		d ++;
	}
	// for(int i = 0; i < top; i++)
	// 	printf("%d\n", divisors[i]);
	// printf("TIME\n");
	int len = 0;
	for(int i = 0; i < top; i++)
		if (divisors[i] > len)
			len = divisors[i];
	len += 1;
	char *sieve = calloc(len, sizeof(char));
	// printf("TIME2\n");
	sieve[0] = sieve[1] = 1;
	for(int i = 2; i*i <= len; i++)
		if(sieve[i] == 0)
			for(int j = i*i; j < len; j += i)
				sieve[j] = 1;	

	long ans = 0;	
	for(int i = 0; i < top; i++)
		if (divisors[i] != x and sieve[divisors[i]] == 0 and ans < divisors[i])
			ans = divisors[i];

	printf("%ld\n", (ans == 0 ? x : ans));
	free(sieve);
	return 0;


}