from collections import namedtuple, deque
from itertools import product
from time import time
import re
import sys
from propagator import derive_words_w, _int_diff

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

IO_CONSTRAINT_ADD2_ID = 0
IO_CONSTRAINT_IF_ID = 1
IO_CONSTRAINT_MAJ_ID = 2
IO_CONSTRAINT_XOR3_ID = 3
IO_CONSTRAINT_ADD3_ID = 4
IO_CONSTRAINT_ADD4_ID = 5
IO_CONSTRAINT_ADD5_ID = 6
IO_CONSTRAINT_ADD6_ID = 7
IO_CONSTRAINT_ADD7_ID = 8


k = [
    0x428A2F98,
    0x71374491,
    0xB5C0FBCF,
    0xE9B5DBA5,
    0x3956C25B,
    0x59F111F1,
    0x923F82A4,
    0xAB1C5ED5,
    0xD807AA98,
    0x12835B01,
    0x243185BE,
    0x550C7DC3,
    0x72BE5D74,
    0x80DEB1FE,
    0x9BDC06A7,
    0xC19BF174,
    0xE49B69C1,
    0xEFBE4786,
    0x0FC19DC6,
    0x240CA1CC,
    0x2DE92C6F,
    0x4A7484AA,
    0x5CB0A9DC,
    0x76F988DA,
    0x983E5152,
    0xA831C66D,
    0xB00327C8,
    0xBF597FC7,
    0xC6E00BF3,
    0xD5A79147,
    0x06CA6351,
    0x14292967,
    0x27B70A85,
    0x2E1B2138,
    0x4D2C6DFC,
    0x53380D13,
    0x650A7354,
    0x766A0ABB,
    0x81C2C92E,
    0x92722C85,
    0xA2BFE8A1,
    0xA81A664B,
    0xC24B8B70,
    0xC76C51A3,
    0xD192E819,
    0xD6990624,
    0xF40E3585,
    0x106AA070,
    0x19A4C116,
    0x1E376C08,
    0x2748774C,
    0x34B0BCB5,
    0x391C0CB3,
    0x4ED8AA4A,
    0x5B9CCA4F,
    0x682E6FF3,
    0x748F82EE,
    0x78A5636F,
    0x84C87814,
    0x8CC70208,
    0x90BEFFFA,
    0xA4506CEB,
    0xBEF9A3F7,
    0xC67178F2,
]

add_prop_rules = {}

Table = namedtuple(
    "Table", ["da", "de", "dw", "ds0", "ds1", "dsigma0", "dsigma1", "dch", "dmaj", "ds"]
)
Equation = namedtuple("Equation", ["x", "y", "diff"])


def otf_2bit_eqs(func, inputs, outputs, names=[]):
    positions = []
    combined = inputs + outputs
    in_length, out_length = len(inputs), len(outputs)
    for i, c in enumerate(combined):
        if c in ["x", "-"]:
            positions.append(i)
    # print(positions)
    n = len(positions)
    # print(n)
    selection = []
    for i in range(pow(2, n)):
        values = [i >> j & 1 for j in range(n)]
        candidate = list(combined)
        for j, position in enumerate(positions):
            value = values[j]
            char = candidate[position]
            candidate[position] = (
                ("u" if value == 1 else "n")
                if char == "x"
                else ("1" if value == 1 else "0")
            )
        candidate_inputs = "".join(candidate[:in_length])
        candidate_outputs = "".join(candidate[in_length:])
        prop_inputs, prop_outputs = otf_prop(add, (candidate_inputs, candidate_outputs))
        # print(prop_inputs, prop_outputs)
        if len(prop_outputs) == 0:
            continue
        selection.append((candidate_inputs, candidate_outputs))
    equations = []
    for s in selection:
        combined = s[0] + s[1]
        # print(s[0], s[1])
        for i in range(len(combined)):
            selector = i + 1
            for j in range(selector, len(combined)):
                eq = Equation(
                    names[i], names[j], 1 if combined[i] != combined[j] else 0
                )
                eq_alt = Equation(eq.y, eq.x, eq.diff)

                if eq.x[0] == "?" or eq.y[0] == "?":
                    continue

                if eq in equations or eq_alt in equations:
                    continue
                equations.append(eq)
    finalized_equations = []
    for eq in equations:
        inv_eq = Equation(eq.x, eq.y, 0 if eq.diff == 1 else 1)
        inv_eq_alt = Equation(inv_eq.y, inv_eq.x, inv_eq.diff)

        # Remove contradicting equations
        if inv_eq in equations or inv_eq_alt in equations:
            continue

        finalized_equations.append(eq)

    # for eq in finalized_equations:
    #     print(eq)
    return finalized_equations


def add(addends):
    sum_ = sum(addends)
    output_bits = 3 if len(addends) > 3 else 2
    output = [sum_ >> i & 1 for i in range(output_bits)]
    output.reverse()
    return output


def load_rules(path):
    rules = {}
    with open(path, "r") as rules_file:
        lines = rules_file.readlines()
        for line in lines:
            segments = line.split()
            id = segments[0]
            rule = segments[1]
            value = segments[2]
            rules[f"{id}{rule}"] = value
    return rules


def load_table(path):
    unknown_diff = "".join(["?" for _ in range(32)])

    da, de, dw = {}, {}, {}
    ds0, ds1, dsigma0, dsigma1 = {}, {}, {}, {}
    dch, dmaj, ds = {}, {}, {}
    with open(path, "r") as characteristics_file:
        lines = characteristics_file.readlines()
        for line in lines:
            segments = line.split()
            i = int(segments[0])
            da[i] = segments[1]
            de[i] = segments[2]
            if i >= 0:
                dw[i] = segments[3]
                ds0[i] = unknown_diff
                ds1[i] = unknown_diff
                dsigma0[i] = unknown_diff
                dsigma1[i] = unknown_diff
                dch[i] = unknown_diff
                dmaj[i] = unknown_diff
                ds[i] = unknown_diff
    return Table(da, de, dw, ds0, ds1, dsigma0, dsigma1, dch, dmaj, ds)


def print_table(table):
    da, de, dw = table.da, table.de, table.dw
    for i in da:
        print(i, da[i], de[i], end=" ")
        if i >= 0:
            print(
                dw[i],
                table.ds0[i],
                table.ds1[i],
                table.dsigma0[i],
                table.dsigma1[i],
                table.dmaj[i],
                table.dch[i],
            )
        else:
            print()


# def _int_diff(word):
#     n = len(word)
#     value = 0
#     for i in range(n):
#         gc_bit = word[n - 1 - i]
#         if gc_bit not in ["u", "n", "-", "1", "0"]:
#             return value, True
#         value += (1 if gc_bit == "u" else -1 if gc_bit == "n" else 0) * pow(2, i)
#     return value % pow(2, n), False


def derive_words_step(word_x, word_y, constant):
    table = {
        "?": ["n", "u", "-"],
        "x": ["n", "u"],
        "3": ["0", "u"],
        "5": ["0", "n"],
        "7": ["0", "u", "n"],
        "A": ["u", "1"],
        "B": ["u", "-"],
        "C": ["n", "1"],
        "D": ["n", "-"],
        "E": ["u", "n", "1"],
    }
    conforms = lambda x, y, table: True if x in table[y] else False
    flatten = lambda l: [item for sublist in l for item in sublist]

    possible_gcs = set(
        flatten([table[x] if x in table else [x] for x in word_x + word_y])
    )

    subject = word_x + word_y
    holes = []
    for i, c in enumerate(subject):
        if c in table:
            holes.append(i)
    n = len(holes)
    combos = product(possible_gcs, repeat=n)
    combos_count = int(pow(len(possible_gcs), n))
    if combos_count >= 1e6:
        return "".join(word_x), "".join(word_y)
    m = int(len(subject) / 2)
    matches = []
    for combo in combos:
        new_subject = list(subject)
        skip = False
        for i, hole in enumerate(holes):
            original_v = subject[hole]
            if not conforms(combo[i], original_v, table):
                skip = True
                break
            new_subject[hole] = combo[i]
        if not skip:
            new_subject = "".join(new_subject)
            parts = (new_subject[:m], new_subject[m:])
            (c1, _), (c2, _) = _int_diff(parts[0], m), _int_diff(parts[1], m)
            c = (c1 + c2) & (pow(2, m) - 1)
            if c == constant:
                matches.append(new_subject)
    commons = list(subject)
    for i in holes:
        col = []
        for match_ in matches:
            col.append(match_[i])
        if len(set(col)) == 1:
            commons[i] = col[0]
    commons = "".join(commons)
    derived_word_x, derived_word_y = commons[:m], commons[m:]
    return derived_word_x, derived_word_y


def propagate_addition(table, row, name, vars_):
    underived_indices = []
    for i, var in enumerate(vars_):
        int_diff, err = _int_diff(var)
        # integers will be on the LHS
        if not err:
            vars_[i] = (1 if i == 0 or (name == "add_a" and i == 2) else -1) * int_diff
        else:
            underived_indices.append(i)
    underived_count = len(underived_indices)
    if underived_count == 1 or underived_count == 2:
        addends = [
            0 if i in underived_indices else addend for i, addend in enumerate(vars_)
        ]
        constant = sum(addends) % pow(2, 32)
        #! Debug
        # print(row, name, constant)

        if underived_count == 1:
            index = underived_indices[0]
            underived_var = vars_[index]
            if index == 0 or (name == "add_a" and index == 2):
                constant *= -1

            derived_vars, err = derive_words_w([underived_var], constant)
            assert not err
        else:
            for index in underived_indices:
                if index == 0 or (name == "add_a" and index == 2):
                    constant *= -1
            underived_vars = [vars_[x] for x in underived_indices]
            # print("Addition:", f"{name}_{row}", underived_vars, constant)
            derived_vars, err = derive_words_w(underived_vars, constant)
            assert not err
        for i, index in enumerate(underived_indices):
            value = derived_vars[i]
            match name:
                case "add_w":
                    match index:
                        case 0:
                            table.dw[row] = value
                        case 1:
                            table.ds1[row] = value
                        case 2:
                            table.dw[row - 7] = value
                        case 3:
                            table.ds0[row] = value
                        case 4:
                            table.dw[row - 16] = value
                case "add_e":
                    match index:
                        case 0:
                            table.de[row] = value
                        case 1:
                            table.da[row - 4] = value
                        case 2:
                            table.de[row - 4] = value
                        case 3:
                            table.dsigma1[row] = value
                        case 4:
                            table.dch[row] = value
                        case 5:
                            table.dw[row] = value
                case "add_a":
                    match index:
                        case 0:
                            table.da[row] = value
                        case 1:
                            table.de[row] = value
                        case 2:
                            table.da[row - 4] = value
                        case 3:
                            table.dsigma0[row] = value
                        case 4:
                            table.dmaj[row] = value


def propagate(table, rules):
    for i in table.dw:
        if i >= 16:
            # s0
            word = table.dw[i - 15][::-1]
            s0 = [None] * 32
            for j in range(32):
                x, y, z = (
                    word[(j + 7) % 32],
                    word[(j + 18) % 32],
                    (word[(j + 3) % 32] if j <= 28 else "0"),
                )
                rule = f"{IO_CONSTRAINT_XOR3_ID}{x}{y}{z}"
                if rule in rules:
                    value = rules[rule]
                    s0[31 - j] = value
            table.ds0[i] = "".join(s0)

            # s1
            word = table.dw[i - 2][::-1]
            s1 = [None] * 32
            for j in range(32):
                x, y, z = (
                    word[(j + 17) % 32],
                    word[(j + 19) % 32],
                    (word[(j + 10) % 32] if j <= 21 else "0"),
                )
                rule = f"{IO_CONSTRAINT_XOR3_ID}{x}{y}{z}"
                if rule in rules:
                    value = rules[rule]
                    s1[31 - j] = value
            table.ds1[i] = "".join(s1)

            # add_W
            add_w_vars = [
                table.dw[i],
                table.ds1[i],
                table.dw[i - 7],
                table.ds0[i],
                table.dw[i - 16],
            ]
            propagate_addition(table, i, "add_w", add_w_vars)

        # sigma1
        sigma1 = [None] * 32
        for j in range(32):
            rule = f"{IO_CONSTRAINT_XOR3_ID}{table.de[i - 1][(j - 6) % 32]}{table.de[i - 1][(j - 11) % 32]}{table.de[i - 1][(j - 25) % 32]}"
            if rule in rules:
                value = rules[rule]
                sigma1[j] = value
        table.dsigma1[i] = "".join(sigma1)

        # ch
        ch = [None] * 32
        for j in range(32):
            rule = f"{IO_CONSTRAINT_IF_ID}{table.de[i - 1][j]}{table.de[i - 2][j]}{table.de[i - 3][j]}"
            if rule in rules:
                value = rules[rule]
                ch[j] = value
        table.dch[i] = "".join(ch)

        # add_e
        add_e_vars = [
            table.de[i],
            table.da[i - 4],
            table.de[i - 4],
            table.dsigma1[i],
            table.dch[i],
            table.dw[i],
        ]
        propagate_addition(table, i, "add_e", add_e_vars)

        # sigma0
        sigma0 = [None] * 32
        for j in range(32):
            rule = f"{IO_CONSTRAINT_XOR3_ID}{table.da[i - 1][(j - 2) % 32]}{table.da[i - 1][(j - 13) % 32]}{table.da[i - 1][(j - 22) % 32]}"
            if rule in rules:
                value = rules[rule]
                sigma0[j] = value
        table.dsigma0[i] = "".join(sigma0)

        # maj
        maj = [None] * 32
        for j in range(32):
            rule = f"{IO_CONSTRAINT_MAJ_ID}{table.da[i - 1][j]}{table.da[i - 2][j]}{table.da[i - 3][j]}"
            if rule in rules:
                value = rules[rule]
                maj[j] = value
        table.dmaj[i] = "".join(maj)

        # add_a
        add_a_vars = [
            table.da[i],
            table.de[i],
            table.da[i - 4],
            table.dsigma0[i],
            table.dmaj[i],
        ]
        propagate_addition(table, i, "add_a", add_a_vars)

        # Propagate modular addition
        # _, ds, _ = otf_prop_add_words([
        #     table.dsigma1[i],
        #     table.dch[i],
        #     "".join([str(k[i] >> i & 1) for i in range(32, -1, -1)]),
        # ], "?" * 32)
        # table.ds[i] = ds
        # otf_prop_add_words([
        #     table.da[i - 4],
        #     table.de[i - 4],
        #     s,
        #     table.dw[i],
        # ], table.de[i])
        # otf_prop_add_words([
        #     table.de[i - 4],
        #     s,
        #     table.dw[i],
        #     table.dsigma0[i],
        #     table.dmaj[i],
        # ], table.da[i])


def otf_prop(func, vars):
    gc_to_bin = (
        lambda gc: (0, 0)
        if gc == "0"
        else (1, 1)
        if gc == "1"
        else (1, 0)
        if gc == "u"
        else (0, 1)
    )
    bin_to_gc = (
        lambda bin: "0"
        if bin == (0, 0)
        else "1"
        if bin == (1, 1)
        else "u"
        if bin == (1, 0)
        else "n"
    )
    flatten_pairs = lambda pairs: (
        [pair[0] for pair in pairs],
        [pair[1] for pair in pairs],
    )

    input_vars, output_vars = vars[0], vars[1]
    n, m = len(input_vars), len(output_vars)
    symbols = {
        "?": ["u", "n", "1", "0"],
        "-": ["1", "0"],
        "x": ["u", "n"],
        "0": ["0"],
        "u": ["u"],
        "n": ["n"],
        "1": ["1"],
        "3": ["0", "u"],
        "5": ["0", "n"],
        "7": ["0", "u", "n"],
        "A": ["u", "1"],
        "B": ["1", "u", "0"],
        "C": ["n", "1"],
        "D": ["0", "n", "1"],
        "E": ["u", "n", "1"],
    }
    input_gcs = set()
    for var in input_vars:
        for var_ in symbols[var]:
            input_gcs.add(var_)
    combos = product(input_gcs, repeat=n)
    possibilities = [set() for _ in range(n + m)]
    for combo in combos:
        # Input must conform to that given
        skip = False
        for i, var in enumerate(combo):
            if var not in symbols[input_vars[i]]:
                skip = True
        if skip:
            continue

        bin_input_vars = flatten_pairs([gc_to_bin(var) for var in combo])
        bin_outputs = [func(inputs) for inputs in bin_input_vars]
        gc_outputs = [
            bin_to_gc((bin_outputs[0][i], bin_outputs[1][i]))
            for i in range(len(bin_outputs[0]))
        ]

        # Output must conform to that given
        skip = False
        for i, gc in enumerate(gc_outputs):
            if gc not in symbols[output_vars[i]]:
                skip = True
        if skip:
            continue

        # rules.append(("".join(combo), "".join(gc_outputs)))
        merged_io = "".join(combo) + "".join(gc_outputs)
        for i, gc in enumerate(merged_io):
            possibilities[i].add(gc)
        # print("".join(combo), "".join(gc_outputs))
    propagation = []
    for p in possibilities:
        for symbol in symbols:
            if set(symbols[symbol]) == p:
                propagation.append(symbol)
    propagation = "".join(propagation)
    return propagation[:n], propagation[n:]


def otf_prop_add_words(words, sum, n=32, model={}):
    imported_carries = "carries" in model
    if imported_carries:
        high_carries, low_carries = model["carries"]
    else:
        high_carries, low_carries = ["0"] * n, ["0"] * n

    # if imported_carries:
    #     print(high_carries, low_carries)

    m = len(words)
    has_high_carry = m > 3
    prop_words = [[]] * m
    for i in range(m):
        prop_words[i] = [None] * n
    output_prop = [None] * n
    steps = []
    for i in range(n):
        gcs = [words[j][n - i - 1] for j in range(m)]
        if i > 0:
            gcs.append(low_carries[i - 1])
        if i > 1 and has_high_carry:
            gcs.append(high_carries[i - 2])

        # Look up in the cache, otherwise do it OTF
        key = ("".join(gcs), f"{sum[n - i - 1]}")
        if key in add_prop_rules:
            inputs_prop = key[0]
            outputs_prop = add_prop_rules[key]
        else:
            inputs_prop, outputs_prop = otf_prop(
                add, (gcs, ("??" if has_high_carry else "?") + f"{sum[n - i - 1]}")
            )
            add_prop_rules[key] = outputs_prop
            # print(
            #     "Cache",
            #     key,
            #     inputs_prop,
            #     len(inputs_prop),
            #     len(outputs_prop),
            #     outputs_prop,
            # )

        assert len(outputs_prop) > 0
        steps.append((inputs_prop, outputs_prop))
        if has_high_carry and not imported_carries:
            high_carries[i] = outputs_prop[0]
        if not imported_carries:
            low_carries[i] = outputs_prop[1] if has_high_carry else outputs_prop[0]
        for k, gc in enumerate(inputs_prop[:m]):
            prop_words[k][n - i - 1] = gc
        # print(outputs_prop[-1], end="")
        output_prop[n - i - 1] = outputs_prop[-1]
    # print()
    for i in range(m):
        prop_words[i] = "".join(prop_words[i])
    return prop_words, "".join(output_prop), (steps, (high_carries, low_carries))


def set_iv(table):
    iv_a = [0x6A09E667, 0xBB67AE85, 0x3C6EF372, 0xA54FF53A]
    iv_e = [0x510E527F, 0x9B05688C, 0x1F83D9AB, 0x5BE0CD19]
    set_word = lambda x: "".join([str(x >> i & 1) for i in range(32)])[::-1]
    for i in range(-4, 0):
        table.da[i] = set_word(iv_a[i])
        table.de[i] = set_word(iv_e[i])
    # for i in range(-4, 0):
    #     table.da[i] = "-" * 32
    #     table.de[i] = "-" * 32


def get_equations(rel_matrix, var_names):
    vars_count = len(var_names)
    rel_matrix_i = -1
    equations = []
    for _ in range(1):
        for i in range(vars_count):
            for j in range(i + 1, vars_count):
                rel_matrix_i += 1
                if (
                    rel_matrix[rel_matrix_i] == "2"
                    or vars_count - 1 == i
                    or vars_count - 1 == j
                ):
                    continue
                are_equal = rel_matrix[rel_matrix_i] == "1"
                equation = Equation(var_names[i], var_names[j], 0 if are_equal else 1)
                equations.append(equation)
                # print(var_names[i], "=" if are_equal else "=/=", var_names[j])
    return equations


def derive_equations(table, rules):
    equations = []
    for i in table.dw:
        k_ = "".join([str(k[i] >> j & 1) for j in range(31, -1, -1)])
        _, dt, (dt_steps, _) = otf_prop_add_words(
            [table.de[i - 4], table.dsigma1[i], table.dch[i], k_, table.dw[i]],
            "?" * 32,
        )

        # print(f"Addition {i}")
        _, _, (de_steps, _) = otf_prop_add_words(
            [
                table.da[i - 4],
                dt,
                # table.de[i - 4],
                # table.dsigma1[i],
                # table.dch[i],
                # k_,
                # table.dw[i],
            ],
            table.de[i],
        )

        #!Debug
        # _, _, (de_steps, _) = otf_prop_add_words(
        #     [
        #         table.da[i - 4],
        #         table.de[i - 4],
        #         table.dsigma1[i],
        #         table.dch[i],
        #         k_,
        #         table.dw[i],
        #     ],
        #     "?" * 32
        #     # table.de[i],
        # )

        for x, (inputs, outputs) in enumerate(de_steps):
            if (
                sum(
                    [
                        1 if c not in ["1", "0", "n", "u"] else 0
                        for c in inputs + outputs
                    ]
                )
                > 4
            ):
                continue
            if len(inputs) != 3:
                inputs += "0" * (3 - len(inputs))
            eqs = otf_2bit_eqs(
                add,
                inputs,
                outputs,
                names=[
                    f"A_{i-4},{x}",
                    # f"E_{i-4},{x}",
                    f"T_{i},{x}",
                    f"?_{i},{x}",
                    # f"?_{i},{x}",
                    # f"W_{i},{x}",
                    # f"?_{i},{x}",
                    # f"?_{i},{x}",
                    # f"?_{i},{x}",
                    f"?_{i},{x}",
                    f"E_{i},{x}",
                ],
            )

            # t = {}
            # for eq in eqs:
            #     x = eq.x[0]
            #     y = eq.y[0]
            #     if x == "T":
            #         t[eq.x] = set()
            #     if y == "T":
            #         t[eq.y] = set()
            # for eq in eqs:
            #     x = eq.x[0]
            #     y = eq.y[0]
            #     if x == "T":
            #         t[eq.x].add((eq.y, eq.diff))
            #     if y == "T":
            #         t[eq.y].add((eq.x, eq.diff))
            # for key in t:
            #     value = t[key]
            #     if len(value) > 1:
            #         print(value)
            equations.extend(eqs)

        for x, (inputs, outputs) in enumerate(dt_steps):
            if (
                sum(
                    [
                        1 if c not in ["1", "0", "n", "u"] else 0
                        for c in inputs + outputs
                    ]
                )
                > 4
            ):
                continue
            if len(inputs) != 7:
                inputs += "0" * (7 - len(inputs))
            eqs = otf_2bit_eqs(
                add,
                inputs,
                outputs,
                names=[
                    f"E_{i-4},{x}",
                    f"?_{i},{x}",
                    f"?_{i},{x}",
                    f"?_{i},{x}",
                    f"W_{i},{x}",
                    f"?_{i},{x}",
                    f"?_{i},{x}",
                    f"?_{i},{x}",
                    f"?_{i},{x}",
                    f"T_{i},{x}",
                ],
            )
            equations.extend(eqs)

            # t = {}
            # for eq in eqs:
            #     x = eq.x[0]
            #     y = eq.y[0]
            #     if x == "T":
            #         t[eq.x] = set()
            #     if y == "T":
            #         t[eq.y] = set()
            # for eq in eqs:
            #     x = eq.x[0]
            #     y = eq.y[0]
            #     if x == "T":
            #         t[eq.x].add((eq.y, eq.diff))
            #     if y == "T":
            #         t[eq.y].add((eq.x, eq.diff))
            # for key in t:
            #     value = t[key]
            #     if len(value) > 1:
            #         print(value)
            equations.extend(eqs)

        # !Debug
        # for x, (inputs, outputs) in enumerate(de_steps):
        #     if (
        #         sum(
        #             [
        #                 1 if c not in ["1", "0", "n", "u"] else 0
        #                 for c in inputs + outputs
        #             ]
        #         )
        #         > 4
        #     ):
        #         continue
        #     if len(inputs) != 8:
        #         inputs += "0" * (8 - len(inputs))
        #     eqs = otf_2bit_eqs(
        #         add,
        #         inputs,
        #         outputs,
        #         names=[
        #             f"A_{i-4},{x}",
        #             f"E_{i-4},{x}",
        #             f"?_{i},{x}",
        #             f"?_{i},{x}",
        #             f"?_{i},{x}",
        #             f"W_{i},{x}",
        #             f"?_{i},{x}",
        #             f"?_{i},{x}",
        #             f"?_{i},{x}",
        #             f"?_{i},{x}",
        #             f"E_{i},{x}",
        #         ],
        #     )
        #     equations.extend(eqs)

        #     # t = {}
        #     # for eq in eqs:
        #     #     x = eq.x[0]
        #     #     y = eq.y[0]
        #     #     if x == "T":
        #     #         t[eq.x] = set()
        #     #     if y == "T":
        #     #         t[eq.y] = set()
        #     # for eq in eqs:
        #     #     x = eq.x[0]
        #     #     y = eq.y[0]
        #     #     if x == "T":
        #     #         t[eq.x].add((eq.y, eq.diff))
        #     #     if y == "T":
        #     #         t[eq.y].add((eq.x, eq.diff))
        #     # for key in t:
        #     #     value = t[key]
        #     #     if len(value) > 1:
        #     #         print(value)
        #     # equations.extend(eqs)

        for j in range(32):
            if i >= 16:
                # s0
                indices = [(j + 7) % 32, (j + 18) % 32, j + 3]
                word_index = i - 15
                s0_i1 = table.dw[word_index][31 - indices[0]]
                s0_i2 = table.dw[word_index][31 - indices[1]]
                s0_i3 = table.dw[word_index][31 - indices[2]]
                s0 = table.ds0[i][31 - j]
                key = (
                    f"{TWO_BIT_CONSTRAINT_XOR3_ID}{s0_i1}{s0_i2}{s0_i3}{s0}"
                    if j <= 29
                    else f"{TWO_BIT_CONSTRAINT_XOR2_ID}{s0_i1}{s0_i2}{s0}"
                )
                if key in rules:
                    value = rules[key]
                    equations.extend(
                        get_equations(
                            value,
                            [
                                f"W_{word_index},{indices[0]}",
                                f"W_{word_index},{indices[1]}",
                                f"W_{word_index},{indices[2]}",
                                f"s0_{i},{j}",
                            ]
                            if j <= 29
                            else [
                                f"W_{word_index},{indices[0]}",
                                f"W_{word_index},{indices[1]}",
                                f"s0_{i},{j}",
                            ],
                        )
                    )

                # s1
                indices = [(j + 17) % 32, (j + 19) % 32, j + 10]
                word_index = i - 2
                s1_i1 = table.dw[word_index][31 - indices[0]]
                s1_i2 = table.dw[word_index][31 - indices[1]]
                s1_i3 = table.dw[word_index][31 - indices[2]]
                s1 = table.ds1[i][31 - j]
                key = (
                    f"{TWO_BIT_CONSTRAINT_XOR3_ID}{s1_i1}{s1_i2}{s1_i3}{s1}"
                    if j <= 21
                    else f"{TWO_BIT_CONSTRAINT_XOR2_ID}{s1_i1}{s1_i2}{s1}"
                )
                if key in rules:
                    value = rules[key]
                    equations.extend(
                        get_equations(
                            value,
                            [
                                f"W_{word_index},{indices[0]}",
                                f"W_{word_index},{indices[1]}",
                                f"W_{word_index},{indices[2]}",
                                f"s1_{i},{j}",
                            ]
                            if j <= 29
                            else [
                                f"W_{word_index},{indices[0]}",
                                f"W_{word_index},{indices[1]}",
                                f"s1_{i},{j}",
                            ],
                        )
                    )

            # sigma0
            indices = [(j + 2) % 32, (j + 13) % 32, (j + 22) % 32]
            word_index = i - 1
            sigma0_i1 = table.da[word_index][31 - indices[0]]
            sigma0_i2 = table.da[word_index][31 - indices[1]]
            sigma0_i3 = table.da[word_index][31 - indices[2]]
            sigma0 = table.dsigma0[i][31 - j]
            key = (
                f"{TWO_BIT_CONSTRAINT_XOR3_ID}{sigma0_i1}{sigma0_i2}{sigma0_i3}{sigma0}"
            )
            if key in rules:
                value = rules[key]
                equations.extend(
                    get_equations(
                        value,
                        [
                            f"A_{word_index},{indices[0]}",
                            f"A_{word_index},{indices[1]}",
                            f"A_{word_index},{indices[2]}",
                            f"sigma0_{i},{j}",
                        ],
                    )
                )

            # sigma1
            indices = [(j + 6) % 32, (j + 11) % 32, (j + 25) % 32]
            word_index = i - 1
            sigma1_i1 = table.de[word_index][31 - indices[0]]
            sigma1_i2 = table.de[word_index][31 - indices[1]]
            sigma1_i3 = table.de[word_index][31 - indices[2]]
            sigma1 = table.dsigma1[i][31 - j]
            key = (
                f"{TWO_BIT_CONSTRAINT_XOR3_ID}{sigma1_i1}{sigma1_i2}{sigma1_i3}{sigma1}"
            )
            if key in rules:
                value = rules[key]
                equations.extend(
                    get_equations(
                        value,
                        [
                            f"E_{word_index},{indices[0]}",
                            f"E_{word_index},{indices[1]}",
                            f"E_{word_index},{indices[2]}",
                            f"sigma1_{i},{j}",
                        ],
                    )
                )

            # maj
            word_indices = [i - 1, i - 2, i - 3]
            maj_i1 = table.da[word_indices[0]][31 - j]
            maj_i2 = table.da[word_indices[1]][31 - j]
            maj_i3 = table.da[word_indices[2]][31 - j]
            maj = table.dmaj[i][31 - j]
            key = f"{TWO_BIT_CONSTRAINT_MAJ_ID}{maj_i1}{maj_i2}{maj_i3}{maj}"
            # maj_found = False
            # if key in rules:
            #     maj_found = True
            # elif maj == 'u' or maj == 'n':
            #     key = (
            #         f"{TWO_BIT_CONSTRAINT_MAJ_ID}{maj_i1}{maj_i2}{maj_i3}x"
            #     )
            #     maj_found = key in rules
            if key in rules:
                value = rules[key]
                equations.extend(
                    get_equations(
                        value,
                        [
                            f"A_{word_indices[0]},{j}",
                            f"A_{word_indices[1]},{j}",
                            f"A_{word_indices[2]},{j}",
                            f"maj_{i},{j}",
                        ],
                    )
                )

            # ch
            word_indices = [i - 1, i - 2, i - 3]
            ch_i1 = table.de[word_indices[0]][31 - j]
            ch_i2 = table.de[word_indices[1]][31 - j]
            ch_i3 = table.de[word_indices[2]][31 - j]
            ch = table.dch[i][31 - j]
            key = f"{TWO_BIT_CONSTRAINT_IF_ID}{ch_i1}{ch_i2}{ch_i3}{ch}"
            if key in rules:
                value = rules[key]
                equations.extend(
                    get_equations(
                        value,
                        [
                            f"E_{word_indices[0]},{j}",
                            f"E_{word_indices[1]},{j}",
                            f"E_{word_indices[2]},{j}",
                            f"ch_{i},{j}",
                        ],
                    )
                )

            # add_e
            # sum_ = table.de[i][31 - j]
            # addends = [
            #     table.da[i - 4][31 - j],
            #     table.de[i - 4][31 - j],
            #     table.dsigma1[i][31 - j],
            #     table.dch[i][31 - j],
            #     str(k[i] >> j & 1),
            #     table.dw[i][31 - j],
            # ]
            # key = f"{TWO_BIT_CONSTRAINT_ADD6_ID}"
            # for addend in addends:
            #     key += addend
            # key += f"??{sum_}"
            # if key in rules:
            #     value = rules[key]
            #     equations.extend(
            #         get_equations(
            #             value,
            #             [
            #                 f"A_{i - 4},{j}",
            #                 f"E_{i - 4},{j}",
            #                 f"?_{i},{j}",
            #                 f"?_{i},{j}",
            #                 f"?_{i},{j}",
            #                 f"W_{i},{j}",
            #                 f"?_{i},{j}",
            #                 f"?_{i},{j}",
            #                 f"E_{i},{j}",
            #             ],
            #         )
            #     )

    return equations


def summarize_2bit_cnds(equations):
    with open("cnds.log", "r") as cnds_file:
        cnds = cnds_file.read()
    found_equations, not_found_equations = [], []
    for equation in equations:
        cnd_a = f"Eq(x='{equation.x}', y='{equation.y}', diff={equation.diff})"
        cnd_b = f"Eq(x='{equation.y}', y='{equation.x}', diff={equation.diff})"
        if cnd_a in cnds or cnd_b in cnds:
            found_equations.append(equation)
        else:
            not_found_equations.append(equation)

    print("Found equations:", len(found_equations))
    for equation in found_equations:
        print(equation)
    print()
    print("Non-existent equations:", len(not_found_equations))
    for equation in not_found_equations:
        print(equation)
    print()

    missing_equations = []
    for cnd in cnds.split("\n"):
        result = re.findall(
            r"Eq\(x='([AWE]_\d+,\d+)', y='([AWE]_\d+,\d+)', diff=(\d)\)", cnd
        )[0]
        x, y, diff = result[0], result[1], int(result[2])
        eq_a = Equation(x, y, diff)
        eq_b = Equation(y, x, diff)
        if eq_a not in equations and eq_b not in equations:
            missing_equations.append(eq_a)
    print("Missing equations:", len(missing_equations))
    for eq in missing_equations:
        print(eq)


def print_steps(table):
    sbs = lambda x: f"{{{x}}}" if x < 0 else x
    k_gc = lambda i: "".join([str(k[i] >> j & 1) for j in range(32)][::-1])

    for i in table.dw:
        print(
            f"W_{i} = M_{i}"
            if i <= 15
            else f"W_{i} = σ_1(W_{i - 2}) + W_{i - 7} + σ_0(W_{i - 15}) + W_{i - 16}"
        )
        print(
            f"[{table.dw[i]}]"
            if i <= 15
            else f"[{table.dw[i]}] = [{table.ds1[i]}] + [{table.dw[i - 7]}] + [{table.ds0[i]}] + [{table.dw[i - 16]}]"
        )

        print(
            f"E_{i} = A_{sbs(i - 4)} + E_{sbs(i - 4)} + Σ_1(E_{sbs(i - 1)}) + IF(E_{sbs(i - 1)}, E_{sbs(i - 2)}, E_{sbs(i - 3)}) + K_{i} + W_{i}"
        )
        print(
            f"[{table.de[i]}] = [{table.da[i - 4]}] + [{table.de[i - 4]}] + [{table.dsigma1[i]}] + [{table.dch[i]}] + [{k_gc(i)}] + [{table.dw[i]}]"
        )

        print(
            f"A_{i} = E_{i} - A_{sbs(i-4)} + Σ_0(A_{sbs(i-1)}) + MAJ_(A_{sbs(i-1)}, A_{sbs(i-2)}, A_{sbs(i-3)})"
        )
        print(
            f"[{table.da[i]}] = [{table.de[i]}] - [{table.da[i - 4]}] + [{table.dsigma0[i]}] + [{table.dmaj[i]}]"
        )


def derive(order):
    two_bit_rules = load_rules("2-bit-rules.txt")
    print(len(two_bit_rules), "2-bit rules")

    prop_rules = load_rules("prop-rules.txt")
    print(len(prop_rules), "prop. rules")

    table = load_table(f"{order}.table")
    start_time = time()
    propagate(table, prop_rules)
    print_table(table)
    # print_steps(table)
    print()
    equations = derive_equations(table, two_bit_rules)
    print("2-bit cnds: {:.2f} seconds".format(time() - start_time), "\n")
    summarize_2bit_cnds(equations)
    print()


if __name__ == "__main__":
    start_time = time()
    derive(27)
    print("Total: {:.2f} seconds".format(time() - start_time))
