import os
from itertools import product
from rules_gen import if_, maj, xor3, bin_add
from collections.abc import Iterable

TWO_BIT_CONSTRAINT_XOR2_ID = 0
TWO_BIT_CONSTRAINT_IF_ID = 1
TWO_BIT_CONSTRAINT_MAJ_ID = 2
TWO_BIT_CONSTRAINT_XOR3_ID = 3
TWO_BIT_CONSTRAINT_ADD2_ID = 4
TWO_BIT_CONSTRAINT_ADD3_ID = 5
TWO_BIT_CONSTRAINT_ADD4_ID = 6
TWO_BIT_CONSTRAINT_ADD5_ID = 7
TWO_BIT_CONSTRAINT_ADD6_ID = 8
TWO_BIT_CONSTRAINT_ADD7_ID = 9

if not os.path.exists("output"):
    os.mkdir("output")
rules_db = open("output/rules-2-bit.db", "wb")


def create_matrix(n):
    matrix = []
    for i in range(n):
        matrix.append([])
        for _ in range(n):
            matrix[i].append("x")
    return matrix


def print_matrix(matrix):
    r = len(matrix)
    c = len(matrix[0])
    print(end="  ")
    for i in range(c):
        print(f"{i}", end=" ")
    print()
    for i in range(r):
        print(f"{i}", end=" ")
        for j in range(c):
            print(matrix[i][j], end=" ")
        print()


def debug_rels(rels_string, n):
    elements = list(range(n))
    matrix = create_matrix(n)
    k = 0
    for i in elements:
        selector = i + 1
        for j in range(selector, n):
            matrix[i][j] = rels_string[k]
            k += 1
    print_matrix(matrix)


def if_w(ops):
    return if_(ops[0], ops[1], ops[2])


def maj_w(ops):
    return maj(ops[0], ops[1], ops[2])


def xor3_w(ops):
    return xor3(ops[0], ops[1], ops[2])


def xor2(ops):
    return ops[0] ^ ops[1]


def bin_add_w(ops):
    carry2, carry1, sum = bin_add(ops)
    if len(ops) >= 4:
        return carry2, carry1, sum
    else:
        return carry1, sum


def rels(rule_key, inputs):
    rels_desc = ""  # relationships descriptor
    n = len(inputs)
    elements = list(range(0, n))
    for i in elements:
        selector = i + 1
        for j in range(selector, n):
            if rule_key[i] in ["-", "x"] and rule_key[j] in ["-", "x"]:
                rels_desc += "1" if inputs[i] == inputs[j] else "0"
            else:
                rels_desc += "2"
    return rels_desc


def conforms_to(vars_f, vars_g, vars_gc):
    if not isinstance(vars_f, Iterable):
        vars_f = [vars_f]
        vars_g = [vars_g]
        vars_gc = vars_gc
    for i in range(len(vars_f)):
        var_f = vars_f[i]
        var_g = vars_g[i]
        var_gc = vars_gc[i]
        if var_gc == "x":
            if var_f == var_g:
                return False
        elif var_gc == "-":
            if var_f != var_g:
                return False
        elif var_gc == "0":
            if var_f != 0 or var_g != 0:
                return False
        elif var_gc == "u":
            if var_f != 1 or var_g != 0:
                return False
        elif var_gc == "n":
            if var_f != 0 or var_g != 1:
                return False
        elif var_gc == "1":
            if var_f != 1 or var_g != 1:
                return False
    return True


def get_rels_colwise(rels):
    rels_colwise = {}
    for rel in rels:
        for i in range(len(rel)):
            if i not in rels_colwise:
                rels_colwise[i] = []
            rels_colwise[i].append(int(rel[i]))
    return rels_colwise


def get_rels_consistency(rels_colwise):
    consistency = []
    for i in range(len(rels_colwise)):
        items = set(rels_colwise[i])
        if len(items) == 1:
            consistency.append(items.pop())
        else:
            consistency.append(2)
    return consistency


# Assumes that there is one output of the function
def gen_2_bit_conds(id, func, inputs_n, outputs_n=1):
    gc_set = ["1", "u", "n", "0", "x", "-"]

    # Generate all the possible rule candidates
    rule_candidates = product(gc_set, repeat=inputs_n)

    # # Try all the candidates
    for rule_candidate in rule_candidates:
        candidates = []
        for i in range(outputs_n):
            entry = ["?"] * outputs_n
            for gc in gc_set:
                entry[i] = gc
                candidates.append(list(rule_candidate) + entry)
        for candidate in candidates:
            rels_f_list, rels_g_list = [], []

            # Try all possible operands
            mask = []
            for c in candidate[:-outputs_n]:
                if c == "x" or c == "-":
                    mask.append(1)
                else:
                    mask.append(0)

            varying_values_n = sum(mask)
            for i in range(pow(2, varying_values_n)):
                ops = []
                k = 0
                for j in range(len(mask)):
                    # If it's masked, we need to fill up the value
                    if mask[j] == 1:
                        value = i >> k & 1
                        ops.append(value)
                        if candidate[j] == "-":
                            ops.append(value)
                        else:
                            ops.append(1 if value == 0 else 0)
                        k += 1
                    elif candidate[j] == "n":
                        ops += [0, 1]
                    elif candidate[j] == "u":
                        ops += [1, 0]
                    elif candidate[j] == "0":
                        ops += [0, 0]
                    elif candidate[j] == "1":
                        ops += [1, 1]

                # f and g are the 2 blocks of SHA-256
                inputs_f, inputs_g = [ops[j * 2] for j in range(inputs_n)], [
                    ops[j * 2 + 1] for j in range(inputs_n)
                ]

                # Ensure that the output matches the candidate
                output_f, output_g = func(inputs_f), func(inputs_g)
                if not conforms_to(output_f, output_g, candidate[inputs_n:]):
                    continue

                # Derive the relationships
                to_list = lambda items: list(items) if type(items) == tuple else [items]
                rels_f = rels(candidate, inputs_f + to_list(output_f))
                rels_f_list.append(rels_f)

                rels_g = rels(candidate, inputs_g + to_list(output_g))
                rels_g_list.append(rels_g)

            if len(rels_f_list) == 0 and len(rels_g_list) == 0:
                continue

            # Go through the rels list column-wise
            rels_colwise_f = get_rels_colwise(rels_f_list)
            rels_colwise_g = get_rels_colwise(rels_g_list)

            # Check the consistency of the rels. column-wise
            consistency_f = get_rels_consistency(rels_colwise_f)
            consistency_g = get_rels_consistency(rels_colwise_g)

            # Skip this rule if there is no consistent column
            if all([x == 2 for x in consistency_f + consistency_g]):
                continue
            key = "".join(candidate)
            value = "".join([str(x) for x in consistency_f + consistency_g])

            # Save the rule to the database
            print(id, key, value)


gen_2_bit_conds(TWO_BIT_CONSTRAINT_XOR2_ID, xor2, 2)
gen_2_bit_conds(TWO_BIT_CONSTRAINT_ADD2_ID, bin_add_w, 2, 2)
gen_2_bit_conds(TWO_BIT_CONSTRAINT_IF_ID, if_w, 3)
gen_2_bit_conds(TWO_BIT_CONSTRAINT_MAJ_ID, maj_w, 3)
gen_2_bit_conds(TWO_BIT_CONSTRAINT_XOR3_ID, xor3_w, 3)
gen_2_bit_conds(TWO_BIT_CONSTRAINT_ADD3_ID, bin_add_w, 3, 2)
gen_2_bit_conds(TWO_BIT_CONSTRAINT_ADD4_ID, bin_add_w, 4, 3)
gen_2_bit_conds(TWO_BIT_CONSTRAINT_ADD5_ID, bin_add_w, 5, 3)
gen_2_bit_conds(TWO_BIT_CONSTRAINT_ADD6_ID, bin_add_w, 6, 3)
gen_2_bit_conds(TWO_BIT_CONSTRAINT_ADD7_ID, bin_add_w, 7, 3)
