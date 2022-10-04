import sys

instances = {}
sat_solvers = []

print("instance ", end="")

for line in sys.stdin:
    if line.startswith("SAT Solver:"):
        sat_solver = line.split(":")[-1].strip()
        print(sat_solver, end=" ")
        sat_solvers.append(sat_solver)
    elif line.startswith("Time:"):
        segments = line.split(" ")
        time = segments[1][:-2]
        instance_index = segments[4].strip()
        # print(sat_solver, time, "-", instance_index)
        if instance_index not in instances:
            instances[instance_index] = []
        instances[instance_index].append(time)

print()

i = 0
while True:
    if str(i) in instances:
        print("instance" + str(i + 1), end=" ")
        for time in instances[str(i)]:
            print(time, end=" ")
        print()
    else:
        break
    i += 1
