#!/usr/bin/env bash
set -euo pipefail

NON_INTERACTIVE="${NON_INTERACTIVE:-0}"
REMOVE_GO="${REMOVE_GO:-0}"
REMOVE_LOGS="${REMOVE_LOGS:-0}"

confirm() {
  local question="$1"
  if [[ "$NON_INTERACTIVE" == "1" ]]; then
    echo "$question [auto: yes]"
    return 0
  fi
  read -r -p "$question [y/N]: " ans
  [[ "${ans,,}" == "y" || "${ans,,}" == "yes" ]]
}

remove_if_exists() {
  local path="$1"
  if [[ -e "$path" ]]; then
    rm -rf "$path"
    echo "Removed $path"
  else
    echo "Not found: $path"
  fi
}

echo "==> Cleo uninstall"
if ! confirm "Proceed with uninstall?"; then
  echo "Cancelled."
  exit 0
fi

remove_if_exists "$HOME/.local/bin/cleo"

if [[ "$REMOVE_GO" == "1" ]]; then
  if confirm "Remove Go toolchain at $HOME/.local/go?"; then
    remove_if_exists "$HOME/.local/go"
  fi
fi

if [[ "$REMOVE_LOGS" == "1" ]]; then
  if confirm "Remove all cleo logs at $HOME/.cleo/logs?"; then
    remove_if_exists "$HOME/.cleo/logs"
  fi
fi

echo "Uninstall complete."
echo "If you added PATH entries manually, remove them from your shell profile."
