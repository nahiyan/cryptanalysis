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
    }
#else
    int d = int_value(s, id);
    if (d == 0) {
        return '-';
    } else if (d == 1) {
        return 'x';
    }
#endif
    return NULL;
}

void from_gc(char& gc, uint8_t* vals)
{
#if DIFF_BITS == 4
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
#else
    if (gc == '-') {
        vals[0] = 0;
    } else if (gc == 'x') {
        vals[0] = 1;
    }
#endif
}

void print_clause(vec<Lit>& clause)
{
    printf("Clause: ");
    for (int i = 0; i < clause.size(); i++) {
        printf("%s%d", sign(clause[i]) ? "-" : "", var(clause[i]) + 1);
        if (i != clause.size() - 1)
            printf(" ");
    }
    printf("\n");
}

bool check_consistency(std::vector<std::pair<int32_t, int32_t>>& equations)
{
    std::map<uint32_t, std::set<int32_t>*> rels;

    for (auto equation : equations) {
        auto var1 = equation.first;
        auto var2 = equation.second;
        auto var1_abs = abs(var1);
        auto var2_abs = abs(var2);
        auto var1_exists = rels.find(var1_abs) == rels.end() ? false : true;
        auto var2_exists = rels.find(var2_abs) == rels.end() ? false : true;

        if (var1_exists && var2_exists) {
            auto var1_inv_exists = rels[var1_abs]->find(-var1) == rels[var1_abs]->end() ? false : true;
            auto var2_inv_exists = rels[var2_abs]->find(-var2) == rels[var2_abs]->end() ? false : true;

            // Ignore if both inverses are found (would be a redudant operation)
            if (var1_inv_exists && var2_inv_exists)
                continue;

            // Try to prevent conflict by inverting one set
            bool invert = false;
            if (var2_inv_exists || var1_inv_exists)
                invert = true;

            // Union the sets
            for (auto item : *rels[var2_abs])
                rels[var1_abs]->insert((invert ? -1 : 1) * item);

            auto& updated_set = rels[var1_abs];
            // If both a var and its inverse is present in the newly updated set, we detected a contradiction
            {
                auto var1_inv_exists = updated_set->find(-var1_abs) == updated_set->end() ? false : true;
                auto var2_inv_exists = updated_set->find(-var2_abs) == updated_set->end() ? false : true;
                auto var1_exists = updated_set->find(var1_abs) == updated_set->end() ? false : true;
                auto var2_exists = updated_set->find(var2_abs) == updated_set->end() ? false : true;

                if ((var1_inv_exists && var1_exists) || (var2_inv_exists && var2_exists)) {
#if DEBUG
                    for (auto equation : equations) {
                        printf("Equation: %d %s %d\n", abs(equation.first) + 1, (equation.first > 0 && equation.second > 0) ? "=" : "/=", abs(equation.second) + 1);
                    }
                    printf("Contradiction detected (%d equations): %d %d\n", equations.size(), abs(var1) + 1, abs(var2) + 1);
#endif
                    return false;
                }
            }

            // Update existing references
            for (auto& item : *updated_set) {
                auto& set = rels[abs(item)];
                if (set == updated_set)
                    continue;

                // Delete last reference
                int counter = 0;
                for (auto& rel : rels) {
                    if (rel.second == rels[abs(item)])
                        counter++;
                }
                if (counter == 1)
                    delete rels[abs(item)];

                rels[abs(item)] = updated_set;
            }
        } else if (var1_exists || var2_exists) {
            // Find an existing set related to any of the variables
            auto& existing_set = var1_exists ? rels[var1_abs] : rels[var2_abs];
            auto var1_inv_in_existing_set = existing_set->find(-var1) == existing_set->end() ? false : true;
            auto var2_inv_in_existing_set = existing_set->find(-var2) == existing_set->end() ? false : true;

            // Invert the lone variable to try to prevent a conflict
            // if (var1_inv_in_existing_set)
            //     var2 *= -1;
            // else if (var2_inv_in_existing_set)
            //     var1 *= -1;

            // Add the var to an existing set
            if (var1_exists)
                rels[var1_abs]->insert(var2);
            else
                rels[var2_abs]->insert(var1);

            // Update existing references
            for (auto& item : *existing_set) {
                auto& set = rels[abs(item)];
                if (set == existing_set)
                    continue;

                // Delete last reference
                int counter = 0;
                for (auto& rel : rels) {
                    if (rel.second == rels[abs(item)])
                        counter++;
                }
                if (counter == 1)
                    delete rels[abs(item)];

                rels[abs(item)] = existing_set;
            }
        } else {
            // Adding novel variables
            auto new_set = new std::set<int32_t> { var1, var2 };
            rels[var1_abs] = new_set;
            rels[var2_abs] = new_set;
        }
    }

    // #if DEBUG
    //     for (auto rel : rels) {
    //         printf("%d: ", rel.first);
    //         auto& set = *rel.second;
    //         for (auto& item : set) {
    //             printf("%d ", item);
    //         }
    //         printf("\n");
    //     }
    // #endif

    return true;
}

bool impose_rule_3i_1o(Minisat::Solver& s, vec<vec<Lit>>& out_refined, int& k, std::string& o, int& x, int& y, int& z, int& w)
{
#if DIFF_BITS == 4
    int vars[DIFF_BITS * 4] = { x,
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
    int vals[DIFF_BITS * 4] = { int_value(s, x),
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
#else
    int vars[DIFF_BITS * 4] = { x, y, z, w };
    int vals[DIFF_BITS * 4] = { int_value(s, x), int_value(s, y), int_value(s, z), int_value(s, w) };
#endif

    uint8_t o1_val[DIFF_BITS];
    from_gc(o[0], o1_val);

    for (int j = 0; j < DIFF_BITS; j++) {
        if (vals[DIFF_BITS * 3 + j] == o1_val[j])
            continue;

        out_refined.push();
        out_refined[k].push(mkLit(vars[j + DIFF_BITS * 3], o1_val[j] == 0));
        for (int i = 0; i < DIFF_BITS * 3; i++)
            out_refined[k].push(mkLit(vars[i], vals[i] == 1));
#if DEBUG
        print_clause(out_refined[k]);
#endif
        k++;
        return true;
    }

    return false;
}

// TODO: Finish this
// bool impose_rule_3i_2o(Minisat::Solver& s, vec<vec<Lit>>& out_refined, int& k, std::string& o, int& x, int& y, int& z, int& v, int& w)
// {
// #if DIFF_BITS == 4
//     // TODO: Implement for diff 4
// #else
//     int vars[DIFF_BITS * 5] = { x, y, z, v, w };
//     int vals[DIFF_BITS * 5] = { int_value(s, x), int_value(s, y), int_value(s, z), int_value(s, v), int_value(s, w) };
// #endif

//     uint8_t o_vals[DIFF_BITS * 2];
//     from_gc(o[0], &o_vals[0]);
//     from_gc(o[1], &o_vals[1]);

//     for (int j = 0; j < 2 * DIFF_BITS; j++) {
//         if (vals[DIFF_BITS * 3 + j] == o_vals[j])
//             continue;

//         out_refined.push();
//         out_refined[k].push(mkLit(vars[j + DIFF_BITS * 3], o_vals[j] == 0));
//         for (int i = 0; i < DIFF_BITS * 3; i++)
//             out_refined[k].push(mkLit(vars[i], vals[i] == 1));
// #if DEBUG
//         print_clause(out_refined[k]);
// #endif
//         k++;
//         return true;
//     }

//     return false;
// }

bool impose_rule_3i_1o_w(Minisat::Solver& s, vec<vec<Lit>>& out_refined, int& k, int id, int& x, int& y, int& z, int& w)
{
    if (k != 0) {
        return false;
    }

    char x_gc = to_gc(s, x), y_gc = to_gc(s, y), z_gc = to_gc(s, z), w_gc = to_gc(s, w);
    if (x_gc == NULL || y_gc == NULL || z_gc == NULL) {
        return false;
    }

    std::string r_key = std::to_string(id);
    r_key.push_back(x_gc);
    r_key.push_back(y_gc);
    r_key.push_back(z_gc);
    auto r_value_it = s.rules.find(r_key);
    // Skip if the rule isn't found
    if (r_value_it == s.rules.end())
        return false;

    auto r_value = r_value_it->second;

    // TODO: Remove this check in the future
    // Skip if the rule's output isn't usable
#if DIFF_BITS == 1
    if (r_value[0] != '-' && r_value[0] != 'x')
        return false;
#endif

    // Skip if the output difference is already correct
    if (w_gc == r_value[0]) {
#if DEBUG
        // printf("%d is correct (%c)\n", w + 1, w_gc);
#endif
        return false;
    }

#if DEBUG
    printf("RKey: %s; RValue: %s; DValue: %c(%d); DId: %d\n", r_key.c_str(), r_value.c_str(), w_gc, w_gc, w + 1);
#endif
    return impose_rule_3i_1o(s, out_refined, k, r_value, x, y, z, w);
}

// TODO: Add XOR rule key construction
void impose_oi_rule(Minisat::Solver& s, vec<vec<Lit>>& out_refined, int& k, int function_id, int* var_ids, int vars_n, int outputs_n)
{
    // Rule key
    int key_size = (outputs_n * 3) + 2;
    char rule_key[key_size];
    rule_key[0] = function_id;
    rule_key[key_size - 1] = NULL;

    // Determine if the GCs are known
    bool gcs_known = true;
    for (int i = 0; i < outputs_n * 3; i += 3) {
        int x = int_value(s, var_ids[i]);
        int x_ = int_value(s, var_ids[i + 1]);
        int dx = int_value(s, var_ids[i + 2]);

        if ((x != -1 && x_ != -1) || ((x != -1 || x_ != -1) && dx != -1)) {

        } else {
            gcs_known = false;
            break;
        }
    }

    // Construct the rule key
    if (gcs_known) {
        for (int i = 0, j = 1; i < outputs_n * 3; i += 3, j++) {
            int x = int_value(s, var_ids[i]);
            int x_ = int_value(s, var_ids[i + 1]);
            int dx = int_value(s, var_ids[i + 2]);
            char gc;
            if (x != -1 && x_ != -1)
                gc = to_gc(x, x_);
            else if ((x == 0 || x_ == 0) && dx == 0)
                gc = '0';
            else if ((x == 1 || x_ == 1) && dx == 0)
                gc = '1';
            else if ((x == 0 || x_ == 0) && dx == 1)
                gc = 'u';
            else
                gc = 'n';
            rule_key[j] = gc;
        }
    } else {
        return;
    }

    // Find the value of the rule (if it exists)
    auto rule_it = s.rules.find(rule_key);
    if (rule_it == s.rules.end())
        return;
    auto rule_value = rule_it->second;

    // TODO: Impose the rule on the input bits (x, x', and dx)
}

// TODO: Reduce redundancy
// The variable IDs provided should include the operands and the output
void add_2_bit_conditions(Minisat::Solver& s, vec<vec<Lit>>& out_refined, int& k, int function_id, int* var_ids, int vars_n, std::vector<std::pair<int32_t, int32_t>>& equations)
{
    assert(DIFF_BITS == 0);
    // Number of variables
    assert(vars_n % 3 == 0 && vars_n > 0); // Must be in triples and non-empty

    // Rule key
    int key_size = (vars_n / 3) + 2;
    char rule_key[key_size];
    rule_key[0] = function_id;
    rule_key[key_size - 1] = NULL;

    // Determine if the XORs are known
    bool xors_known = true;
    for (int i = 2; i < vars_n; i += 3) {
        if (int_value(s, var_ids[i]) == -1) {
            xors_known = false;
            break;
        }
    }

    // Determine if the GCs are known
    bool gcs_known = true;
    for (int i = 0; i < vars_n; i += 3) {
        int x = int_value(s, var_ids[i]);
        int x_ = int_value(s, var_ids[i + 1]);
        int dx = int_value(s, var_ids[i + 2]);
        // TODO: Invert the condition and convert into an if-only statement
        if ((x != -1 && x_ != -1) || ((x != -1 || x_ != -1) && dx != -1)) {

        } else {
            gcs_known = false;
            break;
        }
    }

    // Construct the rule key
    if (gcs_known) {
        for (int i = 0, j = 1; i < vars_n; i += 3, j++) {
            int x = int_value(s, var_ids[i]);
            int x_ = int_value(s, var_ids[i + 1]);
            int dx = int_value(s, var_ids[i + 2]);
            char gc;
            if (x != -1 && x_ != -1)
                gc = to_gc(x, x_);
            else if ((x == 0 || x_ == 0) && dx == 0)
                gc = '0';
            else if ((x == 1 || x_ == 1) && dx == 0)
                gc = '1';
            else if ((x == 0 || x_ == 0) && dx == 1)
                gc = 'u';
            else
                gc = 'n';
            rule_key[j] = gc;
        }
    } else if (xors_known) {
        for (int i = 2, j = 1; i < vars_n; i += 3, j++) {
            char gc = int_value(s, var_ids[i]) == 1 ? 'x' : '-';
            rule_key[j] = gc;
        }
    } else {
        return;
    }

    // Find the value of the rule (if it exists)
    auto rule_it = s.rules.find(rule_key);
    if (rule_it == s.rules.end())
        return;
    auto rule_value = rule_it->second;

    // Derive the equations from the relationships between the input vars
#if DEBUG
    // printf("2-bit %d %s %s\n", function_id, rule_key, rule_value.c_str());
#endif
    std::set<int> visited;
    int equations_added = 0;
    for (int i = 0; i < rule_value.length(); i++) {
        if (rule_value[i] == '2')
            continue;

        int var_id = var_ids[i * 3];
        for (int j = 0; j < vars_n - 3; j += 3) {
            int var2_id = var_ids[j];
            if (var2_id == var_id || visited.find(var2_id) != visited.end())
                continue;

            // DEBUG
            // printf("%d %d %c\n", var_id, var2_id, rule_value[i]);
            if (rule_value[i] == '1') {
                equations.push_back({ var_id, var2_id });
                equations_added++;
            } else if (rule_value[i] == '0') {
                equations.push_back({ var_id, -var2_id });
                equations_added++;
            }
        }

        visited.insert(var_id);
    }

    // Check the consistency of the entire set of equations (includes one added for other functions)
    if (!check_consistency(equations)) {
        // Block the input variables that lead to the contradiction
        if (gcs_known) {
            // out_refined.push();
            for (int i = 0; i < vars_n - 3; i++) {
                // printf("%d\n", int_value(s, var_ids[i]));
                // out_refined[k].push(mkLit(var_ids[i], int_value(s, var_ids[i]) == 1));
                // print_clause(out_refined[k]);
            }
            k++;
        } else if (xors_known) {
            out_refined.push();
            for (int i = 2; i < vars_n - 3; i += 3)
                out_refined[k].push(mkLit(var_ids[i], int_value(s, var_ids[i]) == 1));
            k++;
        }
        // // Block the equations that lead to the contradiction
        // for (int i = 0; i < equations_added; i++) {
        //     auto equation = equations.back();
        //     equations.pop_back();

        //     int var1 = equation.first;
        //     int var2 = equation.second;

        //     if (var1 > 0 && var2 > 0) {
        //         // Block equality
        //         out_refined.push();
        //         out_refined[k].push(mkLit(abs(var2), false));
        //         out_refined[k].push(mkLit(abs(var1), false));
        //         print_clause(out_refined[k]);
        //         k++;
        //         out_refined.push();
        //         out_refined[k].push(mkLit(abs(var1), true));
        //         out_refined[k].push(mkLit(abs(var2), true));
        //         print_clause(out_refined[k]);
        //         k++;
        //     } else {
        //         // Block inequality
        //         out_refined.push();
        //         out_refined[k].push(mkLit(abs(var2), true));
        //         out_refined[k].push(mkLit(abs(var1), false));
        //         print_clause(out_refined[k]);
        //         k++;
        //         out_refined.push();
        //         out_refined[k].push(mkLit(abs(var1), true));
        //         out_refined[k].push(mkLit(abs(var2), false));
        //         print_clause(out_refined[k]);
        //         k++;
        //     }
        // }
        equations.clear();
    }
}

bool sort_by_value(const std::pair<int, lbool>& a, const std::pair<int, lbool>& b)
{
    return a.second == l_Undef;
}

void infer_carries(Minisat::Solver& s, vec<vec<Lit>>& out_refined, int& k, int* var_ids, int vars_n, int carries_n)
{
    int inputs_n = vars_n - carries_n - 1, input_1s_n = 0;
    for (int i = 0; i < inputs_n; i++)
        if (s.value(var_ids[i]) == l_True)
            input_1s_n++;

    std::vector<std::pair<int, lbool>> vars;
    for (int i = 0; i < inputs_n; i++) {
        vars.push_back({ var_ids[i], s.value(var_ids[i]) });
    }

    // // Reorder variables by placing the undefined variables first in the array
    // int vars_reordered[inputs_n];
    // int j = 0;
    // for (int i = 0; i < inputs_n; i++)
    //     if (s.value(var_ids[i]) == l_Undef) {
    //         vars_reordered[j] = var_ids[i];
    //         j++;
    //     }
    // for (int i = 0; i < inputs_n; i++)
    //     if (s.value(var_ids[i]) != l_Undef) {
    //         vars_reordered[j] = var_ids[i];
    //         j++;
    //     }

    // High carry must be 1 if no. of 1s >= 4
    if (carries_n == 2) {
        int high_carry_id = var_ids[inputs_n];
        vars.push_back({ high_carry_id, s.value(high_carry_id) });
        std::sort(vars.begin(), vars.end(), sort_by_value);

        if (input_1s_n >= 4 && s.value(high_carry_id) != l_True) {
            out_refined.push();
            // out_refined[k].push(mkLit(high_carry_id, false));
            for (int i = 0; i < vars.size(); i++) {
                printf("%d: %d\n", vars[i].first + 1, int_value(s, vars[i].first));
                out_refined[k].push(mkLit(vars[i].first, vars[i].first == high_carry_id ? false : vars[i].second == l_True));
            }

            print_clause(out_refined[k]);
            k++;
            printf("Inferred high carry %d %d %d %d\n", inputs_n, input_1s_n, high_carry_id + 1, int_value(s, high_carry_id));
        }
    }

    if (vars_n < 6)
        return;

    // Low carry must be 1 if no. of 1s >= 6
    // int low_carry_id = var_ids[vars_n - carries_n + 1];
    // if (input_1s_n >= 6 && s.value(low_carry_id) != l_True) {
    //     out_refined.push();
    //     out_refined[k].push(mkLit(low_carry_id, false));
    //     for (int i = 0; i < inputs_n; i++)
    //         out_refined[k].push(mkLit(var_ids[i], s.value(var_ids[i]) == l_True));
    //     k++;
    //     printf("Inferred low carry %d\n", input_1s_n);
    // }
}

void add_clauses(Minisat::Solver& s, vec<vec<Lit>>& out_refined)
{
    std::vector<std::pair<int32_t, int32_t>> equations;

    int k = 0;
    for (int i = 0; i < s.steps; i++) {
        int dw_0_base = s.var_map["DW_" + std::to_string(i) + "_g"];
        int da_4_base = s.var_map["DA_" + std::to_string(i + 4) + "_g"];
        int da_3_base = s.var_map["DA_" + std::to_string(i + 3) + "_g"];
        int da_2_base = s.var_map["DA_" + std::to_string(i + 2) + "_g"];
        int da_1_base = s.var_map["DA_" + std::to_string(i + 1) + "_g"];
        int da_0_base = s.var_map["DA_" + std::to_string(i) + "_g"];
        int a_3_base_f = s.var_map["A_" + std::to_string(i + 3) + "_f"];
        int a_3_base_g = s.var_map["A_" + std::to_string(i + 3) + "_g"];
        int de_4_base = s.var_map["DE_" + std::to_string(i + 4) + "_g"];
        int de_3_base = s.var_map["DE_" + std::to_string(i + 3) + "_g"];
        int de_2_base = s.var_map["DE_" + std::to_string(i + 2) + "_g"];
        int de_1_base = s.var_map["DE_" + std::to_string(i + 1) + "_g"];
        int de_0_base = s.var_map["DE_" + std::to_string(i) + "_g"];
        int df1_base = s.var_map["Df1_" + std::to_string(i) + "_g"];
        int df2_base = s.var_map["Df2_" + std::to_string(i) + "_g"];
        int dsigma0_base = s.var_map["Dsigma0_" + std::to_string(i) + "_g"];
        int dsigma1_base = s.var_map["Dsigma1_" + std::to_string(i) + "_g"];
        int sigma0_f_base = s.var_map["Sigma0_" + std::to_string(i) + "_f"];
        int sigma1_f_base = s.var_map["Sigma1_" + std::to_string(i) + "_f"];
        int sigma0_g_base = s.var_map["Sigma0_" + std::to_string(i) + "_g"];
        int sigma1_g_base = s.var_map["Sigma1_" + std::to_string(i) + "_g"];
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
            int dw_0 = dw_0_base + j * DIFF_BITS; // DW[i]
            int da_4 = da_4_base + j * DIFF_BITS; // DA[i+4]
            int da_3 = da_3_base + j * DIFF_BITS; // DA[i+3]
            int da_2 = da_2_base + j * DIFF_BITS; // DA[i+2]
            int da_1 = da_1_base + j * DIFF_BITS; // DA[i+1]
            // int a_3_f = a_3_base_f + j * DIFF_BITS; // A[i+3]
            // int a_3_g = a_3_base_g + j * DIFF_BITS; // A[i+3]'
            int da_0 = da_0_base + j * DIFF_BITS; // DA[i]
            int de_4 = de_4_base + j * DIFF_BITS; // DE[i+4]
            int de_3 = de_3_base + j * DIFF_BITS; // DE[i+3]
            int de_2 = de_2_base + j * DIFF_BITS; // DE[i+2]
            int de_1 = de_1_base + j * DIFF_BITS; // DE[i+1]
            int de_0 = de_0_base + j * DIFF_BITS; // DE[i]
            int df1 = df1_base + j * DIFF_BITS; // Df1 <- IF
            int df2 = df2_base + j * DIFF_BITS; // Df2 <- MAJ
            int dsigma0 = dsigma0_base + j * DIFF_BITS; // DSigma0
            int dsigma1 = dsigma1_base + j * DIFF_BITS; // DSigma1
            int sigma0_f = sigma0_f_base + j * DIFF_BITS; // Sigma0
            int sigma0_g = sigma0_g_base + j * DIFF_BITS; // Sigma0'
            int sigma1_f = sigma1_f_base + j * DIFF_BITS; // Sigma1
            int sigma1_g = sigma1_g_base + j * DIFF_BITS; // Sigma1'
            int ds0 = ds0_base + j * DIFF_BITS; // DS0
            int ds1 = ds1_base + j * DIFF_BITS; // DS1
            int dt = dt_base + j * DIFF_BITS; // DT
            int dk = dk_base + j * DIFF_BITS; // DT
            int dr1_carry = dr1_carry_base + j * DIFF_BITS; // Dr1_carry
            int dr2_carry = dr2_carry_base + j * DIFF_BITS; // Dr2_carry
            int dr2_carry2 = dr2_carry2_base + (j - 1) * DIFF_BITS; // Dr2_carry2
            int dr0_carry = dr0_carry_base + j * DIFF_BITS; // Dr0_carry
            int dr0_carry2 = dr0_carry2_base + j * DIFF_BITS; // Dr0_carry2
            int dw_carry = dw_carry_base + j * DIFF_BITS; // Dw_carry
            int dw_carry2 = dw_carry2_base + j * DIFF_BITS; // Dw_carry2

            // IF
#if DIFF_BITS == 4
            if (impose_rule_3i_1o_w(s, out_refined, k, IO_CONSTRAINT_IF_ID, de_3, de_2, de_1, df1))
                goto END_CALLBACK;
            // MAJ
            if (impose_rule_3i_1o_w(s, out_refined, k, IO_CONSTRAINT_MAJ_ID, da_3, da_2, da_1, df2))
                goto END_CALLBACK;
#endif

            // Sigma0 2-bit conditions
            {
                int x = a_3_base_f + (j + 2) % 32;
                int x_prime = a_3_base_g + (j + 2) % 32;
                int dx = da_3_base + (j + 2) % 32;

                int y = a_3_base_f + (j + 13) % 32;
                int y_prime = a_3_base_g + (j + 13) % 32;
                int dy = da_3_base + (j + 13) % 32;

                int z = a_3_base_f + (j + 22) % 32;
                int z_prime = a_3_base_g + (j + 22) % 32;
                int dz = da_3_base + (j + 22) % 32;

                int var_ids[] = { x, x_prime, dx, y, y_prime, dy, z, z_prime, dz, sigma0_f, sigma0_g, dsigma0 };
                // add_2_bit_conditions(s, out_refined, k, TWO_BIT_CONSTRAINT_XOR3_ID, (int*)var_ids, 12, equations);
            }

            // Sigma1 2-bit conditions
            {
                int x = a_3_base_f + (j + 6) % 32;
                int x_prime = a_3_base_g + (j + 6) % 32;
                int dx = da_3_base + (j + 6) % 32;

                int y = a_3_base_f + (j + 11) % 32;
                int y_prime = a_3_base_g + (j + 11) % 32;
                int dy = da_3_base + (j + 11) % 32;

                int z = a_3_base_f + (j + 25) % 32;
                int z_prime = a_3_base_g + (j + 25) % 32;
                int dz = da_3_base + (j + 25) % 32;

                int var_ids[] = { x, x_prime, dx, y, y_prime, dy, z, z_prime, dz, sigma1_f, sigma1_g, dsigma1 };
                // add_2_bit_conditions(s, out_refined, k, TWO_BIT_CONSTRAINT_XOR3_ID, (int*)var_ids, 12, equations);
            }

            // TODO: s0
            // TODO: s1

            // Compression: 3 to 2
            {
                int op1 = da_0;
                int op2 = dt;
                int op3 = dr1_carry - 1 * DIFF_BITS;
                int o1 = dr1_carry;
                int o2 = de_4;

                // bool out_def_in_undef = int_value(s, o1) != -1 && int_value(s, o2) != -1 && (int_value(s, op1) == -1 || int_value(s, op2) == -1 || int_value(s, op3) == -1);

                // if (j > 0) {
                //     int var_ids[] = { op1, op2, op3, o1, o2 };
                //     infer_carries(s, out_refined, k, var_ids, 5, 1);
                // } else {
                //     int var_ids[] = { op1, op2, o1, o2 };
                //     infer_carries(s, out_refined, k, var_ids, 4, 1);
                // }

                // if (out_def_in_undef)
                //     printf("ADD3: %d %d %d = %d %d\n", int_value(s, op1), int_value(s, op2), int_value(s, op3), int_value(s, o1), int_value(s, o2));

                // if (int_value(s, op1) == -1 || int_value(s, op2) == -1 || int_value(s, op3) == -1 || int_value(s, o1) == -1 || int_value(s, o2) == -1)
                //     printf("%d %d %d = %d %d\n", int_value(s, op1), int_value(s, op2), int_value(s, op3), int_value(s, o1), int_value(s, o2));

                // if (int_value(s, op1) != -1 && int_value(s, op2) != -1) {
                //     if (j > 0 && int_value(s, op3) != -1) {
                //         // comp_3_2(s, out_refined, k, i, j, op1, op2, op3, o1, o2);
                //     } else if (j == 0) {
                //         // comp_2_2(s, out_refined, k, i, j, op1, op2, o1, o2);
                //     }
                // }
            }

            // Compression: 5 to 3
            // g.cnf.diff_add(DA[i + 4], DT[i], Dsigma0[i], Dr2carry[i], Dr2Carry[i],
            // Df2[i]);
            {
                int op1 = dt;
                int op2 = dsigma0;
                int op3 = df2;
                int op4 = dr2_carry - 1 * DIFF_BITS; // t[j - 1]
                int op5 = dr2_carry2 - 2 * DIFF_BITS; // T[j - 2]
                int o1 = dr2_carry2; // T[j]
                int o2 = dr2_carry; // t[j]
                int o3 = da_4; // DA[i+4]

                // if (j > 2) {
                //     int var_ids[] = { op1, op2, op3, op4, op5, o1, o2, o3 };
                //     infer_carries(s, out_refined, k, var_ids, 8, 2);
                // } else if (j == 2 || j == 1) {
                //     int var_ids[] = { op1, op2, op3, op4, o1, o2, o3 };
                //     infer_carries(s, out_refined, k, var_ids, 7, 2);
                // } else {
                //     int var_ids[] = { op1, op2, op3, o1, o2 };
                //     infer_carries(s, out_refined, k, var_ids, 5, 1);
                // }

                // if (int_value(s, o1) != -1 && int_value(s, o2) != -1 && int_value(s, o3) != -1 && (int_value(s, op1) == -1 || int_value(s, op2) == -1 || int_value(s, op3) == -1 || int_value(s, op4) == -1 || int_value(s, op5) == -1))
                //     printf("ADD5: %d %d %d %d %d = %d %d %d\n", int_value(s, op1), int_value(s, op2), int_value(s, op3), int_value(s, op4), int_value(s, op5), int_value(s, o1), int_value(s, o2), int_value(s, o3));

                // if (int_value(s, op1) != -1 && int_value(s, op2) != -1 && int_value(s, op3) != -1) {
                //     if (j > 2 && int_value(s, op4) != -1 && int_value(s, op5) != -1) {
                //         comp_5_3(s, out_refined, k, i, j, op1, op2, op3, op4, op5, o1, o2,
                //             o3);
                //     } else if ((j == 2 || j == 1) && int_value(s, op4) != -1) {
                //         comp_4_3(s, out_refined, k, i, j, op1, op2, op3, op4, o1, o2, o3);
                //     } else if (j == 0) {
                //         comp_3_2(s, out_refined, k, i, j, op1, op2, op3, o2, o3);
                //     }
                // }
            }

            // Compression: 6 to 3
            // g.cnf.diff_add(DW[i], DW[i - 16], Ds0[i], Dwcarry[i], DwCarry[i], DW[i
            // - 7], Ds1[i]);
            if (i >= 16) {
                int op1 = s.var_map["DW_" + std::to_string(i - 16) + "_g"] + j * DIFF_BITS;
                int op2 = ds0;
                int op3 = s.var_map["DW_" + std::to_string(i - 7) + "_g"] + j * DIFF_BITS;
                int op4 = ds1;
                int op5 = dw_carry - 1 * DIFF_BITS; // t[j - 1]
                int op6 = dw_carry2 - 2 * DIFF_BITS; // T[j - 2]
                int o1 = dw_carry2; // T[j]
                int o2 = dw_carry; // t[j]
                int o3 = dw_0; // DW[i]

                // if (j > 1) {
                //     int var_ids[] = { op1, op2, op3, op4, op5, op6, o1, o2, o3 };
                //     infer_carries(s, out_refined, k, var_ids, 9, 2);
                // } else if (j == 1) {
                //     int var_ids[] = { op1, op2, op3, op4, op5, o1, o2, o3 };
                //     infer_carries(s, out_refined, k, var_ids, 8, 2);
                // } else {
                //     int var_ids[] = { op1, op2, op3, op4, o1, o2, o3 };
                //     infer_carries(s, out_refined, k, var_ids, 7, 2);
                // }

                // if (int_value(s, o1) != -1 && int_value(s, o2) != -1 && int_value(s, o3) != -1 && (int_value(s, op1) == -1 || int_value(s, op2) == -1 || int_value(s, op3) == -1 || int_value(s, op4) == -1 || int_value(s, op5) == -1 || int_value(s, op6) == -1))
                //     printf("ADD6: %d %d %d %d %d %d = %d %d %d\n", int_value(s, op1), int_value(s, op2), int_value(s, op3), int_value(s, op4), int_value(s, op5), int_value(s, op6), int_value(s, o1), int_value(s, o2), int_value(s, o3));

                // if (s.value(op1) != l_Undef && s.value(op2) != l_Undef && s.value(op3) != l_Undef && s.value(op4) != l_Undef) {
                //     if (j > 1 && int_value(s, op5) != -1 && int_value(s, op6) != -1) {
                //         comp_6_3(s, out_refined, k, i, j, op1, op2, op3, op4, op5, op6, o1,
                //             o2, o3);
                //     } else if (j == 1 && int_value(s, op5) != -1) {
                //         comp_5_3(s, out_refined, k, i, j, op1, op2, op3, op4, op5, o1, o2,
                //             o3);
                //     } else if (j == 0) {
                //         comp_4_3(s, out_refined, k, i, j, op1, op2, op3, op4, o1, o2, o3);
                //     }
                // }
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
                int op6 = dr0_carry - 1 * DIFF_BITS; // t[j - 1]
                int op7 = dr0_carry2 - 2 * DIFF_BITS; // T[j - 2]
                int o1 = dr0_carry2; // T[j]
                int o2 = dr0_carry; // t[j]
                int o3 = dt; // DT[i]

                if (j > 1) {
                    int var_ids[] = { op1, op2, op3, op4, op5, op6, op7, o1, o2, o3 };
                    // infer_carries(s, out_refined, k, var_ids, 10, 2);
                } else if (j == 1) {
                    // int var_ids[] = { op1, op2, op3, op4, op5, op6, o1, o2, o3 };
                    // infer_carries(s, out_refined, k, var_ids, 9, 2);
                } else {
                    // int var_ids[] = { op1, op2, op3, op4, op5, o1, o2, o3 };
                    // infer_carries(s, out_refined, k, var_ids, 8, 2);
                }

                // if (int_value(s, o1) != -1 && int_value(s, o2) != -1 && int_value(s, o3) != -1 && (int_value(s, op1) == -1 || int_value(s, op2) == -1 || int_value(s, op3) == -1 || int_value(s, op4) == -1 || int_value(s, op5) == -1 || int_value(s, op6) == -1 || int_value(s, op7) == -1))
                //     printf("ADD7: %d %d %d %d %d %d %d = %d %d %d\n", int_value(s, op1), int_value(s, op2), int_value(s, op3), int_value(s, op4), int_value(s, op5), int_value(s, op6), int_value(s, op7), int_value(s, o1), int_value(s, o2), int_value(s, o3));

                // if (int_value(s, op1) != -1 && int_value(s, op2) != -1 && int_value(s, op3) != -1 && int_value(s, op4) != -1 && int_value(s, op5) != -1) {
                //     if (j > 1 && int_value(s, op6) != -1 && int_value(s, op7) != -1) {
                //         comp_7_3(s, out_refined, k, i, j, op1, op2, op3, op4, op5, op6, op7,
                //             o1, o2, o3);
                //     } else if (j == 1 && int_value(s, op6)) {
                //         comp_6_3(s, out_refined, k, i, j, op1, op2, op3, op4, op5, op6, o1,
                //             o2, o3);
                //     } else if (j == 0) {
                //         comp_5_3(s, out_refined, k, i, j, op1, op2, op3, op4, op5, o1, o2,
                //             o3);
                //     }
                // }
            }
        }
    }
END_CALLBACK:
    return;
}