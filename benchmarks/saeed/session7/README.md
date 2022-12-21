# Overview

Saeed's encoder with espresso adder has been used to generate the encoding, which has then been pre-processed with SatELite before fed into March for cubing.

1000 random cubes were tried (seed = 1), and the ones with the following indices have been found to be satisfiable:

- 136538
- 124272
- 30054
- 93000
- 16161

11 of the problems were found to be UNSAT, while the others timed out. The result is in `log.csv`.

However, with a seed value of 2, 2 SAT solutions could be found with 11 UNSAT ones.

# Cubes

March has been provided with the simplified instance and ran for generating the cubeset with a cutoff vars threshold of 3220, producing 253,419 cubes, including 2,226 refuted leaves.

