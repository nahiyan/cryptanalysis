from z3 import *
from derive import _int_diff


def derive_words(word_x, word_y, constant):
    table = {
        "x": [1, -1],
        "B": [0, 1],
        "D": [0, -1],
        "-": [0],
        "u": [1],
        "n": [-1],
        "?": [-1, 0, 1],
        "A": [-1, 0, 1],
        "5": [1, -1],
        "0": [0],
        "1": [0],
    }

    words = (word_x, word_y)
    n = len(word_x)
    vars = z3.Ints(" ".join([str(x) for x in range(n * 2)]))
    s = Solver()
    for i in range(n):
        bit_x, bit_y = word_x[n - 1 - i], word_y[n - 1 - i]
        s.add(Or([vars[i] == x for x in table[bit_x]]))
        s.add(Or([vars[n + i] == x for x in table[bit_y]]))
    s.add(
        sum(
            [
                (
                    (0 if word_x[n - 1 - i] in ["-", "1", "0"] else vars[i])
                    + (0 if word_y[n - 1 - i] in ["-", "1", "0"] else vars[n + i])
                )
                * pow(2, i)
                for i in range(n)
            ]
        )
        % pow(2, n)
        == constant
    )
    solutions = []
    while True:
        result = s.check()
        if result == unsat:
            if len(solutions) == 0:
                print("Failed")

            # check consistency
            derived_words = ([None] * n, [None] * n)
            matrix = {}
            for i in range(2):
                matrix[i] = {}
                for j in range(n):
                    matrix[i][j] = []
                    for solution in solutions:
                        matrix[i][j].append(solution[i][j])

                        # Sanity check
                        (int_diff1, err1), (int_diff2, err2) = _int_diff(
                            "".join(solution[0]), n=n
                        ), _int_diff("".join(solution[1]), n=n)
                        assert (not err1 and not err2) and (
                            (int_diff1 + int_diff2) % pow(2, n)
                        ) == constant

                    gcs = set(matrix[i][j])
                    derived_words[i][j] = list(gcs)[0] if len(gcs) == 1 else words[i][j]

            # Remove any loss of GCs with diff. of 0
            for i, derived_word in enumerate(derived_words):
                for j in range(n):
                    if derived_word[j] == "-" and words[i][j] in ["1", "0"]:
                        derived_words[i][j] = words[i][j]

            # Return the result
            return "".join(derived_words[0]), "".join(derived_words[1])
        model_ = s.model()
        model = {}
        for var in model_:
            model[int(var.name())] = model_[var].as_long()

        solution = ([], [])
        derive_gc = lambda value: "u" if value == 1 else "n" if value == -1 else "-"
        for i in range(n - 1, -1, -1):
            value = model[i]
            solution[0].append(derive_gc(value))
        for i in range(2 * n - 1, n - 1, -1):
            value = model[i]
            solution[1].append(derive_gc(value))
        solution = ("".join(solution[0]), "".join(solution[1]))
        solutions.append(solution)
        print(solution[0], solution[1], sep="\n")

        # block it
        s.add(Or([vars[var] != model_[vars[var]] for var in model]))


def derive_words_new(words, constant, n=32):
    assert len(words) > 0
    n = len(words[0])
    table = {
        "?": [(0, 0), (0, 1), (1, 0), (1, 1)],
        "-": [(0, 0), (1, 1)],
        "x": [(0, 1), (1, 0)],
        "0": [(0, 0)],
        "u": [(1, 0)],
        "n": [(0, 1)],
        "1": [(1, 0)],
        "3": [(0, 0), (1, 0)],
        "5": [(0, 0), (0, 1)],
        "7": [(0, 0), (1, 0), (0, 1)],
        "A": [(1, 0), (1, 1)],
        "B": [(0, 0), (1, 0), (1, 1)],
        "C": [(0, 1), (1, 1)],
        "D": [(0, 0), (0, 1), (1, 1)],
        "E": [(1, 0), (0, 1), (1, 1)],
    }

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

    print(f"Constant: {constant}")
    m = len(words)
    bitvec_pairs = [
        (BitVec(f"w{i}_a", n), BitVec(f"w{i}_b", n)) for i, _ in enumerate(words)
    ]
    addends = [a - b for a, b in bitvec_pairs]
    sum_ = sum(addends)

    s = Solver()
    for i, word in enumerate(words):
        for j in range(n):
            gc = word[n - j - 1]
            vec_x, vec_y = bitvec_pairs[i]
            bit_f = Extract(j, j, vec_x)
            bit_g = Extract(j, j, vec_y)
            assert gc in table
            if gc == "-":
                s.add(bit_f == bit_g)
                continue
            elif gc == "x":
                s.add(bit_f != bit_g)
                continue
            possibilities = table[gc]
            if len(possibilities) == 1:
                p_f, p_g = possibilities[0]
                s.add(And(bit_f == p_f, bit_g == p_g))
            else:
                s.add(
                    Or(
                        [
                            And([bit_f == p_f, bit_g == p_g])
                            for p_f, p_g in possibilities
                        ]
                    )
                )

    s.add(sum_ == constant)

    solutions = []
    while s.check() == sat:
        model = s.model()
        # print(model)
        solutions.append(model)

        # Block it
        block = []
        for var in model:
            block.append(var() != model[var])
        s.add(Or(block))

        block = []
        for x, y in bitvec_pairs:
            x_value, y_value = model[x].as_long(), model[y].as_long()
            d = x_value - y_value
            block.append(x - y != d)
        s.add(Or(block))
    else:
        cases = [[set() for _ in range(n)] for _ in range(m)]
        for solution in solutions:
            for word_index, (x, y) in enumerate(bitvec_pairs):
                x_value, y_value = solution[x].as_long(), solution[y].as_long()
                word = []
                for i in range(n - 1, -1, -1):
                    pair = (x_value >> i & 1, y_value >> i & 1)
                    word.append(
                        "0"
                        if pair == (0, 0)
                        else "1"
                        if pair == (1, 1)
                        else "u"
                        if pair == (1, 0)
                        else "n"
                    )
                for i, gc in enumerate(word):
                    cases[word_index][i].add(gc)
                # word = "".join(word)
                # print(word)
        words = []
        for case in cases:
            word = []
            for i, gcs in enumerate(case):
                for symbol in symbols:
                    if gcs == set(symbols[symbol]):
                        word.append(symbol)
            word = "".join(word)
            words.append(word)

        print("Fail" if len(solutions) == 0 else f"Done: {len(solutions)}")

        return words


if __name__ == "__main__":
    # words = derive_words_new(["-uxxxx-xxx-", "--DD-BBBB--"], 0b00011111110)
    # print(words[0], words[1], sep="\n")
    # words = derive_words_new(["xxxxxxx----x-xxxxx-x-", "DDDDD-D-nn--B-----D-B"], 0b001010110011100010111)
    # print(words[0], words[1], sep="\n")
    words = derive_words_new(
        ["--xxxx-xx--x-", "--B--D-BBBB--"], 0b0100011110110
    )
    print(words[0], words[1], sep="\n")
    words = derive_words_new(
        ["----x-xx---x--x--xx", "DDDD-DD-nn-Du-----D"], 0b1000000000110101110
    )
    print(words[0], words[1], sep="\n")
