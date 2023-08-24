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


def if_w(ops):
    return if_(ops[0], ops[1], ops[2])


def maj_w(ops):
    return maj(ops[0], ops[1], ops[2])


def xor3_w(ops):
    return xor3(ops[0], ops[1], ops[2])


def bin_add_w(ops):
    carry2, carry1, sum = bin_add(ops)
    return sum


def rels(values):
    rels = ""
    visited = set()
    i = 0
    for op in values:
        j = -1
        for op_ in values:
            j += 1
            if i == j or j in visited:
                continue
            rels += "1" if op == op_ else "0"
        visited.add(i)
        i += 1
    return rels


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


# Assumes that there is one output of the function
def gen_2_bit_conds(id, func, inputs_n, outputs_n = 1):
    rules = {}
    gc_set = ["1", "u", "n", "0", "x", "-"]

    # Generate all the possible rule candidates
    rule_candidates = product(gc_set, repeat=inputs_n + 1)

    # Try all the candidates
    for rule_candidate in rule_candidates:
        # Check if the candidate doesn't involve all known bits
        if not any([x == "x" or x == "-" for x in rule_candidate]):
            continue

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
            xs_f, xs_g = [ops[j * 2] for j in range(inputs_n)], [
                ops[j * 2 + 1] for j in range(inputs_n)
            ]

            # Ensure that the output matches the candidate
            w_f, w_g = func(xs_f), func(xs_g)
            if not conforms_to(w_f, w_g, rule_candidate[inputs_n]):
                continue

            
            rels_f = rels(xs_f)
            rels_f_list.append(rels_f)

            rels_g = rels(xs_g)
            rels_g_list.append(rels_g)

        if len(rels_f_list) == 0 and len(rels_g_list) == 0:
            continue

        # Go through the rels list column-wise
        rels_f_column_wise = {}
        for rel_ in rels_f_list:
            for i in range(len(rel_)):
                if i not in rels_f_column_wise:
                    rels_f_column_wise[i] = []
                rels_f_column_wise[i].append(int(rel_[i]))
        rels_g_column_wise = {}
        for rel_ in rels_g_list:
            for i in range(len(rel_)):
                if i not in rels_g_column_wise:
                    rels_g_column_wise[i] = []
                rels_g_column_wise[i].append(int(rel_[i]))

        # Check the consistency of the rels. column-wise
        consistency_f = []
        for i in range(len(rels_f_column_wise)):
            if all(rels_f_column_wise[i]):
                consistency_f.append(1)
            elif not any(rels_f_column_wise[i]):
                consistency_f.append(0)
            else:
                consistency_f.append(2)
        consistency_g = []
        for i in range(len(rels_g_column_wise)):
            if all(rels_g_column_wise[i]):
                consistency_g.append(1)
            elif not any(rels_g_column_wise[i]):
                consistency_g.append(0)
            else:
                consistency_g.append(2)

        # Accept as rule if there's at least 1 consistent column
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
