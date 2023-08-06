import sys

args = sys.argv
log_file = ""
solution_file = ""
if len(args) != 3:
    print("Invalid args, use <log_file> <solution_file>")
log_file = args[1]
solution_file = args[2]

model = []
with open(solution_file) as file:
    for line in file:
        segments = line.split()
        model.extend([1 if int(x) > 0 else 0 for x in segments])

valid = True
clauses_n = 0
with open(log_file) as file:
    for line in file:
        if not line.startswith("Clause: "):
            continue
        clauses_n += 1
        line = line.removeprefix("Clause: ")

        segments = line.split()
        clause = [int(x) for x in segments]
        falsified = True
        for item in clause:
            model_value = model[abs(item) - 1]
            item_value = 1 if item > 0 else 0
            if model_value == item_value:
                falsified = False
                break
        if falsified:
            print(" ".join([str(x) for x in clause]))
            valid = False

if valid:
    print("Valid;", clauses_n, "clauses checked")
