#include "Crypto.h"
#include "NTL/GF2.h"
#include "NTL/mat_GF2.h"
#include "NTL/vec_GF2.h"
#include <algorithm>
#include <chrono>
#include <map>
#include <memory>
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

namespace Crypto {
void load_rule(Solver& solver, FILE*& db, int& id)
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
}

void load_rules(Solver& solver, const char* filename)
{
    FILE* db = fopen(filename, "r");
    char buffer[1];
    int count = 0;
    while (1) {
        int n = fread(buffer, 1, 1, db);
        if (n == 0)
            break;

        int id = buffer[0];
        load_rule(solver, db, id);
        count++;
    }

    printf("Loaded %d rules into %d buckets\n", count, solver.rules.bucket_count());
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

void process_var_map(Solver& solver)
{
    printf("Steps: %d\n", solver.steps);
    printf("Var. map entries: %d\n", solver.var_map.size());
    for (int i = 0; i < solver.steps; i++) {
        // if
        add_to_var_ids(solver, "if_" + std::to_string(i) + "_", solver.var_ids.if_[i], 3, 1);

        // maj
        add_to_var_ids(solver, "maj_" + std::to_string(i) + "_", solver.var_ids.maj[i], 3, 1);

        // sigma0
        add_to_var_ids(solver, "sigma0_" + std::to_string(i) + "_", solver.var_ids.sigma0[i], 3, 1);

        // sigma1
        add_to_var_ids(solver, "sigma1_" + std::to_string(i) + "_", solver.var_ids.sigma1[i], 3, 1);

        if (i >= 16) {
            // s0
            add_to_var_ids(solver, "s0_" + std::to_string(i) + "_", solver.var_ids.s0[i - 16], 3, 1);

            // s1
            add_to_var_ids(solver, "s1_" + std::to_string(i) + "_", solver.var_ids.s1[i - 16], 3, 1);

            // add_w
            add_to_var_ids(solver, "add_w" + std::to_string(i) + "_", solver.var_ids.add_w[i - 16], 6, 3);
        }

        // add_t
        add_to_var_ids(solver, "add_T" + std::to_string(i) + "_", solver.var_ids.add_t[i], 7, 3);

        // add_e
        add_to_var_ids(solver, "add_E" + std::to_string(i + 4) + "_", solver.var_ids.add_e[i], 3, 2);

        // add_a
        add_to_var_ids(solver, "add_A" + std::to_string(i + 4) + "_", solver.var_ids.add_a[i], 5, 3);
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
        return 0;
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
        return 0;
}

int sum(std::vector<int>& v)
{
    int s = 0;
    for (int& v_ : v)
        s += v_;
    return s;
}

int sum(int* v, int n)
{
    int s = 0;
    for (int i = 0; i < n; i++)
        s += v[i];
    return s;
}

NTL::GF2 sum(NTL::vec_GF2& v)
{
    NTL::GF2 sum = NTL::to_GF2(0);
    for (int i = 0; i < v.length(); i++)
        sum += v[i];

    return sum;
}

int sum_dec_from_bin(NTL::vec_GF2& v)
{
    int sum = 0;
    for (int i = 0; i < v.length(); i++)
        sum += NTL::conv<int>(v[i]);

    return sum;
}

void print(std::vector<int>& vector, int offset = 0)
{
    printf("Vector: ");
    for (int i = 0; i < vector.size(); i++) {
        printf("%d", vector[i] + offset);
        if (i != vector.size() - 1)
            printf(" ");
    }
    printf("\n");
}

void print(minisat_clause_t& clause)
{
    printf("Clause: ");
    for (int i = 0; i < clause.size(); i++) {
        printf("%s%d", sign(clause[i]) ? "-" : "", var(clause[i]) + 1);
        if (i != clause.size() - 1)
            printf(" ");
    }
    printf("\n");
}

void print(equation_t equation)
{
    int x = std::get<0>(equation);
    int y = std::get<1>(equation);
    int z = std::get<2>(equation);
    printf("Equation: %d %s %d\n", x, z == 1 ? "=/=" : "=", y);
}

void print(equations_t equations)
{
    for (auto& equation : equations)
        print(equation);
}

void print(NTL::vec_GF2& equation)
{
    printf("Vector [GF(2)]; size: %d: ", equation.length());
    for (int i = 0; i < equation.length(); i++)
        printf("%d ", NTL::conv<int>(equation[i]));
    printf("\n");
}

// Get index of the shortest conflict clause, -1 if no conflict clause is found
int get_shortest_conflict_clause(State& state)
{
    int shortest_index = -1, shortest_length = INT_MAX;
    for (int i = 0; i < state.out_refined.size(); i++) {
        // Skip propagation clauses
        if (state.solver.value(state.out_refined[i][0]) == l_Undef)
            continue;

        int size = state.out_refined[i].size();
        if (size >= shortest_length)
            continue;

        shortest_length = size;
        shortest_index = i;
    }

    return shortest_index;
}

// Checks GF(2) equations and returns conflicting equations (equations that conflicts with previously added ones)
std::shared_ptr<equations_t> check_consistency(std::shared_ptr<equations_t>& equations, bool exhaustive)
{
    auto conflicting_equations = std::make_shared<equations_t>();
    std::map<uint32_t, std::shared_ptr<std::set<int32_t>>> rels;

    for (auto& equation : *equations) {
        auto var1 = std::get<0>(equation) + 1;
        auto var2 = (std::get<2>(equation) == 1 ? -1 : 1) * (std::get<1>(equation) + 1);
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
                    auto confl_eq = std::make_tuple(var1_abs - 1, var2_abs - 1, var2 < 0 ? 1 : 0);
                    conflicting_equations->push_back(confl_eq);
                    if (!exhaustive)
                        return conflicting_equations;
                }
            }

            // Update existing references
            for (auto& item : *updated_set) {
                auto& set = rels[abs(item)];
                if (set == updated_set)
                    continue;
                rels[abs(item)] = updated_set;
            }
        } else if (var1_exists || var2_exists) {
            // Find an existing set related to any of the variables
            auto& existing_set = var1_exists ? rels[var1_abs] : rels[var2_abs];
            auto var1_inv_in_existing_set = existing_set->find(-var1) == existing_set->end() ? false : true;
            auto var2_inv_in_existing_set = existing_set->find(-var2) == existing_set->end() ? false : true;

            // Invert the lone variable to try to prevent a conflict
            if (var1_inv_in_existing_set)
                var2 *= -1;
            else if (var2_inv_in_existing_set)
                var1 *= -1;

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
                rels[abs(item)] = existing_set;
            }
        } else {
            // Adding novel variables
            auto new_set = std::make_shared<std::set<int32_t>>(std::set<int32_t> { var1, var2 });
            rels[var1_abs] = new_set;
            rels[var2_abs] = new_set;
        }
    }

    return conflicting_equations;
}

// The variable IDs provided should include the operands and the output
void add_2_bit_equations(State& state, int operation_id, int function_id, std::vector<int> var_ids)
{
    // Number of variables
    int vars_n = var_ids.size();
    assert(vars_n % 3 == 0 && vars_n > 0); // Must be in triples and non-empty
    int chunks_n = vars_n / 3;

    // Lay out the rule key's foundation
    int key_size = chunks_n + 2;
    char rule_key[key_size];
    rule_key[0] = operation_id;
    rule_key[key_size - 1] = 0;

    // Process chunk-wise (each chunk has 3 bits)
    std::vector<Lit> base_clause;
    bool hasXOrDash = false;
    for (int i = 0, j = 1; i < vars_n; i += 3, j++) {
        // There are 3 possible ways to derive the GC of the chunk: from x and x_, from dx and x or x_, or from dx alone, else we can't
        int& x_id = var_ids[i];
        lbool x_value = state.solver.value(var_ids[i]);

        int& x_prime_id = var_ids[i + 1];
        lbool x_prime_value = state.solver.value(var_ids[i + 1]);

        int& dx_id = var_ids[i + 2];
        lbool dx_value = state.solver.value(var_ids[i + 2]);

        if (x_value != l_Undef && x_prime_value != l_Undef) {
            rule_key[j] = to_gc(x_value, x_prime_value);
            base_clause.push_back(mkLit(x_id, x_value == l_True));
            base_clause.push_back(mkLit(x_prime_id, x_prime_value == l_True));
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

            base_clause.push_back(mkLit(y_id, y_value == l_True));
            base_clause.push_back(mkLit(dx_id, dx_value == l_True));
        } else if (dx_value != l_Undef) {
            rule_key[j] = dx_value == l_True ? 'x' : '-';
            hasXOrDash = true;
            base_clause.push_back(mkLit(dx_id, dx_value == l_True));
        } else {
            // Terminate since we can't derive the rule if we don't know any of {1, u, n, 0, x, -}, and without the rule we can't derive the 2-bit conditions
            return;
        }
    }

    // Find the value of the rule (if it exists)
    auto rule_it = state.solver.rules.find(rule_key);
    if (rule_it == state.solver.rules.end())
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
            bool are_equal = rule_value[rule_i] == '1';

            // Add the equation
            auto equation = equation_t { var1_id, var2_id, are_equal ? 0 : 1 };
            state.equations->push_back(equation);

            // Map the equation variables (if they don't exist)
            if (state.eq_var_map.find(var1_id) == state.eq_var_map.end())
                state.eq_var_map[var1_id] = state.eq_var_map.size();
            if (state.eq_var_map.find(var2_id) == state.eq_var_map.end())
                state.eq_var_map[var2_id] = state.eq_var_map.size();

            // Connect the equation with this function result
            std::vector<int> variables;
            for (Lit& lit : base_clause)
                variables.push_back(var(lit));
            auto func_result = FunctionResult {
                operation_id,
                function_id,
                variables,
            };
            auto eq_func_relation = state.eq_func_rels.find(equation);
            auto it = state.eq_func_rels.find(equation);
            if (it == state.eq_func_rels.end())
                state.eq_func_rels.insert({ equation, { func_result } });
            else
                state.eq_func_rels[equation].push_back(func_result);
        }

        visited.insert(var1_id);
    }
}

// Prepare the vector of all the operands and carries of addition (may also remove operands equal to carries from previous columns)
void prepare_add_vec(std::vector<int>& ids, std::vector<int>& f, std::vector<int>& g, int carries_n, int offset, int carry_removal_n = 0)
{
    int inputs_n = ids.size() - (carries_n + 1) * 3;
    for (int i = 0, j = 0; i < inputs_n; i += 3, j++) {
        if (carry_removal_n > 0 && j == 3 || carry_removal_n == 2 && j == 2)
            continue;

        f.push_back(ids[i] + offset);
        g.push_back(ids[i + 1] + offset);
    }

    // Add the carries but ignore the sum
    for (int i = inputs_n; i < ids.size() - 3; i += 3) {
        f.push_back(ids[i] + offset);
        g.push_back(ids[i + 1] + offset);
    }
}

std::vector<int> prepare_add_2_bit_vec(std::vector<int>& ids, int carries_n, int offset, int carry_removal_n = 0)
{
    std::vector<int> new_vec;
    int inputs_n = ids.size() - (carries_n + 1) * 3;
    for (int i = 0, j = 0; i < inputs_n; i += 3, j++) {
        if (carry_removal_n > 0 && j == 3 || carry_removal_n == 2 && j == 2)
            continue;

        new_vec.push_back(ids[i] + offset);
        new_vec.push_back(ids[i + 1] + offset);
        new_vec.push_back(ids[i + 2] + offset);
    }

    // Add the sum but ignore the carries
    for (int i = inputs_n + (3 * carries_n); i < ids.size(); i += 3) {
        new_vec.push_back(ids[i] + offset);
        new_vec.push_back(ids[i + 1] + offset);
        new_vec.push_back(ids[i + 2] + offset);
    }

    return new_vec;
}

int r_rotate_id(int id, int amount, int offset)
{
    return (id - amount) + (amount + offset) % 32;
}

std::vector<int> prepare_func_vec(std::vector<int>& ids, int offset, int function_id = -1)
{
    std::vector<int> new_vec;
    if (function_id < 0 || function_id > 3)
        for (int i = 0; i < ids.size(); i++)
            new_vec.push_back(ids[i] + offset);
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

    return new_vec;
}

// Create the augmented matrix from equations
void make_aug_matrix(State& state, NTL::mat_GF2& coeff_matrix, NTL::vec_GF2& rhs)
{
    auto variables_n = state.eq_var_map.size();
    auto equations_n = state.equations->size();
    coeff_matrix.SetDims(equations_n, variables_n);
    rhs.SetLength(equations_n);

    // Construct the coefficient matrix
    auto& equations_deref = *state.equations;
    for (int eq_index = 0; eq_index < equations_n; eq_index++) {
        auto& equation = equations_deref[eq_index];
        int& x = state.eq_var_map[std::get<0>(equation)];
        int& y = state.eq_var_map[std::get<1>(equation)];
        for (int col_index = 0; col_index < variables_n; col_index++)
            coeff_matrix[eq_index][col_index] = NTL::to_GF2(col_index == x || col_index == y ? 1 : 0);

        rhs.put(eq_index, std::get<2>(equation));
    }
}

// Detect inconsistencies from nullspace vectors
int find_inconsistency_from_vectors(State& state, NTL::mat_GF2& coeff_matrix, NTL::vec_GF2& rhs, NTL::mat_GF2& nullspace_vectors, NTL::vec_GF2*& inconsistency)
{
    int coeff_n = coeff_matrix.NumCols();
    int inconsistent_eq_n = 0;
    int least_hamming_weight = INT_MAX;
    int nullspace_vectors_n = nullspace_vectors.NumRows();
    int equations_n = state.equations->size();
    for (int index = 0; index < nullspace_vectors_n; index++) {
        auto& nullspace_vector = nullspace_vectors[index];

        // Initialize the values to 0
        NTL::GF2 rhs_sum = NTL::to_GF2(0);
        NTL::vec_GF2 coeff_sums;
        coeff_sums.SetLength(coeff_n);
        for (int x = 0; x < coeff_n; x++)
            coeff_sums[x] = 0;

        // Go through the nullspace vector and add the equations and RHS
        for (int eq_index = 0; eq_index < equations_n; eq_index++) {
            if (nullspace_vector[eq_index] == 0)
                continue;

            // Add the coefficients of the equations
            coeff_sums += coeff_matrix[eq_index];

            // Add the RHS
            rhs_sum += rhs[eq_index];
        }

        // Mismatching RHS sum and coefficients sum is a contradiction
        if (rhs_sum != sum(coeff_sums)) {
            int hamming_weight = 0;
            for (int x = 0; x < equations_n; x++)
                hamming_weight += NTL::conv<int>(nullspace_vector[x]);

            if (hamming_weight < least_hamming_weight) {
                inconsistency = &nullspace_vector;
            }
            inconsistent_eq_n++;
        }
    }

    return inconsistent_eq_n;
}

// Use NTL to find cycles of inconsistent equations
bool block_inconsistency(State& state)
{
    // Make the augmented matrix
    clock_t start_time = std::clock();
    NTL::mat_GF2 coeff_matrix;
    NTL::vec_GF2 rhs;
    make_aug_matrix(state, coeff_matrix, rhs);
    state.solver.stats.two_bit_cpu_time_segments[1] += std::clock() - start_time;

    // Find the basis of the coefficient matrix's left kernel
    start_time = std::clock();
    NTL::mat_GF2 left_kernel_basis;
    NTL::kernel(left_kernel_basis, coeff_matrix);
    auto nullspace_vectors_n = left_kernel_basis.NumRows();
    auto equations_n = left_kernel_basis.NumCols();

    printf("Basis elements: %d, %d\n", nullspace_vectors_n, equations_n);
    state.solver.stats.two_bit_cpu_time_segments[2] += std::clock() - start_time;

    start_time = std::clock();
    // TODO: Add combinations of the basis vectors
    state.solver.stats.two_bit_cpu_time_segments[3] += std::clock() - start_time;

    // Check for inconsistencies
    start_time = std::clock();
    NTL::vec_GF2* inconsistency = NULL;
    auto inconsistent_eq_n = find_inconsistency_from_vectors(state, coeff_matrix, rhs, left_kernel_basis, inconsistency);
    state.solver.stats.two_bit_cpu_time_segments[4] += std::clock() - start_time;

    // Blocking inconsistencies
    if (inconsistency != NULL) {
        start_time = std::clock();
        auto& inconsistency_deref = *inconsistency;
        printf("Found inconsistencies (%d): %d equations\n", inconsistent_eq_n, sum_dec_from_bin(inconsistency_deref));
        print(inconsistency_deref);

        state.out_refined.push();
        std::set<Lit> confl_clause_lits;
        for (int eq_index = 0; eq_index < equations_n; eq_index++) {
            if (inconsistency_deref[eq_index] == 0)
                continue;

            auto& equation = (*state.equations)[eq_index];
            auto results_it = state.eq_func_rels.find(equation);
            assert(results_it != state.eq_func_rels.end());

            // Instances refer to the function instances
            auto& instances = results_it->second;
            print(equation);

            printf("Number of functions for the equation (ID: %d): %d\n", eq_index, instances.size());
            for (auto& instance : instances) {
                printf("Adding to confl. clause: op. %d func. %d\n", instance.operation_id, instance.functon_id);
                print(instance.variables, 1);
                for (auto& var : instance.variables) {
                    auto value = state.solver.value(var);
                    assert(value != l_Undef);
                    // state.out_refined[state.k].push(mkLit(var, value == l_True));
                    confl_clause_lits.insert(mkLit(var, value == l_True));
                }
                state.solver.stats.two_bit_clauses_n[instance.operation_id - TWO_BIT_CONSTRAINT_IF_ID]++;
            }
        }
        for(auto& lit: confl_clause_lits) {
            state.out_refined[state.k].push(lit);
        }
        state.k++;
        print(state.out_refined[state.k - 1]);
        state.solver.stats.two_bit_cpu_time_segments[5] += std::clock() - start_time;

        // Terminate since we already already a few conflict clauses
        return true;
    }

    return false;
}

void infer_carries(State& state, std::vector<int>& var_ids, int carries_n, int function_id)
{
    int inputs_n = var_ids.size() - carries_n;
    int input_1s_n = 0, input_1s_ids[inputs_n];
    int input_0s_n = 0, input_0s_ids[inputs_n];
    int input_us_n = 0, input_us_ids[inputs_n];
    for (int i = 0; i < inputs_n; i++) {
        if (state.solver.value(var_ids[i]) == l_True)
            input_1s_ids[input_1s_n++] = var_ids[i];
        else if (state.solver.value(var_ids[i]) == l_False)
            input_0s_ids[input_0s_n++] = var_ids[i];
        else
            input_us_ids[input_us_n++] = var_ids[i];
    }

    if (carries_n == 2) {
        int high_carry_id = var_ids[inputs_n];
        lbool high_carry_value = state.solver.value(high_carry_id);
        bool inferred = false;

        // High carry must be 1 if no. of 1s >= 4
        if (input_1s_n >= 4 && high_carry_value != l_True) {
            state.out_refined.push();
            state.out_refined[state.k].push(mkLit(high_carry_id));
            for (int i = 0; i < input_1s_n; i++)
                state.out_refined[state.k].push(~mkLit(input_1s_ids[i]));
            state.k++;
            inferred = true;
            // High carry must be 0 if no. of 0s >= 4
        } else if (input_0s_n >= 4 && high_carry_value != l_False) {
            state.out_refined.push();
            state.out_refined[state.k].push(~mkLit(high_carry_id));
            for (int i = 0; i < input_0s_n; i++)
                state.out_refined[state.k].push(mkLit(input_0s_ids[i]));
            state.k++;
            inferred = true;
        }

        if (inferred) {
            printf("Inferred high carry (function: %d, inputs %d, carry_id %d)\n", function_id, inputs_n, high_carry_id + 1);
            print(state.out_refined[state.k - 1]);
            state.solver.stats.carry_infer_high_clauses_n[function_id]++;
        }
    }

    // Low carry must be 1 if no. of 1s >= 6
    int low_carry_id = var_ids[var_ids.size() - 1];
    lbool low_carry_value = state.solver.value(low_carry_id);
    bool inferred = false;
    // TODO: Check if the logic is correct
    if (low_carry_value != l_True && ((input_1s_n >= 6) || (input_1s_n >= 2 && input_1s_n + input_us_n < 4))) {
        state.out_refined.push();
        state.out_refined[state.k].push(mkLit(low_carry_id));
        for (int i = 0; i < input_1s_n; i++)
            state.out_refined[state.k].push(~mkLit(input_1s_ids[i]));
        for (int i = 0; i < input_0s_n; i++)
            state.out_refined[state.k].push(mkLit(input_0s_ids[i]));
        state.k++;
        inferred = true;
    } else if (low_carry_value != l_False && ((input_0s_n >= 6) || (input_1s_n >= 4 && input_1s_n + input_us_n < 6))) {
        state.out_refined.push();
        state.out_refined[state.k].push(~mkLit(low_carry_id));
        for (int i = 0; i < input_1s_n; i++)
            state.out_refined[state.k].push(~mkLit(input_1s_ids[i]));
        for (int i = 0; i < input_0s_n; i++)
            state.out_refined[state.k].push(mkLit(input_0s_ids[i]));
        state.k++;
        inferred = true;
    }

    if (inferred) {
        printf("Inferred low carry (function: %d, inputs %d, carry_id %d)\n", function_id, inputs_n, low_carry_id + 1);
        print(state.out_refined[state.k - 1]);
        state.solver.stats.carry_infer_low_clauses_n[function_id]++;
    }

    // TODO: Implement r1 == 0 ==> sum(inp[]) <= 3
    // TODO: Implement r1 == 1 ==> sum(inp[]) >= 4
}

void add_addition_clauses(State& state, int i, int j, std::vector<int>& ids, int carries_n, int function_id)
{
    std::vector<int> ids_f, ids_g;
    if (j > 1)
        prepare_add_vec(ids, ids_f, ids_g, carries_n, j, 0);
    else if (j == 1)
        prepare_add_vec(ids, ids_f, ids_g, carries_n, j, 1);
    else
        prepare_add_vec(ids, ids_f, ids_g, carries_n, j, 2);

    infer_carries(state, ids_f, carries_n, function_id);
    infer_carries(state, ids_g, carries_n, function_id);
}

void add_addition_2_bit_clauses(State& state, int i, int j, std::vector<int>& ids, int carries_n, int function_id)
{
    std::vector<int> new_vec;
    if (j > 1)
        new_vec = prepare_add_2_bit_vec(ids, carries_n, j, 0);
    else if (j == 1)
        new_vec = prepare_add_2_bit_vec(ids, carries_n, j, 1);
    else
        new_vec = prepare_add_2_bit_vec(ids, carries_n, j, 2);

    int operation_id = 0;
    if ((j > 1 && function_id == 1) || (j == 0 && function_id == 3) || (j == 1 && function_id == 2))
        operation_id = TWO_BIT_CONSTRAINT_ADD5_ID;
    else if ((j == 1 && function_id == 1) || (j == 0 && function_id == 2))
        operation_id = TWO_BIT_CONSTRAINT_ADD4_ID;
    else if ((j == 0 && function_id == 1) || (j > 1 && function_id == 0))
        operation_id = TWO_BIT_CONSTRAINT_ADD3_ID;
    else if ((j == 1 && function_id == 3) || (j > 1 && function_id == 2))
        operation_id = TWO_BIT_CONSTRAINT_ADD6_ID;
    else if (j > 1 && function_id == 3)
        operation_id = TWO_BIT_CONSTRAINT_ADD7_ID;

    // TODO: Add two-bit conditions for ADD2
    // printf("Debug imp: ");
    // for (auto& x: ids_f) {
    //     printf("%d ", x);
    // }
    // printf("\n");

    add_2_bit_equations(state, operation_id, function_id, new_vec);
}

void add_clauses(State& state)
{
    clock_t two_bit_start_time = std::clock();
    // Handle 2-bit conditions
    for (int i = 0; i < state.solver.steps; i++) {
        for (int j = 0; j < 32; j++) {
            // If
            {
                std::vector<int> ids = prepare_func_vec(state.solver.var_ids.if_[i], j);
                add_2_bit_equations(state, TWO_BIT_CONSTRAINT_IF_ID, 0, ids);
            }

            // Maj
            {
                std::vector<int> ids = prepare_func_vec(state.solver.var_ids.maj[i], j);
                add_2_bit_equations(state, TWO_BIT_CONSTRAINT_MAJ_ID, 1, ids);
            }

            // Sigma0
            {
                std::vector<int> ids = prepare_func_vec(state.solver.var_ids.sigma0[i], j, 0);
                add_2_bit_equations(state, TWO_BIT_CONSTRAINT_XOR3_ID, 2, ids);
            }

            // Sigma1
            {
                std::vector<int> ids = prepare_func_vec(state.solver.var_ids.sigma1[i], j, 1);
                add_2_bit_equations(state, TWO_BIT_CONSTRAINT_XOR3_ID, 3, ids);
            }

            if (i >= 16) {
                // S0
                if (j <= 28) {
                    std::vector<int> ids = prepare_func_vec(state.solver.var_ids.s0[i - 16], j, 2);
                    add_2_bit_equations(state, TWO_BIT_CONSTRAINT_XOR3_ID, 4, ids);
                } else {
                    // TODO: Implement XOR2 2-bit conditions
                }

                // S1
                if (j <= 21) {
                    std::vector<int> ids = prepare_func_vec(state.solver.var_ids.s1[i - 16], j, 3);
                    add_2_bit_equations(state, TWO_BIT_CONSTRAINT_XOR3_ID, 5, ids);
                } else {
                    // TODO: Implement XOR2 2-bit conditions
                }

                // Add W
                add_addition_2_bit_clauses(state, i, j, state.solver.var_ids.add_w[i - 16], 2, 2);
            }

            // Add E
            add_addition_2_bit_clauses(state, i, j, state.solver.var_ids.add_e[i], 1, 0);

            // Add A
            add_addition_2_bit_clauses(state, i, j, state.solver.var_ids.add_a[i], 2, 1);

            // Add T
            add_addition_2_bit_clauses(state, i, j, state.solver.var_ids.add_t[i], 2, 3);
        }
    }
    state.solver.stats.two_bit_cpu_time_segments[0] += std::clock() - two_bit_start_time;

    bool is_inconsistent = false;
    {
        auto start_time = std::clock();
        auto confl_equations = check_consistency(state.equations, false);
        is_inconsistent = confl_equations->size() > 0;
        state.solver.stats.incons_set_based_cpu_time += std::clock() - start_time;
    }
    state.solver.stats.two_bit_cpu_time += std::clock() - two_bit_start_time;

    // Block inconsistencies
    if (is_inconsistent) {
        auto start_time = std::clock();
        bool blocked = block_inconsistency(state);
        state.solver.stats.two_bit_cpu_time += std::clock() - two_bit_start_time;
        if (blocked)
            return;
    }

    // Handle carry inference
    clock_t carry_inference_start_time = std::clock();
    for (int i = 0; i < state.solver.steps; i++) {
        for (int j = 0; j < 32; j++) {
            // Add E
            add_addition_clauses(state, i, j, state.solver.var_ids.add_e[i], 1, 0);

            // Add A
            add_addition_clauses(state, i, j, state.solver.var_ids.add_a[i], 2, 1);

            // Add W
            if (i >= 16)
                add_addition_clauses(state, i, j, state.solver.var_ids.add_w[i - 16], 2, 2);

            // Add T
            add_addition_clauses(state, i, j, state.solver.var_ids.add_t[i], 2, 3);
        }
    }
    state.solver.stats.carry_inference_cpu_time += std::clock() - carry_inference_start_time;

    // TODO: Don't add propagation clauses when a conflict clause is detected
    // Post processing
    int index = get_shortest_conflict_clause(state);
    if (index != -1) {
        minisat_clause_t conflict_clause;
        state.out_refined[index].copyTo(conflict_clause);
        state.out_refined.clear();

        printf("Note: Adding only the shortest conflict clause of size %d\n", conflict_clause.size());
        state.out_refined.push();
        for (int count = 0; count < conflict_clause.size(); count++) {
            state.out_refined[0].push(mkLit(var(conflict_clause[count]), state.solver.value((conflict_clause[count])) == l_False));
        }

        print(conflict_clause);
    }
}
}