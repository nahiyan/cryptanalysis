from collections import namedtuple, deque
from itertools import product
from time import time
import sys

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

Table = namedtuple(
    "Table", ["da", "de", "dw", "ds0", "ds1", "dsigma0", "dsigma1", "dch", "dmaj", "dt"]
)

F32 = 0xFFFFFFFF

Equation = namedtuple("Equation", ["x", "y", "diff"])


# def _rotr(x, y):
#     return ((x >> y) | (x << (32 - y))) & F32


def _maj(x, y, z):
    return (x & y) ^ (x & z) ^ (y & z)


def _ch(x, y, z):
    return (x & y) ^ ((~x) & z)


# def _s0(x):
#     return _rotr(x, 7) ^ _rotr(x, 18) ^ (x >> 3)


def s0_gc(gc):
    gc = [1 if x in ["u", "n", "x"] else 0 for x in gc]
    x, y, z = deque(gc), deque(gc), deque(gc)
    x.rotate(7)
    y.rotate(18)
    z.rotate(3)
    for i in range(3):
        z[i] = 0
    output = [None] * 32
    for i in range(32):
        output[i] = str("x" if x[i] ^ y[i] ^ z[i] == 1 else "-")
    return "".join(output)


# def _s1(x):
#     return _rotr(x, 17) ^ _rotr(x, 19) ^ (x >> 10)


def s1_gc(gc):
    gc = [1 if x in ["u", "n", "x"] else 0 for x in gc]
    x, y, z = deque(gc), deque(gc), deque(gc)
    x.rotate(17)
    y.rotate(19)
    z.rotate(10)
    for i in range(10):
        z[i] = 0
    output = [None] * 32
    for i in range(32):
        output[i] = str("x" if x[i] ^ y[i] ^ z[i] == 1 else "-")
    return "".join(output)


# def _sigma0(x):
#     return _rotr(x, 2) ^ _rotr(x, 13) ^ _rotr(x, 22)


def sigma0_gc(gc):
    gc = [1 if x in ["u", "n", "x"] else 0 for x in gc]
    x, y, z = deque(gc), deque(gc), deque(gc)
    x.rotate(2)
    y.rotate(13)
    z.rotate(22)
    output = [None] * 32
    for i in range(32):
        output[i] = str("x" if x[i] ^ y[i] ^ z[i] == 1 else "-")
    return "".join(output)


# def _sigma1(x):
#     return _rotr(x, 6) ^ _rotr(x, 11) ^ _rotr(x, 25)


def sigma1_gc(gc):
    gc = [1 if x in ["u", "n", "x"] else 0 for x in gc]
    x, y, z = deque(gc), deque(gc), deque(gc)
    x.rotate(6)
    y.rotate(11)
    z.rotate(25)
    output = [None] * 32
    for i in range(32):
        output[i] = str("x" if x[i] ^ y[i] ^ z[i] == 1 else "-")
    return "".join(output)


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
    dch, dmaj, dt = {}, {}, {}
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
                dt[i] = unknown_diff
    return Table(da, de, dw, ds0, ds1, dsigma0, dsigma1, dch, dmaj, dt)


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
                # table.dt[i],
            )
        else:
            print()


def _int_diff(gc, n=32):
    value = 0
    for i in range(n):
        gc_bit = gc[n - 1 - i]
        if gc_bit not in ["u", "n", "-", "1", "0"]:
            return value, True
        value += (1 if gc_bit == "u" else -1 if gc_bit == "n" else 0) * pow(2, i)
    return value & (pow(2, n) - 1), False


# def naive_search(str_var, constant):
#     _has_q_marks = lambda var: sum([1 if c == "?" else 0 for c in var]) > 0

#     hole_positions = []
#     for i, c in enumerate(str_var):
#         if c in ["x", "?"]:
#             hole_positions.append(i)
#     hole_count = len(hole_positions)

#     conforms = (
#         lambda x, y: sum(
#             [
#                 1 if y[hole_positions[i]] == "x" and x[i] not in ["u", "n"] else 0
#                 for i, _ in enumerate(x)
#             ]
#         )
#         == 0
#     )
#     combos = product(
#         ["n", "u", "-"] if _has_q_marks(str_var) else ["n", "u"], repeat=hole_count
#     )

#     # if _has_q_marks(str_var):
#     #     print("Debugging q marks")
#     #     for combo in combos:
#     #         print("".join(combo))

#     # print(str_var, _has_q_marks(str_var))
#     for combo in combos:
#         if not conforms(combo, str_var):
#             continue
#         new_str_var = list(str_var)
#         for i, position in enumerate(hole_positions):
#             choice = combo[i]
#             new_str_var[position] = choice
#         new_str_var = "".join(new_str_var)
#         int_diff, err = _int_diff(new_str_var)
#         assert not err
#         # print(new_str_var, int_diff)
#         if int_diff == constant:
#             return new_str_var, False
#     return 0, True


def derive_words_(word_x, word_y, constant):
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
    if combos_count >= 1000000:
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


def derive_words(word_x, word_y, constant, n=32):
    new_words = [], []
    delta_zero_gcs = ["-", "1", "0"]
    new_words_list = []
    new_words_constant = []
    for i in range(n):
        gcs = word_x[n - i - 1], word_y[n - i - 1]
        new_words[0].append(gcs[0])
        new_words[1].append(gcs[1])
        new_words_constant.append(constant >> i & 1)
        if (
            gcs[0] in delta_zero_gcs
            and gcs[1] in delta_zero_gcs
            and new_words_constant[-1] == 0
        ) or i == n - 1:
            new_words[0].reverse()
            new_words[1].reverse()
            new_words_constant.reverse()
            new_words_constant[: len(new_words[0])]
            new_words_constant = int("".join([str(x) for x in new_words_constant]), 2)
            # print("".join(new_words[0]), "".join(new_words[1]), format(new_words_constant, f"0{len(new_words[0])}b"), sep="\n", end="\n\n")
            new_words_list.append((new_words[0], new_words[1], new_words_constant))
            new_words = [], []
            new_words_constant = []

    new_words = [], []
    for word_x, word_y, constant in new_words_list:
        derived_word_x, derived_word_y = derive_words_(word_x, word_y, constant)
        new_words[0].append(derived_word_x)
        new_words[1].append(derived_word_y)
    new_words[0].reverse()
    new_words[1].reverse()
    derived_word_x, derived_word_y = "".join(new_words[0]), "".join(new_words[1])
    return derived_word_x, derived_word_y


def derive_word(word, constant):
    problemetic_gcs = set(["7", "E"])
    adjustable_gcs = set(["x", "n", "5", "C", "D"])
    gcs = set(list(word))
    # table used to derive GCs after manipulation
    table = {
        "x": {0: "n", 1: "u"},
        "3": {0: "0", 1: "u"},
        "5": {0: "n", 1: "0"},
        "A": {0: "1", 1: "u"},
        "B": {0: "-", 1: "u"},
        "C": {0: "n", 1: "1"},
        "D": {0: "n", 1: "-"},
    }

    # detect if derivable
    if not gcs.isdisjoint(problemetic_gcs):
        return "", True
    if "x" in gcs and not gcs == set(["-", "x"]):
        return "", True

    # perform manipulations to keep each bit from influencing others during diff. calc.
    new_constant = constant
    for i in range(32):
        gc = word[31 - i]
        if gc in adjustable_gcs:
            new_constant += pow(2, i)
    if "x" in gcs:
        new_constant = int(new_constant / 2)

    # derive
    derived_word = list(word)
    for i in range(32):
        gc = word[31 - i]
        bit = new_constant >> i & 1
        if gc in table:
            derived_word[31 - i] = table[gc][bit]
    derived_word = "".join(derived_word)

    int_diff, err = _int_diff(derived_word)
    assert not err and int_diff == constant
    return derived_word, False


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
        constant = sum(addends) & F32

        if underived_count == 1:
            index = underived_indices[0]
            underived_var = vars_[index]
            if index == 0 or (name == "add_a" and index == 2):
                constant *= -1

            derived_var, err = derive_word(underived_var, constant)
            if err:
                print("Failed", underived_var, constant)
                return
            derived_vars = [derived_var]
        else:
            for index in underived_indices:
                if index == 0 or (name == "add_a" and index == 2):
                    constant *= -1
            underived_vars = [vars_[x] for x in underived_indices]
            derived_vars = derive_words(underived_vars[0], underived_vars[1], constant)
        
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
            table.ds0[i] = s0_gc(table.dw[i - 15])
            # s1
            table.ds1[i] = s1_gc(table.dw[i - 2])

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
        table.dsigma1[i] = sigma1_gc(table.de[i - 1])

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
        table.dsigma0[i] = sigma0_gc(table.da[i - 1])

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

    return equations


def derive(order):
    two_bit_rules = load_rules("2-bit-rules.txt")
    print(len(two_bit_rules), "2-bit rules")

    prop_rules = load_rules("prop-rules.txt")
    print(len(prop_rules), "prop. rules")

    table = load_table(f"{order}.table")
    # set_iv(table)
    propagate(table, prop_rules)
    print_table(table)
    equations = derive_equations(table, two_bit_rules)

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
    print("Non-existent equations:", len(not_found_equations))
    for equation in not_found_equations:
        print(equation)


start_time = time()
derive(27)
print("{:.2f} seconds".format(time() - start_time))
