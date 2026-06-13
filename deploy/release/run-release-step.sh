#!/usr/bin/env bash
set -euo pipefail

if [[ "$#" -eq 0 ]]; then
  printf 'usage: %s <command> [args...]\n' "$0" >&2
  exit 2
fi

transcript="${CRM_DEPLOY_TRANSCRIPT:-/opt/crm-system/releases/66d2531/deploy-transcript.log}"
stdin_file="${CRM_RELEASE_STDIN_FILE:-}"
mkdir -p "$(dirname "$transcript")"

die() {
  printf 'ERROR: %s\n' "$*" >&2
  exit 1
}

if [[ "${1:-}" == "git" && "${2:-}" == "checkout" ]]; then
  die "production release forbids source checkout"
fi
if [[ "${1:-}" == "npm" && "${2:-}" == "run" && "${3:-}" == "build" ]]; then
  die "production release forbids frontend build"
fi
if [[ "${1:-}" == "docker" && "${2:-}" == "build" ]]; then
  die "production release forbids docker image build"
fi
if [[ "${1:-}" == "docker" && "${2:-}" == "compose" ]]; then
  for arg in "$@"; do
    [[ "$arg" != "build" ]] || die "production release forbids compose build"
    [[ "$arg" != "--build" ]] || die "production release forbids compose up with build"
  done
fi
if [[ "${1:-}" == "docker-compose" ]]; then
  for arg in "$@"; do
    [[ "$arg" != "build" ]] || die "production release forbids compose build"
    [[ "$arg" != "--build" ]] || die "production release forbids compose up with build"
  done
fi

{
  printf '+'
  printf ' %q' "$@"
  if [[ -n "$stdin_file" ]]; then
    [[ -f "$stdin_file" ]] || die "stdin file not found: $stdin_file"
    printf ' < %q\n' "$stdin_file"
    "$@" < "$stdin_file"
  else
    printf '\n'
    "$@"
  fi
} 2>&1 | tee -a "$transcript"
