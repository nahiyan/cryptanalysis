#ifndef Crypto_h
#define Crypto_h

#include "Solver.h"
#include "SolverTypes.h"
#include "mtl/Vec.h"
#include <map>
#include <memory>

using namespace Minisat;

namespace Crypto {
typedef std::tuple<int, int, int> equation_t;
typedef std::vector<equation_t> equations_t;
typedef vec<Lit> minisat_clause_t;
typedef vec<vec<Lit>> minisat_clauses_t;
typedef std::tuple<int, int, int> triple_t;

struct FunctionResult {
    int operation_id;
    int functon_id;
    std::vector<int> variables;
};

struct State {
    minisat_clauses_t& out_refined;
    Solver& solver;
    int k = 0;

    std::shared_ptr<equations_t> equations;
    std::map<int, int> eq_var_map;
    std::map<equation_t, std::vector<FunctionResult>> eq_func_rels;
};

void add_clauses(State& state);
void load_rules(Solver& solver, const char* filename);
void load_rule(Solver& solver, FILE*& db, int& id);
void process_var_map(Solver& solver);
std::shared_ptr<equations_t> check_consistency(std::shared_ptr<equations_t>& equations, bool exhaustive = true);
bool block_inconsistency(State& state);
void print(equations_t equations);
};
#endif