#!/bin/sh

encoding_file=$1
cubes_file=$2
line_number=$3
line=$(sed "${line_number}q;d" "$cubes_file")
cube=${line:2: -2}
cube_length=$(echo $cube | wc -w)
num_clauses=$(head -n 1 $encoding_file | awk '{print $4}')
num_clauses=$((num_clauses+cube_length))
num_vars=$(head -n 1 $encoding_file | awk '{print $3}')

echo "p cnf $num_vars $num_clauses"
tail -n +2 $encoding_file
for lit in $(echo $cube); do
    echo "$lit 0"
done
