#!/usr/bin/env sh

# Change this to your target folder
TARGET_DIR="/app/internal/database/mongo/migrations"

# Create folder if it does not exist
mkdir -p "$TARGET_DIR"

# Ask for base name
printf "Enter base filename (without extension): "
read BASE_NAME

# Current time in milliseconds since Unix epoch
TS_MS=$(date +%s%3N)

# Build final filename: <milliseconds>_<input>.go
FINAL_NAME="${TS_MS}_${BASE_NAME}.go"

# Full path
FULL_PATH="${TARGET_DIR}/${FINAL_NAME}"

# Create the file
: > "$FULL_PATH"

echo "Created file: $FULL_PATH"
