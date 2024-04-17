import argparse

"""
SHA-256 collision verifier adopted from
https://gist.github.com/DavidBuchanan314/aa9ab4265fe402ab86399b5f9da82888.

SHA256 impl follows FIPS 180-4
https://nvlpubs.nist.gov/nistpubs/FIPS/NIST.FIPS.180-4.pdf
"""

F32 = 0xFFFFFFFF


# Section 2.2.2
def rotr(x, n):
    return ((x >> n) | (x << (32 - n))) & F32


# Section 2.2.2
def shr(x, n):
    return x >> n


# Section 4.1.2 (4.2)
def Ch(x, y, z):
    return (x & y) ^ (~x & z)


# Section 4.1.2 (4.3)
def Maj(x, y, z):
    return (x & y) ^ (x & z) ^ (y & z)


# Section 4.1.2 (4.4)
def S0(x):
    return rotr(x, 2) ^ rotr(x, 13) ^ rotr(x, 22)


# Section 4.1.2 (4.5)
def S1(x):
    return rotr(x, 6) ^ rotr(x, 11) ^ rotr(x, 25)


# Section 4.1.2 (4.6)
def s0(x):
    return rotr(x, 7) ^ rotr(x, 18) ^ shr(x, 3)


# Section 4.1.2 (4.7)
def s1(x):
    return rotr(x, 17) ^ rotr(x, 19) ^ shr(x, 10)


"""
Section 4.2.2 - SHA-256 Constants
"""
K = [
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


# Section 5.2.1
def word_iterator(m):
    for i in range(0, len(m), 32 // 8):
        yield int.from_bytes(m[i : i + 32 // 8], "big")


# Section 5.3.3
standard_iv = "6a09e667 bb67ae85 3c6ef372 a54ff53a 510e527f 9b05688c 1f83d9ab 5be0cd19"


# Section 6.2
def sha256(order, blocks, cv):
    # 6.2.1, 1) Initialize H
    H = [
        int(format(int.from_bytes(cv[i * 4 : (i * 4) + 4], "big"), "08x"), 16)
        for i in range(8)
    ]

    # 6.2.2 - SHA-256 Hash Computation
    for block in blocks:

        # 1. Prepare the message schedule
        W = list(word_iterator(block))
        for t in range(16, order):
            W.append((s1(W[t - 2]) + W[t - 7] + s0(W[t - 15]) + W[t - 16]) & F32)

        # 2. Initialize the eight working variables
        a, b, c, d, e, f, g, h = H

        # 3.
        for t in range(order):
            T1 = h + S1(e) + Ch(e, f, g) + K[t] + W[t]
            T2 = S0(a) + Maj(a, b, c)
            h = g
            g = f
            f = e
            e = (d + T1) & F32
            d = c
            c = b
            b = a
            a = (T1 + T2) & F32

        # 4. Calculate the next Hash value
        for i, x in enumerate((a, b, c, d, e, f, g, h)):
            H[i] = (H[i] + x) & F32

    # convert the result to bytes
    M = b""
    for word in H:
        M += word.to_bytes(32 // 8, "big")
    return M


if __name__ == "__main__":
    parser = argparse.ArgumentParser("simple_example")
    parser.add_argument(
        "-s",
        help="The number of steps in the step-reduced SHA-256.",
        type=int,
        required=True,
    )
    parser.add_argument(
        "--m0",
        help="The first message in the pair.",
        type=str,
        required=True,
    )
    parser.add_argument(
        "--m1",
        help="The second message in the pair.",
        type=str,
        required=True,
    )
    parser.add_argument(
        "--c0",
        help="(Optional) The first chaining valueDefaults to the standard IV.",
        type=str,
        default=standard_iv,
    )
    parser.add_argument(
        "--c1",
        help="(Optional) The first chaining valueDefaults to the standard IV.",
        type=str,
        default=standard_iv,
    )
    args = parser.parse_args()

    m0 = bytes.fromhex(args.m0)
    m1 = bytes.fromhex(args.m1)

    c0 = bytes.fromhex(args.c0)
    c1 = bytes.fromhex(args.c1)

    h0 = sha256(args.s, [m0], c0)
    h1 = sha256(args.s, [m1], c1)

    if m0 != m1 and h0 == h1:
        if c0 == bytes.fromhex(standard_iv) and c1 == bytes.fromhex(standard_iv):
            col_type = "(regular)"
        elif c0 != c1:
            col_type = "(free-start)"
        elif c0 == c1:
            col_type = "(semi-free-start)"
        print("Valid", col_type, "collision.")
        exit(0)
    else:
        print("Invalid collision.")
        exit(-1)
