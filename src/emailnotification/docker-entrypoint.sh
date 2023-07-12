#!/bin/sh
# Abort on any error (including if wait-for-it fails).
set -e
# Wait for nats
if [ -n "$NATS_CONNECT" ]; then
/go/src/app/wait-for-it.sh "$NATS_CONNECT" -t 20
fi
# Run the main container command.
exec "$@"