#!/bin/bash

./config.sh --url $GITHUB_RUNNER_REPOSITORY --token $GITHUB_RUNNER_TOKEN --labels="${GITHUB_RUNNER_LABELS}"
./run.sh