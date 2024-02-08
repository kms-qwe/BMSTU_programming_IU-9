#include <stdio.h>
#include <stdbool.h>

int main()
{
    long long A[8], B[8], cA[8][2], cB[8][2];
    for(int i = 0; i < 8; i++)
    {
        scanf("%lld", &A[i]);
    }
    for(int i = 0; i < 8; i++)
    {
        scanf("%lld", &B[i]);
    }
    for(int i = 0; i < 8; i++)
    {
        cA[i][0] = A[i];
        cA[i][1] = 0;
        for( int j = 0;j < 8; j++)
        {
            if (cA[i][0] == A[j])
                cA[i][1] += 1;
        }
    }

    for(int i = 0; i < 8; i++)
    {
        cB[i][0] = B[i];
        cB[i][1] = 0;
        for(int j = 0; j <8; j++)
        {
            if (cB[i][0] == B[j])
                cB[i][1] += 1;
        }
    }

    bool is_permut = true;
    for(int i = 0; i< 8; i++)
    {
        bool a_i_in_B = false;
        for(int j = 0;j<8;j++)
        {
            if (cA[i][0] == cB[j][0] && cA[i][1] == cB[j][1])
            {
                a_i_in_B = true;
                break;
            }
        }
        if (! a_i_in_B)
        {
            is_permut = false;
            break;
        }
    }
    if (is_permut)
        printf("yes");
    else
        printf("no");
    return 0;
}