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

void infer_carries(Solver& s, vec<vec<Lit>>& out_refined, int& k, int* var_ids, int vars_n, int carries_n)
{
    int inputs_n = vars_n - carries_n - 1, input_1s_n = 0, input_1s_ids[6];
    for (int i = 0; i < inputs_n; i++) {
        if (s.value(var_ids[i]) == l_True)
            input_1s_ids[input_1s_n++] = var_ids[i];
        if (input_1s_n == 6)
            break;
    }

    // High carry must be 1 if no. of 1s >= 4
    if (carries_n == 2) {
        int high_carry_id = var_ids[vars_n - 3];
        if (input_1s_n >= 4 && s.value(high_carry_id) != l_True) {
            out_refined.push();
            out_refined[k].push(mkLit(high_carry_id, s.value(high_carry_id) == l_True));
            printf("Debug: %d, %d %d %d %d\n", int_value(s, high_carry_id), int_value(s, input_1s_ids[0]), int_value(s, input_1s_ids[1]), int_value(s, input_1s_ids[2]), int_value(s, input_1s_ids[3]));
            for (int i = 0; i < 4; i++)
                out_refined[k].push(mkLit(input_1s_ids[i], true));
            // assert(falsifiedClause(s, out_refined[k]));
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

void add7_var_ids(Solver& s, int& i, int& j, char block_id, int* var_ids)
{
    // Set the input bits
    for (int index = 0; index < 7; index++)
        var_ids[index] = s.var_map["add_T" + std::to_string(i) + "_x" + std::to_string(index) + "_" + block_id] + j;

    // Set the output bits in order: {high, low} carries, and sum
    var_ids[7] = s.var_map["add_T" + std::to_string(i) + "_z0_" + block_id] + j;
    var_ids[8] = s.var_map["add_T" + std::to_string(i) + "_z1_" + block_id] + j;
    var_ids[9] = s.var_map["T" + std::to_string(i) + "_" + block_id];

    // if (var_ids[7] == 3218 - 1 && j > 1) {
    //     printf("Debug: ");
    //     for (int i = 0; i < 10; i++) {
    //         printf("%d(%d) ", var_ids[i] + 1, int_value(s, var_ids[i]));
    //     }
    //     printf("\n");
    // }
}

void add_clauses(Minisat::Solver& s, vec<vec<Lit>>& out_refined)
{
    int k = 0;
    for (int i = 0; i < s.steps; i++) {
        // int dw_base[1];
        // for (int count = 0; count < 1; count++)
        //     dw_base[count] = s.var_map["DW_" + std::to_string(i + count) + "_g"];
        // int da_base[5];
        // for (int count = 0; count < 5; count++)
        //     da_base[count] = s.var_map["DA_" + std::to_string(i + count) + "_g"];
        // int a_3_base_f = s.var_map["A_" + std::to_string(i + 3) + "_f"];
        // int a_3_base_g = s.var_map["A_" + std::to_string(i + 3) + "_g"];
        // int de_4_base = s.var_map["DE_" + std::to_string(i + 4) + "_g"];
        // int de_3_base = s.var_map["DE_" + std::to_string(i + 3) + "_g"];
        // int de_2_base = s.var_map["DE_" + std::to_string(i + 2) + "_g"];
        // int de_1_base = s.var_map["DE_" + std::to_string(i + 1) + "_g"];
        // int de_0_base = s.var_map["DE_" + std::to_string(i) + "_g"];
        // int df1_base = s.var_map["Df1_" + std::to_string(i) + "_g"];
        // int df2_base = s.var_map["Df2_" + std::to_string(i) + "_g"];
        // int dsigma0_base = s.var_map["Dsigma0_" + std::to_string(i) + "_g"];
        // int dsigma1_base = s.var_map["Dsigma1_" + std::to_string(i) + "_g"];
        // int sigma0_f_base = s.var_map["Sigma0_" + std::to_string(i) + "_f"];
        // int sigma1_f_base = s.var_map["Sigma1_" + std::to_string(i) + "_f"];
        // int sigma0_g_base = s.var_map["Sigma0_" + std::to_string(i) + "_g"];
        // int sigma1_g_base = s.var_map["Sigma1_" + std::to_string(i) + "_g"];
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
            // int dw[] = { dw_base[0] + j }; // DW[i]
            // int da_4 = da_4_base + j; // DA[i+4]
            // int da_3 = da_3_base + j; // DA[i+3]
            // int da_2 = da_2_base + j; // DA[i+2]
            // int da_1 = da_1_base + j; // DA[i+1]
            // // int a_3_f = a_3_base_f + j ; // A[i+3]
            // // int a_3_g = a_3_base_g + j ; // A[i+3]'
            // int da_0 = da_0_base + j; // DA[i]
            // int de_4 = de_4_base + j; // DE[i+4]
            // int de_3 = de_3_base + j; // DE[i+3]
            // int de_2 = de_2_base + j; // DE[i+2]
            // int de_1 = de_1_base + j; // DE[i+1]
            // int de_0 = de_0_base + j; // DE[i]
            // int df1 = df1_base + j; // Df1 <- IF
            // int df2 = df2_base + j; // Df2 <- MAJ
            // int dsigma0 = dsigma0_base + j; // DSigma0
            // int dsigma1 = dsigma1_base + j; // DSigma1
            // int sigma0_f = sigma0_f_base + j; // Sigma0
            // int sigma0_g = sigma0_g_base + j; // Sigma0'
            // int sigma1_f = sigma1_f_base + j; // Sigma1
            // int sigma1_g = sigma1_g_base + j; // Sigma1'
            // int ds0 = ds0_base + j; // DS0
            // int ds1 = ds1_base + j; // DS1
            // int dt = dt_base + j; // DT
            // int dk = dk_base + j; // DT
            // int dr1_carry = dr1_carry_base + j; // Dr1_carry
            // int dr2_carry = dr2_carry_base + j; // Dr2_carry
            // int dr2_carry2 = dr2_carry2_base + (j - 1); // Dr2_carry2
            // int dr0_carry = dr0_carry_base + j; // Dr0_carry
            // int dr0_carry2 = dr0_carry2_base + j; // Dr0_carry2
            // int dw_carry = dw_carry_base + j; // Dw_carry
            // int dw_carry2 = dw_carry2_base + j; // Dw_carry2

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
            // {
            //     int op1 = da_0;
            //     int op2 = dt;
            //     int op3 = dr1_carry - 1;
            //     int o1 = dr1_carry;
            //     int o2 = de_4;

            //     // bool out_def_in_undef = int_value(s, o1) != -1 && int_value(s, o2) != -1 && (int_value(s, op1) == -1 || int_value(s, op2) == -1 || int_value(s, op3) == -1);

            //     // if (j > 0) {
            //     //     int var_ids[] = { op1, op2, op3, o1, o2 };
            //     //     infer_carries(s, out_refined, k, var_ids, 5, 1);
            //     // } else {
            //     //     int var_ids[] = { op1, op2, o1, o2 };
            //     //     infer_carries(s, out_refined, k, var_ids, 4, 1);
            //     // }

            //     // if (out_def_in_undef)
            //     //     printf("ADD3: %d %d %d = %d %d\n", int_value(s, op1), int_value(s, op2), int_value(s, op3), int_value(s, o1), int_value(s, o2));

            //     // if (int_value(s, op1) == -1 || int_value(s, op2) == -1 || int_value(s, op3) == -1 || int_value(s, o1) == -1 || int_value(s, o2) == -1)
            //     //     printf("%d %d %d = %d %d\n", int_value(s, op1), int_value(s, op2), int_value(s, op3), int_value(s, o1), int_value(s, o2));

            //     // if (int_value(s, op1) != -1 && int_value(s, op2) != -1) {
            //     //     if (j > 0 && int_value(s, op3) != -1) {
            //     //         // comp_3_2(s, out_refined, k, i, j, op1, op2, op3, o1, o2);
            //     //     } else if (j == 0) {
            //     //         // comp_2_2(s, out_refined, k, i, j, op1, op2, o1, o2);
            //     //     }
            //     // }
            // }

            // TODO: Compression 4 to 3

            // Compression: 5 to 3
            // g.cnf.diff_add(DA[i + 4], DT[i], Dsigma0[i], Dr2carry[i], Dr2Carry[i],
            // Df2[i]);
            // {
            //     int op1 = dt;
            //     int op2 = dsigma0;
            //     int op3 = df2;
            //     int op4 = dr2_carry - 1; // t[j - 1]
            //     int op5 = dr2_carry2 - 2; // T[j - 2]
            //     int o1 = dr2_carry2; // T[j]
            //     int o2 = dr2_carry; // t[j]
            //     int o3 = da_4; // DA[i+4]

            //     // if (j > 2) {
            //     //     int var_ids[] = { op1, op2, op3, op4, op5, o1, o2, o3 };
            //     //     infer_carries(s, out_refined, k, var_ids, 8, 2);
            //     // } else if (j == 2 || j == 1) {
            //     //     int var_ids[] = { op1, op2, op3, op4, o1, o2, o3 };
            //     //     infer_carries(s, out_refined, k, var_ids, 7, 2);
            //     // } else {
            //     //     int var_ids[] = { op1, op2, op3, o1, o2 };
            //     //     infer_carries(s, out_refined, k, var_ids, 5, 1);
            //     // }

            //     // if (int_value(s, o1) != -1 && int_value(s, o2) != -1 && int_value(s, o3) != -1 && (int_value(s, op1) == -1 || int_value(s, op2) == -1 || int_value(s, op3) == -1 || int_value(s, op4) == -1 || int_value(s, op5) == -1))
            //     //     printf("ADD5: %d %d %d %d %d = %d %d %d\n", int_value(s, op1), int_value(s, op2), int_value(s, op3), int_value(s, op4), int_value(s, op5), int_value(s, o1), int_value(s, o2), int_value(s, o3));

            //     // if (int_value(s, op1) != -1 && int_value(s, op2) != -1 && int_value(s, op3) != -1) {
            //     //     if (j > 2 && int_value(s, op4) != -1 && int_value(s, op5) != -1) {
            //     //         comp_5_3(s, out_refined, k, i, j, op1, op2, op3, op4, op5, o1, o2,
            //     //             o3);
            //     //     } else if ((j == 2 || j == 1) && int_value(s, op4) != -1) {
            //     //         comp_4_3(s, out_refined, k, i, j, op1, op2, op3, op4, o1, o2, o3);
            //     //     } else if (j == 0) {
            //     //         comp_3_2(s, out_refined, k, i, j, op1, op2, op3, o2, o3);
            //     //     }
            //     // }
            // }

            // Compression: 6 to 3
            // g.cnf.diff_add(DW[i], DW[i - 16], Ds0[i], Dwcarry[i], DwCarry[i], DW[i
            // - 7], Ds1[i]);
            // if (i >= 16) {
            //     int op1 = s.var_map["DW_" + std::to_string(i - 16) + "_g"] + j;
            //     int op2 = ds0;
            //     int op3 = s.var_map["DW_" + std::to_string(i - 7) + "_g"] + j;
            //     int op4 = ds1;
            //     int op5 = dw_carry - 1; // t[j - 1]
            //     int op6 = dw_carry2 - 2; // T[j - 2]
            //     int o1 = dw_carry2; // T[j]
            //     int o2 = dw_carry; // t[j]
            //     int o3 = dw_0; // DW[i]

            //     // if (j > 1) {
            //     //     int var_ids[] = { op1, op2, op3, op4, op5, op6, o1, o2, o3 };
            //     //     infer_carries(s, out_refined, k, var_ids, 9, 2);
            //     // } else if (j == 1) {
            //     //     int var_ids[] = { op1, op2, op3, op4, op5, o1, o2, o3 };
            //     //     infer_carries(s, out_refined, k, var_ids, 8, 2);
            //     // } else {
            //     //     int var_ids[] = { op1, op2, op3, op4, o1, o2, o3 };
            //     //     infer_carries(s, out_refined, k, var_ids, 7, 2);
            //     // }

            //     // if (int_value(s, o1) != -1 && int_value(s, o2) != -1 && int_value(s, o3) != -1 && (int_value(s, op1) == -1 || int_value(s, op2) == -1 || int_value(s, op3) == -1 || int_value(s, op4) == -1 || int_value(s, op5) == -1 || int_value(s, op6) == -1))
            //     //     printf("ADD6: %d %d %d %d %d %d = %d %d %d\n", int_value(s, op1), int_value(s, op2), int_value(s, op3), int_value(s, op4), int_value(s, op5), int_value(s, op6), int_value(s, o1), int_value(s, o2), int_value(s, o3));

            //     // if (s.value(op1) != l_Undef && s.value(op2) != l_Undef && s.value(op3) != l_Undef && s.value(op4) != l_Undef) {
            //     //     if (j > 1 && int_value(s, op5) != -1 && int_value(s, op6) != -1) {
            //     //         comp_6_3(s, out_refined, k, i, j, op1, op2, op3, op4, op5, op6, o1,
            //     //             o2, o3);
            //     //     } else if (j == 1 && int_value(s, op5) != -1) {
            //     //         comp_5_3(s, out_refined, k, i, j, op1, op2, op3, op4, op5, o1, o2,
            //     //             o3);
            //     //     } else if (j == 0) {
            //     //         comp_4_3(s, out_refined, k, i, j, op1, op2, op3, op4, o1, o2, o3);
            //     //     }
            //     // }
            // }

            // Compression: 7 to 3
            // g.cnf.diff_add(DT[i], DE[i], Dsigma1[i], Dr0carry[i], Dr0Carry[i], Df1[i], DK[i], DW[i]);
            {
                int var_ids_f[10], var_ids_g[10]; // 7 inputs + 3 outputs
                add7_var_ids(s, i, j, 'f', var_ids_f);
                add7_var_ids(s, i, j, 'g', var_ids_g);
                if (j > 1) {
                    infer_carries(s, out_refined, k, var_ids_f, 10, 2);
                    // if (out_refined.size() > 0) {
                    //     printf("i,j = %d,%d\n", i, j);
                    //     goto END_CALLBACK;
                    // }
                    infer_carries(s, out_refined, k, var_ids_g, 10, 2);
                    // if (out_refined.size() > 0) {
                    //     printf("i,j = %d,%d\n", i, j);
                    //     goto END_CALLBACK;
                    // }
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