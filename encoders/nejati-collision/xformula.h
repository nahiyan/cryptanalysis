#ifndef _EXTENDED_FORMULA_H_
#define _EXTENDED_FORMULA_H_

#include "formula.h"

class xFormula : public Formula {
    public:
    xFormula(string name = "");
    virtual ~xFormula();

    // t: first carry bit, T: second carry bit
    // (ai + bi + ci + di + ei + Ti-2 + ti-1 --> Ti ti zi)
    void add(int *z, int *a, int *b, int *t, int *T = NULL, int *c = NULL, int *d = NULL, int *e = NULL);
    void xor5(int *z, int *a, int *b, int *c, int *d, int *e, int n = 32);
    void xor6(int *z, int *a, int *b, int *c, int *d, int *e, int *f, int n = 32);
    void xor7(int *z, int *a, int *b, int *c, int *d, int *e, int *f, int *g, int n = 32);

    void diff_add(int *z, int *a, int *b, int *t, int *T = NULL, int *c = NULL, int *d = NULL, int *e = NULL);

    void comp(int z, int *v, int n, int t, int T = -1);
};

#endif
