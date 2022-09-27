import os

# Goal 1: Generate the encoding with and without XOR clauses
# Goal 2: Try variants with the adders: counter_chain and dot_matrix
# Goal 3: Try various rounds from 16-48
# TODO: Goal 4: Try various hashes for the inversion
# TODO: Add options to pick the encoder: Saeed, Oleg, etc.

xor_options = [True, False]

# TODO: Implement passing the hashes to the encoder
hashes = ["f f f f f f f f f f f f f f f f f f f f f f f f f f f f f f f f",
          "0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0"]
adder_types = ["counter_chain", "dot_matrix"]
step_variations = list(range(16, 49))

for xor_option in xor_options:
    xor_flag = xor_option
    for adder_type in adder_types:
        for steps in step_variations:
            os.system("./encoders/ {} -A {} -r {} -f md4 -a preimage > encodings/saeed/md4_{}_{}_xor{}.cnf".format(
                "--xor" if xor_flag else None,
                adder_type,
                steps,
                steps,
                adder_type,
                "1" if xor_flag == True else "0"
            ))
