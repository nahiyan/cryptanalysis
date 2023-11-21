from sha256 import _hash

# TODO: Make the script dynamic - take the input from a file or stdin

msg_str = """
10000111101001110001100110110001
01011010101101011101110001011000
10011100100010110000001000010010
00100100010001001000000011000000
01000011000101000000001110000000
01010001001110110000001001010000
00100001000000000000000100001000
00010000000000000000000000000010
nnu0u0nn01110n0nn0uu0u11uu001110
nuun1011u10001nun10110nu00uu0011
00101101100000000111110011000001
10111111010101010010011000100100
10101111110111111101010111101011
u11n11u11nnu11n10u1001u010nnnnuu
01000110111100000000010101111111
00110101110001100011101010101111
"""

cv_str = """
10101001000000011110001101010111 00100011101101100001011111101010
11011010110001101000011100001011 00111010001110100100011000010100
11111110001100100001101011101100 01101010011011001011100100000000
00110100100000010100100100111100 10111110001111011100011000110111
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

order = 26
is_collision = _hash(order, msg_f, cv) == _hash(order, msg_g, cv)
print("Verified" if is_collision else "Failed")
