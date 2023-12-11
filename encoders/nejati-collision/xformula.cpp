#include "xformula.h"
#include <assert.h>

xFormula::xFormula(string name)
    : Formula(name)
{
}

xFormula::~xFormula()
{
}

void xFormula::add(int* z, int* a, int* b, int* t, int* T, int* c, int* d, int* e)
{
    assert(multiAdderType == ESPRESSO);
    int n = 32;
    vector<int> addends[n + 5];
    for (int i = 0; i < n; i++) {
        addends[i].push_back(a[i]);
        addends[i].push_back(b[i]);
        if (c != NULL)
            addends[i].push_back(c[i]);
        if (d != NULL)
            addends[i].push_back(d[i]);
        if (e != NULL)
            addends[i].push_back(e[i]);

        int m = addends[i].size() > 3 ? 3 : 2;
        vector<int> sum(m);
        sum[0] = z[i];
        sum[1] = t[i];
        if (m == 3)
            sum[2] = T[i];
        addends[i + 1].push_back(sum[1]);
        if (m == 3)
            addends[i + 2].push_back(sum[2]);

        espresso(addends[i], sum);
    }
}

void xFormula::comp(int z, int* v, int n, int t, int T)
{
    assert(n >= 2 && n <= 7);
    if (n > 3)
        assert(T != -1);

    if (n == 2) {
        xor2(&z, v, v + 1, 1);
        addClause({ v[0], v[1], -t });
    } else if (n == 3) {
        xor3(&z, v, v + 1, v + 2, 1);
        addClause({ v[0], v[1], v[2], -t });
        addClause({ -v[0], -v[1], -v[2], t });
    } else if (n == 4) {
        xor4(&z, v, v + 1, v + 2, v + 3, 1);
        addClause({ v[0], v[1], v[2], v[3], -t });
        addClause({ v[0], v[1], v[2], v[3], -T });
    } else if (n == 5) {
        xor5(&z, v, v + 1, v + 2, v + 3, v + 4, 1);
        addClause({ v[0], v[1], v[2], v[3], v[4], -t });
        addClause({ v[0], v[1], v[2], v[3], v[4], -T });
        addClause({ -v[0], -v[1], -v[2], -v[3], -v[4], -t });
    } else if (n == 6) {
        xor6(&z, v, v + 1, v + 2, v + 3, v + 4, v + 5, 1);
        addClause({ v[0], v[1], v[2], v[3], v[4], v[5], -t });
        addClause({ v[0], v[1], v[2], v[3], v[4], v[5], -T });
    } else if (n == 7) {
        xor7(&z, v, v + 1, v + 2, v + 3, v + 4, v + 5, v + 6, 1);
        addClause({ v[0], v[1], v[2], v[3], v[4], v[5], v[6], -t });
        addClause({ v[0], v[1], v[2], v[3], v[4], v[5], v[6], -T });
        addClause({ -v[0], -v[1], -v[2], -v[3], -v[4], -v[5], -v[6], t });
        addClause({ -v[0], -v[1], -v[2], -v[3], -v[4], -v[5], -v[6], T });
    }
}

void xFormula::diff_add(int* z, int* a, int* b, int* t, int* T, int* c, int* d, int* e)
{
    int n = 32;
    int m = 2;
    if (c)
        m++;
    if (d)
        m++;
    if (e)
        m++;
    int v[10], k;
    for (int j = 0; j < 32; j++) {
        k = 0;
        v[k++] = a[j];
        v[k++] = b[j];
        if (c)
            v[k++] = c[j];
        if (d)
            v[k++] = d[j];
        if (e)
            v[k++] = e[j];
        if (j > 0)
            v[k++] = t[j - 1];
        if (j > 1)
            if ((m == 3 && j >= 3) || (m > 3))
                v[k++] = T[j - 2];

        if (m == 2)
            comp(z[j], v, k, t[j]);
        else
            comp(z[j], v, k, t[j], T[j]);
    }
}

void xFormula::xor5(int* z, int* a, int* b, int* c, int* d, int* e, int n)
{
    for (int i = 0; i < n; i++) {
        if (useXORClauses) {
            addClause(Clause({ -z[i], a[i], b[i], c[i], d[i], e[i] }, true));
        } else {
            addClause({ -z[i], -a[i], -b[i], -c[i], -d[i], e[i] });
            addClause({ z[i], -a[i], -b[i], -c[i], d[i], e[i] });
            addClause({ z[i], -a[i], -b[i], c[i], -d[i], e[i] });
            addClause({ -z[i], -a[i], -b[i], c[i], d[i], e[i] });
            addClause({ z[i], -a[i], b[i], -c[i], -d[i], e[i] });
            addClause({ -z[i], -a[i], b[i], -c[i], d[i], e[i] });
            addClause({ -z[i], -a[i], b[i], c[i], -d[i], e[i] });
            addClause({ z[i], -a[i], b[i], c[i], d[i], e[i] });
            addClause({ z[i], a[i], -b[i], -c[i], -d[i], e[i] });
            addClause({ -z[i], a[i], -b[i], -c[i], d[i], e[i] });
            addClause({ -z[i], a[i], -b[i], c[i], -d[i], e[i] });
            addClause({ z[i], a[i], -b[i], c[i], d[i], e[i] });
            addClause({ -z[i], a[i], b[i], -c[i], -d[i], e[i] });
            addClause({ z[i], a[i], b[i], -c[i], d[i], e[i] });
            addClause({ z[i], a[i], b[i], c[i], -d[i], e[i] });
            addClause({ -z[i], a[i], b[i], c[i], d[i], e[i] });
            addClause({ z[i], -a[i], -b[i], -c[i], -d[i], -e[i] });
            addClause({ -z[i], -a[i], -b[i], -c[i], d[i], -e[i] });
            addClause({ -z[i], -a[i], -b[i], c[i], -d[i], -e[i] });
            addClause({ z[i], -a[i], -b[i], c[i], d[i], -e[i] });
            addClause({ -z[i], -a[i], b[i], -c[i], -d[i], -e[i] });
            addClause({ z[i], -a[i], b[i], -c[i], d[i], -e[i] });
            addClause({ z[i], -a[i], b[i], c[i], -d[i], -e[i] });
            addClause({ -z[i], -a[i], b[i], c[i], d[i], -e[i] });
            addClause({ -z[i], a[i], -b[i], -c[i], -d[i], -e[i] });
            addClause({ z[i], a[i], -b[i], -c[i], d[i], -e[i] });
            addClause({ z[i], a[i], -b[i], c[i], -d[i], -e[i] });
            addClause({ -z[i], a[i], -b[i], c[i], d[i], -e[i] });
            addClause({ z[i], a[i], b[i], -c[i], -d[i], -e[i] });
            addClause({ -z[i], a[i], b[i], -c[i], d[i], -e[i] });
            addClause({ -z[i], a[i], b[i], c[i], -d[i], -e[i] });
            addClause({ z[i], a[i], b[i], c[i], d[i], -e[i] });
        }
    }
}

void xFormula::xor6(int* z, int* a, int* b, int* c, int* d, int* e, int* f, int n)
{
    for (int i = 0; i < n; i++) {
        if (useXORClauses) {
            addClause(Clause({ -z[i], a[i], b[i], c[i], d[i], e[i], f[i] }, true));
        } else {
            addClause({ -z[i], -a[i], -b[i], -c[i], -d[i], e[i], f[i] });
            addClause({ z[i], -a[i], -b[i], -c[i], d[i], e[i], f[i] });
            addClause({ z[i], -a[i], -b[i], c[i], -d[i], e[i], f[i] });
            addClause({ -z[i], -a[i], -b[i], c[i], d[i], e[i], f[i] });
            addClause({ z[i], -a[i], b[i], -c[i], -d[i], e[i], f[i] });
            addClause({ -z[i], -a[i], b[i], -c[i], d[i], e[i], f[i] });
            addClause({ -z[i], -a[i], b[i], c[i], -d[i], e[i], f[i] });
            addClause({ z[i], -a[i], b[i], c[i], d[i], e[i], f[i] });
            addClause({ z[i], a[i], -b[i], -c[i], -d[i], e[i], f[i] });
            addClause({ -z[i], a[i], -b[i], -c[i], d[i], e[i], f[i] });
            addClause({ -z[i], a[i], -b[i], c[i], -d[i], e[i], f[i] });
            addClause({ z[i], a[i], -b[i], c[i], d[i], e[i], f[i] });
            addClause({ -z[i], a[i], b[i], -c[i], -d[i], e[i], f[i] });
            addClause({ z[i], a[i], b[i], -c[i], d[i], e[i], f[i] });
            addClause({ z[i], a[i], b[i], c[i], -d[i], e[i], f[i] });
            addClause({ -z[i], a[i], b[i], c[i], d[i], e[i], f[i] });
            addClause({ z[i], -a[i], -b[i], -c[i], -d[i], -e[i], f[i] });
            addClause({ -z[i], -a[i], -b[i], -c[i], d[i], -e[i], f[i] });
            addClause({ -z[i], -a[i], -b[i], c[i], -d[i], -e[i], f[i] });
            addClause({ z[i], -a[i], -b[i], c[i], d[i], -e[i], f[i] });
            addClause({ -z[i], -a[i], b[i], -c[i], -d[i], -e[i], f[i] });
            addClause({ z[i], -a[i], b[i], -c[i], d[i], -e[i], f[i] });
            addClause({ z[i], -a[i], b[i], c[i], -d[i], -e[i], f[i] });
            addClause({ -z[i], -a[i], b[i], c[i], d[i], -e[i], f[i] });
            addClause({ -z[i], a[i], -b[i], -c[i], -d[i], -e[i], f[i] });
            addClause({ z[i], a[i], -b[i], -c[i], d[i], -e[i], f[i] });
            addClause({ z[i], a[i], -b[i], c[i], -d[i], -e[i], f[i] });
            addClause({ -z[i], a[i], -b[i], c[i], d[i], -e[i], f[i] });
            addClause({ z[i], a[i], b[i], -c[i], -d[i], -e[i], f[i] });
            addClause({ -z[i], a[i], b[i], -c[i], d[i], -e[i], f[i] });
            addClause({ -z[i], a[i], b[i], c[i], -d[i], -e[i], f[i] });
            addClause({ z[i], a[i], b[i], c[i], d[i], -e[i], f[i] });
            addClause({ z[i], -a[i], -b[i], -c[i], -d[i], e[i], -f[i] });
            addClause({ -z[i], -a[i], -b[i], -c[i], d[i], e[i], -f[i] });
            addClause({ -z[i], -a[i], -b[i], c[i], -d[i], e[i], -f[i] });
            addClause({ z[i], -a[i], -b[i], c[i], d[i], e[i], -f[i] });
            addClause({ -z[i], -a[i], b[i], -c[i], -d[i], e[i], -f[i] });
            addClause({ z[i], -a[i], b[i], -c[i], d[i], e[i], -f[i] });
            addClause({ z[i], -a[i], b[i], c[i], -d[i], e[i], -f[i] });
            addClause({ -z[i], -a[i], b[i], c[i], d[i], e[i], -f[i] });
            addClause({ -z[i], a[i], -b[i], -c[i], -d[i], e[i], -f[i] });
            addClause({ z[i], a[i], -b[i], -c[i], d[i], e[i], -f[i] });
            addClause({ z[i], a[i], -b[i], c[i], -d[i], e[i], -f[i] });
            addClause({ -z[i], a[i], -b[i], c[i], d[i], e[i], -f[i] });
            addClause({ z[i], a[i], b[i], -c[i], -d[i], e[i], -f[i] });
            addClause({ -z[i], a[i], b[i], -c[i], d[i], e[i], -f[i] });
            addClause({ -z[i], a[i], b[i], c[i], -d[i], e[i], -f[i] });
            addClause({ z[i], a[i], b[i], c[i], d[i], e[i], -f[i] });
            addClause({ -z[i], -a[i], -b[i], -c[i], -d[i], -e[i], -f[i] });
            addClause({ z[i], -a[i], -b[i], -c[i], d[i], -e[i], -f[i] });
            addClause({ z[i], -a[i], -b[i], c[i], -d[i], -e[i], -f[i] });
            addClause({ -z[i], -a[i], -b[i], c[i], d[i], -e[i], -f[i] });
            addClause({ z[i], -a[i], b[i], -c[i], -d[i], -e[i], -f[i] });
            addClause({ -z[i], -a[i], b[i], -c[i], d[i], -e[i], -f[i] });
            addClause({ -z[i], -a[i], b[i], c[i], -d[i], -e[i], -f[i] });
            addClause({ z[i], -a[i], b[i], c[i], d[i], -e[i], -f[i] });
            addClause({ z[i], a[i], -b[i], -c[i], -d[i], -e[i], -f[i] });
            addClause({ -z[i], a[i], -b[i], -c[i], d[i], -e[i], -f[i] });
            addClause({ -z[i], a[i], -b[i], c[i], -d[i], -e[i], -f[i] });
            addClause({ z[i], a[i], -b[i], c[i], d[i], -e[i], -f[i] });
            addClause({ -z[i], a[i], b[i], -c[i], -d[i], -e[i], -f[i] });
            addClause({ z[i], a[i], b[i], -c[i], d[i], -e[i], -f[i] });
            addClause({ z[i], a[i], b[i], c[i], -d[i], -e[i], -f[i] });
            addClause({ -z[i], a[i], b[i], c[i], d[i], -e[i], -f[i] });
        }
    }
}

void xFormula::xor7(int* z, int* a, int* b, int* c, int* d, int* e, int* f, int* g, int n)
{
    for (int i = 0; i < n; i++) {
        if (useXORClauses) {
            addClause(Clause({ -z[i], a[i], b[i], c[i], d[i], e[i], f[i], g[i] }, true));
        } else {
            addClause({ -z[i], -a[i], -b[i], -c[i], -d[i], e[i], f[i], g[i] });
            addClause({ z[i], -a[i], -b[i], -c[i], d[i], e[i], f[i], g[i] });
            addClause({ z[i], -a[i], -b[i], c[i], -d[i], e[i], f[i], g[i] });
            addClause({ -z[i], -a[i], -b[i], c[i], d[i], e[i], f[i], g[i] });
            addClause({ z[i], -a[i], b[i], -c[i], -d[i], e[i], f[i], g[i] });
            addClause({ -z[i], -a[i], b[i], -c[i], d[i], e[i], f[i], g[i] });
            addClause({ -z[i], -a[i], b[i], c[i], -d[i], e[i], f[i], g[i] });
            addClause({ z[i], -a[i], b[i], c[i], d[i], e[i], f[i], g[i] });
            addClause({ z[i], a[i], -b[i], -c[i], -d[i], e[i], f[i], g[i] });
            addClause({ -z[i], a[i], -b[i], -c[i], d[i], e[i], f[i], g[i] });
            addClause({ -z[i], a[i], -b[i], c[i], -d[i], e[i], f[i], g[i] });
            addClause({ z[i], a[i], -b[i], c[i], d[i], e[i], f[i], g[i] });
            addClause({ -z[i], a[i], b[i], -c[i], -d[i], e[i], f[i], g[i] });
            addClause({ z[i], a[i], b[i], -c[i], d[i], e[i], f[i], g[i] });
            addClause({ z[i], a[i], b[i], c[i], -d[i], e[i], f[i], g[i] });
            addClause({ -z[i], a[i], b[i], c[i], d[i], e[i], f[i], g[i] });
            addClause({ z[i], -a[i], -b[i], -c[i], -d[i], -e[i], f[i], g[i] });
            addClause({ -z[i], -a[i], -b[i], -c[i], d[i], -e[i], f[i], g[i] });
            addClause({ -z[i], -a[i], -b[i], c[i], -d[i], -e[i], f[i], g[i] });
            addClause({ z[i], -a[i], -b[i], c[i], d[i], -e[i], f[i], g[i] });
            addClause({ -z[i], -a[i], b[i], -c[i], -d[i], -e[i], f[i], g[i] });
            addClause({ z[i], -a[i], b[i], -c[i], d[i], -e[i], f[i], g[i] });
            addClause({ z[i], -a[i], b[i], c[i], -d[i], -e[i], f[i], g[i] });
            addClause({ -z[i], -a[i], b[i], c[i], d[i], -e[i], f[i], g[i] });
            addClause({ -z[i], a[i], -b[i], -c[i], -d[i], -e[i], f[i], g[i] });
            addClause({ z[i], a[i], -b[i], -c[i], d[i], -e[i], f[i], g[i] });
            addClause({ z[i], a[i], -b[i], c[i], -d[i], -e[i], f[i], g[i] });
            addClause({ -z[i], a[i], -b[i], c[i], d[i], -e[i], f[i], g[i] });
            addClause({ z[i], a[i], b[i], -c[i], -d[i], -e[i], f[i], g[i] });
            addClause({ -z[i], a[i], b[i], -c[i], d[i], -e[i], f[i], g[i] });
            addClause({ -z[i], a[i], b[i], c[i], -d[i], -e[i], f[i], g[i] });
            addClause({ z[i], a[i], b[i], c[i], d[i], -e[i], f[i], g[i] });
            addClause({ z[i], -a[i], -b[i], -c[i], -d[i], e[i], -f[i], g[i] });
            addClause({ -z[i], -a[i], -b[i], -c[i], d[i], e[i], -f[i], g[i] });
            addClause({ -z[i], -a[i], -b[i], c[i], -d[i], e[i], -f[i], g[i] });
            addClause({ z[i], -a[i], -b[i], c[i], d[i], e[i], -f[i], g[i] });
            addClause({ -z[i], -a[i], b[i], -c[i], -d[i], e[i], -f[i], g[i] });
            addClause({ z[i], -a[i], b[i], -c[i], d[i], e[i], -f[i], g[i] });
            addClause({ z[i], -a[i], b[i], c[i], -d[i], e[i], -f[i], g[i] });
            addClause({ -z[i], -a[i], b[i], c[i], d[i], e[i], -f[i], g[i] });
            addClause({ -z[i], a[i], -b[i], -c[i], -d[i], e[i], -f[i], g[i] });
            addClause({ z[i], a[i], -b[i], -c[i], d[i], e[i], -f[i], g[i] });
            addClause({ z[i], a[i], -b[i], c[i], -d[i], e[i], -f[i], g[i] });
            addClause({ -z[i], a[i], -b[i], c[i], d[i], e[i], -f[i], g[i] });
            addClause({ z[i], a[i], b[i], -c[i], -d[i], e[i], -f[i], g[i] });
            addClause({ -z[i], a[i], b[i], -c[i], d[i], e[i], -f[i], g[i] });
            addClause({ -z[i], a[i], b[i], c[i], -d[i], e[i], -f[i], g[i] });
            addClause({ z[i], a[i], b[i], c[i], d[i], e[i], -f[i], g[i] });
            addClause({ -z[i], -a[i], -b[i], -c[i], -d[i], -e[i], -f[i], g[i] });
            addClause({ z[i], -a[i], -b[i], -c[i], d[i], -e[i], -f[i], g[i] });
            addClause({ z[i], -a[i], -b[i], c[i], -d[i], -e[i], -f[i], g[i] });
            addClause({ -z[i], -a[i], -b[i], c[i], d[i], -e[i], -f[i], g[i] });
            addClause({ z[i], -a[i], b[i], -c[i], -d[i], -e[i], -f[i], g[i] });
            addClause({ -z[i], -a[i], b[i], -c[i], d[i], -e[i], -f[i], g[i] });
            addClause({ -z[i], -a[i], b[i], c[i], -d[i], -e[i], -f[i], g[i] });
            addClause({ z[i], -a[i], b[i], c[i], d[i], -e[i], -f[i], g[i] });
            addClause({ z[i], a[i], -b[i], -c[i], -d[i], -e[i], -f[i], g[i] });
            addClause({ -z[i], a[i], -b[i], -c[i], d[i], -e[i], -f[i], g[i] });
            addClause({ -z[i], a[i], -b[i], c[i], -d[i], -e[i], -f[i], g[i] });
            addClause({ z[i], a[i], -b[i], c[i], d[i], -e[i], -f[i], g[i] });
            addClause({ -z[i], a[i], b[i], -c[i], -d[i], -e[i], -f[i], g[i] });
            addClause({ z[i], a[i], b[i], -c[i], d[i], -e[i], -f[i], g[i] });
            addClause({ z[i], a[i], b[i], c[i], -d[i], -e[i], -f[i], g[i] });
            addClause({ -z[i], a[i], b[i], c[i], d[i], -e[i], -f[i], g[i] });
            addClause({ z[i], -a[i], -b[i], -c[i], -d[i], e[i], f[i], -g[i] });
            addClause({ -z[i], -a[i], -b[i], -c[i], d[i], e[i], f[i], -g[i] });
            addClause({ -z[i], -a[i], -b[i], c[i], -d[i], e[i], f[i], -g[i] });
            addClause({ z[i], -a[i], -b[i], c[i], d[i], e[i], f[i], -g[i] });
            addClause({ -z[i], -a[i], b[i], -c[i], -d[i], e[i], f[i], -g[i] });
            addClause({ z[i], -a[i], b[i], -c[i], d[i], e[i], f[i], -g[i] });
            addClause({ z[i], -a[i], b[i], c[i], -d[i], e[i], f[i], -g[i] });
            addClause({ -z[i], -a[i], b[i], c[i], d[i], e[i], f[i], -g[i] });
            addClause({ -z[i], a[i], -b[i], -c[i], -d[i], e[i], f[i], -g[i] });
            addClause({ z[i], a[i], -b[i], -c[i], d[i], e[i], f[i], -g[i] });
            addClause({ z[i], a[i], -b[i], c[i], -d[i], e[i], f[i], -g[i] });
            addClause({ -z[i], a[i], -b[i], c[i], d[i], e[i], f[i], -g[i] });
            addClause({ z[i], a[i], b[i], -c[i], -d[i], e[i], f[i], -g[i] });
            addClause({ -z[i], a[i], b[i], -c[i], d[i], e[i], f[i], -g[i] });
            addClause({ -z[i], a[i], b[i], c[i], -d[i], e[i], f[i], -g[i] });
            addClause({ z[i], a[i], b[i], c[i], d[i], e[i], f[i], -g[i] });
            addClause({ -z[i], -a[i], -b[i], -c[i], -d[i], -e[i], f[i], -g[i] });
            addClause({ z[i], -a[i], -b[i], -c[i], d[i], -e[i], f[i], -g[i] });
            addClause({ z[i], -a[i], -b[i], c[i], -d[i], -e[i], f[i], -g[i] });
            addClause({ -z[i], -a[i], -b[i], c[i], d[i], -e[i], f[i], -g[i] });
            addClause({ z[i], -a[i], b[i], -c[i], -d[i], -e[i], f[i], -g[i] });
            addClause({ -z[i], -a[i], b[i], -c[i], d[i], -e[i], f[i], -g[i] });
            addClause({ -z[i], -a[i], b[i], c[i], -d[i], -e[i], f[i], -g[i] });
            addClause({ z[i], -a[i], b[i], c[i], d[i], -e[i], f[i], -g[i] });
            addClause({ z[i], a[i], -b[i], -c[i], -d[i], -e[i], f[i], -g[i] });
            addClause({ -z[i], a[i], -b[i], -c[i], d[i], -e[i], f[i], -g[i] });
            addClause({ -z[i], a[i], -b[i], c[i], -d[i], -e[i], f[i], -g[i] });
            addClause({ z[i], a[i], -b[i], c[i], d[i], -e[i], f[i], -g[i] });
            addClause({ -z[i], a[i], b[i], -c[i], -d[i], -e[i], f[i], -g[i] });
            addClause({ z[i], a[i], b[i], -c[i], d[i], -e[i], f[i], -g[i] });
            addClause({ z[i], a[i], b[i], c[i], -d[i], -e[i], f[i], -g[i] });
            addClause({ -z[i], a[i], b[i], c[i], d[i], -e[i], f[i], -g[i] });
            addClause({ -z[i], -a[i], -b[i], -c[i], -d[i], e[i], -f[i], -g[i] });
            addClause({ z[i], -a[i], -b[i], -c[i], d[i], e[i], -f[i], -g[i] });
            addClause({ z[i], -a[i], -b[i], c[i], -d[i], e[i], -f[i], -g[i] });
            addClause({ -z[i], -a[i], -b[i], c[i], d[i], e[i], -f[i], -g[i] });
            addClause({ z[i], -a[i], b[i], -c[i], -d[i], e[i], -f[i], -g[i] });
            addClause({ -z[i], -a[i], b[i], -c[i], d[i], e[i], -f[i], -g[i] });
            addClause({ -z[i], -a[i], b[i], c[i], -d[i], e[i], -f[i], -g[i] });
            addClause({ z[i], -a[i], b[i], c[i], d[i], e[i], -f[i], -g[i] });
            addClause({ z[i], a[i], -b[i], -c[i], -d[i], e[i], -f[i], -g[i] });
            addClause({ -z[i], a[i], -b[i], -c[i], d[i], e[i], -f[i], -g[i] });
            addClause({ -z[i], a[i], -b[i], c[i], -d[i], e[i], -f[i], -g[i] });
            addClause({ z[i], a[i], -b[i], c[i], d[i], e[i], -f[i], -g[i] });
            addClause({ -z[i], a[i], b[i], -c[i], -d[i], e[i], -f[i], -g[i] });
            addClause({ z[i], a[i], b[i], -c[i], d[i], e[i], -f[i], -g[i] });
            addClause({ z[i], a[i], b[i], c[i], -d[i], e[i], -f[i], -g[i] });
            addClause({ -z[i], a[i], b[i], c[i], d[i], e[i], -f[i], -g[i] });
            addClause({ z[i], -a[i], -b[i], -c[i], -d[i], -e[i], -f[i], -g[i] });
            addClause({ -z[i], -a[i], -b[i], -c[i], d[i], -e[i], -f[i], -g[i] });
            addClause({ -z[i], -a[i], -b[i], c[i], -d[i], -e[i], -f[i], -g[i] });
            addClause({ z[i], -a[i], -b[i], c[i], d[i], -e[i], -f[i], -g[i] });
            addClause({ -z[i], -a[i], b[i], -c[i], -d[i], -e[i], -f[i], -g[i] });
            addClause({ z[i], -a[i], b[i], -c[i], d[i], -e[i], -f[i], -g[i] });
            addClause({ z[i], -a[i], b[i], c[i], -d[i], -e[i], -f[i], -g[i] });
            addClause({ -z[i], -a[i], b[i], c[i], d[i], -e[i], -f[i], -g[i] });
            addClause({ -z[i], a[i], -b[i], -c[i], -d[i], -e[i], -f[i], -g[i] });
            addClause({ z[i], a[i], -b[i], -c[i], d[i], -e[i], -f[i], -g[i] });
            addClause({ z[i], a[i], -b[i], c[i], -d[i], -e[i], -f[i], -g[i] });
            addClause({ -z[i], a[i], -b[i], c[i], d[i], -e[i], -f[i], -g[i] });
            addClause({ z[i], a[i], b[i], -c[i], -d[i], -e[i], -f[i], -g[i] });
            addClause({ -z[i], a[i], b[i], -c[i], d[i], -e[i], -f[i], -g[i] });
            addClause({ -z[i], a[i], b[i], c[i], -d[i], -e[i], -f[i], -g[i] });
            addClause({ z[i], a[i], b[i], c[i], d[i], -e[i], -f[i], -g[i] });
        }
    }
}

void xFormula::basic_rules(int dx[32][4], int x[32], int x_[32])
{
    for (int i = 0; i < 32; i++) {
        // Define the 4-bit differences in terms of each difference bit
        addClause({x[i], x_[i], dx[i][0]});
        addClause({-x[i], x_[i], dx[i][1]});
        addClause({x[i], -x_[i], dx[i][2]});
        addClause({-x[i], -x_[i], dx[i][3]});

        // '-' -> x xnor x'
        addClause({ -dx[i][0], dx[i][1], dx[i][2], -dx[i][3], -x[i], x_[i] });
        addClause({ -dx[i][0], dx[i][1], dx[i][2], -dx[i][3], x[i], -x_[i] });

        // 'x' -> x xor x'
        addClause({ dx[i][0], -dx[i][1], -dx[i][2], dx[i][3], -x[i], -x_[i] });
        addClause({ dx[i][0], -dx[i][1], -dx[i][2], dx[i][3], x[i], x_[i] });

        // '0' -> ~x and ~x'
        addClause({ -dx[i][0], dx[i][1], dx[i][2], dx[i][3], -x[i] });
        addClause({ -dx[i][0], dx[i][1], dx[i][2], dx[i][3], -x_[i] });

        // 'u' -> x and ~x'
        addClause({ dx[i][0], -dx[i][1], dx[i][2], dx[i][3], x[i] });
        addClause({ dx[i][0], -dx[i][1], dx[i][2], dx[i][3], -x_[i] });

        // 'n' -> ~x and x'
        addClause({ dx[i][0], dx[i][1], -dx[i][2], dx[i][3], -x[i] });
        addClause({ dx[i][0], dx[i][1], -dx[i][2], dx[i][3], x_[i] });

        // '1' -> x and x'
        addClause({ dx[i][0], dx[i][1], dx[i][2], -dx[i][3], x[i] });
        addClause({ dx[i][0], dx[i][1], dx[i][2], -dx[i][3], x_[i] });

        // '3' -> ~x'
        addClause({ -dx[i][0], -dx[i][1], dx[i][2], dx[i][3], -x_[i] });

        // '5 -> ~x
        addClause({ -dx[i][0], dx[i][1], -dx[i][2], dx[i][3], -x[i] });

        // 'A' -> x
        addClause({ dx[i][0], -dx[i][1], dx[i][2], -dx[i][3], x[i] });

        // 'C' -> x'
        addClause({ dx[i][0], dx[i][1], -dx[i][2], -dx[i][3], x_[i] });

        // '7' -> ~x or ~x'
        addClause({ -dx[i][0], -dx[i][1], -dx[i][2], dx[i][3], -x[i], -x_[i] });

        // 'B' -> x or ~x'
        addClause({ -dx[i][0], -dx[i][1], dx[i][2], -dx[i][3], x[i], -x_[i] });

        // 'D' -> ~x or x'
        addClause({ -dx[i][0], dx[i][1], -dx[i][2], -dx[i][3], -x[i], x_[i] });

        // 'E' -> x or x'
        addClause({ dx[i][0], -dx[i][1], -dx[i][2], -dx[i][3], x[i], x_[i] });
    }
}

void xFormula::impose_rule(vector<int (*)[32][4]> inputs, vector<int (*)[32][4]> outputs, pair<string, string> rule)
{
    auto get_value = [](char diff) -> vector<int> {
        switch (diff) {
        case '?':
            return { 1, 1, 1, 1 };
            break;
        case '-':
            return { 1, 0, 0, 1 };
            break;
        case 'x':
            return { 0, 1, 1, 0 };
            break;
        case '0':
            return { 1, 0, 0, 0 };
            break;
        case 'u':
            return { 0, 1, 0, 0 };
            break;
        case 'n':
            return { 0, 0, 1, 0 };
            break;
        case '1':
            return { 0, 0, 0, 1 };
            break;
        case '3':
            return { 1, 1, 0, 0 };
            break;
        case '5':
            return { 1, 0, 1, 0 };
            break;
        case '7':
            return { 1, 1, 1, 0 };
            break;
        case 'A':
            return { 0, 1, 0, 1 };
            break;
        case 'B':
            return { 1, 1, 0, 1 };
            break;
        case 'C':
            return { 0, 0, 1, 1 };
            break;
        case 'D':
            return { 1, 0, 1, 1 };
            break;
        case 'E':
            return { 0, 1, 1, 1 };
            break;
        default:
            return { 0, 0, 0, 0 };
        }
    };

    string inputs_diff = rule.first, outputs_diff = rule.second;

    // Avoid encoding rules where we propagate to only '?'s
    int q_count = 0;
    for (auto& diff : outputs_diff) {
        if (diff == '?')
            q_count++;
    }
    if (q_count == outputs_diff.size())
        return;

    for (int i = 0; i < 32; i++) {
        vector<int> base_clause;
        for (int x = 0; x < inputs.size(); x++) {
            vector<int> values = get_value(inputs_diff[x]);
            assert(values.size() == 4);
            for (int k = 0; k < 4; k++)
                base_clause.push_back((values[k] == 1 ? -1 : 1) * (*inputs[x])[i][k]);
        }

        for (int x = 0; x < outputs.size(); x++) {
            vector<int> values = get_value(outputs_diff[x]);
            for (int k = 0; k < 4; k++) {
                vector<int> clause(base_clause);
                clause.push_back((values[k] == 1 ? 1 : -1) * (*outputs[x])[i][k]);
                addClause(clause);
            }
        }
    }
}