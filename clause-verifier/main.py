import sys
import subprocess


def get_solution(file_path):
    assignment = {}

    try:
        # Run the grep command and capture the output
        result = subprocess.run(
            ["grep", "^v.*", file_path], stdout=subprocess.PIPE, text=True, check=True
        )

        lines = result.stdout[:-3].split("\n")
        for line in lines:
            lits = line[2:].split(" ")
            for lit_ in lits:
                if lit_[:1] == "v" or len(lit_.strip()) == 0:
                    continue
                lit = int(lit_)
                assert lit != 0
                assignment[abs(lit)] = True if lit > 0 else False

        return assignment

    except subprocess.CalledProcessError as e:
        # Handle errors if the grep command fails
        print(f"Error: {e}")


def get_clauses(file_path):
    try:
        # Run the grep command and capture the output
        result = subprocess.run(
            [
                "grep",
                "-E",
                "^Reason clause:.*|^Blocking clause:.*",
                file_path,
            ],
            stdout=subprocess.PIPE,
            text=True,
            check=True,
        )

        lines = result.stdout.split("\n")
        clauses = []
        for line_ in lines:
            clause_type = "reason" if line_.startswith("Reason clause:") else "blocking" if line_.startswith("Blocking clause:") else ""
            if clause_type == "":
                continue
            line = line_.replace("Reason clause: ", "").replace("Blocking clause: ", "").strip()
            lits = [int(lit_) for lit_ in line.split(" ")]
            clause = {"lits": lits, "type": clause_type}
            clauses.append(clause)

        # Print the output
        return clauses

    except subprocess.CalledProcessError as e:
        # Handle errors if the grep command fails
        print(f"Error: {e}")
        return []

def is_satisfied(clause, assignment):
    for lit in clause:
        var = abs(lit)
        value = True if lit > 0 else False
        if assignment[var] == value:
            return True
    return False

def print_clause(clause, assignment):
    print("Reason: " if clause["type"] == "reason" else "Blocking: ", end="")
    for lit in clause["lits"]:
        var = abs(lit)
        value = True if lit > 0 else False
        print(str(lit) + "(" + ("✓" if assignment[var] == value else "𐄂") + ")", end=" ")
    print()

def verify_clauses(sol_log_path, solver_log_path):
    assignment = get_solution(sol_log_path)
    clauses = get_clauses(solver_log_path)
    stats = {"reason": 0, "blocking": 0}
    for clause in clauses:
        lits = clause["lits"]
        type_ = clause["type"]
        stats[type_] += 1
        if type_ == "reason":
            if not is_satisfied(lits, assignment):
                print_clause(clause, assignment)
        elif type_ == "blocking":
            if not is_satisfied(lits, assignment):
                print_clause(clause, assignment)
    print(stats)


sol_log_path = sys.argv[1]
solver_log_path = sys.argv[2]
verify_clauses(sol_log_path, solver_log_path)
