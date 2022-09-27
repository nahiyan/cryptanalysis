/*
 * The expected input format is the output of minisat style output.
 * SAT/UNSAT
 * +-1 +-2 ...
 */

#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <math.h>
#include <vector>
#include "util.h"
using namespace std;

void extract(unsigned *v, int n, vector<int>& sol, int offset)
{
    int k = offset;
    for( int i=0; i<n; i++ )
    {
        v[i] = 0;
        for( int j=0; j<32; j++ )
            v[i] = v[i] | (sol[k++] << j);
    }
}

int main(int argc, char** argv)
{
    int rounds = 48;
    int initBlock = 1;
    if ( argc >= 2 )
        rounds = atoi(argv[1]);
    if ( argc >= 3 )
        initBlock = atoi(argv[2]);
    /*if ( argc < 2 )
    {
        printf("USAGE: ./verify-sha1 {number_of_rounds} {initial_block(0/1)}\n");
        return 1;
    }*/

    char res[20];
    int r = scanf("%s", res);
    if ( r != 0 && res[0] != 'S' )
    {
        printf("Result is UNSAT!\n");
        return 1;
    }

    /* Reading the solution */
    vector<int> solution;
    int x;
    while( scanf("%d", &x) != EOF )
    {
        solution.push_back( x>0 ? 1 : 0 );
    }

    /* Extracting the message words */
    unsigned w[16];
    extract(w, 16, solution, 0);

    /* Extracting chaining value */
    unsigned cv[4];
    extract(cv, 4, solution, 16 * 32);

    /* Extracting the correct hash output */
    unsigned hash[4];
    extract(hash, 4, solution, 16 * 32 + 4 * 32);

    md4_comp(w, cv, rounds);
    bool matched = true;
    for( int i=0; i<4; i++ )
        if ( cv[i] != hash[i] )
            matched = false;

    printf("Target: ");
    for( int i=0; i<4; i++ )
        printf("%08x ", hash[i]);
    printf("\n");

    printf("Solution: ");
    for( int i=0; i<16; i++ )
        printf("%08x ", w[i]);
    printf("\n");

    if ( matched )
        printf("Solution's hash matches the target!\n");
    else
    {
        printf("Solution's hash DOES NOT match the target: \n");
        for( int i=0; i<4; i++ )
            printf("%08x ", cv[i]);
        printf("\n");
    }



    return 0;
}
