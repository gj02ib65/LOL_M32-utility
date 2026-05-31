#!/bin/sh
# M32 Utility — Container Startup Script
#
# This script runs before Node-RED starts. It copies flows.json from the
# read-only bind-mounted source directory into /data/ where Node-RED can
# read and write it freely (no bind-mount rename restrictions).
#
# On every restart, the latest flows from git are picked up automatically.

SOURCE=/opt/m32-flows-src/flows.json
DEST=/data/flows.json

if [ -f "$SOURCE" ]; then
    cp "$SOURCE" "$DEST"
    echo "M32: Loaded flows from $SOURCE → $DEST"
else
    echo "M32: No source flows at $SOURCE — using existing $DEST (or Node-RED default)"
fi

# Delegate to the official Node-RED entrypoint
exec /usr/src/node-red/entrypoint.sh "$@"
