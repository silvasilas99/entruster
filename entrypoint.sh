#!/bin/sh
# entrypoint.sh — waits for the .fabric_ready sentinel written by init.sh
# before launching the API server.  This prevents the gRPC connection from
# loading stale certs while the Fabric network is being (re)initialized.
#
# The sentinel lives at /app/shared/.fabric_ready which is bind-mounted from
# the project root on the host via docker-compose.yml.

set -e

READY_MARKER="/app/shared/.fabric_ready"

echo "⏳ Waiting for Fabric network to be ready..."
echo "   Watching for sentinel: $READY_MARKER"

TIMEOUT=300   # 5 minutes maximum
ELAPSED=0
INTERVAL=3

while [ "$ELAPSED" -lt "$TIMEOUT" ]; do
    if [ -f "$READY_MARKER" ]; then
        echo "✅ Fabric network ready (sentinel found after ${ELAPSED}s) — starting API server."
        exec ./server
    fi
    sleep "$INTERVAL"
    ELAPSED=$((ELAPSED + INTERVAL))
done

echo "❌ Timed out waiting for Fabric network after ${TIMEOUT}s. Exiting."
exit 1
