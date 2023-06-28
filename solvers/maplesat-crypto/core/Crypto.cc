#include "Crypto.h"

int int_value(Minisat::Solver &s, int var) {
  auto value = s.value(var);
  return value == l_True ? 1 : value == l_False ? 0 : -1;
}

void add_impl2(Minisat::Solver &s,
               Minisat::vec<Minisat::vec<Minisat::Lit>> &out_refined, int &k,
               int i, int j, const char *name, int a, int b, int c, bool v1,
               bool v2, bool v3) {
  printf("Add %s; i,j = %d,%d; vars: %d %d %d; vals: %d %d %d\n", name, i, j,
         a + 1, b + 1, c + 1, int_value(s, a), int_value(s, b),
         int_value(s, c));
  out_refined.push();
  out_refined[k].push(mkLit(a, v1));
  out_refined[k].push(mkLit(b, v2));
  out_refined[k].push(mkLit(c, v3));
  k++;
}
void add_impl3(Minisat::Solver &s,
               Minisat::vec<Minisat::vec<Minisat::Lit>> &out_refined, int &k,
               int i, int j, const char *name, int a, int b, int c, int d,
               bool v1, bool v2, bool v3, bool v4) {
  printf("Add %s; i,j = %d,%d; vars: %d %d %d %d; vals: %d %d %d %d\n", name, i,
         j, a + 1, b + 1, c + 1, d + 1, int_value(s, a), int_value(s, b),
         int_value(s, c), int_value(s, d));
  out_refined.push();
  out_refined[k].push(mkLit(a, v1));
  out_refined[k].push(mkLit(b, v2));
  out_refined[k].push(mkLit(c, v3));
  out_refined[k].push(mkLit(d, v4));
  k++;
}

void add_impl4(Minisat::Solver &s,
               Minisat::vec<Minisat::vec<Minisat::Lit>> &out_refined, int &k,
               int i, int j, const char *name, int a, int b, int c, int d,
               int e, bool v1, bool v2, bool v3, bool v4, bool v5) {
  printf("Add %s; i,j = %d,%d; vars: %d %d %d %d %d; vals: %d %d %d %d %d\n",
         name, i, j, a + 1, b + 1, c + 1, d + 1, e + 1, int_value(s, a),
         int_value(s, b), int_value(s, c), int_value(s, d), int_value(s, e));
  out_refined.push();
  out_refined[k].push(mkLit(a, v1));
  out_refined[k].push(mkLit(b, v2));
  out_refined[k].push(mkLit(c, v3));
  out_refined[k].push(mkLit(d, v4));
  out_refined[k].push(mkLit(e, v5));
  k++;
}

void add_impl5(Minisat::Solver &s,
               Minisat::vec<Minisat::vec<Minisat::Lit>> &out_refined, int &k,
               int i, int j, const char *name, int a, int b, int c, int d,
               int e, int f, bool v1, bool v2, bool v3, bool v4, bool v5,
               bool v6) {
  printf("Add %s; i,j = %d,%d; vars: %d %d %d %d %d %d; vals: %d %d %d %d %d "
         "%d\n",
         name, i, j, a + 1, b + 1, c + 1, d + 1, e + 1, f + 1, int_value(s, a),
         int_value(s, b), int_value(s, c), int_value(s, d), int_value(s, e),
         int_value(s, f));
  out_refined.push();
  out_refined[k].push(mkLit(a, v1));
  out_refined[k].push(mkLit(b, v2));
  out_refined[k].push(mkLit(c, v3));
  out_refined[k].push(mkLit(d, v4));
  out_refined[k].push(mkLit(e, v5));
  out_refined[k].push(mkLit(f, v6));
  printf("Info: %d %d %d %d %d %d\n", !v1, !v2, !v3, !v4, !v5, v6);
  k++;
}

void add_impl6(Minisat::Solver &s,
               Minisat::vec<Minisat::vec<Minisat::Lit>> &out_refined, int &k,
               int i, int j, const char *name, int a, int b, int c, int d,
               int e, int f, int g, bool v1, bool v2, bool v3, bool v4, bool v5,
               bool v6, bool v7) {
  printf("Added %s; i,j = %d,%d; vars: %d %d %d %d %d %d %d; vals: %d %d %d %d "
         "%d %d %d\n",
         name, i, j, a + 1, b + 1, c + 1, d + 1, e + 1, f + 1, g + 1,
         int_value(s, a), int_value(s, b), int_value(s, c), int_value(s, d),
         int_value(s, e), int_value(s, f), int_value(s, g));
  out_refined.push();
  out_refined[k].push(mkLit(a, v1));
  out_refined[k].push(mkLit(b, v2));
  out_refined[k].push(mkLit(c, v3));
  out_refined[k].push(mkLit(d, v4));
  out_refined[k].push(mkLit(e, v5));
  out_refined[k].push(mkLit(f, v6));
  out_refined[k].push(mkLit(g, v7));
  k++;
}

void add_impl7(Minisat::Solver &s,
               Minisat::vec<Minisat::vec<Minisat::Lit>> &out_refined, int &k,
               int i, int j, const char *name, int a, int b, int c, int d,
               int e, int f, int g, int h, bool v1, bool v2, bool v3, bool v4,
               bool v5, bool v6, bool v7, bool v8) {
  printf("Added %s; i,j = %d,%d; vars: %d %d %d %d %d %d %d %d; vals: %d %d %d "
         "%d %d %d %d %d\n",
         name, i, j, a + 1, b + 1, c + 1, d + 1, e + 1, f + 1, g + 1, h + 1,
         int_value(s, a), int_value(s, b), int_value(s, c), int_value(s, d),
         int_value(s, e), int_value(s, f), int_value(s, g), int_value(s, h));
  out_refined.push();
  out_refined[k].push(mkLit(a, v1));
  out_refined[k].push(mkLit(b, v2));
  out_refined[k].push(mkLit(c, v3));
  out_refined[k].push(mkLit(d, v4));
  out_refined[k].push(mkLit(e, v5));
  out_refined[k].push(mkLit(f, v6));
  out_refined[k].push(mkLit(g, v7));
  out_refined[k].push(mkLit(h, v8));
  k++;
}

void comp_7_3(Minisat::Solver &s,
              Minisat::vec<Minisat::vec<Minisat::Lit>> &out_refined, int &k,
              int &i, int &j, int &op1, int &op2, int &op3, int &op4, int &op5,
              int &op6, int &op7, int &o1, int &o2, int &o3) {
  int xs = 0, nxs = 0;
  xs += s.value(op1) == l_True ? 1 : 0;
  xs += s.value(op2) == l_True ? 1 : 0;
  xs += s.value(op3) == l_True ? 1 : 0;
  xs += s.value(op4) == l_True ? 1 : 0;
  xs += s.value(op5) == l_True ? 1 : 0;
  xs += s.value(op6) == l_True ? 1 : 0;
  xs += s.value(op7) == l_True ? 1 : 0;
  nxs += s.value(op1) == l_False ? 1 : 0;
  nxs += s.value(op2) == l_False ? 1 : 0;
  nxs += s.value(op3) == l_False ? 1 : 0;
  nxs += s.value(op4) == l_False ? 1 : 0;
  nxs += s.value(op5) == l_False ? 1 : 0;
  nxs += s.value(op6) == l_False ? 1 : 0;
  nxs += s.value(op7) == l_False ? 1 : 0;

  bool o1_nf = int_value(s, o1) != 0;
  bool o1_nt = int_value(s, o1) != 1;
  bool o2_nf = int_value(s, o2) != 0;
  bool o2_nt = int_value(s, o2) != 1;
  bool o3_nf = int_value(s, o3) != 0;
  bool o3_nt = int_value(s, o3) != 1;

  if (nxs == 7 && xs == 0 && o1_nf) {
    // ----- -> ---
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 1", op1, op2, op3, op4, op5,
              op6, op7, o1, true, true, true, true, true, true, true, false);
  }
  if (nxs == 7 && xs == 0 && o2_nf) {
    // ----- -> ---
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 2", op1, op2, op3, op4, op5,
              op6, op7, o2, true, true, true, true, true, true, true, false);
  }
  if (nxs == 7 && xs == 0 && o3_nf) {
    // ----- -> ---
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 3", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, true, true, true, true, true, false);
  }

  if (nxs == 6 && xs == 1 && o3_nt) {
    // ------x -> ??x
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 4", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, true, true, true, true, true,
              true); // 0111111
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 4", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, true, true, true, true, true,
              true); // 1011111
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 4", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, false, true, true, true, true,
              true); // 1101111
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 4", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, true, false, true, true, true,
              true); // 1110111
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 4", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, true, true, false, true, true,
              true); // 1111011
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 4", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, true, true, true, false, true,
              true); // 1111101
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 4", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, true, true, true, true, false,
              true); // 1111110
  } else if (nxs == 5 && xs == 2 && o3_nf) {
    // -----xx -> ??-
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 5", op1, op2, op3, op4, op5,
              op6, op7, o3, false, false, true, true, true, true, true, false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 5", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, false, true, true, true, true, false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 5", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, true, false, true, true, true, false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 5", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, true, true, false, true, true, false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 5", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, true, true, true, false, true, false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 5", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, true, true, true, true, false, false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 5", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, false, true, true, true, true, false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 5", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, true, false, true, true, true, false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 5", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, true, true, false, true, true, false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 5", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, true, true, true, false, true, false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 5", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, true, true, true, true, false, false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 5", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, false, false, true, true, true, false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 5", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, false, true, false, true, true, false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 5", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, false, true, true, false, true, false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 5", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, false, true, true, true, false, false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 5", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, true, false, false, true, true, false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 5", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, true, false, true, false, true, false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 5", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, true, false, true, true, false, false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 5", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, true, true, false, false, true, false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 5", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, true, true, false, true, false, false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 5", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, true, true, true, false, false, false);
  } else if (nxs == 4 && xs == 3 && o3_nt) {
    // ----xxx -> ??x
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, false, false, false, true, true, true, true, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, false, false, true, false, true, true, true, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, false, false, true, true, false, true, true, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, false, false, true, true, true, false, true, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, false, false, true, true, true, true, false, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, false, false, true, true, true, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, false, true, false, true, true, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, false, true, true, false, true, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, false, true, true, true, false, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, true, false, false, true, true, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, true, false, true, false, true, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, true, false, true, true, false, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, true, true, false, false, true, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, true, true, false, true, false, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, true, true, true, false, false, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, false, false, true, true, true, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, false, true, false, true, true, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, false, true, true, false, true, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, false, true, true, true, false, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, true, false, false, true, true, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, true, false, true, false, true, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, true, false, true, true, false, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, true, true, false, false, true, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, true, true, false, true, false, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, true, true, true, false, false, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, false, false, false, true, true, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, false, false, true, false, true, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, false, false, true, true, false, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, false, true, false, false, true, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, false, true, false, true, false, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, false, true, true, false, false, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, true, false, false, false, true, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, true, false, false, true, false, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, true, false, true, false, false, true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 6", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, true, true, false, false, false, true);
  } else if (nxs == 3 && xs == 4 && o3_nf) {
    // ---xxxx -> ??-
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, false, false, false, false, true, true, true,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, false, false, false, true, false, true, true,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, false, false, false, true, true, false, true,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, false, false, false, true, true, true, false,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, false, false, true, false, false, true, true,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, false, false, true, false, true, false, true,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, false, false, true, false, true, true, false,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, false, false, true, true, false, false, true,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, false, false, true, true, false, true, false,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, false, false, true, true, true, false, false,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, false, false, false, true, true,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, false, false, true, false, true,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, false, false, true, true, false,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, false, true, false, false, true,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, false, true, false, true, false,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, false, true, true, false, false,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, true, false, false, false, true,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, true, false, false, true, false,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, true, false, true, false, false,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, true, true, false, false, false,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, false, false, false, true, true,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, false, false, true, false, true,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, false, false, true, true, false,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, false, true, false, false, true,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, false, true, false, true, false,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, false, true, true, false, false,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, true, false, false, false, true,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, true, false, false, true, false,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, true, false, true, false, false,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, true, true, false, false, false,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, false, false, false, false, true,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, false, false, false, true, false,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, false, false, true, false, false,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, false, true, false, false, false,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 7", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, true, false, false, false, false,
              false);
  } else if (nxs == 2 && xs == 5 && o3_nt) {
    // --xxxxx -> ??x
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 8", op1, op2, op3, op4, op5,
              op6, op7, o3, false, false, false, false, false, true, true,
              true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 8", op1, op2, op3, op4, op5,
              op6, op7, o3, false, false, false, false, true, false, true,
              true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 8", op1, op2, op3, op4, op5,
              op6, op7, o3, false, false, false, false, true, true, false,
              true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 8", op1, op2, op3, op4, op5,
              op6, op7, o3, false, false, false, true, false, false, true,
              true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 8", op1, op2, op3, op4, op5,
              op6, op7, o3, false, false, false, true, false, true, false,
              true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 8", op1, op2, op3, op4, op5,
              op6, op7, o3, false, false, false, true, true, false, false,
              true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 8", op1, op2, op3, op4, op5,
              op6, op7, o3, false, false, true, false, false, false, true,
              true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 8", op1, op2, op3, op4, op5,
              op6, op7, o3, false, false, true, false, false, true, false,
              true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 8", op1, op2, op3, op4, op5,
              op6, op7, o3, false, false, true, false, true, false, false,
              true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 8", op1, op2, op3, op4, op5,
              op6, op7, o3, false, false, true, true, false, false, false,
              true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 8", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, false, false, false, false, true,
              true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 8", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, false, false, false, true, false,
              true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 8", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, false, false, true, false, false,
              true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 8", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, false, true, false, false, false,
              true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 8", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, true, false, false, false, false,
              true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 8", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, false, false, false, false, true,
              true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 8", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, false, false, false, true, false,
              true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 8", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, false, false, true, false, false,
              true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 8", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, false, true, false, false, false,
              true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 8", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, true, false, false, false, false,
              true);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 8", op1, op2, op3, op4, op5,
              op6, op7, o3, true, true, false, false, false, false, false,
              true);
  } else if (nxs == 1 && xs == 6 && o3_nf) {
    // -xxxxxx -> ??-
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 9", op1, op2, op3, op4, op5,
              op6, op7, o3, false, false, false, false, false, false, true,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 9", op1, op2, op3, op4, op5,
              op6, op7, o3, false, false, false, false, false, true, false,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 9", op1, op2, op3, op4, op5,
              op6, op7, o3, false, false, false, false, true, false, false,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 9", op1, op2, op3, op4, op5,
              op6, op7, o3, false, false, false, true, false, false, false,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 9", op1, op2, op3, op4, op5,
              op6, op7, o3, false, false, true, false, false, false, false,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 9", op1, op2, op3, op4, op5,
              op6, op7, o3, false, true, false, false, false, false, false,
              false);
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 9", op1, op2, op3, op4, op5,
              op6, op7, o3, true, false, false, false, false, false, false,
              false);
  }

  if (nxs == 0 && xs == 7 && o1_nt) {
    // xxxxxxx -> xxx
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 10", op1, op2, op3, op4, op5,
              op6, op7, o1, false, false, false, false, false, false, false,
              true);
  }
  if (nxs == 0 && xs == 7 && o2_nt) {
    // xxxxxxx -> xxx
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 11", op1, op2, op3, op4, op5,
              op6, op7, o1, false, false, false, false, false, false, false,
              true);
  }
  if (nxs == 0 && xs == 7 && o3_nt) {
    // xxxxxxx -> xxx
    add_impl7(s, out_refined, k, i, j, "COMP_7_3 12", op1, op2, op3, op4, op5,
              op6, op7, o1, false, false, false, false, false, false, false,
              true);
  }
}

void comp_6_3(Minisat::Solver &s,
              Minisat::vec<Minisat::vec<Minisat::Lit>> &out_refined, int &k,
              int &i, int &j, int &op1, int &op2, int &op3, int &op4, int &op5,
              int &op6, int &o1, int &o2, int &o3) {
  int xs = 0, nxs = 0;
  xs += s.value(op1) == l_True ? 1 : 0;
  xs += s.value(op2) == l_True ? 1 : 0;
  xs += s.value(op3) == l_True ? 1 : 0;
  xs += s.value(op4) == l_True ? 1 : 0;
  xs += s.value(op5) == l_True ? 1 : 0;
  xs += s.value(op6) == l_True ? 1 : 0;
  nxs += s.value(op1) == l_False ? 1 : 0;
  nxs += s.value(op2) == l_False ? 1 : 0;
  nxs += s.value(op3) == l_False ? 1 : 0;
  nxs += s.value(op4) == l_False ? 1 : 0;
  nxs += s.value(op5) == l_False ? 1 : 0;
  nxs += s.value(op6) == l_False ? 1 : 0;

  bool o1_nf = int_value(s, o1) != 0;
  bool o2_nf = int_value(s, o2) != 0;
  bool o3_nf = int_value(s, o3) != 0;
  bool o3_nt = int_value(s, o3) != 1;

  if (xs == 0 && nxs == 6 && o1_nf) {
    // ------ -> ---
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 1", op1, op2, op3, op4, op5,
              op6, o1, true, true, true, true, true, true, false);
  }
  if (xs == 0 && nxs == 6 && o2_nf) {
    // ------ -> ---
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 2", op1, op2, op3, op4, op5,
              op6, o2, true, true, true, true, true, true, false);
  }
  if (xs == 0 && nxs == 6 && o3_nf) {
    // ------ -> ---
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 3", op1, op2, op3, op4, op5,
              op6, o3, true, true, true, true, true, true, false);
  }

  if (nxs == 5 && xs == 1 && o3_nt) {
    // -----x -> ??x
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 4", op1, op2, op3, op4, op5,
              op6, o3, false, true, true, true, true, true, true);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 4", op1, op2, op3, op4, op5,
              op6, o3, true, false, true, true, true, true, true);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 4", op1, op2, op3, op4, op5,
              op6, o3, true, true, false, true, true, true, true);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 4", op1, op2, op3, op4, op5,
              op6, o3, true, true, true, false, true, true, true);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 4", op1, op2, op3, op4, op5,
              op6, o3, true, true, true, true, false, true, true);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 4", op1, op2, op3, op4, op5,
              op6, o3, true, true, true, true, true, false, true);
  } else if (nxs == 4 && xs == 2 && o3_nf) {
    // ----xx -> ??-
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 5", op1, op2, op3, op4, op5,
              op6, o3, false, false, true, true, true, true, false);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 5", op1, op2, op3, op4, op5,
              op6, o3, false, true, false, true, true, true, false);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 5", op1, op2, op3, op4, op5,
              op6, o3, false, true, true, false, true, true, false);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 5", op1, op2, op3, op4, op5,
              op6, o3, false, true, true, true, false, true, false);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 5", op1, op2, op3, op4, op5,
              op6, o3, false, true, true, true, true, false, false);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 5", op1, op2, op3, op4, op5,
              op6, o3, true, false, false, true, true, true, false);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 5", op1, op2, op3, op4, op5,
              op6, o3, true, false, true, false, true, true, false);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 5", op1, op2, op3, op4, op5,
              op6, o3, true, false, true, true, false, true, false);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 5", op1, op2, op3, op4, op5,
              op6, o3, true, false, true, true, true, false, false);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 5", op1, op2, op3, op4, op5,
              op6, o3, true, true, false, false, true, true, false);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 5", op1, op2, op3, op4, op5,
              op6, o3, true, true, false, true, false, true, false);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 5", op1, op2, op3, op4, op5,
              op6, o3, true, true, false, true, true, false, false);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 5", op1, op2, op3, op4, op5,
              op6, o3, true, true, true, false, false, true, false);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 5", op1, op2, op3, op4, op5,
              op6, o3, true, true, true, false, true, false, false);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 5", op1, op2, op3, op4, op5,
              op6, o3, true, true, true, true, false, false, false);
  } else if (nxs == 3 && xs == 3 && o3_nt) {
    // ---xxx -> ??x
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 6", op1, op2, op3, op4, op5,
              op6, o3, false, false, false, true, true, true, true);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 6", op1, op2, op3, op4, op5,
              op6, o3, false, false, true, false, true, true, true);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 6", op1, op2, op3, op4, op5,
              op6, o3, false, false, true, true, false, true, true);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 6", op1, op2, op3, op4, op5,
              op6, o3, false, false, true, true, true, false, true);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 6", op1, op2, op3, op4, op5,
              op6, o3, false, true, false, false, true, true, true);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 6", op1, op2, op3, op4, op5,
              op6, o3, false, true, false, true, false, true, true);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 6", op1, op2, op3, op4, op5,
              op6, o3, false, true, false, true, true, false, true);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 6", op1, op2, op3, op4, op5,
              op6, o3, false, true, true, false, false, true, true);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 6", op1, op2, op3, op4, op5,
              op6, o3, false, true, true, false, true, false, true);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 6", op1, op2, op3, op4, op5,
              op6, o3, false, true, true, true, false, false, true);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 6", op1, op2, op3, op4, op5,
              op6, o3, true, false, false, false, true, true, true);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 6", op1, op2, op3, op4, op5,
              op6, o3, true, false, false, true, false, true, true);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 6", op1, op2, op3, op4, op5,
              op6, o3, true, false, false, true, true, false, true);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 6", op1, op2, op3, op4, op5,
              op6, o3, true, false, true, false, false, true, true);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 6", op1, op2, op3, op4, op5,
              op6, o3, true, false, true, false, true, false, true);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 6", op1, op2, op3, op4, op5,
              op6, o3, true, false, true, true, false, false, true);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 6", op1, op2, op3, op4, op5,
              op6, o3, true, true, false, false, false, true, true);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 6", op1, op2, op3, op4, op5,
              op6, o3, true, true, false, false, true, false, true);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 6", op1, op2, op3, op4, op5,
              op6, o3, true, true, false, true, false, false, true);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 6", op1, op2, op3, op4, op5,
              op6, o3, true, true, true, false, false, false, true);
  } else if (nxs == 2 && xs == 4 && o3_nf) {
    // --xxxx -> ??-
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 7", op1, op2, op3, op4, op5,
              op6, o3, false, false, false, false, true, true, false);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 7", op1, op2, op3, op4, op5,
              op6, o3, false, false, false, true, false, true, false);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 7", op1, op2, op3, op4, op5,
              op6, o3, false, false, false, true, true, false, false);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 7", op1, op2, op3, op4, op5,
              op6, o3, false, false, true, false, false, true, false);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 7", op1, op2, op3, op4, op5,
              op6, o3, false, false, true, false, true, false, false);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 7", op1, op2, op3, op4, op5,
              op6, o3, false, false, true, true, false, false, false);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 7", op1, op2, op3, op4, op5,
              op6, o3, false, true, false, false, false, true, false);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 7", op1, op2, op3, op4, op5,
              op6, o3, false, true, false, false, true, false, false);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 7", op1, op2, op3, op4, op5,
              op6, o3, false, true, false, true, false, false, false);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 7", op1, op2, op3, op4, op5,
              op6, o3, false, true, true, false, false, false, false);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 7", op1, op2, op3, op4, op5,
              op6, o3, true, false, false, false, false, true, false);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 7", op1, op2, op3, op4, op5,
              op6, o3, true, false, false, false, true, false, false);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 7", op1, op2, op3, op4, op5,
              op6, o3, true, false, false, true, false, false, false);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 7", op1, op2, op3, op4, op5,
              op6, o3, true, false, true, false, false, false, false);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 7", op1, op2, op3, op4, op5,
              op6, o3, true, true, false, false, false, false, false);
  } else if (nxs == 1 && xs == 5 && o3_nt) {
    // -xxxxx -> ??x
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 8", op1, op2, op3, op4, op5,
              op6, o3, false, false, false, false, false, true, true);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 8", op1, op2, op3, op4, op5,
              op6, o3, false, false, false, false, true, false, true);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 8", op1, op2, op3, op4, op5,
              op6, o3, false, false, false, true, false, false, true);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 8", op1, op2, op3, op4, op5,
              op6, o3, false, false, true, false, false, false, true);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 8", op1, op2, op3, op4, op5,
              op6, o3, false, true, false, false, false, false, true);
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 8", op1, op2, op3, op4, op5,
              op6, o3, true, false, false, false, false, false, true);
  }

  if (nxs == 0 && xs == 6 && o3_nf) {
    // xxxxxx -> ??-
    add_impl6(s, out_refined, k, i, j, "COMP_6_3 9", op1, op2, op3, op4, op5,
              op6, o3, false, false, false, false, false, false, false);
  }
}

void comp_5_3(Minisat::Solver &s,
              Minisat::vec<Minisat::vec<Minisat::Lit>> &out_refined, int &k,
              int &i, int &j, int &op1, int &op2, int &op3, int &op4, int &op5,
              int &o1, int &o2, int &o3) {
  int xs = 0, nxs = 0;
  xs += s.value(op1) == l_True ? 1 : 0;
  xs += s.value(op2) == l_True ? 1 : 0;
  xs += s.value(op3) == l_True ? 1 : 0;
  xs += s.value(op4) == l_True ? 1 : 0;
  xs += s.value(op5) == l_True ? 1 : 0;
  nxs += s.value(op1) == l_False ? 1 : 0;
  nxs += s.value(op2) == l_False ? 1 : 0;
  nxs += s.value(op3) == l_False ? 1 : 0;
  nxs += s.value(op4) == l_False ? 1 : 0;
  nxs += s.value(op5) == l_False ? 1 : 0;

  bool o1_nf = int_value(s, o1) != 0;
  bool o2_nf = int_value(s, o2) != 0;
  bool o3_nf = int_value(s, o3) != 0;
  bool o3_nt = int_value(s, o3) != 1;

  if (nxs == 5 && xs == 0 && o1_nf) {
    // ----- -> ---
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 1", op1, op2, op3, op4, op5,
              o1, true, true, true, true, true, false);
  }
  if (nxs == 5 && xs == 0 && o2_nf) {
    // ----- -> ---
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 2", op1, op2, op3, op4, op5,
              o2, true, true, true, true, true, false);
  }
  if (nxs == 5 && xs == 0 && o3_nf) {
    // ----- -> ---
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 3", op1, op2, op3, op4, op5,
              o3, true, true, true, true, true, false);
  }

  if (nxs == 4 && xs == 1 && o3_nt) {
    // ----x -> ??x
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 4", op1, op2, op3, op4, op5,
              o3, true, true, true, true, false, true);
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 4", op1, op2, op3, op4, op5,
              o3, true, true, true, false, true, true);
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 4", op1, op2, op3, op4, op5,
              o3, true, true, false, true, true, true);
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 4", op1, op2, op3, op4, op5,
              o3, true, false, true, true, true, true);
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 4", op1, op2, op3, op4, op5,
              o3, false, true, true, true, true, true);
  } else if (nxs == 3 && xs == 2 && o3_nf) {
    // ---xx -> ??-
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 5", op1, op2, op3, op4, op5,
              o3, true, true, true, false, false, false); // 11100
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 5", op1, op2, op3, op4, op5,
              o3, true, true, false, true, false, false); // 11010
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 5", op1, op2, op3, op4, op5,
              o3, true, true, false, false, true, false); // 11001
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 5", op1, op2, op3, op4, op5,
              o3, true, false, true, true, false, false); // 10110
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 5", op1, op2, op3, op4, op5,
              o3, true, false, true, false, true, false); // 10101
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 5", op1, op2, op3, op4, op5,
              o3, true, false, false, true, true, false); // 10011
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 5", op1, op2, op3, op4, op5,
              o3, false, true, true, true, false, false); // 01110
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 5", op1, op2, op3, op4, op5,
              o3, false, true, true, false, true, false); // 01101
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 5", op1, op2, op3, op4, op5,
              o3, false, true, false, true, true, false); // 01011
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 5", op1, op2, op3, op4, op5,
              o3, false, false, true, true, true, false); // 00111
  } else if (nxs == 2 && xs == 3 && o3_nt) {
    // --xxx -> ??x
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 6", op1, op2, op3, op4, op5,
              o3, true, true, false, false, false, true); // 11000
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 6", op1, op2, op3, op4, op5,
              o3, true, false, true, false, false, true); // 10100
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 6", op1, op2, op3, op4, op5,
              o3, true, false, false, true, false, true); // 10010
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 6", op1, op2, op3, op4, op5,
              o3, true, false, false, false, true, true); // 10001
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 6", op1, op2, op3, op4, op5,
              o3, false, true, true, false, false, true); // 01100
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 6", op1, op2, op3, op4, op5,
              o3, false, true, false, true, false, true); // 01010
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 6", op1, op2, op3, op4, op5,
              o3, false, true, false, false, true, true); // 01001
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 6", op1, op2, op3, op4, op5,
              o3, false, false, true, true, false, true); // 00110
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 6", op1, op2, op3, op4, op5,
              o3, false, false, true, false, true, true); // 00101
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 6", op1, op2, op3, op4, op5,
              o3, false, false, false, true, true, true); // 00011
  } else if (nxs == 1 && xs == 4 && o3_nf) {
    // -xxxx -> ??-
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 7", op1, op2, op3, op4, op5,
              o3, true, false, false, false, false, false); // 10000
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 7", op1, op2, op3, op4, op5,
              o3, false, true, false, false, false, false); // 01000
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 7", op1, op2, op3, op4, op5,
              o3, false, false, true, false, false, false); // 00100
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 7", op1, op2, op3, op4, op5,
              o3, false, false, false, true, false, false); // 00010
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 7", op1, op2, op3, op4, op5,
              o3, false, false, false, false, true, false); // 00001
  }

  if (xs == 5 && nxs == 0 && o2_nf) {
    // xxxxx -> ?-x
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 8", op1, op2, op3, op4, op5,
              o2, false, false, false, false, false, false);
  }
  if (xs == 5 && nxs == 0 && o3_nt) {
    // xxxxx -> ?-x
    add_impl5(s, out_refined, k, i, j, "COMP_5_3 9", op1, op2, op3, op4, op5,
              o2, false, false, false, false, false, true);
  }
}

void comp_4_3(Minisat::Solver &s,
              Minisat::vec<Minisat::vec<Minisat::Lit>> &out_refined, int &k,
              int &i, int &j, int &op1, int &op2, int &op3, int &op4, int &o1,
              int &o2, int &o3) {
  int xs = 0, nxs = 0;
  xs += s.value(op1) == l_True ? 1 : 0;
  xs += s.value(op2) == l_True ? 1 : 0;
  xs += s.value(op3) == l_True ? 1 : 0;
  xs += s.value(op4) == l_True ? 1 : 0;
  nxs += s.value(op1) == l_False ? 1 : 0;
  nxs += s.value(op2) == l_False ? 1 : 0;
  nxs += s.value(op3) == l_False ? 1 : 0;
  nxs += s.value(op4) == l_False ? 1 : 0;

  bool o1_nf = int_value(s, o1) != 0;
  bool o2_nf = int_value(s, o2) != 0;
  bool o3_nf = int_value(s, o3) != 0;
  bool o3_nt = int_value(s, o3) != 1;

  if (nxs == 4 && xs == 0 && o1_nf) {
    // ---- -> ---
    add_impl4(s, out_refined, k, i, j, "COMP_5_3 1", op1, op2, op3, op4, o1,
              true, true, true, true, false);
  }
  if (nxs == 4 && xs == 0 && o2_nf) {
    // ---- -> ---
    add_impl4(s, out_refined, k, i, j, "COMP_5_3 2", op1, op2, op3, op4, o2,
              true, true, true, true, false);
  }
  if (nxs == 4 && xs == 0 && o3_nf) {
    // ---- -> ---
    add_impl4(s, out_refined, k, i, j, "COMP_5_3 3", op1, op2, op3, op4, o3,
              true, true, true, true, false);
  }

  if (nxs == 3 && xs == 1 && o3_nt) {
    // ---x -> ??x
    add_impl4(s, out_refined, k, i, j, "COMP_5_3 4", op1, op2, op3, op4, o3,
              false, true, true, true, true);
    add_impl4(s, out_refined, k, i, j, "COMP_5_3 4", op1, op2, op3, op4, o3,
              true, false, true, true, true);
    add_impl4(s, out_refined, k, i, j, "COMP_5_3 4", op1, op2, op3, op4, o3,
              true, true, false, true, true);
    add_impl4(s, out_refined, k, i, j, "COMP_5_3 4", op1, op2, op3, op4, o3,
              true, true, true, false, true);
  } else if (nxs == 2 && xs == 2 && o3_nf) {
    add_impl4(s, out_refined, k, i, j, "COMP_5_3 5", op1, op2, op3, op4, o3,
              false, false, true, true, false);
    add_impl4(s, out_refined, k, i, j, "COMP_5_3 5", op1, op2, op3, op4, o3,
              false, true, false, true, false);
    add_impl4(s, out_refined, k, i, j, "COMP_5_3 5", op1, op2, op3, op4, o3,
              false, true, true, false, false);
    add_impl4(s, out_refined, k, i, j, "COMP_5_3 5", op1, op2, op3, op4, o3,
              true, false, false, true, false);
    add_impl4(s, out_refined, k, i, j, "COMP_5_3 5", op1, op2, op3, op4, o3,
              true, false, true, false, false);
    add_impl4(s, out_refined, k, i, j, "COMP_5_3 5", op1, op2, op3, op4, o3,
              true, true, false, false, false);
  } else if (nxs == 1 && xs == 3 && o3_nt) {
    add_impl4(s, out_refined, k, i, j, "COMP_5_3 6", op1, op2, op3, op4, o3,
              false, false, false, true, true);
    add_impl4(s, out_refined, k, i, j, "COMP_5_3 6", op1, op2, op3, op4, o3,
              false, false, true, false, true);
    add_impl4(s, out_refined, k, i, j, "COMP_5_3 6", op1, op2, op3, op4, o3,
              false, true, false, false, true);
    add_impl4(s, out_refined, k, i, j, "COMP_5_3 6", op1, op2, op3, op4, o3,
              true, false, false, false, true);
  }

  if (nxs == 0 && xs == 4 && o3_nf) {
    add_impl4(s, out_refined, k, i, j, "COMP_5_3 7", op1, op2, op3, op4, o3,
              true, true, true, true, false);
  }
}

void comp_3_2(Minisat::Solver &s,
              Minisat::vec<Minisat::vec<Minisat::Lit>> &out_refined, int &k,
              int &i, int &j, int &op1, int &op2, int &op3, int &o1, int &o2) {
  int xs = 0, nxs = 0;
  xs += s.value(op1) == l_True ? 1 : 0;
  xs += s.value(op2) == l_True ? 1 : 0;
  xs += s.value(op3) == l_True ? 1 : 0;
  nxs += s.value(op1) == l_False ? 1 : 0;
  nxs += s.value(op2) == l_False ? 1 : 0;
  nxs += s.value(op3) == l_False ? 1 : 0;

  int o1_nf = int_value(s, o1) != 0;
  int o1_nt = int_value(s, o1) != 1;
  int o2_nf = int_value(s, o2) != 0;
  int o2_nt = int_value(s, o2) != 1;

  if (nxs == 3 && xs == 0 && o1_nf) {
    // --- -> --
    add_impl3(s, out_refined, k, i, j, "COMP_3_2 1", op1, op2, op3, o1, true,
              true, true, false);
  }
  if (nxs == 3 && xs == 0 && o2_nf) {
    // --- -> --
    add_impl3(s, out_refined, k, i, j, "COMP_3_2 2", op1, op2, op3, o2, true,
              true, true, false);
  }

  if (nxs == 2 && xs == 1 && o2_nt) {
    // --x -> ?x
    add_impl3(s, out_refined, k, i, j, "COMP_3_2 3", op1, op2, op3, o2, true,
              true, false, true);
    add_impl3(s, out_refined, k, i, j, "COMP_3_2 3", op1, op2, op3, o2, true,
              false, true, true);
    add_impl3(s, out_refined, k, i, j, "COMP_3_2 3", op1, op2, op3, o2, false,
              true, true, true);
  } else if (nxs == 1 && xs == 2 && o2_nf) {
    // -xx -> ?-
    add_impl3(s, out_refined, k, i, j, "COMP_3_2 4", op1, op2, op3, o2, true,
              false, false, false);
    add_impl3(s, out_refined, k, i, j, "COMP_3_2 4", op1, op2, op3, o2, false,
              true, false, false);
    add_impl3(s, out_refined, k, i, j, "COMP_3_2 4", op1, op2, op3, o2, false,
              false, true, false);
  }

  if (xs == 3 && nxs == 0 && o1_nt) {
    // xxx -> xx
    add_impl3(s, out_refined, k, i, j, "COMP_3_2 5", op1, op2, op3, o1, false,
              false, false, true);
  }
  if (xs == 3 && nxs == 0 && o2_nt) {
    // xxx -> xx
    add_impl3(s, out_refined, k, i, j, "COMP_3_2 6", op1, op2, op3, o2, false,
              false, false, true);
  }
}

void comp_2_2(Minisat::Solver &s,
              Minisat::vec<Minisat::vec<Minisat::Lit>> &out_refined, int &k,
              int &i, int &j, int &op1, int &op2, int &o1, int &o2) {
  int xs = 0, nxs = 0;
  xs += s.value(op1) == l_True ? 1 : 0;
  xs += s.value(op2) == l_True ? 1 : 0;
  nxs += s.value(op1) == l_False ? 1 : 0;
  nxs += s.value(op2) == l_False ? 1 : 0;

  int o1_nf = int_value(s, o1) != 0;
  int o2_nf = int_value(s, o2) != 0;
  int o2_nt = int_value(s, o2) != 1;

  if (nxs == 2 && xs == 0 && o1_nf) {
    // -- -> --
    add_impl2(s, out_refined, k, i, j, "COMP_2_2 1", op1, op2, o1, true, true,
              false);
  }
  if (nxs == 2 && xs == 0 && o2_nf) {
    // -- -> --
    add_impl2(s, out_refined, k, i, j, "COMP_2_2 2", op1, op2, o2, true, true,
              false);
  }

  if (nxs == 1 && xs == 1 && o2_nt) {
    // -x -> ?x
    add_impl2(s, out_refined, k, i, j, "COMP_2_2 3", op1, op2, o2, true, false,
              true);
    add_impl2(s, out_refined, k, i, j, "COMP_2_2 3", op1, op2, o2, false, true,
              true);
  } else if (nxs == 0 && xs == 2 && o2_nf) {
    // xx -> ?-
    add_impl2(s, out_refined, k, i, j, "COMP_2_2 4", op1, op2, o2, false, false,
              false);
  }
}

void xor3_impl(Minisat::Solver &s,
               Minisat::vec<Minisat::vec<Minisat::Lit>> &out_refined, int &k,
               int i, int j, int a, int b, int c, int o, const char *name) {
  if (s.value(a) != l_Undef && s.value(b) != l_Undef && s.value(c) != l_Undef) {
    int xs = 0, nxs = 0;
    if (s.value(a) == l_True)
      xs++;
    if (s.value(b) == l_True)
      xs++;
    if (s.value(c) == l_True)
      xs++;

    if (s.value(a) == l_False)
      nxs++;
    if (s.value(b) == l_False)
      nxs++;
    if (s.value(c) == l_False)
      nxs++;
    auto o_ = s.value(o);

    if (xs == 0 && nxs == 3 && o_ != l_False) {
      // --- -> -
      add_impl3(s, out_refined, k, i, j, name, a, b, c, o, true, true, true,
                false);
    } else if (xs == 1 && nxs == 2 && o_ != l_True) {
      // --x -> x
      add_impl3(s, out_refined, k, i, j, name, a, b, c, o, true, true, false,
                true);
      // -x- -> x
      add_impl3(s, out_refined, k, i, j, name, a, b, c, o, true, false, true,
                true);
      // x-- -> x
      add_impl3(s, out_refined, k, i, j, name, a, b, c, o, false, true, true,
                true);
    } else if (xs == 2 && nxs == 1 && o_ != l_False) {
      // -xx -> -
      add_impl3(s, out_refined, k, i, j, name, a, b, c, o, true, false, false,
                false);
      // x-x -> -
      add_impl3(s, out_refined, k, i, j, name, a, b, c, o, false, true, false,
                false);
      // xx- -> -
      add_impl3(s, out_refined, k, i, j, name, a, b, c, o, false, false, true,
                false);
    } else if (xs == 3 && nxs == 0 && o_ != l_True) {
      // xxx -> x
      add_impl3(s, out_refined, k, i, j, name, a, b, c, o, false, false, false,
                true);
    }
  }
}

void xor2_impl(Minisat::Solver &s,
               Minisat::vec<Minisat::vec<Minisat::Lit>> &out_refined, int &k,
               int i, int j, int op1, int op2, int o, const char *name) {
  if (s.value(op1) != l_Undef && s.value(op2) != l_Undef) {
    int xs = 0, nxs = 0;
    xs += s.value(op1) == l_True ? 1 : 0;
    xs += s.value(op2) == l_True ? 1 : 0;
    nxs += s.value(op1) == l_False ? 1 : 0;
    nxs += s.value(op2) == l_False ? 1 : 0;
    auto o_nf = s.value(o) != l_False;
    auto o_nt = s.value(o) != l_True;

    if (xs == 0 && nxs == 2 && o_nf) {
      // -- -> -
      add_impl2(s, out_refined, k, i, j, name, op1, op2, o, true, true, false);
    } else if (nxs == 1 && xs == 1 && o_nt) {
      // -x -> x
      add_impl2(s, out_refined, k, i, j, name, op1, op2, o, true, false, true);
      add_impl2(s, out_refined, k, i, j, name, op1, op2, o, false, true, true);
    } else if (nxs == 0 && xs == 2 && o_nf) {
      // xx -> -
      add_impl2(s, out_refined, k, i, j, name, op1, op2, o, false, false,
                false);
    }
  }
}

void add_clauses(Minisat::Solver &s, vec<vec<Lit>> &out_refined) {
  int k = 0;
  for (int i = 0; i < s.steps; i++) {
    int dw_0_base = s.var_map["DW_" + std::to_string(i) + "_g"];
    int da_4_base = s.var_map["DA_" + std::to_string(i + 4) + "_g"];
    int da_0_base = s.var_map["DA_" + std::to_string(i) + "_g"];
    int de_4_base = s.var_map["DE_" + std::to_string(i + 4) + "_g"];
    int de_0_base = s.var_map["DE_" + std::to_string(i) + "_g"];
    int df1_base = s.var_map["Df1_" + std::to_string(i) + "_g"];
    int df2_base = s.var_map["Df2_" + std::to_string(i) + "_g"];
    int dsigma0_base = s.var_map["Dsigma0_" + std::to_string(i) + "_g"];
    int dsigma1_base = s.var_map["Dsigma1_" + std::to_string(i) + "_g"];
    int ds0_base = s.var_map["Ds0_" + std::to_string(i) + "_g"];
    int ds1_base = s.var_map["Ds1_" + std::to_string(i) + "_g"];
    int dt_base = s.var_map["DT_" + std::to_string(i) + "_g"];
    int dk_base = s.var_map["DK_" + std::to_string(i) + "_g"];
    int dr1_carry_base = s.var_map["Dr1_carry_" + std::to_string(i) + "_g"];
    int dr2_carry_base = s.var_map["Dr2_carry_" + std::to_string(i) + "_g"];
    int dr2_carry2_base = s.var_map["Dr2_Carry_" + std::to_string(i) + "_g"];
    int dr0_carry_base = s.var_map["Dr0_carry_" + std::to_string(i) + "_g"];
    int dr0_carry2_base = s.var_map["Dr0_Carry_" + std::to_string(i) + "_g"];
    int dw_carry_base = s.var_map["Dw_carry_" + std::to_string(i) + "_g"];
    int dw_carry2_base = s.var_map["Dw_Carry_" + std::to_string(i) + "_g"];
    for (int j = 0; j < 32; j++) {
      int dw_0 = dw_0_base + j;                 // DW[i]
      int da_4 = da_4_base + j;                 // DA[i+4]
      int da_0 = da_0_base + j;                 // DA[i]
      int de_4 = de_4_base + j;                 // DE[i+4]
      int de_0 = de_0_base + j;                 // DE[i]
      int df1 = df1_base + j;                   // DIF
      int df2 = df2_base + j;                   // DMAJ
      int dsigma0 = dsigma0_base + j;           // DSigma0
      int dsigma1 = dsigma1_base + j;           // DSigma1
      int ds0 = ds0_base + j;                   // DS0
      int ds1 = ds1_base + j;                   // DS1
      int dt = dt_base + j;                     // DT
      int dk = dk_base + j;                     // DT
      int dr1_carry = dr1_carry_base + j;       // Dr1_carry
      int dr2_carry = dr2_carry_base + j;       // Dr2_carry
      int dr2_carry2 = dr2_carry2_base + j - 1; // Dr2_carry2
      int dr0_carry = dr0_carry_base + j;       // Dr0_carry
      int dr0_carry2 = dr0_carry2_base + j;     // Dr0_carry2
      int dw_carry = dw_carry_base + j;         // Dw_carry
      int dw_carry2 = dw_carry2_base + j;       // Dw_carry2

      // // IF
      // if (s.value(de_3) == l_False && s.value(de_2) == l_True &&
      // s.value(de_1) == l_True && s.value(df1) == l_Undef) {
      //     // -xx -> x
      //     ADD_IMPL3("IF", de_3, de_2, de_1, df1, true, false, false, true);
      // } else if (s.value(de_3) == l_False && s.value(de_2) == l_False &&
      // s.value(de_1) == l_False && s.value(df1) == l_Undef) {
      //     // --- -> -
      //     ADD_IMPL3("IF", de_3, de_2, de_1, df1, true, true, true, false);
      // }

      // // MAJ
      // if (s.value(da_3) == l_True && s.value(da_2) == l_True && s.value(da_1)
      // == l_True && s.value(df2) == l_Undef) {
      //     // xxx -> x
      //     ADD_IMPL3("MAJ", da_3, da_2, da_1, df2, false, false, false, true);
      // } else if (s.value(da_3) == l_False && s.value(da_2) == l_False &&
      // s.value(da_1) == l_False && s.value(df2) == l_Undef) {
      //     // --- -> -
      //     ADD_IMPL3("MAJ", da_3, da_2, da_1, df2, true, true, true, false);
      // }

      // Sigma0
      int da_3_base = s.var_map["DA_" + std::to_string(i + 3) + "_g"];
      int da_3_r1 = da_3_base + (j + 2) % 32;
      int da_3_r2 = da_3_base + (j + 13) % 32;
      int da_3_r3 = da_3_base + (j + 22) % 32;
      xor3_impl(s, out_refined, k, i, j, da_3_r1, da_3_r2, da_3_r3, dsigma0,
                "sigma0");

      // Sigma1
      int de_3_base = s.var_map["DE_" + std::to_string(i + 3) + "_g"];
      int de_3_r1 = de_3_base + (j + 6) % 32;
      int de_3_r2 = de_3_base + (j + 11) % 32;
      int de_3_r3 = de_3_base + (j + 25) % 32;
      xor3_impl(s, out_refined, k, i, j, de_3_r1, de_3_r2, de_3_r3, dsigma1,
                "sigma1");

      // S0
      if (i >= 16) {
        int dw_m15_base = s.var_map["DW_" + std::to_string(i - 15) + "_g"];
        int dw_s0_r1 = dw_m15_base + (j + 7) % 32;
        int dw_s0_r2 = dw_m15_base + (j + 18) % 32;
        int dw_s0_r3 = dw_m15_base + j + 3;
        if (j < 29) {
          xor3_impl(s, out_refined, k, i, j, dw_s0_r1, dw_s0_r2, dw_s0_r3, ds0,
                    "s0 1");
        } else {
          xor2_impl(s, out_refined, k, i, j, dw_s0_r1, dw_s0_r2, ds0, "s0 2");
        }
      }

      // S1
      if (i >= 16) {
        int dw_m2_base = s.var_map["DW_" + std::to_string(i - 2) + "_g"];
        int dw_s1_r1 = dw_m2_base + (j + 17) % 32;
        int dw_s1_r2 = dw_m2_base + (j + 19) % 32;
        int dw_s1_r3 = dw_m2_base + j + 10;
        if (j < 22) {
          xor3_impl(s, out_refined, k, i, j, dw_s1_r1, dw_s1_r2, dw_s1_r3, ds1,
                    "s1 2");
        } else {
          xor2_impl(s, out_refined, k, i, j, dw_s1_r1, dw_s1_r2, ds1, "s1 1");
        }
      }

      // Compression: 3 to 2
      {
        int op1 = da_0;
        int op2 = dt;
        int op3 = dr1_carry - 1;
        int o1 = dr1_carry;
        int o2 = de_4;

        if (int_value(s, op1) != -1 && int_value(s, op2) != -1) {
          if (j > 0 && int_value(s, op3) != -1) {
            comp_3_2(s, out_refined, k, i, j, op1, op2, op3, o1, o2);
          } else if (j == 0) {
            comp_2_2(s, out_refined, k, i, j, op1, op2, o1, o2);
          }
        }
      }

      // Compression: 5 to 3
      // g.cnf.diff_add(DA[i + 4], DT[i], Dsigma0[i], Dr2carry[i], Dr2Carry[i],
      // Df2[i]);
      {
        int op1 = dt;
        int op2 = dsigma0;
        int op3 = df2;
        int op4 = dr2_carry - 1;  // t[j - 1]
        int op5 = dr2_carry2 - 2; // T[j - 2]
        int o1 = dr2_carry2;      // T[j]
        int o2 = dr2_carry;       // t[j]
        int o3 = da_4;            // DA[i+4]

        if (int_value(s, op1) != -1 && int_value(s, op2) != -1 &&
            int_value(s, op3) != -1) {
          if (j > 2 && int_value(s, op4) != -1 && int_value(s, op5) != -1) {
            comp_5_3(s, out_refined, k, i, j, op1, op2, op3, op4, op5, o1, o2,
                     o3);
          } else if ((j == 2 || j == 1) && int_value(s, op4) != -1) {
            comp_4_3(s, out_refined, k, i, j, op1, op2, op3, op4, o1, o2, o3);
          } else if (j == 0) {
            comp_3_2(s, out_refined, k, i, j, op1, op2, op3, o2, o3);
          }
        }
      }

      // Compression: 6 to 3
      // g.cnf.diff_add(DW[i], DW[i - 16], Ds0[i], Dwcarry[i], DwCarry[i], DW[i
      // - 7], Ds1[i]);
      if (i >= 16) {
        int op1 = s.var_map["DW_" + std::to_string(i - 16) + "_g"] + j;
        int op2 = ds0;
        int op3 = s.var_map["DW_" + std::to_string(i - 7) + "_g"] + j;
        int op4 = ds1;
        int op5 = dw_carry - 1;  // t[j - 1]
        int op6 = dw_carry2 - 2; // T[j - 2]
        int o1 = dw_carry2;      // T[j]
        int o2 = dw_carry;       // t[j]
        int o3 = dw_0;           // DW[i]

        // printf("%d %d\n", op2, op4);

        if (s.value(op1) != l_Undef && s.value(op2) != l_Undef &&
            s.value(op3) != l_Undef && s.value(op4) != l_Undef) {
          if (j > 1 && int_value(s, op5) != -1 && int_value(s, op6) != -1) {
            comp_6_3(s, out_refined, k, i, j, op1, op2, op3, op4, op5, op6, o1,
                     o2, o3);
          } else if (j == 1 && int_value(s, op5) != -1) {
            comp_5_3(s, out_refined, k, i, j, op1, op2, op3, op4, op5, o1, o2,
                     o3);
          } else if (j == 0) {
            comp_4_3(s, out_refined, k, i, j, op1, op2, op3, op4, o1, o2, o3);
          }
        }
      }

      // Compression: 7 to 3
      // g.cnf.diff_add(DT[i], DE[i], Dsigma1[i], Dr0carry[i], Dr0Carry[i],
      // Df1[i], DK[i], DW[i]);
      {
        int op1 = de_0;
        int op2 = dsigma1;
        int op3 = df1;
        int op4 = dk;
        int op5 = dw_0;
        int op6 = dr0_carry - 1;  // t[j - 1]
        int op7 = dr0_carry2 - 2; // T[j - 2]
        int o1 = dr0_carry2;      // T[j]
        int o2 = dr0_carry;       // t[j]
        int o3 = dt;              // DT[i]
        if (int_value(s, op1) != -1 && int_value(s, op2) != -1 &&
            int_value(s, op3) != -1 && int_value(s, op4) != -1 &&
            int_value(s, op5) != -1) {
          if (j > 1 && int_value(s, op6) != -1 && int_value(s, op7) != -1) {
            comp_7_3(s, out_refined, k, i, j, op1, op2, op3, op4, op5, op6, op7,
                     o1, o2, o3);
          } else if (j == 1 && int_value(s, op6)) {
            comp_6_3(s, out_refined, k, i, j, op1, op2, op3, op4, op5, op6, o1,
                     o2, o3);
          } else if (j == 0) {
            comp_5_3(s, out_refined, k, i, j, op1, op2, op3, op4, op5, o1, o2,
                     o3);
          }
        }
      }
    }
  }

  if (k > 0) {
    printf("k = %d\n", k);
  }
}