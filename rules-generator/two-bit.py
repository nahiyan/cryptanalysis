import os

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


def to_bytearray(func_id, entries):
    data = [func_id]
    data += [ord(x) for x in entries]
    return bytearray(data)


def if_(ops):
    x = ops[0]
    y = ops[1]
    z = ops[2]
    return (x & y) ^ (x & z) ^ z


def maj(ops):
    x = ops[0]
    y = ops[1]
    z = ops[2]
    return (x & y) ^ (y & z) ^ (x & z)


def xor3(ops):
    x = ops[0]
    y = ops[1]
    z = ops[2]
    return x ^ y ^ z


def bin_add(ops):
    carries = []
    bitsum = 0

    for op in ops:
        bitsum += op >> 0 & 1
    sum = bitsum & 1
    carries.append(bitsum >> 1 & 1)
    carries.append(bitsum >> 2 & 1)
    return carries[1], carries[0], sum


def bin_add_w(ops):
    carry2, carry1, sum = bin_add(ops)
    return sum


def gc(x, x_):
    return (
        "0"
        if x == 0 and x_ == 0
        else "1"
        if x == 1 and x_ == 1
        else "u"
        if x == 1 and x_ == 0
        else "n"
    )


def rels(values):
    rels = ""
    visited = []
    i = 0
    for op in values:
        j = -1
        for op_ in values:
            j += 1
            if i == j or j in visited:
                continue
            rels += "1" if op == op_ else "0"
        visited.append(i)
        i += 1
    return rels


def gen_2_bit_conds(id, func, n_ops):
    print(id, func)
    rules = {}
    rules_xors = {}
    for i in range(pow(2, n_ops * 2)):
        ops = [i >> j & 1 for j in range(n_ops * 2)]
        ops_gc = [gc(ops[j * 2], ops[j * 2 + 1]) for j in range(n_ops)]

        w = func([ops[j * 2] for j in range(n_ops)])
        w_ = func([ops[j * 2 + 1] for j in range(n_ops)])
        w_gc = gc(w, w_)

        r = rels([ops[j * 2] for j in range(n_ops)])
        key = "".join(ops_gc) + w_gc
        rules[key] = r

        # XOR rules
        xors = [("1" if c == "n" or c == "u" else "0") for c in key]
        rules_xors_key = "".join(xors)
        if rules_xors_key not in rules_xors:
            rules_xors[rules_xors_key] = []
        rules_xors[rules_xors_key].append(r)

    for rule_xor in rules_xors:
        rels_ = rules_xors[rule_xor]
        r_vals = {}
        for r in rels_:
            i = 0
            for r_ in r:
                if i not in r_vals:
                    r_vals[i] = []
                r_vals[i].append(int(r_))
                i += 1
        consistency = []
        value = []
        for r_val in r_vals:
            all_1 = all(r_vals[r_val])
            all_0 = not any(r_vals[r_val])
            consistency.append(all_1 or all_0)
            if all_1:
                value.append("1")
            elif all_0:
                value.append("0")
            else:
                value.append("2")
        value = "".join(value)
        if any(consistency):
            key = "".join(["x" if c == "1" else "-" for c in rule_xor])
            rules[key] = value

    for key in rules:
        print(key, rules[key])
        rules_db.write(to_bytearray(id, key + rules[key]))


gen_2_bit_conds(TWO_BIT_CONSTRAINT_IF_ID, if_, 3)
# gen_2_bit_conds(TWO_BIT_CONSTRAINT_MAJ_ID, maj, 3)
# gen_2_bit_conds(TWO_BIT_CONSTRAINT_XOR3_ID, xor3, 3)
# gen_2_bit_conds(TWO_BIT_CONSTRAINT_ADD3_ID, bin_add_w, 3)
# gen_2_bit_conds(TWO_BIT_CONSTRAINT_ADD4_ID, bin_add_w, 4)
# gen_2_bit_conds(TWO_BIT_CONSTRAINT_ADD5_ID, bin_add_w, 5)
# gen_2_bit_conds(TWO_BIT_CONSTRAINT_ADD6_ID, bin_add_w, 6)
# gen_2_bit_conds(TWO_BIT_CONSTRAINT_ADD7_ID, bin_add_w, 7)
