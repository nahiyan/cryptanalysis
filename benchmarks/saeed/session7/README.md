# Overview

Saeed's encoder with espresso adder has been used to generate the encoding, which has then been simplified using CaDiCaL and then pre-processed with SatELite before being fed into March for producing the cubes.

1000 random cubes were tried, and the ones with the following indices have been found to be satisfiable:

- 136538
- 124272
- 30054
- 93000
- 16161

11 of the problems were found to be UNSAT, while the others timed out.

# Simplification

CaDiCaL has been used for simplifying the instance with 8 passes, each for 100 seconds, while SatELite was used for pre-processing after the in-processing. The reason why it's done so is because even after simplification with CaDiCaL, the problem was still too hard to find SAT solutions.

# Cubes

March has been provided with the simplified instance and ran for generating the cubeset with a cutoff vars threshold of 3220, producing 253,419 cubes, including 2,226 refuted leaves.

