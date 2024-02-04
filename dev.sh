#!/bin/bash

# Function to run commands in background
run_in_background() {
    DEV_MODE=true "$@" &
}

# Trap Ctrl+C to kill background processes
trap 'kill $(jobs -p); echo "Ctrl+C pressed. Exiting."' INT

# Run commands in background
run_in_background pnpm tw:watch
run_in_background pnpm viter
run_in_background air

# Wait for all background jobs to finish
wait

echo "All commands completed!"

