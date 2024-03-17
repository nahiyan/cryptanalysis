#!/bin/bash

file=$1
output_dir=$2
min_cubes=$3
decrement=$4
offset=$5

basename=$(basename "$file")

# Skip if already cubed
if [ -f "$output_dir/$basename.cubes" ]; then
    echo "Skipped"
    exit
fi

# Free vars
echo $basename
free_vars=$(march_cu "$file" -d 1 |
    grep "free variables" |
    awk '{print $NF}' 2>&1)
free_vars=$((free_vars + offset))

# Loop until we get the min. cubes
while true; do
    echo $free_vars
    cubes=$(march_cu "$file" -n $free_vars |
        grep "number of cubes" |
        awk '{print int($5)}' 2>&1)
    echo $cubes
    if [ $cubes -ge $min_cubes ]; then
        march_cu "$file" -n $free_vars -o $output_dir/$basename.cubes
        break
    fi
    free_vars=$(expr "$free_vars" - $decrement)
done
