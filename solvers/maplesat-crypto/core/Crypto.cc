#include "Crypto.h"
#include <algorithm>
#include <map>
#include <set>
#include <vector>
#define DEBUG true
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
    } else if (id == TWO_BIT_CONSTRAINT_ADD4_ID) {
        key_size = 5;
        val_size = 6;
    } else if (id == TWO_BIT_CONSTRAINT_ADD5_ID) {
        key_size = 6;
        val_size = 10;
    } else if (id == TWO_BIT_CONSTRAINT_ADD6_ID) {
        key_size = 7;
        val_size = 15;
    } else if (id == TWO_BIT_CONSTRAINT_ADD7_ID) {
        key_size = 8;
        val_size = 21;
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

    printf("Loaded %d rules into %d buckets\n", count, solver.rules.bucket_count());

    // DEBUG
    // char key[] = {19, 'n', '1', '1', 'n'};
    // printf("Found: %s\n", solver.rules.find(key)->second.c_str());
    // exit(0);
}

void add_to_var_ids(Solver& solver, std::string prefix, std::string suffix, std::vector<int>& var_ids, int inputs_n, int outputs_n)
{
    for (int i = 0; i < inputs_n; i++)
        var_ids.push_back(solver.var_map[prefix + "x" + std::to_string(i) + suffix]);

    for (int i = 0; i < outputs_n; i++)
        var_ids.push_back(solver.var_map[prefix + "z" + std::to_string(i) + suffix]);
}

void processVarMap(Solver& solver)
{
    printf("Varmap: %d %d\n", solver.var_map.size(), solver.steps);
    for (int i = 0; i < solver.steps; i++) {
        // add_w
        add_to_var_ids(solver, "add_w" + std::to_string(i) + "_", "_f", solver.var_ids_.add_w_f[i], 6, 2);
        add_to_var_ids(solver, "add_w" + std::to_string(i) + "_", "_g", solver.var_ids_.add_w_g[i], 6, 2);

        // add_t
        add_to_var_ids(solver, "add_T" + std::to_string(i) + "_", "_f", solver.var_ids_.add_t_f[i], 7, 2);
        add_to_var_ids(solver, "add_T" + std::to_string(i) + "_", "_g", solver.var_ids_.add_t_g[i], 7, 2);

        // add_e
        add_to_var_ids(solver, "add_E" + std::to_string(i + 4) + "_", "_f", solver.var_ids_.add_e_f[i], 3, 1);
        add_to_var_ids(solver, "add_E" + std::to_string(i + 4) + "_", "_g", solver.var_ids_.add_e_g[i], 3, 1);

        // add_a
        add_to_var_ids(solver, "add_A" + std::to_string(i + 4) + "_", "_f", solver.var_ids_.add_a_f[i], 5, 2);
        add_to_var_ids(solver, "add_A" + std::to_string(i + 4) + "_", "_g", solver.var_ids_.add_a_g[i], 5, 2);
    }
}

int int_value(Minisat::Solver& s, int var)
{
    auto value = s.value(var);
    return value == l_True ? 1 : value == l_False ? 0
                                                  : -1;
}

char to_gc(lbool x, lbool x_prime)
{
    if (x == l_False && x_prime == l_False)
        return '0';
    else if (x == l_True && x_prime == l_False)
        return 'u';
    else if (x == l_False && x_prime == l_True)
        return 'n';
    else if (x == l_True && x_prime == l_True)
        return '1';
    else
        return NULL;
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
    int d = int_value(s, id);
    if (d == 0) {
        return '-';
    } else if (d == 1) {
        return 'x';
    }
    return NULL;
}

void from_gc(char& gc, uint8_t* vals)
{
    if (gc == '-') {
        vals[0] = 0;
    } else if (gc == 'x') {
        vals[0] = 1;
    }
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

// TODO: Get rid of unwanted vector
// TODO: Reduce redundancy
// The variable IDs provided should include the operands and the output
void add_2_bit_conditions(Minisat::Solver& s, vec<vec<Lit>>& out_refined, int& k, int function_id, int* var_ids, int vars_n)
{
    // Number of variables
    assert(vars_n % 3 == 0 && vars_n > 0); // Must be in triples and non-empty
    int chunks_n = vars_n / 3;

    // Lay out the rule key's foundation
    int key_size = chunks_n + 2;
    char rule_key[key_size];
    rule_key[0] = function_id;
    rule_key[key_size - 1] = NULL;

    // Process chunk-wise (each chunk has 3 bits)
    vec<Lit> base_clause;
    bool hasXOrDash = false;
    for (int i = 0, j = 1; i < vars_n; i += 3, j++) {
        // There are 3 possible ways to derive the GC of the chunk: from x and x_, from dx and x or x_, or from dx alone, else we can't
        int& x_id = var_ids[i];
        lbool x_value = s.value(var_ids[i]);

        int& x_prime_id = var_ids[i + 1];
        lbool x_prime_value = s.value(var_ids[i + 1]);

        int& dx_id = var_ids[i + 2];
        lbool dx_value = s.value(var_ids[i + 2]);

        // printf("Values: %d %d %d; IDs: %d %d %d\n", int_value(s, x_id), int_value(s, x_prime_id), int_value(s, dx_id), x_id + 1, x_prime_id + 1, dx_id + 1);

        if (x_value != l_Undef && x_prime_value != l_Undef) {
            rule_key[j] = to_gc(x_value, x_prime_value);
            base_clause.push(mkLit(x_id, x_value == l_True));
            base_clause.push(mkLit(x_prime_id, x_prime_value == l_True));
        } else if (dx_value != l_Undef && (x_value != l_Undef || x_prime_value != l_Undef)) {
            // y is x or x' that is defined
            int y_id;
            lbool y_value;
            bool x_defined = false;
            if (x_value != l_Undef) {
                y_id = x_id;
                y_value = x_value;
                x_defined = true;
            } else {
                y_id = x_prime_id;
                y_value = x_prime_value;
            }

            if (dx_value == l_False)
                rule_key[j] = x_defined && y_value == l_False ? '0' : '1';
            else
                rule_key[j] = x_defined && y_value == l_False ? 'n' : 'u';

            base_clause.push(mkLit(y_id, y_value == l_True));
            base_clause.push(mkLit(dx_id, dx_value == l_True));
        } else if (dx_value != l_Undef) {
            rule_key[j] = dx_value == l_True ? 'x' : '-';
            hasXOrDash = true;
            base_clause.push(mkLit(dx_id, dx_value == l_True));
        } else {
            // Terminate since we can't derive the rule if we don't know any of {1, u, n, 0, x, -}, and without the rule we can't derive the 2-bit conditions
            return;
        }
    }

    if (!hasXOrDash)
        return;

    // Find the value of the rule (if it exists)
    auto rule_it = s.rules.find(rule_key);
    if (rule_it == s.rules.end())
        return;
    auto rule_value = rule_it->second;
    // printf("Found for key %s: %s\n", rule_key, rule_value.c_str());

    // Derive the relationships between the x and x_ of the chunks and enforce them through clauses
    std::set<int> visited;
    int rule_i = -1;
    for (int i = 0; i < vars_n - 3; i += 3) {
        int var1_id = var_ids[i];
        for (int j = 0; j < vars_n - 3; j += 3) {
            int var2_id = var_ids[j];
            // printf("Trying %d\n", var2_id + 1);
            if (var2_id == var1_id || visited.find(var2_id) != visited.end())
                continue;
            rule_i++;
            if (rule_value[rule_i] == '2')
                continue;
            // printf("Passed %d with %d; %c in %d\n", var1_id + 1, var2_id + 1, rule_value[rule_i], rule_i);

            // Inferred variables should be undefined
            lbool var1_value = s.value(var1_id);
            lbool var2_value = s.value(var2_id);

            // Skip if both the values are defined
            if (var1_value != l_Undef && var2_value != l_Undef)
                continue;

            // Skip if both the values are undefined
            if (var1_value == l_Undef && var2_value == l_Undef)
                continue;

            // printf("1175 = %d\n", int_value(s, 1175 - 1));
            // printf("1110 = %d\n", int_value(s, 1110 - 1));

            printf("\nUsing key %s: %s\n", rule_key, rule_value.c_str());

            printf("Related vars: %d and %d; values: %d and %d\n", var1_id + 1, var2_id + 1, int_value(s, var1_id), int_value(s, var2_id));

            printf("DEBUG: ");
            for (int l = 0; l < base_clause.size(); l++) {
                printf("%d = %d, ", var(base_clause[l]) + 1, int_value(s, var(base_clause[i])));
            }
            printf("\n");

            // DEBUG
            if (rule_value[rule_i] == '1') {
                out_refined.push();
                if (var1_value == l_Undef) {
                    out_refined[k].push(mkLit(var1_id, true));
                    out_refined[k].push(mkLit(var2_id, false));
                } else {
                    out_refined[k].push(mkLit(var2_id, false));
                    out_refined[k].push(mkLit(var1_id, true));
                }
                for (int l = 0; l < base_clause.size(); l++) {
                    if (var(base_clause[l]) == var1_id || var(base_clause[l]) == var2_id)
                        continue;
                    out_refined[k].push(base_clause[l]);
                }
                k++;

                out_refined.push();
                if (var1_value == l_Undef) {
                    out_refined[k].push(mkLit(var1_id, false));
                    out_refined[k].push(mkLit(var2_id, true));
                } else {
                    out_refined[k].push(mkLit(var2_id, true));
                    out_refined[k].push(mkLit(var1_id, false));
                }
                for (int l = 0; l < base_clause.size(); l++) {
                    if (var(base_clause[l]) == var1_id || var(base_clause[l]) == var2_id)
                        continue;
                    out_refined[k].push(base_clause[l]);
                }
                k++;

                print_clause(out_refined[k - 1]);
                print_clause(out_refined[k - 2]);
            } else {
                out_refined.push();
                if (var1_value == l_Undef) {
                    out_refined[k].push(mkLit(var1_id, true));
                    out_refined[k].push(mkLit(var2_id, true));
                } else {
                    out_refined[k].push(mkLit(var2_id, true));
                    out_refined[k].push(mkLit(var1_id, true));
                }
                for (int l = 0; l < base_clause.size(); l++) {
                    if (var(base_clause[l]) == var1_id || var(base_clause[l]) == var2_id)
                        continue;
                    out_refined[k].push(base_clause[l]);
                }
                k++;

                out_refined.push();
                if (var1_value == l_Undef) {
                    out_refined[k].push(mkLit(var1_id, false));
                    out_refined[k].push(mkLit(var2_id, false));
                } else {
                    out_refined[k].push(mkLit(var2_id, false));
                    out_refined[k].push(mkLit(var1_id, false));
                }
                for (int l = 0; l < base_clause.size(); l++) {
                    if (var(base_clause[l]) == var1_id || var(base_clause[l]) == var2_id)
                        continue;
                    out_refined[k].push(base_clause[l]);
                }
                k++;

                print_clause(out_refined[k - 1]);
                print_clause(out_refined[k - 2]);
            }
        }

        visited.insert(var1_id);
    }
}

// Prepare the vector of all the operands and carries of addition (may also remove operands equal to carries from previous columns)
std::vector<int> prepare_add_vec(std::vector<int>& ids, int amount, int carry_removal_n = 0)
{
    std::vector<int> new_vec;
    for (int i = 0; i < ids.size(); i++) {
        if (carry_removal_n > 0 && i == 3 || carry_removal_n == 2 && i == 2)
            continue;

        new_vec.push_back(ids[i] + amount);
    }
    return new_vec;
}

void infer_carries(Solver& s, vec<vec<Lit>>& out_refined, int& k, std::vector<int>& var_ids, int vars_n, int carries_n)
{
    int inputs_n = vars_n - carries_n - 1, input_1s_n = 0, input_1s_ids[inputs_n];
    for (int i = 0; i < inputs_n; i++) {
        if (s.value(var_ids[i]) == l_True)
            input_1s_ids[input_1s_n++] = var_ids[i];
        // if (input_1s_n == 6)
        //     break;
    }

    // High carry must be 1 if no. of 1s >= 4
    if (carries_n == 2) {
        int high_carry_id = var_ids[vars_n - 3];
        if (input_1s_n >= 4 && s.value(high_carry_id) != l_True) {
            out_refined.push();
            out_refined[k].push(mkLit(high_carry_id, s.value(high_carry_id) == l_True));
            for (int i = 0; i < 4; i++)
                out_refined[k].push(mkLit(input_1s_ids[i], true));

            // assert(falsifiedClause(s, out_refined[k]));
            printf("Debug: %d, %d %d %d %d\n", int_value(s, high_carry_id), int_value(s, input_1s_ids[0]), int_value(s, input_1s_ids[1]), int_value(s, input_1s_ids[2]), int_value(s, input_1s_ids[3]));
            print_clause(out_refined[k]);
            k++;
            printf("Inferred high carry (inputs %d, carry_id %d)\n", inputs_n, high_carry_id + 1);
        }
    }

    if (vars_n < 6)
        return;

    // Low carry must be 1 if no. of 1s >= 6
    int low_carry_id = var_ids[vars_n - 2];
    if (input_1s_n >= 6 && s.value(low_carry_id) != l_True) {
        out_refined.push();
        out_refined[k].push(mkLit(low_carry_id, false));
        for (int i = 0; i < 6; i++)
            out_refined[k].push(mkLit(input_1s_ids[i], true));
        print_clause(out_refined[k]);
        k++;
        printf("Inferred low carry (inputs %d, carry_id %d)\n", inputs_n, low_carry_id + 1);
    }
}

void add_addition_clauses(Solver& s, vec<vec<Lit>>& out_refined, int& k, int i, int j, std::vector<int>& f, std::vector<int>& g, int vars_n, int carries_n)
{
    if (j > 1) {
        std::vector<int> ids_f = prepare_add_vec(f, j);
        infer_carries(s, out_refined, k, ids_f, vars_n, carries_n);
        // if (out_refined.size() > 0) {
        //     printf("i,j = %d,%d\n", i, j);
        //     goto END_CALLBACK;
        // }
        std::vector<int> ids_g = prepare_add_vec(g, j);
        infer_carries(s, out_refined, k, ids_g, vars_n, carries_n);
    } else if (j == 1) {
        std::vector<int> ids_f = prepare_add_vec(f, j);
        infer_carries(s, out_refined, k, ids_f, vars_n - 1, carries_n);

        std::vector<int> ids_g = prepare_add_vec(g, j);
        infer_carries(s, out_refined, k, ids_g, vars_n - 1, carries_n);
    } else {
        std::vector<int> ids_f = prepare_add_vec(f, j);
        infer_carries(s, out_refined, k, ids_f, vars_n - 2, carries_n);

        std::vector<int> ids_g = prepare_add_vec(g, j);
        infer_carries(s, out_refined, k, ids_g, vars_n - 2, carries_n);
    }
}

void add_clauses(Minisat::Solver& s, vec<vec<Lit>>& out_refined)
{
    int k = 0;
    for (int i = 0; i < s.steps; i++) {
        for (int j = 0; j < 32; j++) {
            // IF

            // printf("DEBUG: %d %d <- %d %d, %d %d, %d %d\n", int_value(s, 1095), int_value(s, 1113), int_value(s, 1109 - 1), int_value(s, 10179 - 1), int_value(s, 18550 - 1), int_value(s, 18532 - 1), int_value(s, 3438 - 1), int_value(s, 12508 - 1));

            // Sigma0 2-bit conditions
            // {
            //     int x = a_3_base_f + (j + 2) % 32;
            //     int x_prime = a_3_base_g + (j + 2) % 32;
            //     int dx = da_3_base + (j + 2) % 32;

            //     int y = a_3_base_f + (j + 13) % 32;
            //     int y_prime = a_3_base_g + (j + 13) % 32;
            //     int dy = da_3_base + (j + 13) % 32;

            //     int z = a_3_base_f + (j + 22) % 32;
            //     int z_prime = a_3_base_g + (j + 22) % 32;
            //     int dz = da_3_base + (j + 22) % 32;

            //     int var_ids[] = { x, x_prime, dx, y, y_prime, dy, z, z_prime, dz, sigma0_f, sigma0_g, dsigma0 };
            //     // printf("DEBUG: %d %d <- %d %d, %d %d, %d\n", int_value(s, 1861 - 1), int_value(s, 1881 - 1), int_value(s, 19681 - 1), int_value(s, 19692 - 1), int_value(s, 10951 - 1), int_value(s, 7222 - 1), int_value(s, 16292 - 1));
            //     // add_2_bit_conditions(s, out_refined, k, TWO_BIT_CONSTRAINT_XOR3_ID, (int*)var_ids, 12);
            //     // if (out_refined.size() > 0)
            //     //     goto END_CALLBACK;
            // }

            // Sigma1 2-bit conditions
            // {
            //     int x = a_3_base_f + (j + 6) % 32;
            //     int x_prime = a_3_base_g + (j + 6) % 32;
            //     int dx = da_3_base + (j + 6) % 32;

            //     int y = a_3_base_f + (j + 11) % 32;
            //     int y_prime = a_3_base_g + (j + 11) % 32;
            //     int dy = da_3_base + (j + 11) % 32;

            //     int z = a_3_base_f + (j + 25) % 32;
            //     int z_prime = a_3_base_g + (j + 25) % 32;
            //     int dz = da_3_base + (j + 25) % 32;

            //     int var_ids[] = { x, x_prime, dx, y, y_prime, dy, z, z_prime, dz, sigma1_f, sigma1_g, dsigma1 };
            //     // add_2_bit_conditions(s, out_refined, k, TWO_BIT_CONSTRAINT_XOR3_ID, (int*)var_ids, 12);
            // }

            // TODO: s0
            // TODO: s1

            // Compression: 3 to 2
            add_addition_clauses(s, out_refined, k, i, j, s.var_ids_.add_e_f[i], s.var_ids_.add_e_g[i], 5, 1);

            // Compression: 5 to 3
            // g.cnf.diff_add(DA[i + 4], DT[i], Dsigma0[i], Dr2carry[i], Dr2Carry[i], Df2[i]);
            add_addition_clauses(s, out_refined, k, i, j, s.var_ids_.add_a_f[i], s.var_ids_.add_a_g[i], 7, 2);

            // Compression: 6 to 3
            // g.cnf.diff_add(DW[i], DW[i - 16], Ds0[i], Dwcarry[i], DwCarry[i], DW[i - 7], Ds1[i]);
            if (i >= 16) {
                add_addition_clauses(s, out_refined, k, i, j, s.var_ids_.add_w_f[i - 16], s.var_ids_.add_w_g[i - 16], 8, 2);
            }

            // Compression: 7 to 3
            // g.cnf.diff_add(DT[i], DE[i], Dsigma1[i], Dr0carry[i], Dr0Carry[i], Df1[i], DK[i], DW[i]);
            add_addition_clauses(s, out_refined, k, i, j, s.var_ids_.add_t_f[i], s.var_ids_.add_t_g[i], 10, 2);
        }
    }
END_CALLBACK:
    return;
}