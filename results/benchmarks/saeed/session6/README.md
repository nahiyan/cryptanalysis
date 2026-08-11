# Overview

Saeed's encoder with espresso adder has been used to generate the encoding, which has then been simplified using CaDiCaL before being fed into March for producing cubes.

The cubes with the following indices have been successfully solved:

- 22252
- 16486
- 35193
- 15607

The others timed out or were found unsatisfiable.

# Simplification

CaDiCaL has been used for simplifying the instance with 8 passes, each for 100 seconds.

# Cubes

March has been provided with the simplified instance and ran for generating the cubeset with a cutoff vars threshold of 2570, producing 35,277 cubes, including 630 refuted leaves. Out of the 35,277 sub-problems in total, only 8,730 could be attempted, of which 229 resulted in UNSAT, while 4 yielded SAT. The rest 8,497 timed out.

# Slurm Statistics

```
State: TIMEOUT (exit code 0)
Nodes: 2
Cores per node: 40
CPU Utilized: 39-21:52:25
CPU Efficiency: 49.89% of 80-00:02:40 core-walltime
Job Wall-clock time: 1-00:00:02
Memory Utilized: 7.62 GB
Memory Efficiency: 11.91% of 64.00 GB
```