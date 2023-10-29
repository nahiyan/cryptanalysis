#ifndef _SHA256_X_H_
#define _SHA256_X_H_

#include "xformula.h"

class SHA256 {
    public:
        SHA256(
                int rnds = 64,
                bool initBlock = true);

        void encode();

        int w[64][32];
        int in[8][32];
        int out[8][32];

        int A[70][32];
        int E[70][32];
        int T[70][32];

        int s0[64][32];
        int s1[64][32];
        int sigma0[64][32];
        int sigma1[64][32];

        int f1[64][32];
        int f2[64][32];

        int wcarry[64][32], wCarry[64][32];
        int r0carry[64][32], r0Carry[64][32];
        int r1carry[64][32];
        int r2carry[64][32], r2Carry[64][32];
        int ocarry[8][32];

        int rounds;
        bool initialBlock;
        void fixOutput(unsigned target[8]);

        xFormula cnf;

        void Sigma0(int *z, int *x, std::string prefix = "");
        void Sigma1(int *z, int *x, std::string prefix = "");
};

#endif
