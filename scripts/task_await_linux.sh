#!/bin/sh

if [ "$#" -lt 1 ]; then
    echo "Usage: $0 <sleep_duration_in_milliseconds>"
    exit 1
fi

SLEEP_DURATION=$1

# Converte milissegundos para segundos
SLEEP_DURATION_SEC=$(echo "scale=3; $SLEEP_DURATION / 1000" | bc)

echo "Awaiting for $SLEEP_DURATION milliseconds..."
sleep "$SLEEP_DURATION_SEC"

echo "Done awaiting!"