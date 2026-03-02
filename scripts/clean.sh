#!/usr/bin/env bash
set -euo pipefail

STEP="clean"
export STEP
# shellcheck disable=SC1091
source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"

ensure_log_dir
run_logged "Unable to clean logs directory." find "$LOG_DIR" -type f -name "*.log" ! -name "clean.log" -delete
