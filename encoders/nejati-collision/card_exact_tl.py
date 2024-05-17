#!/usr/bin/python3
import sys

global top_id


# Return the number of leaves under node k in the complete binary tree with N nodes
def num_leaves_under(N, k):
    if k >= N:
        return 0
    if 2 * k + 1 >= N:
        return 1
    return num_leaves_under(N, 2 * k + 1) + num_leaves_under(N, 2 * k + 2)


# Generate clauses which say that exactly k of the variables in X are set to true
# using the totalizer encoding of Bailleux and Boufkhad
def card(X, k):
    global numvars
    N = len(X)
    R = [["F" for j in range(N + 2)] for i in range(2 * N - 1)]
    for i in range(2 * N - 1):
        R[i][0] = "T"
    c = numvars
    for i in range(N - 1):
        t = num_leaves_under(2 * N - 1, i)
        for j in range(t):
            c += 1
            R[i][j + 1] = c
    for i in range(N - 1, 2 * N - 1):
        R[i][1] = X[i - N + 1]
    # print("c "+str(R))

    numvars = c
    newcl = []

    for i in range(N - 1):
        m = num_leaves_under(2 * N - 1, i)
        for sigma in range(m + 1):
            # Solve alpha + beta = sigma
            for alpha in range(sigma + 1):
                beta = sigma - alpha
                newcl.append(
                    "-{0} -{1} {2} 0".format(
                        R[2 * i + 1][alpha], R[2 * i + 2][beta], R[i][sigma]
                    )
                )
                newcl.append(
                    "{0} {1} -{2} 0".format(
                        R[2 * i + 1][alpha + 1], R[2 * i + 2][beta + 1], R[i][sigma + 1]
                    )
                )

    for i in range(1, k + 1):
        newcl.append("{0} 0".format(R[0][i]))
    for i in range(k + 1, N + 1):
        newcl.append("-{0} 0".format(R[0][i]))

    return newcl


if len(sys.argv) < 4:
    print("Need N, k, and top ID")
    exit()

N = int(sys.argv[1])  # Number of variables in cardinality constraint
k = int(sys.argv[2])  # Number of variables to set true
top_id = int(sys.argv[3])  # Top variable ID

numvars = N
clauses = card(range(1, N + 1), k)

# If "-T" appears in a clause, remove it
i = 0
while i < len(clauses):
    clauses[i] = clauses[i].replace("-T ", "")
    i += 1

# If "-F" or "T" appears in a clause, remove the clause
i = 0
while i < len(clauses):
    if "-F" in clauses[i]:
        clauses.remove(clauses[i])
    elif "T" in clauses[i]:
        clauses.remove(clauses[i])
    else:
        i += 1

# If "F" appears in a clause, remove it
i = 0
while i < len(clauses):
    clauses[i] = clauses[i].replace("F ", "")
    i += 1

# print("p cnf {} {}".format(numvars, len(clauses)))
for C in clauses:
    # print(C.split(" "))
    lits = list(map(lambda x: int(x), C.split(" ")))
    for lit in lits:
        if lit == 0:
            print(lit)
        elif abs(lit) > N:
            print((-1 if lit < 0 else 1) * (abs(lit) + top_id), end=" ")
        else:
            print(lit, end=" ")

    # print(lits)
    # print(C)
