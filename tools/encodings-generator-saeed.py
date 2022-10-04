import os
import sys

# Goal 1: Generate the encoding with and without XOR clauses
# Goal 2: Try variants with the adders: counter_chain and dot_matrix
# Goal 3: Try various rounds from 16-48
# Goal 4: Try various hashes for the inversion
# TODO: Add options to pick the encoder: Saeed, Oleg, etc.

xor_options = [True, False]

hashes = ["ffffffffffffffffffffffffffffffff",
          "00000000000000000000000000000000"]
adder_types = ["counter_chain", "dot_matrix"]
step_variations = list(range(16, 49))

# Check if the executable exists
encoder_path = "encoders/saeed/crypto/main"
if not os.path.exists(encoder_path):
    sys.exit("Failed to find the encoder in the 'encoders/saeed/crypto' directory. Can you ensure that you compiled it?")

for hash in hashes:
    for xor_option in xor_options:
        xor_flag = "--xor" if xor_flag else None
        for adder_type in adder_types:
            for steps in step_variations:
                os.system("{} {} -A {} -r {} -f md4 -a preimage -t {} > encodings/saeed/md4_{}_{}_xor{}_{}.cnf".format(
                    encoder_path,
                    xor_flag,
                    adder_type,
                    steps,
                    hash,
                    steps,
                    adder_type,
                    "1" if xor_flag == True else "0",
                    hash
                ))
