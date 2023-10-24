from itertools import product


def _int_diff(word):
    n = len(word)
    value = 0
    for i in range(n):
        gc_bit = word[n - 1 - i]
        if gc_bit not in ["u", "n", "-", "1", "0"]:
            return value, True
        value += (1 if gc_bit == "u" else -1 if gc_bit == "n" else 0) * pow(2, i)
    return value % pow(2, n), False


def naive_derive_step(word_x, word_y, constant):
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
    # if combos_count >= 1e10:
    #     return "".join(word_x), "".join(word_y)
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
            (c1, _), (c2, _) = _int_diff(parts[0]), _int_diff(parts[1])
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


def naive_derive(word_x, word_y, constant, n=32):
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
        derived_word_x, derived_word_y = naive_derive_step(word_x, word_y, constant)
        new_words[0].append(derived_word_x)
        new_words[1].append(derived_word_y)
    new_words[0].reverse()
    new_words[1].reverse()
    derived_word_x, derived_word_y = "".join(new_words[0]), "".join(new_words[1])
    return derived_word_x, derived_word_y


def split_words(words, constant):
    _is_trivial = lambda col, bit: all([c in ["0", "1", "-"] for c in col]) and bit == 0

    n = len(words[0])
    new_words_list = []
    new_words_constant = []
    new_words = []
    for i in range(n):
        column = [word[n - 1 - i] for word in words]
        new_words.append(column)
        bit = constant >> i & 1
        new_words_constant.append(bit)
        if _is_trivial(column, bit) or i == n - 1:
            for new_word in new_words:
                new_word = str(list(new_word)[::-1])
            new_words_constant.reverse()
            new_words_constant[: len(new_words[0])]
            new_words_constant = int("".join([str(x) for x in new_words_constant]), 2)
            new_words_list.append((new_words, new_words_constant))
            new_words = []
            new_words_constant = []

    m = len(words)
    pieces = []
    for cols, c in new_words_list:
        words = [[] for _ in range(m)]
        for col in cols:
            for word_index in range(m):
                words[word_index].append(col[word_index])
        words_ = []
        for w_l in words:
            words_.append("".join(w_l[::-1]))
        pieces.append((words_, c))

    return pieces


def adjust_gcs(word, constant, adjustable_gcs_=None):
    adjustable_gcs = (
        set(["x", "n", "5", "C", "D", "?"])
        if adjustable_gcs_ == None
        else adjustable_gcs_
    )

    gcs = set(list(word))

    # TODO: Do this only for a special case, otherwise it's going to be a problem
    # if gcs.issubset(set(["-", "1", "0", "u", "n", "x"])):
    #     adjustable_gcs.add("u")

    n = len(word)
    for i in range(n):
        gc = word[n - 1 - i]
        if gc in adjustable_gcs:
            constant += pow(2, i)

    # TODO: Reconsider it
    # if gcs.issubset(set(["x", "n", "u"])) and 'u' in adjustable_gcs:
    #     constant /= 2
    #     constant = int(constant)

    return constant % pow(2, n)


# derive with the assumption that we have ? with grounded bits
def derive_q(word, constant):
    new_constant = adjust_gcs(word, constant)

    n = len(word)
    staged_values = []
    staged_bits = []
    for i in range(n):
        gc = word[n - 1 - i]
        bit = new_constant >> i & 1
        clean_state = (
            True
            if (gc == "u" and bit == 1)
            or (gc == "n" and bit == 0)
            or (gc in ["-", "1", "0"] and bit == 0)
            else False
        )
        if not clean_state:
            staged_values.append(gc)
            staged_bits.append(bit)
    staged_values = staged_values[::-1]
    staged_bits = staged_bits[::-1]

    if set(staged_values) == set("?"):
        new_word = [
            "-"
            if gc == "?" and new_constant >> (n - 1 - i) & 1 == 0
            else "u"
            if gc == "?" and new_constant >> (n - 1 - i) & 1 == 1
            else gc
            for i, gc in enumerate(word)
        ]
        new_word = "".join(new_word)
    else:
        # TODO: Cover other cases
        new_word = word

    diff, err = _int_diff(new_word)
    assert not err and diff == constant
    return new_word, False


def derive_word(word, constant):
    problemetic_gcs = set(["7", "E"])

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
    n = len(word)
    gcs = set(list(word))

    # does it require derivation?
    if gcs.issubset(set(["u", "n", "1", "0", "-"])):
        return word, False

    # detect if derivable
    if not gcs.isdisjoint(problemetic_gcs):
        return None, True
    # if "x" in gcs and not gcs == set(["-", "x"]):
    #     return None, True

    # derive '?' if all other bits are grounded
    if gcs.issubset(set(["u", "n", "-", "0", "1", "?"])):
        return derive_q(word, constant)

    # perform manipulations to keep each bit from influencing others during diff. calc.
    new_constant = adjust_gcs(word, constant)

    # derive
    derived_word = list(word)
    for i in range(n):
        gc = word[n - 1 - i]
        bit = new_constant >> i & 1
        if gc in table:
            derived_word[n - 1 - i] = table[gc][bit]
    derived_word = "".join(derived_word)

    int_diff, err = _int_diff(derived_word)
    assert not err and int_diff == constant
    return derived_word, False


def derive_word_w(word, constant):
    pieces = split_words([word], constant)
    derivation = []
    for words, c in pieces:
        word = words[0]
        derived_piece, err = derive_word(word, c)
        assert not err
        derivation.append(derived_piece)
    derivation_str = "".join(derivation[::-1])
    diff, _ = _int_diff(derivation_str)
    assert diff == constant
    return derivation_str, False


# def derive_words(words, constant):
#     if len(words) != 2:
#         return [], True

#     n = len(words[0])
#     m = len(words)

#     all_grounded = all(
#         [
#             True if set(word).issubset(["n", "u", "1", "0", "-"]) else False
#             for word in words
#         ]
#     )
#     if all_grounded:
#         return words, False

#     # adjustments gcs
#     new_constant = constant
#     for word in words:
#         for i in range(n):
#             if word[n - i - 1] in ["x", "n", "5", "C", "D"]:
#                 new_constant += pow(2, i)

#     # spread out the variables
#     table = {
#         "D": [0, 1],
#         "x": [0, 2],
#         "B": [0, 1],
#         "-": [0],
#         "0": [0],
#         "1": [0],
#         "n": [0],
#         "u": [1],
#     }
#     variables = []
#     for i in range(m):
#         variables.append([None] * n)
#         for j in range(n):
#             gc = words[i][j]
#             if gc not in table:
#                 return words, constant
#             variables[i][j] = table[gc]

#     for i in range(len(variables)):
#         for j in range(n):
#             if j != n - 1:
#                 if variables[i][j] == [0]:
#                     for k in range(len(variables)):
#                         if variables[k][j + 1] == [0, 2]:
#                             variables[k][j + 1] = [0]
#                             variables[i][j] = [0, 1]
#             # if variables[i][j] == [0, 2]:
#             #     variables[i][j] = [0]
#             #     if j != n - 1:
#             #         if variables[i][j - 1] == [0]:
#             #             variables[i][j - 1] = [0, 1]
#             #         elif i > 0 and variables[i - 1][j - 1] == [0]:
#             #             variables[i - 1][j - 1] = [0]
#             #         elif i < len(variables) - 1 == [0]:
#             #             variables[i + 1][j - 1] = [0]

#     new_words = [list(words[i]) for i in range(m)]
#     for i in range(n):
#         sum = new_constant >> i & 1
#         # Left column
#         left, right = [None] * m, [None] * m
#         if i != n - 1:
#             left = (words[0][n - 2 - i], words[1][n - 2 - i])
#         # Right column
#         if i > 0:
#             right = (words[0][n - i], words[1][n - i])
#         current = (variables[0][n - 1 - i], variables[1][n - 1 - i])
#         is_last = True if i == n - 1 else False
#         gcs = set([words[j][n - 1 - i] for j in range(m)])
#         print(gcs, current, left, right, sum)
#         # if current[0] == [0, 1] and current[1] == [0, 1]:
#         #     continue
#         if gcs == set(["B", "-"]) and [0] in current and [0, 1] in current:
#             for j in range(m):
#                 if words[j][n - i - 1] == "B":
#                     new_words[j][n - i - 1] = "-" if sum == 0 else "u"
#         elif gcs == set(["B", "x"]) and [0] in current and [0, 1] in current:
#             for j in range(m):
#                 if words[j][n - i - 1] == "B":
#                     new_words[j][n - i - 1] = "-" if sum == 0 else "u"
#         elif gcs == set(["D", "x"]) and [0, 1] in current and [0, 1] in current:
#             for j in range(m):
#                 if words[j][n - i - 1] == "B":
#                     new_words[j][n - i - 1] = "-" if sum == 0 else "u"
#         elif gcs == set(["D", "-"]) and [0] in current and [0, 1] in current:
#             for j in range(m):
#                 if words[j][n - i - 1] == "D":
#                     new_words[j][n - i - 1] = "-" if sum == 1 else "n"
#         elif gcs == set(["D", "n"]) and [0] in current and [0, 1] in current:
#             for j in range(m):
#                 if words[j][n - i - 1] == "D":
#                     new_words[j][n - i - 1] = "-" if sum == 1 else "n"
#         elif (gcs == set(["-"]) or gcs == set(["-", "x"])) and [0] in current and [0, 1] in current and (right != [None, None]):
#             for j in range(m):
#                 if right[j] == "x":
#                     new_words[j][n - i] = "u" if sum == 1 else "n"
#         elif [0] in current and [0, 1] in current and (right != [None, None]):
#             for j in range(m):
#                 if right[j] == "x" and sum == 1:
#                     new_words[j][n - i] = "u"
#         elif [0] in current and not [0, 1] in current and (right != [None, None]):
#             for j in range(m):
#                 if right[j] == "B" and sum == 1:
#                     new_words[j][n - i] = "u"


#     print(words, constant, new_constant)
#     return ["".join(word) for word in new_words], False


def gen_variables(words, constant):
    n = len(words[0])
    vars = [[] for _ in range(n)]
    for word in words:
        for i in range(n - 1, -1, -1):
            gc = word[i]
            is_msb = i == 0
            if gc == "x" and not is_msb:
                vars[i - 1].append([0, 1])
            elif gc in ["?", "7", "E"]:
                vars[i].append([0, 1])
                if not is_msb:
                    vars[i - 1].append([0, 1])
            elif gc in [
                "3",
                "5",
                "A",
                "B",
                "C",
                "D",
            ]:
                vars[i].append([0, 1])

    # Deal with the constants and treat them as derived variables
    for i in range(n - 1, -1, -1):
        sum = 0
        for word in words:
            gc = word[i]
            is_msb = i == 0
            sum += 1 if word[i] == "u" else 0
        assert sum in [2, 1, 0]
        if sum == 2:
            vars[i - 1].append([0, 1])
        elif sum == 1:
            vars[i].append([sum])

    return vars


def is_congruent(a, b, m):
    return (a - b) % m == 0


def _does_overflow(stash):
    bits, vars_colwise = stash["bits"], stash["vars"]
    n = len(bits)
    value = sum([bit * pow(2, i) for i, bit in enumerate(bits)])
    max_value = sum([pow(2, i) for i in range(n)])

    vars, consts = [], []
    for order, col in enumerate(vars_colwise):
        for var in col:
            if var == [0, 1]:
                vars.append({"value": -1, "order": order})
            else:
                consts.append({"value": var[0], "order": order})

    m = pow(2, n)
    for i in range(pow(2, len(vars)), -1, -1):
        for j in range(len(vars)):
            vars[j]["value"] = i >> j & 1
        candidate_value = sum(
            [var["value"] * pow(2, var["order"]) for var in vars]
        ) + sum([const["value"] * pow(2, const["order"]) for const in consts])
        if candidate_value <= max_value:
            continue
        # TODO: Ensure that the candidate conforms to the constraints
        has_solution = is_congruent(candidate_value, value, m)
        if has_solution:
            return True
    return False


def brute_force(vars_colwise, constant, min_gt=-1):
    vars, consts = [], []
    for order, col in enumerate(vars_colwise):
        for var in col:
            if var == [0, 1]:
                vars.append({"value": -1, "order": order})
            else:
                consts.append({"value": var[0], "order": order})

    n = len(vars)
    solutions = []
    for i in range(pow(2, n)):
        values = [i >> j & 1 for j in range(n)]
        var_index = 0
        for j, value in enumerate(values):
            vars[j]["value"] = value
        sum_ = sum([var["value"] * pow(2, var["order"]) for var in vars]) + sum(
            [const["value"] * pow(2, const["order"]) for const in consts]
        )
        if min_gt == -1 and sum_ != constant:
            continue
        if min_gt != -1 and sum_ <= min_gt:
            continue
        solutions.append(values)

    pattern = [set() for _ in range(n)]
    for solution in solutions:
        for i, var in enumerate(solution):
            pattern[i].add(var)

    return pattern


def apply_grounding(words, vars_colwise, values):
    # print(values)
    # print(vars_colwise)
    # Remove constants from colwise vars
    new_vars_colwise = []
    for col in vars_colwise:
        new_col = []
        for var in col:
            if var == [0, 1]:
                new_col.append(var)
        new_vars_colwise.append(new_col)

    n = len(words[0])
    var_index = 0
    words_ = [list(word) for word in words]
    for i in range(n - 1, -1, -1):
        current_col = new_vars_colwise[i]
        next_col = new_vars_colwise[i - 1] if i != 0 else []
        new_var_index = var_index + len(current_col)
        # print(i, len(next_col), next_col, new_var_index)
        next_col_values = [values[k + new_var_index] for k in range(len(next_col))]
        current_col_values = [values[k + var_index] for k in range(len(current_col))]
        # print(new_var_index)
        for j, word in enumerate(words):
            gc = word[i]
            next_col_uniq_values = list(set(next_col_values))
            current_col_uniq_values = list(set(current_col_values))
            if (
                gc == "x"
                and len(next_col_uniq_values) == 1
                and next_col_uniq_values[0] != -1
            ):
                words_[j][i] = "u" if next_col_uniq_values[0] == 1 else "n"
            elif (
                gc in ["3", "5", "A", "B", "C", "D"]
                and len(current_col_uniq_values) == 1
                and current_col_uniq_values[0] != -1
            ):
                if gc == "3":
                    words_[j][i] = "0" if current_col_uniq_values[0] == 0 else "u"
                elif gc == "5":
                    words_[j][i] = "n" if current_col_uniq_values[0] == 0 else "0"
                elif gc == "A":
                    words_[j][i] = "1" if current_col_uniq_values[0] == 0 else "u"
                elif gc == "B":
                    words_[j][i] = "u" if current_col_uniq_values[0] == 1 else "-"
                elif gc == "C":
                    words_[j][i] = "n" if current_col_uniq_values[0] == 0 else "1"
                elif gc == "D":
                    words_[j][i] = "n" if current_col_uniq_values[0] == 0 else "-"
            # Handle 7, E and ?
            elif gc in ["7", "E", "?"]:
                if (
                    len(current_col_uniq_values) == 1
                    and current_col_uniq_values[0] != -1
                ):
                    words_[j][i] = (
                        ("-" if gc == "?" else "0" if gc == "7" else "1")
                        if current_col_uniq_values[0] == 1
                        else "n"
                    )
                elif (
                    len(next_col_uniq_values) == 1
                    and next_col_uniq_values[0] != -1
                    and next_col_uniq_values[0] == 1
                ):
                    words_[j][i] = "u"
            # print(current_col_uniq_values, next_col_uniq_values)

        var_index = new_var_index

    derived_words = ["".join(word_) for word_ in words_]
    # for w in derived_words:
    #     print(w)
    return derived_words


def derive_words(words, adj_constant):
    derived_words, err = words, False

    # Checks
    n = len(words[0])
    m = len(words)

    # Skip if all words are grounded
    all_grounded = all(
        [
            True if set(word).issubset(["n", "u", "1", "0", "-"]) else False
            for word in words
        ]
    )
    if all_grounded:
        return words, False

    # Generate variables
    vars_colwise = gen_variables(words, adj_constant)

    # print("Adjustment", words, format(adj_constant, "0b"))
    # for entry in vars_colwise[::-1]:
    #     print(entry)
    # print("End")

    # Linear scan
    segments = []
    stash = {"vars": [], "bits": []}
    last_does_overflow = False
    for i in range(n - 1, -1, -1):
        # Skip if there's nothing in the stash and there's no variable either
        if len(stash["vars"]) == 0 and len(vars_colwise[i]) == 0:
            continue

        bit = adj_constant >> (n - i - 1) & 1
        stash["vars"].append(vars_colwise[i])
        stash["bits"].append(bit)

        # print(stash["vars"], stash["bits"])

        # Check if it overflows
        does_overflow = _does_overflow(stash)
        # print(i, bit, vars_colwise[i], does_overflow)
        segment_ends = False
        # if (
        #     does_overflow
        #     and i != 0
        #     and (
        #         len(vars_colwise[i - 1]) == 0
        #         or (len(vars_colwise[i - 1]) == 1 and vars_colwise[i - 1][0] == [1])
        #     )
        # ):
        #     if len(vars_colwise[i - 1]) != 0:
        #         stash["vars"].append(vars_colwise[i - 1])
        #     stash["bits"].append(adj_constant >> (n - i) & 1)
        #     # TODO: Optional: Set the variables since we know them now
        #     segment_ends = True

        # If it doesn't overflow, cut it off
        if not does_overflow:
            segment_ends = True

        # Flush the stash
        if i == 0:
            segment_ends = True
            # If it can overflow, brute force should search for solutions with sum > highest possible value
            if does_overflow:
                last_does_overflow = True

        if segment_ends:
            segments.append((stash["vars"], stash["bits"]))
            stash["bits"], stash["vars"] = [], []

    # Derive the bits
    vars_values = {}
    var_index = 0
    for s_i, (vars_colwise_, bits) in enumerate(segments):
        sum_ = sum([bit * pow(2, i) for i, bit in enumerate(bits)])
        is_last = s_i == len(segments) - 1

        vars_count, cons_count = 0, 0
        for col in vars_colwise_:
            for var in col:
                if var == [0, 1]:
                    vars_count += 1
                else:
                    cons_count += 1

        # print(vars_colwise_, sum_)

        # If the variable doesn't require any search
        # if vars_count == 1:
        #     # TODO: Generalize the process
        #     if sum_ == 0 and cons_count == 2:
        #         vars_values[var_index] = 1
        #     elif sum_ == 1 and cons_count == 1:
        #         vars_values[var_index] = 0
        #     elif sum_ == 1 and cons_count == 0:
        #         vars_values[var_index] = 1
        #     elif sum_ == 2 and cons_count == 1:
        #         vars_values[var_index] = 1
        #     else:
        #         vars_values[var_index] = sum_

        #     var_index += 1
        # else:
        if is_last and last_does_overflow:
            min_gt = sum([pow(2, i) for i in range(len(bits))])
            propagated_vars = brute_force(vars_colwise_, -1, min_gt=min_gt)
        else:
            propagated_vars = brute_force(vars_colwise_, sum_)
        local_index = 0
        for i, col in enumerate(vars_colwise_):
            for var in col:
                # Ignore constants
                if var != [0, 1]:
                    continue
                value_ = propagated_vars[local_index]
                value = list(value_)[0] if len(value_) == 1 else -1
                vars_values[var_index] = value
                local_index += 1
                var_index += 1
    # print(len(vars_colwise), len(words[0]), words, adj_constant)
    # print(vars_values)

    # Update the words with the derived bits
    derived_words = apply_grounding(words, vars_colwise, vars_values)

    return derived_words, err


def derive_words_w(words, constant):
    adj_constant = constant
    for word in words:
        adj_constant = adjust_gcs(word, adj_constant)

    pieces = split_words(words, adj_constant)
    derivation = []
    for _ in range(len(words)):
        derivation.append([])
    for words, c in pieces:
        new_words, err = derive_words(words, c)
        assert not err
        # print("End of the piece")
        # print("New words", new_words)
        for i, word in enumerate(new_words):
            derivation[i] = list(word) + derivation[i]
    diff_sum = 0
    for i in range(len(derivation)):
        derivation[i] = "".join(derivation[i])
        diff, _ = _int_diff(derivation[i])
        diff_sum = (diff_sum + diff) % pow(2, len(derivation[i]))
    return derivation, False


# Tests

# Adjustment of GCs
assert 2196 == adjust_gcs("--xxxx-xx--x", adjust_gcs("--B--D-BBBB-", 1147))

# Congruency
assert True == is_congruent(2, 0, 2)
assert False == is_congruent(2, 1, 2)
assert True == is_congruent(16, 2, 2)
assert False == is_congruent(16, 1, 16)

# Overflow detector
assert True == _does_overflow({"vars": [[[0, 1], [0, 1]]], "bits": [0]})
assert False == _does_overflow({"vars": [[[0, 1], [0, 1]], [[0, 1]]], "bits": [0, 1]})
assert True == _does_overflow({"vars": [[[0, 1], [0, 1]], [[0, 1]]], "bits": [0, 0]})
assert False == _does_overflow({"vars": [[[0, 1]]], "bits": [0]})
assert True == _does_overflow(
    {"vars": [[[0, 1], [0, 1]], [[0, 1]], [[0, 1]]], "bits": [0, 0, 0]}
)
assert True == _does_overflow(
    {"vars": [[[0, 1], [0, 1]], [[0, 1]], [[0, 1], [0, 1]]], "bits": [0, 0, 1]}
)
# assert False == _does_overflow(
#     {"vars": [[[0, 1], [0, 1]], [[0, 1]], [[0, 1], [0, 1]]], "bits": [0, 1, 0]}
# )
assert True == _does_overflow(
    {
        "vars": [[[0, 1], [0, 1]], [[0, 1]], [[0, 1], [0, 1]], [[0, 1], [0, 1]]],
        "bits": [0, 0, 1, 1],
    }
)
assert True == _does_overflow(
    {"vars": [[[0, 1], [0, 1]], [[0, 1], [0, 1]], [[0, 1], [1]]], "bits": [0, 0, 0]}
)
assert False == _does_overflow(
    {"vars": [[[0, 1], [0, 1]], [[0, 1], [0, 1]], [[0, 1], [1]]], "bits": [0, 0, 0, 1]}
)

# Brute force
unknown_var = set([0, 1])
one_var = set([1])
assert brute_force([[[0, 1], [0, 1]], [[0, 1]]], 2) == [
    unknown_var,
    unknown_var,
    unknown_var,
]
assert brute_force([[[0, 1], [0, 1]], [[0, 1]]], 4) == [one_var, one_var, one_var]
assert brute_force([[[0, 1], [0, 1]], [[0, 1], [0, 1]]], 0, min_gt=3) == [
    unknown_var,
    unknown_var,
    unknown_var,
    unknown_var,
]
assert brute_force([[[0, 1], [0, 1]], [[0, 1], [0, 1]], [[0, 1], [1]]], 8) == [
    unknown_var,
    unknown_var,
    unknown_var,
    unknown_var,
    unknown_var,
]

# Derive words
# assert (["--uunu-nx--x", "--u--n-B-BB-"], False) == derive_words_w(
#     ["--xxxx-xx--x", "--B--D-BBBB-"], 1147
# )
