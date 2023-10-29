#ifndef _MD4_H_
#define _MD4_H_

#include "hash.h"

class MD4 : public MDHash {
public:
    MD4(int rnds = 48, int dobbertin = 0, int bits = 32, bool initBlock = true);

    void encode();

    Word q[52];
    bool dobbertin;
    int bits;
};

#endif
