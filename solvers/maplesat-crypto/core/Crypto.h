#ifndef Crypto_h
#define Crypto_h

#include "Solver.h"
#include "SolverTypes.h"
#include "mtl/Vec.h"

using namespace Minisat;

void add_clauses(Minisat::Solver& s, vec<vec<Lit>>& out_refined);
void loadRules(Minisat::Solver& solver, const char* filename);
void loadRule(Minisat::Solver& solver, FILE*& db, int& id);
void process_var_map(Solver& solver);

#endif