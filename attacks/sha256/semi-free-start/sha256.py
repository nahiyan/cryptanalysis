#!/usr/bin/python3

__base__ = "https://github.com/thomdixon/pysha2/blob/master/sha2/sha256.py"
__author__ = "Lukas Prokop"
__license__ = "MIT"

import copy
import struct
import binascii

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

_h = [
    0x6A09E667,
    0xBB67AE85,
    0x3C6EF372,
    0xA54FF53A,
    0x510E527F,
    0x9B05688C,
    0x1F83D9AB,
    0x5BE0CD19,
]


def _rotr(x, y):
    return ((x >> y) | (x << (32 - y))) & F32


def _maj(x, y, z):
    return (x & y) ^ (x & z) ^ (y & z)


def _ch(x, y, z):
    return (x & y) ^ ((~x) & z)


class SHA256:
    _output_size = 8
    blocksize = 1
    block_size = 64
    digest_size = 32
    order = 64

    def __init__(self, m=None, order=64, h=_h):
        self._counter = 0
        self._cache = b""
        self._k = copy.deepcopy(_k)
        self._h = copy.deepcopy(h)
        self.order = order

        self.update(m)

    def _compress(self, c):
        w = [0] * 64
        w[0:16] = struct.unpack("!16L", c)

        for i in range(16, self.order):
            s0 = _rotr(w[i - 15], 7) ^ _rotr(w[i - 15], 18) ^ (w[i - 15] >> 3)
            s1 = _rotr(w[i - 2], 17) ^ _rotr(w[i - 2], 19) ^ (w[i - 2] >> 10)
            w[i] = (w[i - 16] + s0 + w[i - 7] + s1) & F32

        a, b, c, d, e, f, g, h = self._h

        for i in range(self.order):
            s0 = _rotr(a, 2) ^ _rotr(a, 13) ^ _rotr(a, 22)
            t2 = s0 + _maj(a, b, c)
            s1 = _rotr(e, 6) ^ _rotr(e, 11) ^ _rotr(e, 25)
            t1 = h + s1 + _ch(e, f, g) + self._k[i] + w[i]

            h = g
            g = f
            f = e
            e = (d + t1) & F32
            d = c
            c = b
            b = a
            a = (t1 + t2) & F32

        for i, (x, y) in enumerate(zip(self._h, [a, b, c, d, e, f, g, h])):
            self._h[i] = (x + y) & F32

    def update(self, m):
        if not m:
            return

        self._cache += m
        self._counter += len(m)

        while len(self._cache) >= 64:
            self._compress(self._cache[:64])
            self._cache = self._cache[64:]

    def digest(self):
        r = copy.deepcopy(self)
        # r.update(_pad(self._counter))
        data = [struct.pack("!L", i) for i in r._h[: self._output_size]]
        return b"".join(data)

    def hexdigest(self):
        return binascii.hexlify(self.digest()).decode("ascii")


if __name__ == "__main__":

    def check(order, msg, sig, h=None):
        msg = bytes.fromhex(msg)
        m = SHA256(msg, order) if h == None else SHA256(msg, order, h)
        # m.update(msg.encode('ascii'))
        print(m.hexdigest())
        print(m.hexdigest() == sig)

    tests = {
        # "":
        #     'e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855',
        # "a":
        #     'ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb',
        # "abc":
        #     'ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad',
        # "message digest":
        #     'f7846f55cf23e14eebeab5b4e1550cad5b509e3348fbc4efa3a1413d393cb650',
        # "abcdefghijklmnopqrstuvwxyz":
        #     '71c480df93d6ae2f1efad1447c66c9525e316218cf51fc8d9ed832f2daf18b73',
        # "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789":
        #     'db4bfcbd4da0cd85a60c3c37d3fbd8805c77f15fc6b1fdfe614ee0a7c8fdb4c0',
        # ("12345678901234567890123456789012345678901234567890123456789"
        #  "012345678901234567890"):
        #     'f371bc4a311f2b009eef952dd83ca80e2b60026c8e935592d0f9c308453c813e',
        # "00baf6626abc2df808da36a518c69f09b0d2ed0a79421ccfde4f559d2e42128b":
        #     'b835e56173be2b5b7177d71bf02850dc578ac855ac60f91a108eec253bd5a543',
        # "c9edzaekfwjksbz0lewkugcjxnxlmthypeaxnok3ny8eexmzsndspzw0vt2eiprz".encode('ascii'): "7e24d5acc17297ecf1978c9642056f7c5bfc43114e74d426ce978a4e25973944",
        "633965647a61656b66776a6b73627a306c65776b7567636a786e786c6d746879706561786e6f6b336e79386565786d7a736e6473707a7730767432656970727a": "7e24d5acc17297ecf1978c9642056f7c5bfc43114e74d426ce978a4e25973944"
    }

    # for inp, out in tests.items():
    #     check(inp, out)

    h = [
        0xC993C1BC,
        0x4685E40F,
        0x41270246,
        0x3D4A8BD1,
        0x24723AF6,
        0xD700757C,
        0x12E6468F,
        0xE8BA4416,
    ]

    check(
        28,
        "b28c2f27e4fce3e2cb749ca9ccf67d4c528911f24184e0e810bd0c845207030d1c07a10cca819dc8e537f9e714a994dff0313d1bc8a153485785be5a33b22d56",
        "0f1569031514c450850318400d02954e23e883d4bbf3f1af092b29dac3e055fe",
        h,
    )
    check(
        28,
        "b28c2f27e4fce3e2cb749ca9ccf67d4c528911f24184e0e810bd0c845207030decb2f1acc6599930e537f9e714a994dff0313d1bcae15bd85785be5a33b22d56",
        "0f1569031514c450850318400d02954e23e883d4bbf3f1af092b29dac3e055fe",
        h,
    )
