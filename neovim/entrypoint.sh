#!/bin/bash

set -eu
tmux new -s golangide -d && tmux ls

exec "$@"
