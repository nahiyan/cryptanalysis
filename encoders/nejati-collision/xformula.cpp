#include "xformula.h"
#include <assert.h>
#include <cstddef>

inline vector<int> get_values_4bit(char diff)
{
    assert(diff != '#');
    switch (diff) {
    case '?':
        // ! IMPORTANT: It may be beneficial to not enforce '?'
        return {};
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
        assert(false);
        return {};
    }
};

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

void xFormula::comp_4bit(Rules& rules, int z, int v[10], int n, int t, int T)
{
    assert(n >= 2 && n <= 7);
    if (n > 3)
        assert(T != -1);

    unordered_map<string, string>* add_rules;
    if (n == 2)
        add_rules = &rules.add2;
    else if (n == 3)
        add_rules = &rules.add3;
    else if (n == 4)
        add_rules = &rules.add4;
    else if (n == 5)
        add_rules = &rules.add5;
    else if (n == 6)
        add_rules = &rules.add6;
    else if (n == 7)
        add_rules = &rules.add7;

    for (auto& differential : *add_rules) {
        string lhs = differential.first, rhs = differential.second;
        // cout << lhs << " " << rhs << endl;
        assert(lhs.size() == n);
        assert(rhs.size() == 3);

        vector<int> antecedent;
        int i = 0;
        for (char& c : lhs) {
            assert(c == '-' || c == 'x');
            if (c == '-') {
                antecedent.push_back(v[i] + 1);
                antecedent.push_back(v[i] + 2);
            } else if (c == 'x') {
                antecedent.push_back(v[i] + 0);
                antecedent.push_back(v[i] + 3);
            }
            i++;
        }
        i = -1;
        for (char& c : rhs) {
            ++i;

            if (i == 0 && (n <= 3 || T == -1))
                continue;

            assert(n <= 3 ? i != 0 : true);
            assert(T == -1 ? i != 0 : true);

            vector<int> values = get_values_4bit(c);
            if (values.size() == 0)
                continue;

            int diff;
            if (i == 0)
                diff = T;
            else if (i == 1)
                diff = t;
            else if (i == 2)
                diff = z;

            // Output
            for (int j = 0; j < 4; j++) {
                vector<int> clause(antecedent);
                if (values[j] == 1)
                    continue;
                clause.push_back((values[j] == 1 ? 1 : -1) * (diff + j));
                addClause(clause);
            }
        }
    }
}

void xFormula::diff_4bit_add(Rules& rules, int z[32], int a[32], int b[32], int t[32], int T[32], int c[32], int d[32], int e[32])
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
        if ((m == 3 && j >= 3) || (j >= 2 && m > 3))
            v[k++] = T[j - 2];

        if (m == 2)
            comp_4bit(rules, z[j], v, k, t[j]);
        else
            comp_4bit(rules, z[j], v, k, t[j], T[j]);
    }
}

void xFormula::comp_1bit(int z, int* v, int n, int t, int T)
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

void xFormula::diff_1bit_add(Rules& rules, int z[32], int a[32], int b[32], int t[32], int T[32], int c[32], int d[32], int e[32])
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
            comp_1bit(z[j], v, k, t[j]);
        else
            comp_1bit(z[j], v, k, t[j], T[j]);
    }
}

void xFormula::diff_add(Rules& rules, int z[32], int a[32], int b[32], int t[32], int T[32], int c[32], int d[32], int e[32])
{
#if IS_4bit
    diff_4bit_add(rules, z, a, b, t, T, c, d, e);
#else
    diff_1bit_add(rules, z, a, b, t, T, c, d, e);
#endif
}

void xFormula::basic_4bit_rules(int dx[32], int x[32], int x_[32])
{
    for (int i = 0; i < 32; i++) {
        // * (0, 0) -> '0'
        addClause({ x[i], x_[i], dx[i] + 0 });

        // * (1, 0) -> 'u'
        addClause({ -x[i], x_[i], dx[i] + 1 });

        // * (0, 1) -> 'n'
        addClause({ x[i], -x_[i], dx[i] + 2 });

        // * (1, 1) -> '1'
        addClause({ -x[i], -x_[i], dx[i] + 3 });

        // * (0, ?) -> can't be 'u' or '1'
        addClause({ x[i], -(dx[i] + 1) });
        addClause({ x[i], -(dx[i] + 3) });

        // * (?, 0) -> can't be 'n' or '1'
        addClause({ x_[i], -(dx[i] + 2) });
        addClause({ x_[i], -(dx[i] + 3) });

        // * (1, ?) -> can't be '0' or 'n'
        addClause({ -x[i], -(dx[i] + 0) });
        addClause({ -x[i], -(dx[i] + 2) });

        // * (?, 1) -> can't be '0' or 'u'
        addClause({ -x_[i], -(dx[i] + 0) });
        addClause({ -x_[i], -(dx[i] + 1) });

        // * '-' -> x xnor x'
        addClause({ dx[i] + 1, dx[i] + 2, -x[i], x_[i] });
        addClause({ dx[i] + 1, dx[i] + 2, x[i], -x_[i] });

        // * 'x' -> x xor x'
        addClause({ dx[i] + 0, dx[i] + 3, -x[i], -x_[i] });
        addClause({ dx[i] + 0, dx[i] + 3, x[i], x_[i] });

        // * '0' -> ~x and ~x'
        addClause({ dx[i] + 1, dx[i] + 2, dx[i] + 3, -x[i] });
        addClause({ dx[i] + 1, dx[i] + 2, dx[i] + 3, -x_[i] });

        // * 'u' -> x and ~x'
        addClause({ dx[i] + 0, dx[i] + 2, dx[i] + 3, x[i] });
        addClause({ dx[i] + 0, dx[i] + 2, dx[i] + 3, -x_[i] });

        // * 'n' -> ~x and x'
        addClause({ dx[i] + 0, dx[i] + 1, dx[i] + 3, -x[i] });
        addClause({ dx[i] + 0, dx[i] + 1, dx[i] + 3, x_[i] });

        // * '1' -> x and x'
        addClause({ dx[i] + 0, dx[i] + 1, dx[i] + 2, x[i] });
        addClause({ dx[i] + 0, dx[i] + 1, dx[i] + 2, x_[i] });

        // * '3' -> ~x'
        // If it can't be 'n' and '1' -> ~x'
        addClause({ dx[i] + 2, dx[i] + 3, -x_[i] });

        // * '5' -> ~x
        // If it can't be 'u' and '1' -> ~x
        addClause({ dx[i] + 1, dx[i] + 3, -x[i] });

        // * '7' -> ~x or ~x'
        addClause({ dx[i] + 3, -x[i], -x_[i] });

        // * 'A' -> x
        // If it can't be '0' and 'n' -> x
        addClause({ dx[i] + 0, dx[i] + 2, x[i] });

        // * 'B' -> x or ~x'
        addClause({ dx[i] + 2, x[i], -x_[i] });

        // * 'C' -> x'
        // If it can't be '0' and 'u' -> x'
        addClause({ dx[i] + 0, dx[i] + 1, x_[i] });

        // * 'D' -> ~x or x'
        addClause({ dx[i] + 1, -x[i], x_[i] });

        // * 'E' -> x or x'
        addClause({ dx[i] + 0, x[i], x_[i] });

        // * Can't be a '#'
        addClause({ dx[i] + 0, dx[i] + 1, dx[i] + 2, dx[i] + 3 });
    }
}

void xFormula::basic_1bit_rules(int dx[32], int x[32], int x_[32])
{
    xor2(dx, x, x_, 32);
}

void xFormula::basic_rules(int dx[32], int x[32], int x_[32])
{
#if IS_4bit
    basic_4bit_rules(dx, x, x_);
#else
    basic_1bit_rules(dx, x, x_);
#endif
}

inline void xFormula::impose_4bit_rule(vector<int> input_ids, vector<int> output_ids, pair<string, string> rule)
{
    string inputs_diff = rule.first, outputs_diff = rule.second;

    vector<int> antecedent;
    for (int x = 0; x < input_ids.size(); x++) {
        if (inputs_diff[x] == '?')
            continue;
        int base_id = input_ids[x];
        if (base_id == 0)
            continue;

        vector<int> values = get_values_4bit(inputs_diff[x]);
        for (int k = 0; k < 4; k++) {
            if (values[k] == 1)
                continue;

            int id = base_id + k;
            antecedent.push_back((values[k] == 1 ? -1 : 1) * id);
        }
    }

    for (int x = 0; x < output_ids.size(); x++) {
        if (outputs_diff[x] == '?')
            continue;
        vector<int> values = get_values_4bit(outputs_diff[x]);
        if (values.size() == 0)
            continue;

        int base_id = output_ids[x];
        assert(base_id > 0);
        for (int k = 0; k < 4; k++) {
            if (values[k] == 1)
                continue;
            vector<int> clause(antecedent);
            clause.push_back((values[k] == 1 ? 1 : -1) * (base_id + k));
            // printf("Clause (%d, %d): ", base_id, x);
            // for (auto& lit : clause)
            //     printf("%d ", lit);
            // printf("\n");
            addClause(clause);
        }
    }
}

inline void xFormula::impose_1bit_rule(vector<int> input_ids, vector<int> output_ids, pair<string, string> rule)
{
    string inputs_diff = rule.first, outputs_diff = rule.second;

    vector<int> antecedent;
    for (int x = 0; x < input_ids.size(); x++) {
        if (inputs_diff[x] == '?')
            continue;
        int id = input_ids[x];
        if (id == 0)
            continue;
        assert(inputs_diff[x] == '-' || inputs_diff[x] == 'x');
        antecedent.push_back((inputs_diff[x] == '-' ? 1 : -1) * id);
    }

    for (int x = 0; x < output_ids.size(); x++) {
        if (outputs_diff[x] == '?')
            continue;
        assert(outputs_diff[x] == '-' || outputs_diff[x] == 'x');
        assert(output_ids[x] > 0);
        vector<int> clause(antecedent);
        clause.push_back((outputs_diff[x] == '-' ? -1 : 1) * output_ids[x]);
        addClause(clause);
    }
}

void xFormula::impose_rule(vector<int> inputs, vector<int> outputs, pair<string, string> rule)
{
#if IS_4bit
    impose_4bit_rule(inputs, outputs, rule);
#else
    impose_1bit_rule(inputs, outputs, rule);
#endif
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