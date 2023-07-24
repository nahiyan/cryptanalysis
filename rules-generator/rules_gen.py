import os
from itertools import permutations, combinations

IO_CONSTRAINT_ADD2_ID = 0
IO_CONSTRAINT_IF_ID = 1
IO_CONSTRAINT_MAJ_ID = 2
IO_CONSTRAINT_XOR3_ID = 3
IO_CONSTRAINT_ADD3_ID = 4
IO_CONSTRAINT_ADD4_ID = 5
IO_CONSTRAINT_ADD5_ID = 6
IO_CONSTRAINT_ADD6_ID = 7
IO_CONSTRAINT_ADD7_ID = 8

OI_CONSTRAINT_IF_ID = 9
OI_CONSTRAINT_MAJ_ID = 10
OI_CONSTRAINT_XOR3_ID = 11
OI_CONSTRAINT_ADD3_ID = 12
OI_CONSTRAINT_ADD4_ID = 13
OI_CONSTRAINT_ADD5_ID = 14
OI_CONSTRAINT_ADD6_ID = 15
OI_CONSTRAINT_ADD7_ID = 16

if not os.path.exists("output"):
    os.mkdir("output")
rules_db = open("output/rules-io-oi.db", "wb")


def uniq(data):
    return [list(x) for x in set(tuple(x) for x in data)]


def if_(x, y, z):
    return (x & y) ^ (x & z) ^ z


def maj(x, y, z):
    return (x & y) ^ (y & z) ^ (x & z)


def xor3(x, y, z):
    return x ^ y ^ z


def bin_add(ops):
    carries = []
    bitsum = 0

    for op in ops:
        bitsum += op & 1
    sum = bitsum & 1
    carries.append(bitsum >> 1 & 1)
    carries.append(bitsum >> 2 & 1)
    return carries[1], carries[0], sum


def gc(x, y):
    if x == 0 and y == 0:
        return "0"
    if x == 1 and y == 0:
        return "u"
    if x == 0 and y == 1:
        return "n"
    if x == 1 and y == 1:
        return "1"


sq = {"0", "n", "u", "1"}
sh = {"0", "1"}
sx = {"n", "u"}
s3 = {"0", "u"}
s5 = {"0", "n"}
s7 = {"0", "u", "n"}
sa = {"u", "1"}
sb = {"0", "u", "1"}
sc = {"n", "1"}
sd = {"0", "n", "1"}
se = {"u", "n", "1"}


def to_sym(s):
    if s == sq:
        return "?"
    elif s == sh:
        return "-"
    elif s == sx:
        return "x"
    elif s == {"u"}:
        return "u"
    elif s == {"n"}:
        return "n"
    elif s == {"1"}:
        return "1"
    elif s == {"0"}:
        return "0"
    elif s == s3:
        return "3"
    elif s == s5:
        return "5"
    elif s == s7:
        return "7"
    elif s == sa:
        return "A"
    elif s == sb:
        return "B"
    elif s == sc:
        return "C"
    elif s == sd:
        return "D"
    elif s == se:
        return "E"


def from_sym(s):
    if s == "?":
        return sq
    elif s == "-":
        return sh
    elif s == "x":
        return sx
    elif s == "n":
        return {"n"}
    elif s == "u":
        return {"u"}
    elif s == "1":
        return {"1"}
    elif s == "0":
        return {"0"}
    elif s == "3":
        return s3
    elif s == "5":
        return s5
    elif s == "7":
        return s7
    elif s == "A":
        return sa
    elif s == "B":
        return sb
    elif s == "C":
        return sc
    elif s == "D":
        return sd
    elif s == "E":
        return se


def to_bytearray(func_id, entries):
    data = [func_id]
    data += [ord(x) for x in entries]
    return bytearray(data)


def gen_weight_dists(n):
    weight_dists = []
    for i in range(pow(10, n) + 1):
        i_s = str(i)

        count = 0
        for c in i_s:
            count += int(c)

        if count != n:
            continue

        for _ in range(n - len(i_s)):
            i_s = "0" + i_s

        weight_dists.append([int(x) for x in i_s])
    return weight_dists


# TODO: Add support for IF, MAJ, and XOR3
def oi_rules_gen(id, n):
    print("ADD" + str(n))
    output_freq = {}
    rules = {}
    for i in range(pow(2, n)):
        ops = [i >> (n - 1 - x) & 1 for x in range(n)]

        o1, o2, o3 = bin_add(ops)
        num1s = 0
        for op in ops:
            if op == 1:
                num1s += 1

        if num1s >= 4:
            assert o1 == 1

        if num1s >= 6:
            assert o2 == 1

        for j in range(pow(2, n)):
            ops_ = [j >> (n - 1 - x) & 1 for x in range(n)]

            o1_, o2_, o3_ = bin_add(ops_)
            num1s_ = 0
            for op in ops_:
                if op == 1:
                    num1s_ += 1

            ops_xor = [ops[x] ^ ops_[x] for x in range(n)]
            ops_gc = [gc(ops[x], ops_[x]) for x in range(n)]

            o1_gc = gc(o1, o1_)
            o2_gc = gc(o2, o2_)
            o3_gc = gc(o3, o3_)

            num1sop = 0
            for op in ops_xor:
                if op == 1:
                    num1sop += 1

            input_ = ""
            for op_gc in ops_gc:
                input_ += op_gc
            output = f"{o1_gc}{o2_gc}{o3_gc}"
            if output not in output_freq:
                output_freq[output] = 1
                rules[output] = input_
            else:
                output_freq[output] += 1
    for output in output_freq:
        if output_freq[output] == 1:
            print(output, rules[output])
            rules_db.write(to_bytearray(id, output + rules[output]))


def gen_io_rules(n, allowed_input_syms=None, allowed_out_syms=None):
    symbols = [
        "0",
        "n",
        "u",
        "1",
        "?",
        "-",
        "x",
        "3",
        "5",
        "7",
        "A",
        "B",
        "C",
        "D",
        "E",
    ]
    if allowed_input_syms == None:
        allowed_input_syms = symbols
    # TODO: Enforce it
    if allowed_out_syms == None:
        allowed_out_syms = symbols
    combos = []
    weight_dists = gen_weight_dists(n)
    for combo in combinations(symbols, n):
        for weight_dist in weight_dists:
            items = []
            for i in range(n):
                if combo[i] not in allowed_input_syms:
                    break
                for _ in range(weight_dist[i]):
                    items.append(combo[i])

            if len(items) != n:
                continue
            combos.append(items)

    perms = []
    for i in uniq(combos):
        perm = permutations(i, n)
        for x in uniq(perm):
            perms.append(x)

    rules = []
    for i in perms:
        s = []
        for j in range(n):
            s.append(from_sym(i[j]))

        if n == 3:
            out_syms_if = set()
            out_syms_maj = set()
            out_syms_xor3 = set()
            carry2_syms = set()
            carry1_syms = set()
            sum_syms = set()
            for a in s[0]:
                for b in s[1]:
                    for c in s[2]:
                        ops_s = [a, b, c]
                        ops = []
                        for op_s in ops_s:
                            if op_s == "u":
                                ops.append(1)
                                ops.append(0)
                            elif op_s == "0":
                                ops.append(0)
                                ops.append(0)
                            elif op_s == "1":
                                ops.append(1)
                                ops.append(1)
                            elif op_s == "n":
                                ops.append(0)
                                ops.append(1)
                        if_1 = if_(ops[0], ops[2], ops[4])
                        if_2 = if_(ops[1], ops[3], ops[5])
                        maj_1 = maj(ops[0], ops[2], ops[4])
                        maj_2 = maj(ops[1], ops[3], ops[5])
                        xor3_1 = xor3(ops[0], ops[2], ops[4])
                        xor3_2 = xor3(ops[1], ops[3], ops[5])
                        out_syms_if.add(gc(if_1, if_2))
                        out_syms_maj.add(gc(maj_1, maj_2))
                        out_syms_xor3.add(gc(xor3_1, xor3_2))
                        carry2_1, carry1_1, sum_1 = bin_add([ops[0], ops[2], ops[4]])
                        carry2_2, carry1_2, sum_2 = bin_add([ops[1], ops[3], ops[5]])
                        carry2_syms.add(gc(carry2_1, carry2_2))
                        carry1_syms.add(gc(carry1_1, carry1_2))
                        sum_syms.add(gc(sum_1, sum_2))
            print("IF:  ", i[0], i[1], i[2], "->", to_sym(out_syms_if))
            print("MAJ: ", i[0], i[1], i[2], "->", to_sym(out_syms_maj))
            print("XOR3:", i[0], i[1], i[2], "->", to_sym(out_syms_xor3))
            print(
                "ADD3:",
                i[0],
                i[1],
                i[2],
                "->",
                to_sym(carry1_syms),
                to_sym(sum_syms),
            )
            rules_db.write(to_bytearray(IO_CONSTRAINT_IF_ID, i + [to_sym(out_syms_if)]))
            rules_db.write(
                to_bytearray(IO_CONSTRAINT_MAJ_ID, i + [to_sym(out_syms_maj)])
            )
            rules_db.write(
                to_bytearray(IO_CONSTRAINT_XOR3_ID, i + [to_sym(out_syms_xor3)])
            )
            # TODO: Add rules to the DB for all other operations
            rules_db.write(
                to_bytearray(
                    IO_CONSTRAINT_ADD3_ID,
                    i + [to_sym(carry1_syms), to_sym(sum_syms)],
                )
            )
        elif n == 4:
            carry2_syms = set()
            carry1_syms = set()
            sum_syms = set()
            for a in s[0]:
                for b in s[1]:
                    for c in s[2]:
                        for d in s[3]:
                            ops_s = [a, b, c, d]
                            ops = []
                            for op_s in ops_s:
                                if op_s == "u":
                                    ops.append(1)
                                    ops.append(0)
                                elif op_s == "0":
                                    ops.append(0)
                                    ops.append(0)
                                elif op_s == "1":
                                    ops.append(1)
                                    ops.append(1)
                                elif op_s == "n":
                                    ops.append(0)
                                    ops.append(1)
                            carry2_1, carry1_1, sum_1 = bin_add(
                                [ops[0], ops[2], ops[4], ops[6]]
                            )
                            carry2_2, carry1_2, sum_2 = bin_add(
                                [ops[1], ops[3], ops[5], ops[7]]
                            )
                            carry2_syms.add(gc(carry2_1, carry2_2))
                            carry1_syms.add(gc(carry1_1, carry1_2))
                            sum_syms.add(gc(sum_1, sum_2))
            print(
                "ADD4:",
                i[0],
                i[1],
                i[2],
                i[3],
                "->",
                to_sym(carry2_syms),
                to_sym(carry1_syms),
                to_sym(sum_syms),
            )
        elif n == 5:
            carry2_syms = set()
            carry1_syms = set()
            sum_syms = set()
            for a in s[0]:
                for b in s[1]:
                    for c in s[2]:
                        for d in s[3]:
                            for e in s[4]:
                                ops_s = [a, b, c, d, e]
                                ops = []
                                for op_s in ops_s:
                                    if op_s == "u":
                                        ops.append(1)
                                        ops.append(0)
                                    elif op_s == "0":
                                        ops.append(0)
                                        ops.append(0)
                                    elif op_s == "1":
                                        ops.append(1)
                                        ops.append(1)
                                    elif op_s == "n":
                                        ops.append(0)
                                        ops.append(1)
                                carry2_1, carry1_1, sum_1 = bin_add(
                                    [ops[0], ops[2], ops[4], ops[6], ops[8]]
                                )
                                carry2_2, carry1_2, sum_2 = bin_add(
                                    [ops[1], ops[3], ops[5], ops[7], ops[9]]
                                )
                                carry2_syms.add(gc(carry2_1, carry2_2))
                                carry1_syms.add(gc(carry1_1, carry1_2))
                                sum_syms.add(gc(sum_1, sum_2))
            print(
                "ADD5:",
                i[0],
                i[1],
                i[2],
                i[3],
                i[4],
                "->",
                to_sym(carry2_syms),
                to_sym(carry1_syms),
                to_sym(sum_syms),
            )
        elif n == 6:
            carry2_syms = set()
            carry1_syms = set()
            sum_syms = set()
            for a in s[0]:
                for b in s[1]:
                    for c in s[2]:
                        for d in s[3]:
                            for e in s[4]:
                                for f in s[5]:
                                    ops_s = [a, b, c, d, e, f]
                                    ops = []
                                    for op_s in ops_s:
                                        if op_s == "u":
                                            ops.append(1)
                                            ops.append(0)
                                        elif op_s == "0":
                                            ops.append(0)
                                            ops.append(0)
                                        elif op_s == "1":
                                            ops.append(1)
                                            ops.append(1)
                                        elif op_s == "n":
                                            ops.append(0)
                                            ops.append(1)
                                    carry2_1, carry1_1, sum_1 = bin_add(
                                        [
                                            ops[0],
                                            ops[2],
                                            ops[4],
                                            ops[6],
                                            ops[8],
                                            ops[10],
                                        ]
                                    )
                                    carry2_2, carry1_2, sum_2 = bin_add(
                                        [
                                            ops[1],
                                            ops[3],
                                            ops[5],
                                            ops[7],
                                            ops[9],
                                            ops[11],
                                        ]
                                    )
                                    carry2_syms.add(gc(carry2_1, carry2_2))
                                    carry1_syms.add(gc(carry1_1, carry1_2))
                                    sum_syms.add(gc(sum_1, sum_2))
            print(
                "ADD6:",
                i[0],
                i[1],
                i[2],
                i[3],
                i[4],
                i[5],
                "->",
                to_sym(carry2_syms),
                to_sym(carry1_syms),
                to_sym(sum_syms),
            )
        elif n == 7:
            carry2_syms = set()
            carry1_syms = set()
            sum_syms = set()
            for a in s[0]:
                for b in s[1]:
                    for c in s[2]:
                        for d in s[3]:
                            for e in s[4]:
                                for f in s[5]:
                                    for g in s[6]:
                                        ops_s = [a, b, c, d, e, f, g]
                                        ops = []
                                        for op_s in ops_s:
                                            if op_s == "u":
                                                ops.append(1)
                                                ops.append(0)
                                            elif op_s == "0":
                                                ops.append(0)
                                                ops.append(0)
                                            elif op_s == "1":
                                                ops.append(1)
                                                ops.append(1)
                                            elif op_s == "n":
                                                ops.append(0)
                                                ops.append(1)
                                        carry2_1, carry1_1, sum_1 = bin_add(
                                            [
                                                ops[0],
                                                ops[2],
                                                ops[4],
                                                ops[6],
                                                ops[8],
                                                ops[10],
                                                ops[12],
                                            ]
                                        )
                                        carry2_2, carry1_2, sum_2 = bin_add(
                                            [
                                                ops[1],
                                                ops[3],
                                                ops[5],
                                                ops[7],
                                                ops[9],
                                                ops[11],
                                                ops[13],
                                            ]
                                        )
                                        carry2_syms.add(gc(carry2_1, carry2_2))
                                        carry1_syms.add(gc(carry1_1, carry1_2))
                                        sum_syms.add(gc(sum_1, sum_2))
            print(
                "ADD7:",
                i[0],
                i[1],
                i[2],
                i[3],
                i[4],
                i[5],
                i[6],
                "->",
                to_sym(carry2_syms),
                to_sym(carry1_syms),
                to_sym(sum_syms),
            )


def __main__():
    # gen_io_rules(4, ["x", "-"])
    for i in range(3, 8):
        id = (
            OI_CONSTRAINT_ADD3_ID
            if i == 3
            else OI_CONSTRAINT_ADD4_ID
            if i == 4
            else OI_CONSTRAINT_ADD5_ID
            if i == 5
            else OI_CONSTRAINT_ADD6_ID
            if i == 6
            else OI_CONSTRAINT_ADD7_ID
        )
        oi_rules_gen(id, i)
