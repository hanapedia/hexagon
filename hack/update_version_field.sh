#!/bin/bash

# Check if two arguments are given
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 [Directory Path] [New Version]"
    exit 1
fi

# Assign arguments to variables
directory_path=$1
new_version=$2

# Check if the directory exists
if [ ! -d "$directory_path" ]; then
    echo "Directory does not exist: $directory_path"
    exit 1
fi

# Iterate over all YAML files in the directory
for file in "$directory_path"/*.yaml; do
    if [ -f "$file" ]; then
        # Using yq to update the version field
        yq eval ".version = \"$new_version\"" "$file" --inplace
    fi
done

echo "Version updated in all YAML files in $directory_path."
