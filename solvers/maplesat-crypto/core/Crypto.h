#ifndef Crypto_h
#define Crypto_h

#include "Solver.h"
#include "SolverTypes.h"
#include "mtl/Vec.h"
#include <map>

using namespace Minisat;

namespace Crypto {
typedef std::pair<int32_t, int32_t> equation_t;
typedef std::vector<equation_t> equations_t;
typedef vec<Lit> minisat_clause_t;
typedef vec<vec<Lit>> minisat_clauses_t;

struct Triple {
    int x, x_, dx;
};

struct FunctionResult {
    int operation_id;
    int functon_id;
    std::vector<Triple> inputs;
    std::vector<Triple> outputs;
};

struct State {
    minisat_clauses_t& out_refined;
    Solver& solver;
    int k = 0;
    bool has_conflict = false;
    minisat_clauses_t conflict_clauses;
    equations_t equations;
    std::unordered_map<int32_t, std::vector<FunctionResult>> var_func_relations;
};

void add_clauses(State& state);
void load_rules(Solver& solver, const char* filename);
void load_rule(Solver& solver, FILE*& db, int& id);
void process_var_map(Solver& solver);
bool check_consistency(equations_t& equations);
};
#endif