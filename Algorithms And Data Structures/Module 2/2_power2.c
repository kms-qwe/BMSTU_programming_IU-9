#include <stdio.h>

int is_power2(int x)
{
	return (x > 0) && ((x & (x - 1)) == 0);
}

int cnt_comb_len_n(int *arr, int n, int len, int start, int sum)
{
	if (n > len - start)
		return 0;

	if (0 == n)
		return is_power2(sum);

	int cnt = 0;
	for (int i = start; i <  len; i++)
	{
		cnt += cnt_comb_len_n(arr, n - 1, len, i + 1, sum + *(arr + i));
	}
	return cnt;

}

int cnt(int *arr, int len)
{
	int cnt = 0;
	for(int comblen = 0; comblen < len; comblen++)
	{
		// for(int comblen = 1; comblen <= len - i; comblen++)
		// {

		// }
		//printf("%d\t%d\n", comblen, cnt_comb_len_n(arr, comblen, len, 0, 0));
		cnt += cnt_comb_len_n(arr, comblen, len, 0, 0);
	}
	return cnt;
}

int main()
{
	int n;
	scanf("%d", &n);
	int arr[n];
	for (int i = 0 ; i < n; i++)
		scanf("%d", &arr[i]);
	// printf("%s\n", "start");
	printf("%d\n", cnt(arr, n));

}