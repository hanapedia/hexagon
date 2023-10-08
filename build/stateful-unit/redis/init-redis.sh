#!/bin/sh

echo "Setting initial data in Redis..."

while IFS= read -r line
do
  redis-cli SET $line
done < "/data/small.txt"

while IFS= read -r line
do
  redis-cli SET $line
done < "/data/medium.txt"

while IFS= read -r line
do
  redis-cli SET $line
done < "/data/large.txt"

echo "Data initialization completed."
