#!/bin/bash

SESSION="smanzy_dev"
# Using $HOME is safer than '~' inside scripts
PROJECT_DIR="$HOME/Develpment/smanzy_site/smanzy_backend"

echo "--- Shutting down development environment ---"

# 1. Stop Docker Containers (Graceful)
if [ -d "$PROJECT_DIR" ]; then
    echo "Stopping Docker containers in $PROJECT_DIR..."
    cd "$PROJECT_DIR"
    docker compose down
else
    echo "⚠️  Warning: Could not find $PROJECT_DIR. Docker containers might still be running."
fi

# 2. Kill the tmux session
if tmux has-session -t $SESSION 2>/dev/null; then
    echo "Killing tmux session: $SESSION..."
    tmux kill-session -t $SESSION
else
    echo "Tmux session '$SESSION' is not running."
fi

echo "--- Shutdown complete ---"
