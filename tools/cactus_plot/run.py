import sys

instances = {}
sat_solvers = []

print("instance ", end="")

for line in sys.stdin:
    segments = line.split(" ")

    sat_solver = segments[10][:-1]
    if sat_solver not in sat_solvers:
        sat_solvers.append(sat_solver)
        print(sat_solver, end=" ")

    sat_solver_index = sat_solvers.index(sat_solver)

    time = segments[1][:-2]
    instance_index = int(segments[4][:-1].strip())
    if instance_index not in instances:
        instances[instance_index] = {}
    instances[instance_index][sat_solver_index] = time

print()

i = 0
while True:
    if i in instances:
        print("instance" + str(i + 1), end=" ")
        for j in range(len(sat_solvers)):
            time = instances[i][j]
            print(time, end=" ")
        print()
    else:
        break
    i += 1
