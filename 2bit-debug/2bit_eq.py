import re
from collections import namedtuple

import collections



Eq = namedtuple("Eq", "x y diff")

def get_eqs(file, pattern):
    eqs = set()
    text = file.read()
    matches = re.findall(pattern, text)
    for match in matches:
        eq = Eq(match[0].replace("_", ""), match[2].replace("_", ""), 1 if match[1] == '"=' or match[1] == "=/=" else 0)
        eqs.add(eq)
    return sorted(eqs)

with open("27_mendel_cnds.txt", "r") as mendel_cnds_file:
    mendel_eqs = get_eqs(mendel_cnds_file, r'([AWE]\d+,\d+)\s(["=]+)\s([AWE]\d+,\d+)')

with open("/tmp/log.txt", "r") as our_eqs_file:
    our_eqs = get_eqs(our_eqs_file, r'([AWE]_\d+,\d+)\s([/=]+)\s([AWE]_\d+,\d+)')
    
for eq in mendel_eqs:
    eq_reversed = Eq(eq.y, eq.x, eq.diff)
    if eq in our_eqs or eq_reversed in our_eqs:
        print(eq, "found=1")
    elif Eq(eq.x, eq.y, not eq.diff) in our_eqs or Eq(eq.y, eq.x, not eq.diff) in our_eqs:
        print(eq, "found=-1")
    else:
        print(eq, "found=0")
