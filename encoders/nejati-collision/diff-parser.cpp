#include "diff-parser.h"


void parse_diff_path(int rounds, FILE* f, vector<string>& A, vector<string>& E, vector<string>& W)
{
    int rnd;
    char a[40], e[40], w[40];
    char tmp1[5], tmp2[5], tmp3[5];

    for( int i=0; i<4; i++ )
    {
        fscanf(f, "%d %s %s %s %s", &rnd, tmp1, a, tmp2, e);
        A.push_back(a);
        E.push_back(e);
    }

    for( int i=0; i<rounds; i++ )
    {
        fscanf(f, "%d %s %s %s %s %s %s", &rnd, tmp1, a, tmp2, e, tmp3, w);
        A.push_back(a);
        E.push_back(e);
        W.push_back(w);
    }
}
