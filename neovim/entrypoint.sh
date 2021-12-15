#!/bin/ash
# entrypoint.sh

set -eu
tmux new -s foo -d && tmux ls

exec "$@"
