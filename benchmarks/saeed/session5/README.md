# Configuration

- Solvers: Kissat
- Steps: 37
- Dobbertin's Attack: 1
- Dobbertin Relaxation Bits: 32
- Cubes: Enabled
- Cube Number of Vars (Cutoff Threshold): 3910 
- Max Time: 10000
- Hashes: 00000000000000000000000000000000
- Adders: Dot matrix
- XOR: 0

# Slurm Statistics

```
Job ID: 32773807
Cluster: beluga
User/Group: nahiyan0/nahiyan0
State: TIMEOUT (exit code 0)
Nodes: 1
Cores per node: 32
CPU Utilized: 3-16:57:49
CPU Efficiency: 99.81% of 3-17:07:44 core-walltime
Job Wall-clock time: 02:47:07
Memory Utilized: 5.08 GB
Memory Efficiency: 63.50% of 8.00 GB
```

# Notes

Previously, Kissat failed to solve a md4-37 instance with the same configuration within 5000 seconds. Cubing not only let it do so, but also score the best time at ~53 seconds! Although, there were thousands of cubes, only around 202 of them were attempted to be solved within that time. Since the integration of cubing into the slurm version isn't tested, the regular method is used with 200 max. concurrent instances.
