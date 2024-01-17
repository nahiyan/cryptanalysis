import sys
import subprocess


def get_solution(file_path):
    try:
        # Run the grep command and capture the output
        result = subprocess.run(
            ["grep", "^v.*", file_path], stdout=subprocess.PIPE, text=True, check=True
        )

        # Print the output
        return result.stdout

    except subprocess.CalledProcessError as e:
        # Handle errors if the grep command fails
        print(f"Error: {e}")


def get_clauses(file_path):
    try:
        # Run the grep command and capture the output
        result = subprocess.run(
            ["grep", "-E", "^Antecedent.*|^Asked for reason of.*", file_path],
            stdout=subprocess.PIPE,
            text=True,
            check=True,
        )

        # Print the output
        return result.stdout

    except subprocess.CalledProcessError as e:
        # Handle errors if the grep command fails
        print(f"Error: {e}")


sol_log_path = sys.argv[1]
solver_log_path = sys.argv[2]

assignment = {}
lines = get_solution(sol_log_path)[:-3].split("\n")
for line in lines:
    lits = line[2:].split(" ")
    for lit_ in lits:
        if lit_[:1] == "v" or len(lit_.strip()) == 0:
            continue
        lit = int(lit_)
        assert lit != 0
        assignment[abs(lit)] = True if lit > 0 else False

antecedents_count = 0
p_lit = 0
lines = get_clauses(solver_log_path).split("\n")
for line in lines:
    if line.startswith("Antecedent:"):
        if p_lit == 0:
            continue
        antecedents_count += 1
        # assert p_lit != 0
        antecedent = line[12:].strip()
        reason_clause = [p_lit]
        lits = antecedent.split(" ")
        for lit in lits:
            reason_clause.append(int(lit))

        is_satisfied = False
        for lit in reason_clause:
            var = abs(lit)
            value = True if lit > 0 else False
            if assignment[var] == value:
                is_satisfied = True
        if not is_satisfied:
            print("Clause: ", end="")
            for lit in reason_clause:
                print(lit, end=" ")
            print()
        p_lit = 0
    elif line.startswith("Asked for reason of"):
        p_lit = int(line.split(" ")[4])

print(antecedents_count, "reason clauses")
