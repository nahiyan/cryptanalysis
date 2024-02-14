#ifndef _EXTENDED_FORMULA_H_
#define _EXTENDED_FORMULA_H_

#include "formula.h"
#include <unordered_map>

struct Rules {
public:
    unordered_map<string, string> ch, maj, xor3, add2, add3, add4, add5, add6, add7;
};

class xFormula : public Formula {

public:
    xFormula(string name = "");
    virtual ~xFormula();

    // t: first carry bit, T: second carry bit
    // (ai + bi + ci + di + ei + Ti-2 + ti-1 --> Ti ti zi)
    void add(int* z, int* a, int* b, int* t, int* T = NULL, int* c = NULL, int* d = NULL, int* e = NULL);
    void xor5(int* z, int* a, int* b, int* c, int* d, int* e, int n = 32);
    void xor6(int* z, int* a, int* b, int* c, int* d, int* e, int* f, int n = 32);
    void xor7(int* z, int* a, int* b, int* c, int* d, int* e, int* f, int* g, int n = 32);

    void diff_1bit_add(Rules& rules, int z[32], int a[32], int b[32], int t[32], int T[32] = NULL, int c[32] = NULL, int d[32] = NULL, int e[32] = NULL);
    void diff_4bit_add(Rules& rules, int z[32], int a[32], int b[32], int t[32], int T[32] = NULL, int c[32] = NULL, int d[32] = NULL, int e[32] = NULL);
    void diff_add(Rules& rules, int z[32], int a[32], int b[32], int t[32], int T[32] = NULL, int c[32] = NULL, int d[32] = NULL, int e[32] = NULL);

    void comp_1bit(int z, int* v, int n, int t, int T = -1);
    void comp_4bit(Rules& rules, int z, int v[10], int n, int t, int T = -1);

    void impose_4bit_rule(vector<int> input_ids, vector<int> output_ids, pair<string, string> rule);
    void impose_1bit_rule(vector<int> input_ids, vector<int> output_ids, pair<string, string> rule);
    void impose_rule(vector<int> input_ids, vector<int> output_ids, pair<string, string> rule);
    void basic_4bit_rules(int dx[32], int x[32], int x_[32]);
    void basic_1bit_rules(int dx[32], int x[32], int x_[32]);
    void basic_rules(int dx[32], int x[32], int x_[32]);
};

#endif
