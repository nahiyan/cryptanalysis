#include "md4.h"

MD4::MD4(int rnds, int dobbertin, int bits, bool initBlock)
    : MDHash(16, 4, rnds, initBlock)
{
    this->dobbertin = dobbertin == 1;
    this->bits = bits;
}

void MD4::encode()
{
    /* Message words */
    for (int i = 0; i < 16; i++)
        cnf.newVars(w[i], 32, "w" + to_string(i));

    for (int i = 0; i < 4; i++)
        cnf.newVars(chain[i], 32);

    for (int i = 0; i < 4; i++)
        cnf.newVars(out[i], 32, "hash" + to_string(i));

    if (initialBlock) {
        cnf.fixedValue(chain[0], 0x67452301); // A
        cnf.fixedValue(chain[1], 0xefcdab89); // B
        cnf.fixedValue(chain[2], 0x98badcfe); // C
        cnf.fixedValue(chain[3], 0x10325476); // D
    }

    /* Round constants */
    Word k[2];
    cnf.newVars(k[0], 32);
    cnf.newVars(k[1], 32);
    cnf.fixedValue(k[0], 0x5a827999);
    cnf.fixedValue(k[1], 0x6ed9eba1);

    /* Round rotations */
    int s[48] = {
        3,
        7,
        11,
        19,
        3,
        7,
        11,
        19,
        3,
        7,
        11,
        19,
        3,
        7,
        11,
        19,
        3,
        5,
        9,
        13,
        3,
        5,
        9,
        13,
        3,
        5,
        9,
        13,
        3,
        5,
        9,
        13,
        3,
        9,
        11,
        15,
        3,
        9,
        11,
        15,
        3,
        9,
        11,
        15,
        3,
        9,
        11,
        15,
    };

    /* Round message word indices */
    int ind[48] = {
        0,
        1,
        2,
        3,
        4,
        5,
        6,
        7,
        8,
        9,
        10,
        11,
        12,
        13,
        14,
        15,
        0,
        4,
        8,
        12,
        1,
        5,
        9,
        13,
        2,
        6,
        10,
        14,
        3,
        7,
        11,
        15,
        0,
        8,
        4,
        12,
        2,
        10,
        6,
        14,
        1,
        9,
        5,
        13,
        3,
        11,
        7,
        15,
    };

    cnf.assign(q[0], chain[0]); // A
    cnf.assign(q[3], chain[1]); // B
    cnf.assign(q[2], chain[2]); // C
    cnf.assign(q[1], chain[3]); // D

    /* Main loop */
    for (int i = 0; i < rounds; i++) {
        Word t;
        cnf.newVars(t, 32);

        Word f;
        cnf.newVars(f, 32);

        if (i < 16) {
            cnf.ch(f, q[i + 3], q[i + 2], q[i + 1]);
            cnf.add3(t, q[i], f, w[ind[i]]);
        } else if (i < 32) {
            cnf.maj3(f, q[i + 3], q[i + 2], q[i + 1]);
            cnf.add4(t, q[i], f, w[ind[i]], k[0]);
        } else {
            cnf.xor3(f, q[i + 3], q[i + 2], q[i + 1]);
            cnf.add4(t, q[i], f, w[ind[i]], k[1]);
        }

        cnf.rotl(q[i + 4], t, s[i]);
    }

    // Dobbertin's constraints
    if (dobbertin) {
        // Dobbertin's constant
        unsigned int k = 0xffffffff;

        int q_indices[12] = {
            16, 20, 24, 28, // q[16] = Q[13] (from [Debapratim De et al. 2007])
            17, 21, 25, 29,
            18, 22, 26, 30
        };
        // Index of q that needs relaxation
        int p = 16;
        for (int& i : q_indices) {
            if (i == p && bits != 32) {
                // TODO: Offset the value of k
                cnf.fixedValue(q[i] + (32 - bits), k, bits);
            } else {
                cnf.fixedValue(q[i], k);
            }
        }
    }

    int R = rounds, r = rounds % 4;
    cnf.add2(out[0], chain[0], q[R + (4 - r) % 4]);
    cnf.add2(out[1], chain[1], q[R + (3 - r) % 4]);
    cnf.add2(out[2], chain[2], q[R + (6 - r) % 4]);
    cnf.add2(out[3], chain[3], q[R + (5 - r) % 4]);
}
