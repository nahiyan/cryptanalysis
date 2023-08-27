import os
from itertools import product
from rules_gen import if_, maj, xor3, bin_add, to_bytearray

TWO_BIT_CONSTRAINT_IF_ID = 17
TWO_BIT_CONSTRAINT_MAJ_ID = 18
TWO_BIT_CONSTRAINT_XOR3_ID = 19
TWO_BIT_CONSTRAINT_ADD3_ID = 20
TWO_BIT_CONSTRAINT_ADD4_ID = 21
TWO_BIT_CONSTRAINT_ADD5_ID = 22
TWO_BIT_CONSTRAINT_ADD6_ID = 23
TWO_BIT_CONSTRAINT_ADD7_ID = 24

if not os.path.exists("output"):
    os.mkdir("output")
rules_db = open("output/rules-2-bit.db", "wb")

def create_matrix(n):
    matrix = []
    for i in range(n):
        matrix.append([])
        for _ in range(n):
            matrix[i].append('x')
    return matrix

def print_matrix(matrix):
    r = len(matrix)
    c = len(matrix[0])
    for i in range(r):
        for j in range(c):
            print(matrix[i][j], end=" ")
        print()

def debug_rels(rels_string, n):
    elements = list(range(0, n))
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


def bin_add_w(ops):
    carry2, carry1, sum = bin_add(ops)
    return sum


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


def conforms_to(x_f, x_g, gc_):
    if gc_ == "x":
        return x_f != x_g
    elif gc_ == "-":
        return x_f == x_g
    elif gc_ == "0":
        return x_f == 0 and x_g == 0
    elif gc_ == "u":
        return x_f == 1 and x_g == 0
    elif gc_ == "n":
        return x_f == 0 and x_g == 1
    elif gc_ == "1":
        return x_f == 1 and x_g == 1


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
    rules = {}
    gc_set = ["1", "u", "n", "0", "x", "-"]

    # Generate all the possible rule candidates
    rule_candidates = product(gc_set, repeat=inputs_n + 1)

    # Try all the candidates
    for rule_candidate in rule_candidates:
        rels_f_list, rels_g_list = [], []
        # Try all possible operands
        mask = []
        for c in rule_candidate:
            if c == "x" or c == "-":
                mask.append(1)
            else:
                mask.append(0)

        varying_values_n = sum(mask)
        for i in range(pow(2, varying_values_n)):
            ops = []
            k = 0
            for j, _ in enumerate(rule_candidate):
                # If it's masked, we need to fill up the value
                if mask[j] == 1:
                    value = i >> k & 1
                    ops.append(value)
                    if rule_candidate[j] == "-":
                        ops.append(value)
                    else:
                        ops.append(1 if value == 0 else 0)
                    k += 1
                elif rule_candidate[j] == "n":
                    ops += [0, 1]
                elif rule_candidate[j] == "u":
                    ops += [1, 0]
                elif rule_candidate[j] == "0":
                    ops += [0, 0]
                elif rule_candidate[j] == "1":
                    ops += [1, 1]

            # f and g are the 2 blocks of SHA-256
            inputs_f, inputs_g = [ops[j * 2] for j in range(inputs_n)], [
                ops[j * 2 + 1] for j in range(inputs_n)
            ]

            # Ensure that the output matches the candidate
            output_f, output_g = func(inputs_f), func(inputs_g)
            if not conforms_to(output_f, output_g, rule_candidate[inputs_n]):
                continue

            # Derive the relationships
            rels_f = rels(rule_candidate, inputs_f)
            rels_f_list.append(rels_f)

            rels_g = rels(rule_candidate, inputs_g)
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
        key = "".join(rule_candidate)
        rules[key] = "".join([str(x) for x in consistency_f + consistency_g])

        # Save the rule to the database
        print(id, key, rules[key])


gen_2_bit_conds(TWO_BIT_CONSTRAINT_IF_ID, if_w, 3)
gen_2_bit_conds(TWO_BIT_CONSTRAINT_MAJ_ID, maj_w, 3)
gen_2_bit_conds(TWO_BIT_CONSTRAINT_XOR3_ID, xor3_w, 3)
gen_2_bit_conds(TWO_BIT_CONSTRAINT_ADD3_ID, bin_add_w, 3)
gen_2_bit_conds(TWO_BIT_CONSTRAINT_ADD4_ID, bin_add_w, 4)
gen_2_bit_conds(TWO_BIT_CONSTRAINT_ADD5_ID, bin_add_w, 5)
gen_2_bit_conds(TWO_BIT_CONSTRAINT_ADD6_ID, bin_add_w, 6)
gen_2_bit_conds(TWO_BIT_CONSTRAINT_ADD7_ID, bin_add_w, 7)
