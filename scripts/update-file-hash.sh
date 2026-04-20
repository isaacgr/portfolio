#!/bin/bash

set -e

BUILD_DIR=$1

#Configuration
TEMPLATE_FILE="$BUILD_DIR/web/views/base.html"
CSS_FILE="$BUILD_DIR/web/static/css/styles.css"
JS_FILE="$BUILD_DIR/web/static/js/main.js"

if [ ! -f "$CSS_FILE" ]; then
    echo "Error: CSS file $CSS_FILE not found."
    exit 1
fi
if [ ! -f "$JS_FILE" ]; then
    echo "Error: JS file $CSS_FILE not found."
    exit 1
fi

# Get modification time in milliseconds
# For Linux (GNU stat)
CSS_VERSION=$(stat -c %Y "$CSS_FILE")000
JS_VERSION=$(stat -c %Y "$JS_FILE")000

# Note: If running on macOS (BSD stat), use:
# VERSION=$(stat -f %m "$CSS_FILE")000

# Define replacement string
CSS_REPLACEMENT="href=\"/static/css/styles.css?v=$CSS_VERSION\""
JS_REPLACEMENT="href=\"/static/js/main.js?v=$JS_VERSION\""

# Use sed to update the file in-place
# The regex matches href="/static/css/styles.css followed by any characters until the closing quote
sed -i "s|href=\"/static/css/styles.css[^\"]*\"|$CSS_REPLACEMENT|g" "$TEMPLATE_FILE"
sed -i "s|href=\"/static/js/main.js[^\"]*\"|$JS_REPLACEMENT|g" "$TEMPLATE_FILE"
