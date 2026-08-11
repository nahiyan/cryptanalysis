from sha256 import _hash

# TODO: Make the script dynamic - take the input from a file or stdin

msg_str = """
11111111111111111111111111111111
11111111111111111111111111111111
11111111111111111111111111111111
u1111111111111111111111111111111
uu1uuuu11111u1uuuuuuuuu11u111111
u1u1111u1uuuuuu1uuu1uuuu1uu1u111
11u0011111111111011111u111111111
u1010111111111111111111111111111
nun101n10110nn01011100u1un011011
11001000000000100000011100000001
00011000111110001010010110010001
n1100101011010000011101001101000
00001110010110000100111011000101
01110011110011111001000100111001
11010010001110101010000110101011
01010110101111111100101110111110
"""

cv_str = """
11111111111111111111111111111111 11111111111111111111111111111111
11111111111111111111111111111111 11111111111111111111111111111111
11111111111111111111111111111111 11111111111111111111111111111111
11111111111111111111111111111111 11111111111111111111111111111111
"""

msg, cvs = [], ([], [])

for line in msg_str.splitlines():
    line = line.strip()
    if len(line) == 0:
        continue
    assert len(line) == 32
    msg.append(line)
for line in cv_str.splitlines():
    line = line.strip()
    if len(line) == 0:
        continue
    assert len(line) == 64 + 1
    words = line.split()
    assert len(words[0]) == 32 and len(words[1]) == 32
    cvs[0].append(int(words[0].strip(), 2))
    cvs[1].append(int(words[1].strip(), 2))
cvs[0].reverse()
cvs[1].reverse()

msg_bin_f, msg_bin_g = [], []

gc_to_bin = (
    lambda gc: [0, 0]
    if gc == "0"
    else [0, 1]
    if gc == "n"
    else [1, 0]
    if gc == "u"
    else [1, 1]
)

msg_hex_f, msg_hex_g = [], []

for i, word in enumerate(msg):
    msg_bin_f.append("")
    msg_bin_g.append("")
    for j in range(32):
        bits = gc_to_bin(word[j])
        msg_bin_f[i] += str(bits[0])
        msg_bin_g[i] += str(bits[1])

    assert len(msg_bin_f[i]) == 32 and len(msg_bin_g[i])
    msg_hex_f.append(format(int(msg_bin_f[i], 2), "08x"))
    msg_hex_g.append(format(int(msg_bin_g[i], 2), "08x"))

msg_f = "".join([m for m in msg_hex_f])
msg_g = "".join([m for m in msg_hex_g])
cv = [word for word in cvs[0] + cvs[1]]

assert len(msg_f) == 128 and len(msg_g) == 128
assert len(cv) == 8

# print(msg_f)
# print(msg_g)
# print(cv)

order = 16
is_collision = _hash(order, msg_f, cv) == _hash(order, msg_g, cv)
print("Verified" if is_collision else "Failed")
