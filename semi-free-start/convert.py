import os
import subprocess
import multiprocessing
import sys
import itertools


F32 = 0xFFFFFFFF

_k = [
    0x428A2F98,
    0x71374491,
    0xB5C0FBCF,
    0xE9B5DBA5,
    0x3956C25B,
    0x59F111F1,
    0x923F82A4,
    0xAB1C5ED5,
    0xD807AA98,
    0x12835B01,
    0x243185BE,
    0x550C7DC3,
    0x72BE5D74,
    0x80DEB1FE,
    0x9BDC06A7,
    0xC19BF174,
    0xE49B69C1,
    0xEFBE4786,
    0x0FC19DC6,
    0x240CA1CC,
    0x2DE92C6F,
    0x4A7484AA,
    0x5CB0A9DC,
    0x76F988DA,
    0x983E5152,
    0xA831C66D,
    0xB00327C8,
    0xBF597FC7,
    0xC6E00BF3,
    0xD5A79147,
    0x06CA6351,
    0x14292967,
    0x27B70A85,
    0x2E1B2138,
    0x4D2C6DFC,
    0x53380D13,
    0x650A7354,
    0x766A0ABB,
    0x81C2C92E,
    0x92722C85,
    0xA2BFE8A1,
    0xA81A664B,
    0xC24B8B70,
    0xC76C51A3,
    0xD192E819,
    0xD6990624,
    0xF40E3585,
    0x106AA070,
    0x19A4C116,
    0x1E376C08,
    0x2748774C,
    0x34B0BCB5,
    0x391C0CB3,
    0x4ED8AA4A,
    0x5B9CCA4F,
    0x682E6FF3,
    0x748F82EE,
    0x78A5636F,
    0x84C87814,
    0x8CC70208,
    0x90BEFFFA,
    0xA4506CEB,
    0xBEF9A3F7,
    0xC67178F2,
]


def _rotr(x, y):
    return ((x >> y) | (x << (32 - y))) & F32


def _maj(x, y, z):
    return (x & y) ^ (x & z) ^ (y & z)


def _ch(x, y, z):
    return (x & y) ^ ((~x) & z)


def sigma0(x):
    return _rotr(x, 2) ^ _rotr(x, 13) ^ _rotr(x, 22)


def sigma1(x):
    return _rotr(x, 6) ^ _rotr(x, 11) ^ _rotr(x, 25)


def s0(x):
    return _rotr(x, 7) ^ _rotr(x, 18) ^ x >> 3


def s1(x):
    return _rotr(x, 17) ^ _rotr(x, 19) ^ x >> 10


combos = []
for i in range(pow(2, 11)):
    x = [i >> j & 1 for j in range(11)]
    y = [index + 8 if i == 1 else 0 for index, i in enumerate(x)]
    z = []
    for item in y:
        if item == 0:
            continue
        z.append(item)
    if len(z) > 0:
        combos.append(z)


def run(command):
    try:
        # Use the subprocess module to run the 'tail' command
        result = subprocess.run(command, stdout=subprocess.PIPE, text=True, check=True)
        return result.stdout
    except Exception:
        return None


def rm_lits(literals, range):
    literals_ = []
    for lit in literals:
        lit_abs = abs(lit)
        if lit_abs >= range[0] and lit_abs <= range[1]:
            continue
        literals_.append(lit)
    return literals_


def try_seeds(strt_seed):
    seed = strt_seed
    while True:
        print(f"Trying with seed = {seed}")
        process = subprocess.Popen(
            ["kissat", f"--seed={seed}", "28-sf-start.cnf"],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
        )

        # Read stdout and stderr
        stdout_data, _ = process.communicate()
        # Check if the process has completed
        process.wait()
        if process.returncode != 10:
            continue
        # with open("output.log", "w") as logfile:
        #     logfile.write(stdout_data)
        msg = get_msg(stdout_data)
        print(msg)
        exitcode = convert(msg)
        if exitcode == 10:
            print("Found", msg)
            break

        seed += 1


def get_msg_lits(msg, strt_index):
    # msg = "14c48440b3c3277fad69812dc3d4dffa7eae690b7f9fe027832aece89a4894581607a45cdb81bdc88786e031d8f2280172b6be5e45a2652ff3fbb17a2ce70f52"
    # mess_words = "14c48440 b3c3277f ad69812d c3d4dffa 7eae690b 7f9fe027 832aece8 9a489458 e6b2f4fc d759b930 8786e031 d8f22801 72b6be5e 47e26dbf f3fbb17a 2ce70f52".split()
    # msg = "b28c2f27e4fce3e2cb749ca9ccf67d4c528911f24184e0e810bd0c845207030d1c07a10cca819dc8e537f9e714a994dff0313d1bc8a153485785be5a33b22d56"
    chunk_size = 8
    msg_words = [msg[i : i + chunk_size] for i in range(0, len(msg), chunk_size)]

    lits = []
    bit_index = 0
    for word_index, word in enumerate(msg_words):
        word_int = int(word, 16)
        for i in range(32):
            bit = word_int >> i & 1
            sign = 1 if bit == 1 else -1 if bit == 0 else "u"
            if word_index >= strt_index:
                lits.append(sign * (bit_index + 1))
            bit_index += 1

    return lits


# def get_a_lits(a, enc_path):
#     chunk_size = 8
#     words = [a[i : i + chunk_size] for i in range(0, len(a), chunk_size)]
#     a_map = get_map(enc_path)

#     lits = []
#     for word_index, word in enumerate(words):
#         word_int = int(word, 16)
#         bit_index = a_map[word_index + 4]
#         for i in range(32):
#             bit = word_int >> i & 1
#             sign = 1 if bit == 1 else -1 if bit == 0 else "u"
#             lits.append(sign * (bit_index + 1))
#             bit_index += 1

#     return lits


def get_word(sol, start_index):
    word = []
    for i in range(32):
        bit = 1 if sol[start_index + i - 1] > 0 else 0
        word.append(bit)
    word.reverse()
    word_b_str = "".join(map(str, word))
    word_int = int(word_b_str, 2)
    return word_int


def get_msg(sol):
    msg_words = []
    for i in range(16):
        msg_words.append(format(get_word(sol, 1 + i * 32), "08x"))
    msg = " ".join(msg_words)
    return msg


def get_lits(name, logdata, enc_path, keys):
    vars = []
    map_ = get_map(enc_path, name)
    lines = logdata.split("\n")
    for line in lines:
        if len(line) == 0:
            continue
        if line[0] == "v":
            words = line[2:].split()
            for word in words:
                word = int(word)
                var_index = abs(word)
                if var_index == 0:
                    continue
                eligible = False
                for key in keys:
                    if var_index >= map_[key] and var_index < map_[key] + 32:
                        eligible = True
                        break
                if eligible:
                    vars.append(word)
    return vars


def get_map(enc_path, name, block="f"):
    map_ = {}
    output = run(["grep", "-E", f"c {name}[0-9]+_{block}", enc_path])
    entries = output.split("\n")
    for entry in entries:
        if len(entry) == 0:
            continue
        segments = entry.split("_")
        if len(segments) == 3:
            index = segments[1]
            value = segments[2].split()[1]
        else:
            index = segments[0].split()[1][1:]
            value = segments[1].split()[1]
        map_[int(index)] = int(value)
    return map_


def get_sol(logdata):
    sol = []
    lines = logdata.split("\n")
    for line in lines:
        if len(line) == 0:
            continue
        if line[0] == "v":
            words = line[2:-1].split()
            words = [int(word) for word in words]
            sol.extend(words)
    return sol


def convert(logdata, enc_path):
    # Get the literals for fixing some of the variables
    x, y = 0, 7
    a_lits = get_lits("A_", logdata, enc_path, list(range(x + 4, y + 4 + 1)))
    # print(a_lits)
    literals = a_lits
    # literals.clear()
    print("Number of unit clauses added:", len(literals))

    # Read the original encoding
    head = run(["head", "-n", "1", enc_path]).split()
    org_vars_count = int(head[2])
    org_clauses_count = int(head[3])

    vars_count = org_vars_count + (
        max([abs(x) for x in literals]) if len(literals) > 0 else 0
    )
    clauses_count = org_clauses_count + len(literals)

    # process = subprocess.Popen(
    #     # ["kissat", "tmp.cnf"],
    #     # ["maplesat", "tmp.cnf"],
    #     ["/home/nahiyan/code/cryptanalysis/solvers/maplesat-crypto/build/maplesat", "-no-pre", "-rnd-init", "-rnd-seed=1"],
    #     # ["/home/nahiyan/code/CDCL-Crypto/src/simp/maplesat_static"],
    #     # stdin=subprocess.PIPE,
    #     stdout=subprocess.PIPE,
    #     stderr=subprocess.PIPE,
    #     text=True,
    # )
    # stdin = process.stdin
    tmp = open("tmp.cnf", "w")
    # stdin.write(f"p cnf {vars_count} {clauses_count}\n")
    tmp.write(f"p cnf {vars_count} {clauses_count}\n")
    with open(enc_path, "r") as original:
        lines = original.readlines()
        for line in lines:
            if line[0] != "p" and line[0] != "c":
                # stdin.write(line)
                tmp.write(line)
    for literal in literals:
        # stdin.write(f"{literal} 0\n")
        tmp.write(f"{literal} 0\n")
    # stdin.flush()
    tmp.flush()

    process = subprocess.Popen(
        ["kissat", "tmp.cnf"],
        # ["maplesat", "tmp.cnf"],
        # ["/home/nahiyan/code/cryptanalysis/solvers/maplesat-crypto/build/maplesat", "-no-pre", "-rnd-init", "-rnd-seed=1", "tmp.cnf"],
        # ["/home/nahiyan/code/CDCL-Crypto/src/simp/maplesat_static"],
        # stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
    )
    stdout, _ = process.communicate()
    # process.stdin.close()
    process.stdout.close()
    process.stderr.close()

    print(stdout)

    exitcode = process.wait()
    print("Exit code", exitcode)
    if exitcode == 10:
        sol = get_sol(stdout)
        print("Solved", get_msg(sol))
        # sol_log = open("solution.log", "a")
        # sol_log.write("Combo: " + str(combo) + "\n" + stdout)
    return exitcode


real_enc, sfs_log = "encodings/real/28-8-step-table.cnf", "logs/28-sfs.log"
with open(sfs_log, "r") as logfile:
    # Get the SFS solution literals
    logdata = logfile.read()
    sfs_sol = get_sol(logdata)

    sfs_a_map = get_map("encodings/sfs/28.cnf", "A_", "f")
    a4 = get_word(sfs_sol, sfs_a_map[4])
    print(format(a4, "08x"))

    # msg = get_msg(sfs_sol)
    # print(msg)

    convert(logdata, real_enc)
    
    # a0 = 0xa54ff53a
    # a1 = 0x3c6ef372
    # a2 = 0xbb67ae85
    # a3 = 0x6a09e667
    # e0 = 0x5be0cd19
    # e1 = 0x1f83d9ab
    # e2 = 0x9b05688c
    # e3 = 0x510e527f
    # w0 = (a4-sigma0(a3)-_maj(a3, a2, a1)-e0-sigma1(e3)-_ch(e3, e2, e1)-_k[0]) & F32
    # assert(format(w0, "08x") == "69a0dca8")