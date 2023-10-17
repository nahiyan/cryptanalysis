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
                    gcs = set(matrix[i][j])
                    derived_words[i][j] = list(gcs)[0] if len(gcs) == 1 else words[i][j]
            # Sanity check
            (int_diff1, err1), (int_diff2, err2) = _int_diff(
                derived_words[0], n=n
            ), _int_diff(derived_words[1], n=n)
            assert (not err1 and not err2) and (int_diff1 + int_diff2) == constant
            
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


if __name__ == "__main__":
    word_x, word_y = derive_words(
        "0nun1x--ux--n--nxxu0x--ux-uuxu0-", ["0"] * 32, 3434604988
    )
    print(word_x, word_y, sep="\n")
