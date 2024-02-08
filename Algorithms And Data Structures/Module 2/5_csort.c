#include <stdio.h>
#include <iso646.h>
#include <stdlib.h>
#include <string.h>

char* strdup (const char* s)
{
  size_t slen = strlen(s);
  char* result = malloc(slen + 1);
  if(result == NULL)
  {
    return NULL;
  }

  memcpy(result, s, slen+1);
  return result;
}

char *concat(char **s, int n)
{
    int N = 0;
    for(int i = 0; i < n; i++)
        N += strlen(s[i]);
    char *res = (char*)malloc(N+1);
    int top = 0;
    for(int i = 0; i < n; i++)
    {
        for(int j =0; (j < N) && (s[i][j] != '\0'); j++)
        {
            memcpy((void *)(res + top), (const void*)&s[i][j], 1);
            top += 1;
        }
    }
    res[top] = '\0';
    return res;
}

int wcount(char *s)
{
    int count = 0, len = strlen(s);
    if (len == 0)
        return 0;
    if (1 - (s[0] == ' ' or s[0] == '\n' or s[0] == '\t'))
        count += 1;
    for (int i = 1; i < len; i++)
    {
        if ((s[i-1] == ' ' or s[i-1] == '\n' or s[i-1] == '\t') 
        and (1 - (s[i] == ' ' or s[i] == '\n' or s[i] == '\t')))
            count += 1;

    }
    return count;

}
void csort(char *src, char *dest)
{	
	//инициализация массива split
	int len = wcount(src);
	char *s[len];
	for(int i = 0; i <  len; i++)
		s[i]  = (char*)calloc(1000/len + 10, 1);
	// заполнение массива сплит
	int i = 0, cnt = 0, len_src = strlen(src);
	while (cnt < len and i < len_src)
	{
		int j = 0;
		while(src[i] != '\t' and src[i] != '\n' 
			and src[i] != ' ' and i < len_src)
		{
			s[cnt][j] = src[i];
			j++;
			i++;
		}
		if(j)
		{
			s[cnt][j] = ' ';
			cnt++;
		}
		i++;
	}
	//сортировка 
	int count[len];
	for(int i = 0; i < len; i++)
		count[i] = 0;

	int j = 0;
	while (j < len - 1)
	{
		int i = j + 1;
		while (i < len)
		{
			if (strlen(s[i]) < strlen(s[j]))
				count[j]++;
			else
				count[i]++;
			i++;
		}
		j++;
	}

	//осуществление перестановки
	char *s2[len];
	for(int i = 0; i <  len; i++)
		s2[count[i]]  = strdup(s[i]);


	char *s3 = concat(s2, len);
	memcpy(dest, s3, strlen(s3));
	dest[strlen(s3) - 1] = '\0';
	for(int i = 0; i < len; i++)
	{
		free(s[i]);
		free(s2[i]);
	}
	free(s3);
}

int main()
{
	char src[1001];
	fgets(src, 1001, stdin);
	char *dest = (char *)calloc(2000, 'f');
	csort(src, dest);
	printf("%s\n", dest);
	free(dest);
	return 0;
}