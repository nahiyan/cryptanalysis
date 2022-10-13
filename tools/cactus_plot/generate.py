import sys
from matplotlib.lines import lineStyles
import matplotlib.pyplot as plt

sat_solvers = []
points = {}
max_instances_solved = 0

# Gather the points in a dictionary from the input
for line in sys.stdin:
    segments = line.split(" ")

    sat_solver = segments[10][:-1]
    if sat_solver not in sat_solvers:
        sat_solvers.append(sat_solver)
        points[sat_solver] = {}

    time = float(segments[1][:-2])
    instance_index = int(segments[4][:-1].strip())
    points[sat_solver][instance_index] = time
    if instance_index + 1 > max_instances_solved:
        max_instances_solved = instance_index + 1

# Organize the points for each line
lines = []
j = -1
for sat_solver in sat_solvers:
    j += 1
    if sat_solver not in points:
        continue

    lines.append([])
    for i in range(0, max_instances_solved):
        if i not in points[sat_solver]:
            continue

        time = points[sat_solver][i]
        lines[j].append(time)


point_styles = ['b+', 'g.', 'r+', 'c*', 'm.', 'ys', 'kp']

plt.figure(dpi=600)
i = 0
for sat_solver in sat_solvers:
    sat_solver_index = sat_solvers.index(sat_solver)

    time_values = []
    instance_count_values = []

    num_instances = 1
    for time in lines[sat_solver_index]:
        time_values.append(time)
        instance_count_values.append(num_instances)
        num_instances += 1

    time_values.sort()

    plt.plot(instance_count_values, time_values,
             point_styles[i], label=sat_solver, linewidth=0.5, linestyle='solid')
    i += 1

plt.grid()
plt.legend()
plt.xlabel("Instances Solved")
plt.ylabel("Time Limit")
plt.savefig("plot.png")
