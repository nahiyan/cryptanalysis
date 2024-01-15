#include "formula.h"
#include <assert.h>
#include <errno.h>
#include <math.h>
#include <stdexcept>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

Formula::Formula(string name)
{
    varID = 0;
    varCnt = 0;
    useXORClauses = false;
    pbMethod = SEQUENTIAL_COUNTER;
    adderType = RIPPLE_CARRY;
    multiAdderType = ESPRESSO;
    formulaName = name;
}

Formula::~Formula()
{
}

void Formula::varName(int* x, string name, int offset)
{
    varNames[name + "_" + formulaName] = x[0] + offset;
}

void Formula::newVars(int* x, int n, string name)
{
    for (int i = 0; i < n; i++)
        x[i] = ++varID;

    if (name != "")
        varNames[name + "_" + formulaName] = x[0];
    varCnt += n;
}
void Formula::newDiff(int x[32][4], string name)
{
    for (int i = 0; i < 32; i++)
        for (int j = 0; j < 4; j++)
            x[i][j] = ++varID;

    if (name != "")
        varNames[name + "_" + formulaName] = x[0][0];
    varCnt += 32 * 4;
}

void Formula::addClause(vector<int> v)
{
    assert(v.size() > 0);
    if (any_of(v.begin(), v.end(), [](int x) { return x == 0; })) {
        fprintf(stderr, "bad vector clause:");
        for (int x : v)
            fprintf(stderr, " %d", x);
        fprintf(stderr, "\n");
        exit(1);
    }
    clauses.push_back(Clause(v));
}

void Formula::addClause(Clause c)
{
    assert(c.lits.size() > 0);
    if (any_of(c.lits.begin(), c.lits.end(), [](int x) { return x == 0; })) {
        fprintf(stderr, "bad clause:");
        for (int x : c.lits)
            fprintf(stderr, " %d", x);
        fprintf(stderr, "\n");
        exit(1);
    }
    clauses.push_back(c);
}

void Formula::fixedValue(int* z, unsigned value, int n)
{
    for (int i = 0; i < n; i++) {
        int x = (value >> i) & 1 ? z[i] : -z[i];
        addClause({ x });
    }
}

void Formula::rotl(int* z, int* x, int p, int n)
{
    for (int i = 0; i < n; i++)
        z[i] = x[(i + n - p) % n];
}

void Formula::eq(int* z, int* x, int n)
{
    for (int i = 0; i < n; i++) {
        addClause({ -z[i], x[i] });
        addClause({ z[i], -x[i] });
    }
}

void Formula::neq(int* z, int* x, int n)
{
    for (int i = 0; i < n; i++) {
        addClause({ z[i], x[i] });
        addClause({ -z[i], -x[i] });
    }
}

void Formula::xor2(int* z, int* x, int* y, int n)
{
    for (int i = 0; i < n; i++) {
        if (useXORClauses) {
            addClause(Clause({ -z[i], x[i], y[i] }, true));
        } else {
            addClause({ -z[i], -x[i], -y[i] });
            addClause({ z[i], -x[i], y[i] });
            addClause({ z[i], x[i], -y[i] });
            addClause({ -z[i], x[i], y[i] });
        }
    }
}

// p -> q
void Formula::implication(vector<int> p, vector<int> q)
{
    vector<int> clause;
    for (auto& var : p)
        clause.push_back(-var);
    for (auto& var : q)
        clause.push_back(var);
    addClause(clause);
}

// TODO: Inject XOR rules if these are difference variables
void Formula::xor3(int* z, int* x, int* y, int* t, int n)
{
    for (int i = 0; i < n; i++) {
        if (useXORClauses) {
            addClause(Clause({ -z[i], x[i], y[i], t[i] }, true));
        } else {
            addClause({ z[i], -x[i], -y[i], -t[i] });
            addClause({ -z[i], -x[i], -y[i], t[i] });
            addClause({ -z[i], -x[i], y[i], -t[i] });
            addClause({ z[i], -x[i], y[i], t[i] });
            addClause({ -z[i], x[i], -y[i], -t[i] });
            addClause({ z[i], x[i], -y[i], t[i] });
            addClause({ z[i], x[i], y[i], -t[i] });
            addClause({ -z[i], x[i], y[i], t[i] });
        }
    }
}

void Formula::xor4(int* z, int* a, int* b, int* c, int* d, int n)
{
    for (int i = 0; i < n; i++) {
        if (useXORClauses) {
            addClause(Clause({ -z[i], a[i], b[i], c[i], d[i] }, true));
        } else {
            addClause({ -z[i], -a[i], -b[i], -c[i], -d[i] });
            addClause({ z[i], -a[i], -b[i], -c[i], d[i] });
            addClause({ z[i], -a[i], -b[i], c[i], -d[i] });
            addClause({ -z[i], -a[i], -b[i], c[i], d[i] });
            addClause({ z[i], -a[i], b[i], -c[i], -d[i] });
            addClause({ -z[i], -a[i], b[i], -c[i], d[i] });
            addClause({ -z[i], -a[i], b[i], c[i], -d[i] });
            addClause({ z[i], -a[i], b[i], c[i], d[i] });
            addClause({ z[i], a[i], -b[i], -c[i], -d[i] });
            addClause({ -z[i], a[i], -b[i], -c[i], d[i] });
            addClause({ -z[i], a[i], -b[i], c[i], -d[i] });
            addClause({ z[i], a[i], -b[i], c[i], d[i] });
            addClause({ -z[i], a[i], b[i], -c[i], -d[i] });
            addClause({ z[i], a[i], b[i], -c[i], d[i] });
            addClause({ z[i], a[i], b[i], c[i], -d[i] });
            addClause({ -z[i], a[i], b[i], c[i], d[i] });
        }
    }
}

void Formula::ch(int* z, int* x, int* y, int* t, int n)
{
    for (int i = 0; i < n; i++) {
        addClause({ -z[i], x[i], t[i] });
        addClause({ -z[i], -x[i], y[i] });
        addClause({ z[i], x[i], -t[i] });
        addClause({ z[i], -x[i], -y[i] });
    }
}

void Formula::maj3(int* z, int* x, int* y, int* t, int n)
{
    for (int i = 0; i < n; i++) {
        addClause({ -z[i], x[i], y[i] });
        addClause({ -z[i], x[i], t[i] });
        addClause({ -z[i], y[i], t[i] });
        addClause({ z[i], -y[i], -t[i] });
        addClause({ z[i], -x[i], -t[i] });
        addClause({ z[i], -x[i], -y[i] });
    }
}

void Formula::espresso(const vector<int>& lhs, const vector<int>& rhs)
{
    static map<pair<unsigned int, unsigned int>, vector<vector<int>>> cache;

    unsigned int n = lhs.size();
    unsigned int m = rhs.size();

    vector<vector<int>> _clauses;
    auto it = cache.find(make_pair(n, m));
    if (it != cache.end()) {
        _clauses = it->second;
    } else {
        int wfd[2], rfd[2];

        /* pipe(): fd[0] is for reading, fd[1] is for writing */

        if (pipe(wfd) == -1)
            throw std::runtime_error("pipe() failed");

        if (pipe(rfd) == -1)
            throw std::runtime_error("pipe() failed");

        pid_t child = fork();
        if (child == 0) {
            if (dup2(wfd[0], STDIN_FILENO) == -1)
                throw std::runtime_error("dup() failed");

            if (dup2(rfd[1], STDOUT_FILENO) == -1)
                throw std::runtime_error("dup() failed");

            if (execlp("espresso", "espresso", 0) == -1)
                throw std::runtime_error("execve() failed");

            exit(EXIT_FAILURE);
        }

        close(wfd[0]);
        close(rfd[1]);

        FILE* eout = fdopen(wfd[1], "w");
        if (!eout)
            throw std::runtime_error("fdopen() failed");

        FILE* ein = fdopen(rfd[0], "r");
        if (!ein)
            throw std::runtime_error("fdopen() failed");

        fprintf(eout, ".i %u\n", n + m);
        fprintf(eout, ".o 1\n");

        for (unsigned int i = 0; i < 1U << n; ++i) {
            for (unsigned int j = 0; j < 1U << m; ++j) {
                for (unsigned int k = n; k--;)
                    fprintf(eout, "%u", 1 - ((i >> k) & 1));
                for (unsigned int k = m; k--;)
                    fprintf(eout, "%u", 1 - ((j >> k) & 1));

                fprintf(eout, " %u\n", __builtin_popcount(i) != j);
            }
        }

        fprintf(eout, ".e\n");
        fflush(eout);

        while (1) {
            char buf[512];
            if (!fgets(buf, sizeof(buf), ein))
                break;

            if (!strncmp(buf, ".i", 2))
                continue;
            if (!strncmp(buf, ".o", 2))
                continue;
            if (!strncmp(buf, ".p", 2))
                continue;
            if (!strncmp(buf, ".e", 2))
                break;

            vector<int> c;
            for (int i = 0; i < n + m; ++i) {
                if (buf[i] == '0')
                    c.push_back(-(i + 1));
                else if (buf[i] == '1')
                    c.push_back(i + 1);
            }

            _clauses.push_back(c);
        }

        fclose(ein);
        fclose(eout);

        while (true) {
            int status;
            pid_t kid = wait(&status);
            if (kid == -1) {
                if (errno == ECHILD)
                    break;
                if (errno == EINTR)
                    continue;

                throw std::runtime_error("wait() failed");
            }

            if (kid == child)
                break;
        }

        cache.insert(make_pair(make_pair(n, m), _clauses));

#ifdef _DUMP_ADDER_CLAUSES_
        FILE* f = fopen("comp_clauses.txt", "a");
        fprintf(f, "%d %d %d\n", n, m, _clauses.size());
        for (vector<int>& c : _clauses) {
            for (int i : c)
                fprintf(f, "%d ", i);
            fprintf(f, "\n");
        }
        fclose(f);
#endif
    }

    for (vector<int>& c : _clauses) {
        Clause cl;
        for (int i : c) {
            int j = abs(i) - 1;
            int var = j < n ? lhs[j] : rhs[m - 1 - (j - n)];
            if (i < 0)
                cl.lits.push_back(-var);
            else
                cl.lits.push_back(var);
        }

        addClause(cl);
    }
}

void Formula::dimacs(int rounds, string fileName, bool header)
{
    FILE* out = fileName == "" ? stdout : fopen(fileName.c_str(), "w");
    if (out == NULL) {
        fprintf(stderr, "Failed to open %s to write!\n", fileName.c_str());
        exit(1);
    }

    if (header)
        fprintf(out, "p cnf %d %d\n", getVarCnt(), getClauseCnt());

    for (Clause c : clauses) {
        if (c.xor_clause)
            fprintf(out, "x ");
        for (int v : c.lits)
            fprintf(out, "%d ", v);
        fprintf(out, "0\n");
    }

    for (auto e : varNames)
        fprintf(out, "c %s %d\n", e.first.c_str(), e.second);

    printf("c order %d\n", rounds);

    fclose(out);
}

int Formula::clauseCheck()
{
    for (Clause c : clauses) {
        for (int v : c.lits) {
            if (abs(v) > getVarCnt()) {
                fprintf(stderr, "Clause check failed: out of bound variable ID (%d)! var_cnt == %d\n", v, getVarCnt());
                abort();
            }
            if (v == 0) {
                fprintf(stderr, "Clause check failed: variable ID is zero!\n");
                abort();
            }
        }
    }
    return 0;
}

vector<Clause> Formula::getClauses()
{
    return clauses;
}

void Formula::AddFormula(Formula& f)
{
    varCnt += f.getVarCnt();
    vector<Clause> c = f.getClauses();
    clauses.insert(clauses.end(), c.begin(), c.end());
    for (auto e : f.varNames)
        varNames[e.first] = e.second;
}
