#include <stdio.h>
#include <stdlib.h>
#include <iso646.h>
#include <string.h>
void conc(char* s1, char* s2, char* s)
{
    // for(int i = 0; i < strlen(s1)+strlen(s2)+1; i++)
    //     s[i] = '\0';
    s[strlen(s2)+strlen(s1)] = '\0';
    int top = 0;
    for(int i = 0; i<strlen(s1); i++)
    {
        s[top] = s1[i];
        top += 1;
    }
    for(int i =0; i <strlen(s2);i++)
    {   
        s[top] = s2[i];
        top += 1;
    }
}
char *fibstr(int n)
{
    char *si_2 = (char*)malloc(2);

    char *si_1 = (char*)malloc(2);
    strcpy(si_2, "a");
    strcpy(si_1, "b");
    if (1 == n)
        return si_2;
    if (2 == n)
        return si_1;

    char *si = (char*)malloc(strlen(si_2)+strlen(si_1)+1);
    conc(si_2, si_1, si);
    printf("%s\n", si);
    int i = 3;
    while (i < n)
    {
        printf("iter %d\n", i);
        // printf("1) %s\n", si_1);
        si_2 = realloc(si_2, strlen(si_1)+1);
        strcpy(si_2, si_1);
        si_1 = realloc(si_1, strlen(si)+1);
        // printf("2) %s\n", si_1);
        strcpy(si_1,si);
        // printf("3) %s\n", si_1);
        si = realloc(si, strlen(si_2)+strlen(si_1)+1);
        if (!si)
            printf("ERROR");
        
        // printf("4) %s %s %s\n",si_2, si_1, si);
        // printf("5) %s\n", si_1);
        conc(si_2, si_1, si);
        // printf("check\n");
        i++;

    }
    free(si_2);
    free(si_1);
    return si;
}
int main()
{
    int n;
    scanf("%d", &n);
    printf("%s\n", fibstr(n));
    printf("check2\n");

    return 0;
}
