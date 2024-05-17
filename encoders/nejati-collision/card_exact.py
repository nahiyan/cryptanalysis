from pysat.card import CardEnc, EncType
import argparse

parser = argparse.ArgumentParser()
parser.add_argument(
    "-l", help="List of variable IDs (space-separated).", type=str, required=True
)
parser.add_argument("-k", help="The value of the bound.", type=int, required=True)
parser.add_argument("-t", help="Top ID used in the encoding.", type=int)

args = parser.parse_args()
var_ids = list(map(lambda x: abs(int(x)), args.l.split()))
k = args.k
top_id = args.t

enc = CardEnc.equals(lits=var_ids, bound=k, encoding=EncType.kmtotalizer, top_id=top_id)
for clause in enc.clauses:
    print(" ".join([str(l) for l in clause]), "0")
