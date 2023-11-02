#include <stdio.h>

long long gorner (long long arr[], long long x, long long len)
{
	long long px = arr[0]*x + arr[1];
	for (int i = 2; i < len; i++){
		px = px*x + arr[i];
	}
	return px;
}

long long gornerdx (long long arr[], long long x, long long len)
{
	for (int i = 0; i < len; i++)
	{
		arr[i] = arr[i]*(len - 1 - i);
	}
	return gorner(arr, x, len-1);
}

int main()
{
	long long n, x;
	scanf("%lld%lld", &n, &x);
	long long arr_a[n+1];
	for (int i = 0; i < n + 1; i++)
	{
		scanf("%lld", &arr_a[i]);
	}

	printf("%lld\n", gorner(arr_a,x,n+1));
	printf("%lld", gornerdx(arr_a,x,n+1));
	return 0;
}