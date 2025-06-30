#!/bin/bash
set -e

# Create output directory
mkdir -p output

# Build the Docker image
docker build -t pom-deb-builder -f packaging/deb/Dockerfile .

# Run the container to build the package
docker run --rm -v "$(pwd)/output:/output" pom-deb-builder

echo "Debian package has been built and placed in the output directory" 