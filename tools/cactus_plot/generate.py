import sys
from matplotlib.lines import lineStyles
import matplotlib.pyplot as plt

time_in_x = False

if len(sys.argv) >= 2 and sys.argv[1] == "invert":
    time_in_x = True

sat_solvers = []
points = {}

# Gather the points in a dictionary from the input
for line in sys.stdin:
    segments = line.split(" ")

    exit_code = int(segments[-1])
    if exit_code != 10:
        continue

    sat_solver = segments[10][:-1]
    if sat_solver not in sat_solvers:
        sat_solvers.append(sat_solver)
        points[sat_solver] = []

    time = float(segments[1][:-2])
    points[sat_solver].append(time)

# Organize the points for each line
lines = []
j = -1
for sat_solver in sat_solvers:
    j += 1
    if sat_solver not in points:
        continue

    lines.append([])
    for time in points[sat_solver]:
        lines[j].append(time)


# Plot the cactus plot
point_styles = ['b+', 'g.', 'r+', 'c*', 'm.', 'ys', 'kp']

plt.figure(dpi=250)
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

    # Control the orientation of the plot
    x = instance_count_values
    y = time_values
    xlabel = "Instances Solved"
    ylabel = "Time Limit"
    if time_in_x:
        x = time_values
        y = instance_count_values
        xlabel = "Time Limit"
        ylabel = "Instances Solved"

    plt.plot(x, y,
             point_styles[i], label=sat_solver, linewidth=0.5, linestyle='solid')
    i += 1

plt.grid()
plt.legend()
plt.ylabel(ylabel)
plt.xlabel(xlabel)
plt.savefig("plot.png")