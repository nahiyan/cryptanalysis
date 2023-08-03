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

void add_to_var_ids(Solver& solver, std::string prefix, std::vector<int>& var_ids, int inputs_n, int outputs_n)
{
    for (int i = 0; i < inputs_n; i++) {
        var_ids.push_back(solver.var_map[prefix + "x" + std::to_string(i) + "_f"]);
        var_ids.push_back(solver.var_map[prefix + "x" + std::to_string(i) + "_g"]);

        std::string diff_name = prefix + "x" + std::to_string(i) + "_g";
        diff_name.insert(diff_name.substr(0, 4) == "add_" ? 4 : 0, "D");
        var_ids.push_back(solver.var_map[diff_name]);
    }

    for (int i = 0; i < outputs_n; i++) {
        std::string z_index = std::to_string(i);
        var_ids.push_back(solver.var_map[prefix + "z" + z_index + "_f"]);
        var_ids.push_back(solver.var_map[prefix + "z" + z_index + "_g"]);

        std::string diff_name = prefix + "z" + z_index + "_g";
        diff_name.insert(diff_name.substr(0, 4) == "add_" ? 4 : 0, "D");
        var_ids.push_back(solver.var_map[diff_name]);
    }
}

void processVarMap(Solver& solver)
{
    printf("Var. map entries: %d\n", solver.var_map.size());
    for (int i = 0; i < solver.steps; i++) {
        // if
        add_to_var_ids(solver, "if_" + std::to_string(i) + "_", solver.var_ids_.if_[i], 3, 1);

        // maj
        add_to_var_ids(solver, "maj_" + std::to_string(i) + "_", solver.var_ids_.maj[i], 3, 1);

        // sigma0
        add_to_var_ids(solver, "sigma0_" + std::to_string(i) + "_", solver.var_ids_.sigma0[i], 3, 1);

        // sigma1
        add_to_var_ids(solver, "sigma1_" + std::to_string(i) + "_", solver.var_ids_.sigma1[i], 3, 1);

        if (i >= 16) {
            // s0
            add_to_var_ids(solver, "s0_" + std::to_string(i) + "_", solver.var_ids_.s0[i - 16], 3, 1);

            // s1
            add_to_var_ids(solver, "s1_" + std::to_string(i) + "_", solver.var_ids_.s1[i - 16], 3, 1);

            // add_w
            add_to_var_ids(solver, "add_w" + std::to_string(i) + "_", solver.var_ids_.add_w[i - 16], 6, 2);
        }

        // add_t
        add_to_var_ids(solver, "add_T" + std::to_string(i) + "_", solver.var_ids_.add_t[i], 7, 2);

        // add_e
        add_to_var_ids(solver, "add_E" + std::to_string(i + 4) + "_", solver.var_ids_.add_e[i], 3, 1);

        // add_a
        add_to_var_ids(solver, "add_A" + std::to_string(i + 4) + "_", solver.var_ids_.add_a[i], 5, 2);
    }
}

int int_value(Minisat::Solver& s, int var)
{
    auto value = s.value(var);
    return value == l_True ? 1 : value == l_False ? 0
                                                  : 2;
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
void add_2_bit_clauses(Minisat::Solver& s, vec<vec<Lit>>& out_refined, int& k, int operation_id, int function_id, std::vector<int> var_ids)
{
    // Number of variables
    int vars_n = var_ids.size();
    assert(vars_n % 3 == 0 && vars_n > 0); // Must be in triples and non-empty
    int chunks_n = vars_n / 3;

    // Lay out the rule key's foundation
    int key_size = chunks_n + 2;
    char rule_key[key_size];
    rule_key[0] = operation_id;
    rule_key[key_size - 1] = NULL;

    // Process chunk-wise (each chunk has 3 bits)
    std::set<Lit> base_clause;
    bool hasXOrDash = false;
    for (int i = 0, j = 1; i < vars_n; i += 3, j++) {
        // There are 3 possible ways to derive the GC of the chunk: from x and x_, from dx and x or x_, or from dx alone, else we can't
        int& x_id = var_ids[i];
        lbool x_value = s.value(var_ids[i]);

        int& x_prime_id = var_ids[i + 1];
        lbool x_prime_value = s.value(var_ids[i + 1]);

        int& dx_id = var_ids[i + 2];
        lbool dx_value = s.value(var_ids[i + 2]);

        // TODO: Consider enforcing the relationship instead of just helping the solver propagate
        if (x_value != l_Undef && x_prime_value != l_Undef) {
            rule_key[j] = to_gc(x_value, x_prime_value);
            base_clause.insert(mkLit(x_id, x_value == l_True));
            base_clause.insert(mkLit(x_prime_id, x_prime_value == l_True));
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

            base_clause.insert(mkLit(y_id, y_value == l_True));
            base_clause.insert(mkLit(dx_id, dx_value == l_True));
        } else if (dx_value != l_Undef) {
            rule_key[j] = dx_value == l_True ? 'x' : '-';
            hasXOrDash = true;
            base_clause.insert(mkLit(dx_id, dx_value == l_True));
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

    // Derive the relationships between the x and x_ of the chunks and enforce them through clauses
    std::set<int> visited;
    int rule_i = -1;
    for (int i = 0; i < vars_n - 3; i += 3) {
        int var1_id = var_ids[i];
        for (int j = 0; j < vars_n - 3; j += 3) {
            int var2_id = var_ids[j];
            if (var2_id == var1_id || visited.find(var2_id) != visited.end())
                continue;
            rule_i++;
            if (rule_value[rule_i] == '2')
                continue;

            // Inferred variables should be undefined
            lbool var1_value = s.value(var1_id);
            lbool var2_value = s.value(var2_id);

            // Skip if both the values are defined
            if (var1_value != l_Undef && var2_value != l_Undef)
                continue;

            // Skip if both the values are undefined
            if (var1_value == l_Undef && var2_value == l_Undef)
                continue;

            // DEBUG
            printf("2-bit conditions met (%d, %d): %s %s ", operation_id, function_id, rule_it->first.c_str(), rule_it->second.c_str());

            bool constrainEquality = rule_value[rule_i] == '1';
            for (int count = 0; count < 2; count++) {
                out_refined.push();

                // Determine the signs of var1 and var2
                bool signs[4];
                if (constrainEquality) {
                    if (count == 0) {
                        signs[0] = true;
                        signs[1] = false;
                        signs[2] = false;
                        signs[3] = true;
                    } else {
                        signs[0] = false;
                        signs[1] = true;
                        signs[2] = true;
                        signs[3] = false;
                    }
                } else if (!constrainEquality) {
                    if (count == 0) {
                        signs[0] = true;
                        signs[1] = true;
                        signs[2] = true;
                        signs[3] = true;
                    } else {
                        signs[0] = false;
                        signs[1] = false;
                        signs[2] = false;
                        signs[3] = false;
                    }
                }

                // Push var1 and var2
                if (var1_value == l_Undef) {
                    out_refined[k].push(mkLit(var1_id, signs[0]));
                    out_refined[k].push(mkLit(var2_id, signs[1]));
                } else {
                    out_refined[k].push(mkLit(var2_id, signs[2]));
                    out_refined[k].push(mkLit(var1_id, signs[3]));
                }

                // Push the other lits
                for (auto& lit : base_clause) {
                    if (var(lit) == var1_id || var(lit) == var2_id)
                        continue;
                    out_refined[k].push(lit);
                }
                k++;
            }

            for (int count = 0; count < out_refined[k - 1].size(); count++)
                printf(count == 0 ? "%d " : "%d", int_value(s, var(out_refined[k - 1][count])));
            printf("\n");
            print_clause(out_refined[k - 1]);
            print_clause(out_refined[k - 2]);

            s.stats.two_bit_clauses_n[operation_id - TWO_BIT_CONSTRAINT_IF_ID] += 2;
        }

        visited.insert(var1_id);
    }
}

// Prepare the vector of all the operands and carries of addition (may also remove operands equal to carries from previous columns)
void prepare_add_vec(std::vector<int>& ids, std::vector<int>& f, std::vector<int>& g, int carries_n, int amount, int carry_removal_n = 0)
{
    int inputs_n = ids.size() - carries_n * 3;
    for (int i = 0, j = 0; i < inputs_n; i += 3, j++) {
        if (carry_removal_n > 0 && j == 3 || carry_removal_n == 2 && j == 2)
            continue;

        f.push_back(ids[i] + amount);
        g.push_back(ids[i + 1] + amount);
    }

    for (int i = inputs_n; i < ids.size(); i += 3) {
        f.push_back(ids[i] + amount);
        g.push_back(ids[i + 1] + amount);
    }
}

int r_rotate_id(int id, int amount, int offset)
{
    return (id - amount) + (amount + offset) % 32;
}

std::vector<int> prepare_func_vec(std::vector<int>& ids, int offset, int function_id = -1)
{
    std::vector<int> new_vec;
    if (function_id == 0 || function_id == 1)
        for (int i = 0; i < ids.size(); i++) {
            int r_rotate_amount = 0;
            if (i >= 0 && i <= 2)
                r_rotate_amount = function_id == 0 ? 2 : 6;
            else if (i >= 3 && i <= 5)
                r_rotate_amount = function_id == 0 ? 13 : 11;
            else if (i >= 6 && i <= 8)
                r_rotate_amount = function_id == 0 ? 22 : 25;
            new_vec.push_back(r_rotate_id(ids[i], r_rotate_amount, offset));
        }
    else if (function_id == 2 || function_id == 3)
        for (int i = 0; i < ids.size(); i++) {
            int r_rotate_amount = 0;
            if (i >= 0 && i <= 2)
                r_rotate_amount = function_id == 2 ? 7 : 17;
            else if (i >= 3 && i <= 5)
                r_rotate_amount = function_id == 2 ? 18 : 19;
            new_vec.push_back(r_rotate_id(ids[i], r_rotate_amount, offset));
        }
    else
        for (int i = 0; i < ids.size(); i++)
            new_vec.push_back(ids[i] + offset);

    return new_vec;
}

void infer_carries(Solver& s, vec<vec<Lit>>& out_refined, int& k, std::vector<int>& var_ids, int carries_n, int function_id)
{
    int inputs_n = var_ids.size() - carries_n;
    int input_1s_n = 0, input_1s_ids[inputs_n];
    int input_0s_n = 0, input_0s_ids[inputs_n];
    int input_us_n = 0, input_us_ids[inputs_n];
    for (int i = 0; i < inputs_n; i++) {
        if (s.value(var_ids[i]) == l_True)
            input_1s_ids[input_1s_n++] = var_ids[i];
        else if (s.value(var_ids[i]) == l_False)
            input_0s_ids[input_0s_n++] = var_ids[i];
        else
            input_us_ids[input_us_n++] = var_ids[i];
    }

    if (carries_n == 2) {
        int high_carry_id = var_ids[inputs_n];
        lbool high_carry_value = s.value(high_carry_id);
        bool inferred = false;

        // High carry must be 1 if no. of 1s >= 4
        if (input_1s_n >= 4 && high_carry_value != l_True) {
            out_refined.push();
            out_refined[k].push(mkLit(high_carry_id));
            for (int i = 0; i < input_1s_n; i++)
                out_refined[k].push(~mkLit(input_1s_ids[i]));
            k++;
            inferred = true;
            // High carry must be 0 if no. of 0s >= 4
        } else if (input_0s_n >= 4 && high_carry_value != l_False) {
            out_refined.push();
            out_refined[k].push(~mkLit(high_carry_id));
            for (int i = 0; i < input_0s_n; i++)
                out_refined[k].push(mkLit(input_0s_ids[i]));
            k++;
            inferred = true;
        }

        if (inferred) {
            printf("Inferred high carry (function: %d, inputs %d, carry_id %d)\n", function_id, inputs_n, high_carry_id + 1);
            print_clause(out_refined[k - 1]);
            s.stats.carry_infer_high_clauses_n[function_id]++;
        }
    }

    // TODO: Fix bug with the following clause injection process
    // Low carry must be 1 if no. of 1s >= 6
    int low_carry_id = var_ids[var_ids.size() - 1];
    lbool low_carry_value = s.value(low_carry_id);
    bool inferred = false;
    // if (low_carry_value != l_True && input_1s_n >= 6) {
    //     out_refined.push();
    //     out_refined[k].push(mkLit(low_carry_id));
    //     for (int i = 0; i < input_1s_n; i++)
    //         out_refined[k].push(~mkLit(input_1s_ids[i]));
    //     for (int i = 0; i < input_0s_n; i++)
    //         out_refined[k].push(mkLit(input_0s_ids[i]));
    //     k++;
    //     inferred = true;
    // } else if (low_carry_value != l_False && input_0s_n >= 6) {
    //     out_refined.push();
    //     out_refined[k].push(~mkLit(low_carry_id));
    //     for (int i = 0; i < input_1s_n; i++)
    //         out_refined[k].push(~mkLit(input_1s_ids[i]));
    //     for (int i = 0; i < input_0s_n; i++)
    //         out_refined[k].push(mkLit(input_0s_ids[i]));
    //     k++;
    //     inferred = true;
    // }

    // TODO: Check if the logic is correct
    if (low_carry_value != l_True && ((input_1s_n >= 6) || (input_1s_n >= 2 && input_1s_n + input_us_n < 4))) {
        out_refined.push();
        out_refined[k].push(mkLit(low_carry_id));
        for (int i = 0; i < input_1s_n; i++)
            out_refined[k].push(~mkLit(input_1s_ids[i]));
        for (int i = 0; i < input_0s_n; i++)
            out_refined[k].push(mkLit(input_0s_ids[i]));
        k++;
        inferred = true;
    } else if (low_carry_value != l_False && ((input_0s_n >= 6) || (input_1s_n >= 4 && input_1s_n + input_us_n < 6))) {
        out_refined.push();
        out_refined[k].push(~mkLit(low_carry_id));
        for (int i = 0; i < input_1s_n; i++)
            out_refined[k].push(~mkLit(input_1s_ids[i]));
        for (int i = 0; i < input_0s_n; i++)
            out_refined[k].push(mkLit(input_0s_ids[i]));
        k++;
        inferred = true;
    }

    if (inferred) {
        printf("Inferred low carry (function: %d, inputs %d, carry_id %d)\n", function_id, inputs_n, low_carry_id + 1);
        print_clause(out_refined[k - 1]);
        s.stats.carry_infer_low_clauses_n[function_id]++;
    }

    // TODO: Implement r1 == 0 ==> sum(inp[]) <= 3
    // TODO: Implement r1 == 1 ==> sum(inp[]) >= 4
}

void add_addition_clauses(Solver& s, vec<vec<Lit>>& out_refined, int& k, int i, int j, std::vector<int>& ids, int carries_n, int function_id)
{
    std::vector<int> ids_f, ids_g;
    if (j > 1)
        prepare_add_vec(ids, ids_f, ids_g, carries_n, j, 0);
    else if (j == 1)
        prepare_add_vec(ids, ids_f, ids_g, carries_n, j, 1);
    else
        prepare_add_vec(ids, ids_f, ids_g, carries_n, j, 2);

    infer_carries(s, out_refined, k, ids_f, carries_n, function_id);
    infer_carries(s, out_refined, k, ids_g, carries_n, function_id);
}

void add_clauses(Minisat::Solver& s, vec<vec<Lit>>& out_refined)
{
    int k = 0;
    for (int i = 0; i < s.steps; i++) {
        for (int j = 0; j < 32; j++) {
            // If
            {
                std::vector<int> ids = prepare_func_vec(s.var_ids_.if_[i], j);
                add_2_bit_clauses(s, out_refined, k, TWO_BIT_CONSTRAINT_IF_ID, 0, ids);
            }

            // Maj
            {
                std::vector<int> ids = prepare_func_vec(s.var_ids_.maj[i], j);
                add_2_bit_clauses(s, out_refined, k, TWO_BIT_CONSTRAINT_MAJ_ID, 1, ids);
            }

            // sigma0
            {
                std::vector<int> ids = prepare_func_vec(s.var_ids_.sigma0[i], j, 0);
                add_2_bit_clauses(s, out_refined, k, TWO_BIT_CONSTRAINT_XOR3_ID, 2, ids);
            }

            // sigma1
            {
                std::vector<int> ids = prepare_func_vec(s.var_ids_.sigma1[i], j, 1);
                add_2_bit_clauses(s, out_refined, k, TWO_BIT_CONSTRAINT_XOR3_ID, 3, ids);
            }

            if (i >= 16) {
                // s0
                if (j <= 28) {
                    std::vector<int> ids = prepare_func_vec(s.var_ids_.s0[i - 16], j, 2);
                    add_2_bit_clauses(s, out_refined, k, TWO_BIT_CONSTRAINT_XOR3_ID, 4, ids);
                }

                // s1
                if (j <= 21) {
                    std::vector<int> ids = prepare_func_vec(s.var_ids_.s1[i - 16], j, 3);
                    add_2_bit_clauses(s, out_refined, k, TWO_BIT_CONSTRAINT_XOR3_ID, 5, ids);
                }
            }

            // Add E
            add_addition_clauses(s, out_refined, k, i, j, s.var_ids_.add_e[i], 1, 0);

            // Add A
            add_addition_clauses(s, out_refined, k, i, j, s.var_ids_.add_a[i], 2, 1);

            // Add W
            if (i >= 16) {
                add_addition_clauses(s, out_refined, k, i, j, s.var_ids_.add_w[i - 16], 2, 2);
            }

            // Add T
            add_addition_clauses(s, out_refined, k, i, j, s.var_ids_.add_t[i], 2, 3);
        }
    }
}