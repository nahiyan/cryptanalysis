#include "../core/Crypto.h"
#include "../core/Solver.h"
#include "../core/SolverTypes.h"
#include <stdio.h>
#include <map>
#include <memory>
#include <set>

#define X 3
#define Y 8
#define Z 7
#define U 6

#define A 1
#define B 2
#define C 3
#define D 4
#define E 5
#define F 6
#define G 7
#define H 8
#define O1 9
#define O2 10
#define O3 11
#define O4 12

using namespace Crypto;
using namespace std;

void test_inconsistency_blocker()
{
    minisat_clauses_t clauses;
    Minisat::Solver solver;
    solver.newVar(false);
    solver.newVar(false);
    solver.newVar(false);
    solver.newVar(false);
    solver.newVar(false);
    solver.newVar(false);
    solver.newVar(false);
    solver.newVar(false);
    solver.newVar(false);
    solver.newVar(false);
    solver.newVar(false);
    solver.newVar(false);
    solver.newVar(false);
    solver.addClause(mkLit(A, true));
    solver.addClause(mkLit(B, true));
    solver.addClause(mkLit(C, true));
    solver.addClause(mkLit(D, true));
    solver.addClause(mkLit(E, true));
    solver.addClause(mkLit(F, true));
    solver.addClause(mkLit(G, false));
    solver.addClause(mkLit(H, true));
    solver.addClause(mkLit(O1, true));
    solver.addClause(mkLit(O2, true));
    solver.addClause(mkLit(O3, true));
    solver.addClause(mkLit(O4, false));
    solver.solve();

    printf("Values: %d %d %d %d %d %d %d %d %d %d %d %d\n", solver.value(A), solver.value(B), solver.value(C), solver.value(D), solver.value(E), solver.value(F), solver.value(G), solver.value(H), solver.value(O1), solver.value(O2), solver.value(O3), solver.value(O4));

    auto equations = make_shared<equations_t>();
    equations->push_back(equation_t { X, Y, 0 });
    equations->push_back(equation_t { Y, Z, 0 });
    equations->push_back(equation_t { Z, U, 0 });
    equations->push_back(equation_t { U, X, 1 });

    auto evm = map<int, int>();
    evm[X] = 0;
    evm[Y] = 1;
    evm[Z] = 2;
    evm[U] = 3;

    auto efr = map<equation_t, vector<FunctionResult>>();
    efr[(*equations)[0]] = { FunctionResult {
        19,
        0,
        { G, A, B, O1 },
    } };
    efr[(*equations)[1]] = { FunctionResult {
        18,
        0,
        { C, B, F, O2 },
    } };
    efr[(*equations)[2]] = { FunctionResult {
        19,
        0,
        { H, D, C, O3 },
    } };
    efr[(*equations)[3]] = { FunctionResult {
        18,
        0,
        { D, A, E, O4 },
    } };

    auto state = State {
        out_refined : clauses,
        solver : solver,
        equations : equations,
        eq_var_map : evm,
        eq_func_rels : efr,
    };

    auto is_blocked = block_inconsistency(state);
    assert(is_blocked);
}

void test_inconsistency_checker()
{
    // Test an inconsistent system
    {
        auto equations = make_shared<equations_t>();
        equations->push_back(equation_t { X, Y, 0 });
        equations->push_back(equation_t { Y, Z, 0 });
        equations->push_back(equation_t { Z, U, 0 });
        equations->push_back(equation_t { U, X, 1 });
        equations->push_back(equation_t { U, 100, 1 });
        equations->push_back(equation_t { U, 101, 0 });
        equations->push_back(equation_t { 101, 100, 0 });
        auto conflict_equations = check_consistency(equations);
        assert((*conflict_equations)[0] == std::make_tuple(6, 3, 1));
        assert((*conflict_equations)[1] == std::make_tuple(101, 100, 0));
    }

    // Test an consistent system
    {
        auto equations = make_shared<equations_t>();
        equations->push_back(equation_t { X, Y, 0 });
        equations->push_back(equation_t { Y, Z, 0 });
        equations->push_back(equation_t { Z, U, 0 });
        equations->push_back(equation_t { U, X, 0 });
        equations->push_back(equation_t { U, 100, 0 });
        equations->push_back(equation_t { U, 101, 0 });
        equations->push_back(equation_t { 101, 100, 0 });
        auto conflict_equations = check_consistency(equations);
        assert(conflict_equations.get()->size() == 0);
    }
}

int main()
{
    // test_inconsistency_blocker();
    test_inconsistency_checker();

    return 0;
}