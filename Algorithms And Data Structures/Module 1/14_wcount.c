#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <iso646.h>

int wcount(char *s)
{
    int count = 0;
    if (1 - (s[0] == ' ' or s[0] == '\n' or s[0] == '\t'))
        count += 1;
    for (int i = 1; i < strlen(s); i++)
    {
        if ((s[i-1] == ' ' or s[i-1] == '\n' or s[i-1] == '\t') 
        and (1 - (s[i] == ' ' or s[i] == '\n' or s[i] == '\t')))
            count += 1;

    }
    return count;

}
int main()
{
    char s[2000];
    fgets(s, 2000, stdin);
    printf("%d", wcount(s));

    return 0;
}