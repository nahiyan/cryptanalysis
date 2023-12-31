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

void xFormula::comp(Rules& rules, int z[4], int v[10][4], int n, int t[4], int T[4])
{
    assert(n >= 2 && n <= 7);
    if (n > 3)
        assert(T != NULL);

    unordered_map<string, string>* rules_;
    if (n == 2)
        rules_ = &rules.add2;
    else if (n == 3)
        rules_ = &rules.add3;
    else if (n == 4)
        rules_ = &rules.add4;
    else if (n == 5)
        rules_ = &rules.add5;
    else if (n == 6)
        rules_ = &rules.add6;
    else if (n == 7)
        rules_ = &rules.add7;

    for (auto& differential : *rules_) {
        string lhs = differential.first, rhs = differential.second;
        vector<int> base_clause;
        int i = 0;
        for (char& c : lhs) {
            assert(c == '-' || c == 'x');
            if (c == '-') {
                base_clause.push_back(-v[i][0]);
                base_clause.push_back(v[i][1]);
                base_clause.push_back(v[i][2]);
                base_clause.push_back(-v[i][3]);
            } else if (c == 'x') {
                base_clause.push_back(v[i][0]);
                base_clause.push_back(-v[i][1]);
                base_clause.push_back(-v[i][2]);
                base_clause.push_back(v[i][3]);
            }
            i++;
        }
        i = -1;
        for (char& c : rhs) {
            ++i;

            if (i == 0 && (n <= 3 || T == NULL))
                continue;

            int signs[4];
            if (c == '-') {
                signs[0] = 1;
                signs[1] = -1;
                signs[2] = -1;
                signs[3] = 1;
            } else if (c == 'x') {
                signs[0] = -1;
                signs[1] = 1;
                signs[2] = 1;
                signs[3] = -1;
            } else if (c == '0') {
                signs[0] = 1;
                signs[1] = -1;
                signs[2] = -1;
                signs[3] = -1;
            } else if (c == '7') {
                signs[0] = 1;
                signs[1] = 1;
                signs[2] = 1;
                signs[3] = -1;
            } else {
                continue;
            }

            int* diff;
            if (i == 0)
                diff = T;
            else if (i == 1)
                diff = t;
            else if (i == 2)
                diff = z;

            for (int j = 0; j < 4; j++) {
                vector<int> clause(base_clause);
                clause.push_back(signs[j] * diff[j]);
                addClause(clause);
            }
        }
    }
}

void xFormula::diff_add(Rules& rules, int z[32][4], int a[32][4], int b[32][4], int t[32][4], int T[32][4], int c[32][4], int d[32][4], int e[32][4])
{
    auto set = [](int x[4], int y[4]) {
        for (int i = 0; i < 4; i++)
            x[i] = y[i];
    };

    int n = 32;
    int m = 2;
    if (c)
        m++;
    if (d)
        m++;
    if (e)
        m++;
    int v[10][4], k;
    for (int j = 0; j < 32; j++) {
        k = 0;
        set(v[k++], a[j]);
        set(v[k++], b[j]);
        if (c)
            set(v[k++], c[j]);
        if (d)
            set(v[k++], d[j]);
        if (e)
            set(v[k++], e[j]);
        if (j > 0)
            set(v[k++], t[j - 1]);
        if ((m == 3 && j >= 3) || (j >= 2 && m > 3))
            set(v[k++], T[j - 2]);

        if (m == 2)
            comp(rules, z[j], v, k, t[j]);
        else
            comp(rules, z[j], v, k, t[j], T[j]);
    }
}

void xFormula::basic_rules(int dx[32][4], int x[32], int x_[32])
{
    for (int i = 0; i < 32; i++) {
        // Define the 4-bit differences in terms of each difference bit
        addClause({ x[i], x_[i], dx[i][0] });
        addClause({ -x[i], x_[i], dx[i][1] });
        addClause({ x[i], -x_[i], dx[i][2] });
        addClause({ -x[i], -x_[i], dx[i][3] });

        // // '-' -> x xnor x'
        // addClause({ -dx[i][0], dx[i][1], dx[i][2], -dx[i][3], -x[i], x_[i] });
        // addClause({ -dx[i][0], dx[i][1], dx[i][2], -dx[i][3], x[i], -x_[i] });

        // // 'x' -> x xor x'
        // addClause({ dx[i][0], -dx[i][1], -dx[i][2], dx[i][3], -x[i], -x_[i] });
        // addClause({ dx[i][0], -dx[i][1], -dx[i][2], dx[i][3], x[i], x_[i] });

        // // '0' -> ~x and ~x'
        // addClause({ -dx[i][0], dx[i][1], dx[i][2], dx[i][3], -x[i] });
        // addClause({ -dx[i][0], dx[i][1], dx[i][2], dx[i][3], -x_[i] });

        // // 'u' -> x and ~x'
        // addClause({ dx[i][0], -dx[i][1], dx[i][2], dx[i][3], x[i] });
        // addClause({ dx[i][0], -dx[i][1], dx[i][2], dx[i][3], -x_[i] });

        // // 'n' -> ~x and x'
        // addClause({ dx[i][0], dx[i][1], -dx[i][2], dx[i][3], -x[i] });
        // addClause({ dx[i][0], dx[i][1], -dx[i][2], dx[i][3], x_[i] });

        // // '1' -> x and x'
        // addClause({ dx[i][0], dx[i][1], dx[i][2], -dx[i][3], x[i] });
        // addClause({ dx[i][0], dx[i][1], dx[i][2], -dx[i][3], x_[i] });

        // // '3' -> ~x'
        // addClause({ -dx[i][0], -dx[i][1], dx[i][2], dx[i][3], -x_[i] });

        // // '5 -> ~x
        // addClause({ -dx[i][0], dx[i][1], -dx[i][2], dx[i][3], -x[i] });

        // // 'A' -> x
        // addClause({ dx[i][0], -dx[i][1], dx[i][2], -dx[i][3], x[i] });

        // // 'C' -> x'
        // addClause({ dx[i][0], dx[i][1], -dx[i][2], -dx[i][3], x_[i] });

        // // '7' -> ~x or ~x'
        // addClause({ -dx[i][0], -dx[i][1], -dx[i][2], dx[i][3], -x[i], -x_[i] });

        // // 'B' -> x or ~x'
        // addClause({ -dx[i][0], -dx[i][1], dx[i][2], -dx[i][3], x[i], -x_[i] });

        // // 'D' -> ~x or x'
        // addClause({ -dx[i][0], dx[i][1], -dx[i][2], -dx[i][3], -x[i], x_[i] });

        // // 'E' -> x or x'
        // addClause({ dx[i][0], -dx[i][1], -dx[i][2], -dx[i][3], x[i], x_[i] });
    }
}

void xFormula::impose_rule(vector<int (*)[32][4]> inputs, vector<int (*)[32][4]> outputs, pair<string, string> rule)
{
    auto get_value = [](char diff) -> vector<int> {
        switch (diff) {
        case '?':
            // ! IMPORTANT: It may be beneficial to not enforce '?'
            return {};
            // return { 1, 1, 1, 1 };
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

    for (int i = 0; i < 32; i++) {
        vector<int> base_clause;
        for (int x = 0; x < inputs.size(); x++) {
            vector<int> values = get_value(inputs_diff[x]);
            for (int k = 0; k < 4; k++)
                base_clause.push_back((values[k] == 1 ? -1 : 1) * (*inputs[x])[i][k]);
        }

        int q_count = 0;
        vector<vector<int>> clauses;
        for (int x = 0; x < outputs.size(); x++) {
            vector<int> values = get_value(outputs_diff[x]);
            if (values.size() == 0) {
                q_count++;
                continue;
            }

            for (int k = 0; k < 4; k++) {
                vector<int> clause(base_clause);
                clause.push_back((values[k] == 1 ? 1 : -1) * (*outputs[x])[i][k]);
                clauses.push_back(clause);
            }
        }

        if (q_count == outputs.size())
            continue;

        for (auto& clause : clauses)
            addClause(clause);
    }
}
