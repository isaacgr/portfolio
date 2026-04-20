#!/bin/bash

set -e

get_version_from_go_version_file() {
    VERSION=$(grep -o 'const Version = "[^"]*"' version/version.go | sed 's/const Version = "\(.*\)"/\1/')
    if [ -z "$VERSION" ]; then
        msg "ERROR: Cannot parse version from version/version.go"
        msg "ERROR: Please ensure the file contains: const Version = \"version-string\""
        exit 1
    fi
    # Strip the 'v' prefix if it exists
    VERSION=$(echo "$VERSION" | sed 's/^v//')
}

main() {
    get_version_from_go_version_file
    # Abort if version is missing
    if [ -z "$VERSION" ]; then
        echo "Failed to determine version"
        exit 1
    fi

    # Echo completed version to caller
    # shellcheck disable=SC3037
    echo -n "$VERSION"
}

main
