#include "Crypto.h"
#include <algorithm>
#include <map>
#include <set>
#include <vector>
#define DEBUG true
#define DIFF_BITS 1
#define IO_CONSTRAINT_ADD2_ID 0
#define IO_CONSTRAINT_IF_ID 1
#define IO_CONSTRAINT_MAJ_ID 2
#define IO_CONSTRAINT_XOR3_ID 3
#define IO_CONSTRAINT_ADD3_ID 4
#define IO_CONSTRAINT_ADD4_ID 5
#define IO_CONSTRAINT_ADD5_ID 6
#define IO_CONSTRAINT_ADD6_ID 7
#define IO_CONSTRAINT_ADD7_ID 8
#define OI_CONSTRAINT_IF_ID 9
#define OI_CONSTRAINT_MAJ_ID 10
#define OI_CONSTRAINT_XOR3_ID 11
#define OI_CONSTRAINT_ADD3_ID 12
#define OI_CONSTRAINT_ADD4_ID 13
#define OI_CONSTRAINT_ADD5_ID 14
#define OI_CONSTRAINT_ADD6_ID 15
#define OI_CONSTRAINT_ADD7_ID 16
#define TWO_BIT_CONSTRAINT_IF_ID 17
#define TWO_BIT_CONSTRAINT_MAJ_ID 18
#define TWO_BIT_CONSTRAINT_XOR3_ID 19
#define TWO_BIT_CONSTRAINT_ADD3_ID 20
#define TWO_BIT_CONSTRAINT_ADD4_ID 21
#define TWO_BIT_CONSTRAINT_ADD5_ID 22
#define TWO_BIT_CONSTRAINT_ADD6_ID 23
#define TWO_BIT_CONSTRAINT_ADD7_ID 24

// TODO: Add support for 4-bit diff.

void loadRule(Minisat::Solver& solver, FILE*& db, int& id)
{
    int key_size = 0, val_size = 0;
    // Note: Put one extra char for the ID
    if (id >= IO_CONSTRAINT_IF_ID && id <= IO_CONSTRAINT_XOR3_ID) {
        key_size = 3;
        val_size = 1;
    } else if (id == IO_CONSTRAINT_ADD3_ID) {
        key_size = 3;
        val_size = 2;
    } else if (id >= TWO_BIT_CONSTRAINT_IF_ID && id <= TWO_BIT_CONSTRAINT_XOR3_ID) {
        key_size = 4;
        val_size = 3;
    } else if (id == OI_CONSTRAINT_ADD7_ID) {
        key_size = 3;
        val_size = 7;
    } else if (id == OI_CONSTRAINT_ADD6_ID) {
        key_size = 3;
        val_size = 6;
    } else if (id == OI_CONSTRAINT_ADD5_ID) {
        key_size = 3;
        val_size = 5;
    } else if (id == OI_CONSTRAINT_ADD4_ID) {
        key_size = 3;
        val_size = 4;
    } else if (id == OI_CONSTRAINT_ADD3_ID) {
        key_size = 2;
        val_size = 3;
    } else if (id >= TWO_BIT_CONSTRAINT_IF_ID && id <= TWO_BIT_CONSTRAINT_ADD3_ID) {
        key_size = 4;
        val_size = 3;
    }

    int size = key_size + val_size;
    char buffer[size];
    int n = fread(buffer, size, 1, db);
    if (n == 0)
        return;

    char key[key_size + 1], value[val_size];
    key[0] = id;
    for (int i = 0; i < key_size; i++) {
        key[i + 1] = buffer[i];
    }
    key[key_size + 1] = 0;
    int j = 0;
    for (int i = key_size; i < size; i++) {
        value[j] = buffer[i];
        j++;
    }
    value[val_size] = 0;

    solver.rules.insert({ key, value });

    // DEBUG
    // printf("Rule: %d %s: %s\n", id, key, value);
    // fflush(stdout);
}

void loadRules(Minisat::Solver& solver, const char* filename)
{
    FILE* db = fopen(filename, "r");
    char buffer[1];
    int count = 0;
    while (1) {
        int n = fread(buffer, 1, 1, db);
        if (n == 0)
            break;

        int id = buffer[0];
        loadRule(solver, db, id);
        count++;
    }

    printf("Loaded %d rules\n", count);

    // DEBUG
    // char key[] = {19, 'n', '1', '1', 'n'};
    // printf("Found: %s\n", solver.rules.find(key)->second.c_str());
    // exit(0);
}

int int_value(Minisat::Solver& s, int var)
{
    auto value = s.value(var);
    return value == l_True ? 1 : value == l_False ? 0
                                                  : -1;
}

char to_gc(int x, int x_prime)
{
    if (x == 0 && x_prime == 0)
        return '0';
    else if (x == 1 && x_prime == 0)
        return 'u';
    else if (x == 0 && x_prime == 1)
        return 'n';
    else if (x == 1 && x_prime == 1)
        return '1';
    else
        return NULL;
}

char to_gc(Minisat::Solver& s, int& id)
{
#if DIFF_BITS == 4
    int d[4] = { int_value(s, id), int_value(s, id + 1), int_value(s, id + 2), int_value(s, id + 3) };
    if (d[0] == 1 && d[1] == 1 && d[2] == 1 && d[3] == 1) {
        return '?';
    } else if (d[0] == 1 && d[1] == 0 && d[2] == 0 && d[3] == 1) {
        return '-';
    } else if (d[0] == 0 && d[1] == 1 && d[2] == 1 && d[3] == 0) {
        return 'x';
    } else if (d[0] == 1 && d[1] == 0 && d[2] == 0 && d[3] == 0) {
        return '0';
    } else if (d[0] == 0 && d[1] == 1 && d[2] == 0 && d[3] == 0) {
        return 'u';
    } else if (d[0] == 0 && d[1] == 0 && d[2] == 1 && d[3] == 0) {
        return 'n';
    } else if (d[0] == 0 && d[1] == 0 && d[2] == 0 && d[3] == 1) {
        return '1';
    } else if (d[0] == 1 && d[1] == 1 && d[2] == 0 && d[3] == 0) {
        return '3';
    } else if (d[0] == 1 && d[1] == 0 && d[2] == 1 && d[3] == 0) {
        return '5';
    } else if (d[0] == 1 && d[1] == 1 && d[2] == 1 && d[3] == 0) {
        return '7';
    } else if (d[0] == 0 && d[1] == 1 && d[2] == 0 && d[3] == 1) {
        return 'A';
    } else if (d[0] == 1 && d[1] == 1 && d[2] == 0 && d[3] == 1) {
        return 'B';
    } else if (d[0] == 0 && d[1] == 0 && d[2] == 1 && d[3] == 1) {
        return 'C';
    } else if (d[0] == 1 && d[1] == 0 && d[2] == 1 && d[3] == 1) {
        return 'D';
    } else if (d[0] == 0 && d[1] == 1 && d[2] == 1 && d[3] == 1) {
        return 'E';
    } else {
        // printf("%d: %d %d %d %d\n", id + 1, d[0], d[1], d[2], d[3]);
        return NULL;
    }
}

void from_gc(char& gc, uint8_t* vals) {
    if (gc == '-') {
        vals[0] = 1;
        vals[1] = 0;
        vals[2] = 0;
        vals[3] = 1;
    } else if (gc == 'x') {
        vals[0] = 0;
        vals[1] = 1;
        vals[2] = 1;
        vals[3] = 0;
    } else if (gc == '?') {
        vals[0] = 1;
        vals[1] = 1;
        vals[2] = 1;
        vals[3] = 1;
    } else if (gc == '0') {
        vals[0] = 1;
        vals[1] = 0;
        vals[2] = 0;
        vals[3] = 0;
    } else if (gc == 'u') {
        vals[0] = 0;
        vals[1] = 1;
        vals[2] = 0;
        vals[3] = 0;
    } else if (gc == 'n') {
        vals[0] = 0;
        vals[1] = 0;
        vals[2] = 1;
        vals[3] = 0;
    } else if (gc == '1') {
        vals[0] = 0;
        vals[1] = 0;
        vals[2] = 0;
        vals[3] = 1;
    } else if (gc == '3') {
        vals[0] = 1;
        vals[1] = 1;
        vals[2] = 0;
        vals[3] = 0;
    } else if (gc == '5') {
        vals[0] = 1;
        vals[1] = 0;
        vals[2] = 1;
        vals[3] = 0;
    } else if (gc == '7') {
        vals[0] = 1;
        vals[1] = 1;
        vals[2] = 1;
        vals[3] = 0;
    } else if (gc == 'A') {
        vals[0] = 0;
        vals[1] = 1;
        vals[2] = 0;
        vals[3] = 1;
    } else if (gc == 'B') {
        vals[0] = 1;
        vals[1] = 1;
        vals[2] = 0;
        vals[3] = 1;
    } else if (gc == 'C') {
        vals[0] = 0;
        vals[1] = 0;
        vals[2] = 1;
        vals[3] = 1;
    } else if (gc == 'D') {
        vals[0] = 1;
        vals[1] = 0;
        vals[2] = 1;
        vals[3] = 1;
    } else if (gc == 'E') {
        vals[0] = 0;
        vals[1] = 1;
        vals[2] = 1;
        vals[3] = 1;
    }
}

void print_clause(vec<Lit>& clause)
{
    for (int i = 0; i < clause.size(); i++) {
        printf("%s%d ", sign(clause[i]) ? "-" : "", var(clause[i]) + 1);
    }
    printf("\n");
}

void impose_rule_3_1(Minisat::Solver& s, vec<vec<Lit>>& out_refined, int& k, std::string& o, int& x, int& y, int& z, int& w)
{
    int vars[16] = { x,
        x + 1,
        x + 2,
        x + 3,
        y,
        y + 1,
        y + 2,
        y + 3,
        z,
        z + 1,
        z + 2,
        z + 3,
        w,
        w + 1,
        w + 2,
        w + 3 };
    int vals[16] = { int_value(s, x),
        int_value(s, x + 1),
        int_value(s, x + 2),
        int_value(s, x + 3),
        int_value(s, y),
        int_value(s, y + 1),
        int_value(s, y + 2),
        int_value(s, y + 3),
        int_value(s, z),
        int_value(s, z + 1),
        int_value(s, z + 2),
        int_value(s, z + 3),
        int_value(s, w),
        int_value(s, w + 1),
        int_value(s, w + 2),
        int_value(s, w + 3) };

    uint8_t o1v[4];
    from_gc(o[0], o1v);

    for (int j = 0; j < 4; j++) {
        if (vals[12 + j] == o1v[j])
            continue;

        out_refined.push();
        out_refined[k].push(mkLit(vars[j + 12], o1v[j] == 0));
        for (int i = 0; i < 12; i++)
            out_refined[k].push(mkLit(vars[i], vals[i] == 1));
        k++;

        break;
    }
}

void impose_rule_3_1_w(Minisat::Solver& s, vec<vec<Lit>>& out_refined, int& k, int id, int& x, int& y, int& z, int& w)
{
    if (k != 0) {
        return;
    }

    char x_gc = to_gc(s, x), y_gc = to_gc(s, y), z_gc = to_gc(s, z), w_gc = to_gc(s, w);
    if (x_gc == NULL || y_gc == NULL || z_gc == NULL) {
        return;
    }

    std::string r_key = std::to_string(id);
    r_key.push_back(x_gc);
    r_key.push_back(y_gc);
    r_key.push_back(z_gc);
    auto r_value_it = s.rules.find(r_key);
    if (r_value_it == s.rules.end()) // Rule not found
        return;

    auto r_value = r_value_it->second;
    if (w_gc == r_value[0]) { // Output difference is already correct
        // printf("%d is correct (%c)\n", w + 1, w_gc);
        return;
    }

    printf("RKey: %s; RValue: %s; DValue: %c(%d); DId: %d\n", r_key.c_str(), r_value.c_str(), w_gc, w_gc, w + 1);
    impose_rule_3_1(s, out_refined, k, r_value, x, y, z, w);
}

void add_clauses(Minisat::Solver& s, vec<vec<Lit>>& out_refined)
{
    int k = 0;
    for (int i = 0; i < s.steps; i++) {
        // int dw_0_base = s.var_map["DW_" + std::to_string(i) + "_g"];
        // int da_4_base = s.var_map["DA_" + std::to_string(i + 4) + "_g"];
        int da_3_base = s.var_map["DA_" + std::to_string(i + 3) + "_g"];
        int da_2_base = s.var_map["DA_" + std::to_string(i + 2) + "_g"];
        int da_1_base = s.var_map["DA_" + std::to_string(i + 1) + "_g"];
        // int da_0_base = s.var_map["DA_" + std::to_string(i) + "_g"];
        // int de_4_base = s.var_map["DE_" + std::to_string(i + 4) + "_g"];
        int de_3_base = s.var_map["DE_" + std::to_string(i + 3) + "_g"];
        int de_2_base = s.var_map["DE_" + std::to_string(i + 2) + "_g"];
        int de_1_base = s.var_map["DE_" + std::to_string(i + 1) + "_g"];
        // int de_0_base = s.var_map["DE_" + std::to_string(i) + "_g"];
        int df1_base = s.var_map["Df1_" + std::to_string(i) + "_g"];
        int df2_base = s.var_map["Df2_" + std::to_string(i) + "_g"];
        // int dsigma0_base = s.var_map["Dsigma0_" + std::to_string(i) + "_g"];
        // int dsigma1_base = s.var_map["Dsigma1_" + std::to_string(i) + "_g"];
        // int ds0_base = s.var_map["Ds0_" + std::to_string(i) + "_g"];
        // int ds1_base = s.var_map["Ds1_" + std::to_string(i) + "_g"];
        // int dt_base = s.var_map["DT_" + std::to_string(i) + "_g"];
        // int dk_base = s.var_map["DK_" + std::to_string(i) + "_g"];
        // int dr1_carry_base = s.var_map["Dr1_carry_" + std::to_string(i) + "_g"];
        // int dr2_carry_base = s.var_map["Dr2_carry_" + std::to_string(i) + "_g"];
        // int dr2_carry2_base = s.var_map["Dr2_Carry_" + std::to_string(i) + "_g"];
        // int dr0_carry_base = s.var_map["Dr0_carry_" + std::to_string(i) + "_g"];
        // int dr0_carry2_base = s.var_map["Dr0_Carry_" + std::to_string(i) + "_g"];
        // int dw_carry_base = s.var_map["Dw_carry_" + std::to_string(i) + "_g"];
        // int dw_carry2_base = s.var_map["Dw_Carry_" + std::to_string(i) + "_g"];

        for (int j = 0; j < 32; j++) {
            // int dw_0 = dw_0_base + j;                 // DW[i]
            // int da_4 = da_4_base + j; // DA[i+4]
            int da_3 = da_3_base + j * 4; // DA[i+3]
            int da_2 = da_2_base + j * 4; // DA[i+2]
            int da_1 = da_1_base + j * 4; // DA[i+1]
            // int de_4 = de_4_base + j * 4; // DE[i+4]
            int de_3 = de_3_base + j * 4; // DE[i+3]
            int de_2 = de_2_base + j * 4; // DE[i+2]
            int de_1 = de_1_base + j * 4; // DE[i+1]
            // int de_0 = de_0_base + j * 4; // DE[i]
            int df1 = df1_base + j * 4; // Df1 <- IF
            int df2 = df2_base + j * 4; // Df2 <- MAJ
            // int dsigma0 = dsigma0_base + j;           // DSigma0
            // int dsigma1 = dsigma1_base + j;           // DSigma1
            // int ds0 = ds0_base + j;                   // DS0
            // int ds1 = ds1_base + j;                   // DS1
            // int dt = dt_base + j;                     // DT
            // int dk = dk_base + j;                     // DT
            // int dr1_carry = dr1_carry_base + j;       // Dr1_carry
            // int dr2_carry = dr2_carry_base + j;       // Dr2_carry
            // int dr2_carry2 = dr2_carry2_base + j - 1; // Dr2_carry2
            // int dr0_carry = dr0_carry_base + j;       // Dr0_carry
            // int dr0_carry2 = dr0_carry2_base + j;     // Dr0_carry2
            // int dw_carry = dw_carry_base + j;         // Dw_carry
            // int dw_carry2 = dw_carry2_base + j;       // Dw_carry2

            // IF
            impose_rule_3_1_w(s, out_refined, k, 1, de_3, de_2, de_1, df1);
            // MAJ
            impose_rule_3_1_w(s, out_refined, k, 2, da_3, da_2, da_1, df2);
            // Sigma0
            // impose_rule_3_1_w(s, out_refined, k, 3, da_3, da_2, da_1, df2);
        }
    }
}