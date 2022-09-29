import os
import time

# Variations
xor_options = [True, False]
hashes = ["ffffffffffffffffffffffffffffffff",
          "00000000000000000000000000000000"]
adder_types = ["counter_chain", "dot_matrix"]
step_variations = list(range(16, 49))

# sat_solvers = ["cryptominisat", "kissat", "cadical", "glucose", "maplesat"]
sat_solvers = ["cryptominisat"]

# Constants
CRYPTOMINISAT = "cryptominisat"
KISSAT = "kissat"
CADICAL = "cadical"
GLUCOSE = "glucose"
MAPLESAT = "maplesat"


# Max time (seconds)
MAX_TIME = 5000


def cryptominisat(filepath, max_time):
    solution_file_path = filepath[-3] + ".sol"
    os.system(
        "cryptominisat5 -t 16 --verb 0 --maxtime {} {} > solutions/{}.sol".format(max_time, filepath, solution_file_path))


for sat_solver in sat_solvers:
    # Record the start time
    start_time = time.time()

    # Numbner of instances solved for this specific SAT solver
    i = 0
    for hash in hashes:
        for xor_option in xor_options:
            xor_flag = xor_option
            for adder_type in adder_types:
                for steps in step_variations:
                    # Invoke the SAT solver
                    filepath = "encodings/saeed/md4_{}_{}_xor{}_{}.cnf".format(
                        steps, adder_type, xor_flag, hash)
                    match sat_solver:
                        case CRYPTOMINISAT:
                            cryptominisat(filepath, MAX_TIME)

                    i += 1
    print("{} solved {} instances in {} seconds",
          sat_solver, i, time.time() - start_time)
