#!/bin/bash

# Function to generate the mock
generate_mock() {
    local path=$1

    # Extract the path components
    path_with_extension="${path}.go"
    package_name=$(basename "$(dirname "$path")")
    destination="${path}_mock.go"

    # Check if source file exists
    if [ ! -f "$path_with_extension" ]; then
        echo "Error: Source file $path_with_extension does not exist."
        return 1
    fi

    # Run mockgen command
    echo "Generating mock for $path_with_extension in package $package_name..."
    mockgen -source="$path_with_extension" -destination="$destination" -package="$package_name"

    # Check if command succeeded
    if [ $? -eq 0 ]; then
        echo "Mock generated successfully at $destination"
        return 0
    else
        echo "Failed to generate mock"
        return 1
    fi
}

# Main interactive script
clear
echo "========================================"
echo "          MOCK GENERATOR TOOL          "
echo "========================================"
echo

# Check if path is provided as argument
if [ ! -z "$1" ]; then
    path="$1"
else
    # Ask for path interactively
    read -p "Enter the path to interface file (without .go extension): " path
    echo
fi

# Confirm with the user
echo "You are about to generate a mock for: $path.go"
read -p "Proceed? (Y/n): " confirm

if [[ "$confirm" == "" || "$confirm" =~ ^[Yy] ]]; then
    generate_mock "$path"
else
    echo "Operation cancelled."
    exit 0
fi