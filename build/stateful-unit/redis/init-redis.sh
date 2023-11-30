#!/bin/sh

echo "Setting initial data in Redis..."

while IFS=":" read -r key value
do
  redis-cli SET $key "$value"
done < "/data/small.txt"

while IFS=":" read -r key value
do
  redis-cli SET $key "$value"
done < "/data/medium.txt"

while IFS=":" read -r key value
do
  redis-cli SET $key "$value"
done < "/data/large.txt"

echo "Data initialization completed."
