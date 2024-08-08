#!/bin/bash

# Define the base directories
CONFIG_DIR="./dev/config"
MANIFEST_DIR="./dev/manifest"

DEFAULT_DOCKER_USER="hexagonbenchmark"

if [ -z "$1" ]; then
    echo "No docker user provided."
    exit 1
fi

# Access the first argument
docker_user="$1"

# Get the list of directories under the config directory
for dir in $(ls -d $CONFIG_DIR/*/); do
    # Extract the directory name without the path
    dir_name=$(basename $dir)

    # Construct the paths
    config_path="$CONFIG_DIR/$dir_name"
    manifest_path="$MANIFEST_DIR/$dir_name/generated"

    mkdir -p $manifest_path

    rm -f $manifest_path/*
    ./bin/hexctl generate -f $config_path -o $manifest_path
done
