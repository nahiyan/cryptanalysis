from rich.console import Console
from rich.table import Table
import sys

args = sys.argv
if len(args) != 2:
    sys.exit("Missing the proof file's path")

variable_ranges = [(1, 512)]
variables = {}
proof_file_path = args[1]
proof = open(proof_file_path, "r")
for line in proof:
    segments = line.split()
    if len(segments) == 2:
        variables[abs(int(segments[0]))] = "1" if int(segments[0]) >= 1 else "0"

console = Console()
counter = 0
for variable, size in variable_ranges:
    table = Table(show_header=True, header_style="bold magenta", expand=True)
    console.print("Variable: " + str(variable), " Size: " + str(size))
    table.add_column("ID")
    for i in range(min(size, 32)):
        table.add_column(str(i + 1))

    rows = []
    row = []
    for i in range(size):
        row.append(variables[variable + i] if variables.get(variable + i) else "")
        if (len(row) == 32) or i == size - 1:
            rows.append(row)
            row = []
    row_index = 0
    for row in rows:
        table.add_row(str(row_index), *row)
        row_index += 1

    console.print(table)
