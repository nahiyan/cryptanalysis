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

    n = len(word)
    for i in range(n):
        gc = word[n - 1 - i]
        if gc in adjustable_gcs:
            constant += pow(2, i)

    return constant % pow(2, n)


def gen_variables(words):
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
            vars[i - 1].append([1])
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
                elif len(next_col_uniq_values) == 1 and next_col_uniq_values[0] == 1:
                    words_[j][i] = "u"
            # print(current_col_uniq_values, next_col_uniq_values)

        var_index = new_var_index

    derived_words = ["".join(word_) for word_ in words_]
    # for w in derived_words:
    #     print(w)
    return derived_words


def derive_words(words, adj_constant, adjusted=True):
    derived_words, err = words, False
    n, m = len(words[0]), len(words)

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
    vars_colwise = gen_variables(words)
    vars_count = 0
    for col in vars_colwise:
        for c in col:
            if c != [0, 1]:
                continue
            vars_count += 1

    if not adjusted:
        for word in words:
            adj_constant = adjust_gcs(word, adj_constant)

    # print("Adjustment", words, format(adj_constant, "0b"))
    # for entry in vars_colwise[::-1]:
    #     print(entry)
    # print("End")

    # Linear scan
    segments = []
    stash = {"vars": [], "bits": []}
    # last_does_overflow = False
    overflow_brute_force_indices = []
    for i in range(n - 1, -1, -1):
        # Skip if there's nothing in the stash and there's no variable either
        if len(stash["vars"]) == 0 and len(vars_colwise[i]) == 0:
            continue

        bit = adj_constant >> (n - i - 1) & 1
        stash["vars"].append(vars_colwise[i])
        stash["bits"].append(bit)

        # print(stash["vars"], stash["bits"])

        # if vars_count > 20:
        #     break

        # Check if it overflows
        does_overflow = _does_overflow(stash)
        # print(i, bit, vars_colwise[i], does_overflow)
        segment_ends = False

        if bit == 0 and all([word[i] in ["1", "0", "-"] for word in words]):
            segment_ends = True
            if does_overflow:
                overflow_brute_force_indices.append(len(segments))

        # If it doesn't overflow, cut it off
        if not does_overflow:
            segment_ends = True

        # Flush the stash
        if i == 0:
            segment_ends = True
            # If it can overflow, brute force should search for solutions with sum > highest possible value
            if does_overflow:
                if does_overflow:
                    overflow_brute_force_indices.append(len(segments))
                # last_does_overflow = True

        if segment_ends:
            segments.append((stash["vars"], stash["bits"]))
            stash["bits"], stash["vars"] = [], []

    # Derive the bits
    vars_values = {}
    var_index = 0
    for s_i, (vars_colwise_, bits) in enumerate(segments):
        sum_ = sum([bit * pow(2, i) for i, bit in enumerate(bits)])
        if s_i in overflow_brute_force_indices:
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

    # Fill in missing values
    for i in range(len(vars_values), vars_count):
        vars_values[var_index] = -1
        var_index += 1;
    assert (len(vars_values) == vars_count);

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

# Generation of variables
assert gen_variables(["-uxxu-xx1u---x-00x", "--?0-?0--u?A-???5-"]) == [
    [],
    [[0, 1], [0, 1], [1]],
    [[0, 1], [0, 1]],
    [],
    [[0, 1], [1]],
    [[0, 1], [0, 1]],
    [[0, 1]],
    [],
    [[1]],
    [[0, 1]],
    [[0, 1]],
    [[0, 1]],
    [[0, 1], [0, 1]],
    [[0, 1], [0, 1]],
    [[0, 1], [0, 1]],
    [[0, 1]],
    [[0, 1], [0, 1]],
    [],
]

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
assert (["--uunu-nx--x", "--u--n-B-BB-"], False) == derive_words_w(
    ["--xxxx-xx--x", "--B--D-BBBB-"], 1147
)
