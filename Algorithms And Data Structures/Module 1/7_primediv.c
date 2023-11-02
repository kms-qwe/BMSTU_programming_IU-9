#include <stdio.h>
#include <stdbool.h>
#define N 65536
int main(int argc, char const *argv[])
{
	int x;
	scanf("%d", &x);
	int sieve[N] = {0};
	sieve[0] = 1;
	sieve[1] = 1;
	for (int i = 2; i*i <= N; i++)
	{
		if (sieve[i] == 0)
		{
			for(int j = i*i; j < N; j += i)
			{
				sieve[j] = 1;
			}
		}
	}

	bool flag = true;
	for(int i = N - 1; i >= 0; i--)
	{
		if(sieve[i] == 0 && x % i == 0)
		{
			printf("%d\n", i);
			flag = false;
			break;
		}
	}
	if (flag)
		printf("%d\n", x);


	return 0;
}