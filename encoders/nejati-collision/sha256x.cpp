#include "sha256x.h"

SHA256::SHA256(int rnds, bool initBlock)
{
    rounds = rnds;
    initialBlock = initBlock;
}

void SHA256::encode()
{
    for (int i = 0; i < rounds; i++)
        cnf.newVars(w[i], 32, "w" + to_string(i));

    for (int i = 0; i < 8; i++)
        cnf.newVars(in[i], 32);

    for (int i = 0; i < 8; i++)
        cnf.newVars(out[i], 32, "hash" + to_string(i));

    for (int i = 0; i < rounds; i++) {
        cnf.newVars(A[i + 4], 32);
        cnf.newVars(E[i + 4], 32);
    }

    /* Message expansion */
    for (int i = 16; i < rounds; i++) {
        cnf.newVars(s0[i]);
        cnf.newVars(s1[i]);

        int r1[32], r2[32];
        cnf.rotr(r1, w[i - 15], 7);
        cnf.rotr(r2, w[i - 15], 18);
        cnf.xor2(s0[i] + 29, r1 + 29, r2 + 29, 3);
        cnf.xor3(s0[i], r1, r2, w[i - 15] + 3, 29);

        cnf.rotr(r1, w[i - 2], 17);
        cnf.rotr(r2, w[i - 2], 19);
        cnf.xor2(s1[i] + 22, r1 + 22, r2 + 22, 10);
        cnf.xor3(s1[i], r1, r2, w[i - 2] + 10, 22);

        cnf.add(w[i], w[i - 16], s0[i], wcarry[i], wCarry[i], w[i - 7], s1[i]);
    }

    /* Round constants */
    unsigned rnd_const[] = {
        0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5, 0x3956c25b, 0x59f111f1, 0x923f82a4, 0xab1c5ed5,
        0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3, 0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174,
        0xe49b69c1, 0xefbe4786, 0x0fc19dc6, 0x240ca1cc, 0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da,
        0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7, 0xc6e00bf3, 0xd5a79147, 0x06ca6351, 0x14292967,
        0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13, 0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85,
        0xa2bfe8a1, 0xa81a664b, 0xc24b8b70, 0xc76c51a3, 0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070,
        0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5, 0x391c0cb3, 0x4ed8aa4a, 0x5b9cca4f, 0x682e6ff3,
        0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208, 0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2
    };

    int k[64][32];
    for (int i = 0; i < rounds; i++) {
        cnf.newVars(k[i]);
        cnf.fixedValue(k[i], rnd_const[i]);
    }

    /* Initialization vector */
    if (initialBlock) {
        cnf.fixedValue(in[0], 0x6a09e667);
        cnf.fixedValue(in[1], 0xbb67ae85);
        cnf.fixedValue(in[2], 0x3c6ef372);
        cnf.fixedValue(in[3], 0xa54ff53a);
        cnf.fixedValue(in[4], 0x510e527f);
        cnf.fixedValue(in[5], 0x9b05688c);
        cnf.fixedValue(in[6], 0x1f83d9ab);
        cnf.fixedValue(in[7], 0x5be0cd19);
    }

    cnf.assign(A[3], in[0]);
    cnf.assign(A[2], in[1]);
    cnf.assign(A[1], in[2]);
    cnf.assign(A[0], in[3]);
    cnf.assign(E[3], in[4]);
    cnf.assign(E[2], in[5]);
    cnf.assign(E[1], in[6]);
    cnf.assign(E[0], in[7]);

    /* Main loop */
    for (int i = 0; i < rounds; i++) {
        cnf.newVars(sigma0[i]);
        cnf.newVars(sigma1[i]);
        Sigma0(sigma0[i], A[i + 3]);
        Sigma1(sigma1[i], E[i + 3]);

        cnf.newVars(f1[i]);
        cnf.newVars(f2[i]);
        cnf.ch(f1[i], E[i + 3], E[i + 2], E[i + 1]);
        cnf.maj3(f2[i], A[i + 3], A[i + 2], A[i + 1]);

        cnf.newVars(T[i]);

        cnf.add(T[i], E[i], sigma1[i], r0carry[i], r0Carry[i], f1[i], k[i], w[i]);

        cnf.add(E[i + 4], A[i], T[i], r1carry[i]);

        cnf.add(A[i + 4], T[i], sigma0[i], r2carry[i], r2Carry[i], f2[i]);
    }

    /* Final addition */
    cnf.add(out[0], in[0], A[rounds + 3], ocarry[0]);
    cnf.add(out[1], in[1], A[rounds + 2], ocarry[1]);
    cnf.add(out[2], in[2], A[rounds + 1], ocarry[2]);
    cnf.add(out[3], in[3], A[rounds], ocarry[3]);
    cnf.add(out[4], in[4], E[rounds + 3], ocarry[4]);
    cnf.add(out[5], in[5], E[rounds + 2], ocarry[5]);
    cnf.add(out[6], in[6], E[rounds + 1], ocarry[6]);
    cnf.add(out[7], in[7], E[rounds], ocarry[7]);
}

void SHA256::fixOutput(unsigned target[8])
{
    for (int i = 0; i < 8; i++)
        cnf.fixedValue(out[i], target[i]);
}

void SHA256::Sigma0(int* z, int* x)
{
    int r1[32], r2[32], r3[32];
    cnf.rotr(r1, x, 2);
    cnf.rotr(r2, x, 13);
    cnf.rotr(r3, x, 22);
    cnf.xor3(z, r1, r2, r3);
}

void SHA256::Sigma1(int* z, int* x)
{
    int r1[32], r2[32], r3[32];
    cnf.rotr(r1, x, 6);
    cnf.rotr(r2, x, 11);
    cnf.rotr(r3, x, 25);
    cnf.xor3(z, r1, r2, r3);
}