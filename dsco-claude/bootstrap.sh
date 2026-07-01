#!/bin/sh
# dsco-claude bootstrap — download the bundle from GitHub and install the skills
# into Claude Code, without a manual checkout. POSIX sh: Linux, macOS, WSL, Git
# Bash. On native Windows PowerShell use bootstrap.ps1.
#
# One-liner:
#   curl -fsSL https://raw.githubusercontent.com/byte4ever/dsco/master/dsco-claude/bootstrap.sh | sh
#   wget -qO-  https://raw.githubusercontent.com/byte4ever/dsco/master/dsco-claude/bootstrap.sh | sh
#
# Pin a version / pass install options (note the `-s --`):
#   curl -fsSL <url>/bootstrap.sh | sh -s -- --ref v1.4.0 --copy
#
# Env vars (alternative to flags):
#   DSCO_CLAUDE_REF   git ref to fetch (branch / tag / sha; default: master)
#   DSCO_CLAUDE_HOME  where to place the bundle (default: ~/.dsco-claude)
#
# Flags consumed here: --ref <ref>, --home <dir>. Any other flag (e.g. --copy,
# --project) is forwarded to the bundle's install.sh.

set -eu

REPO="byte4ever/dsco"
REF="${DSCO_CLAUDE_REF:-master}"
HOME_DIR="${DSCO_CLAUDE_HOME:-$HOME/.dsco-claude}"
INSTALL_ARGS=""

while [ $# -gt 0 ]; do
  case "$1" in
    --ref) REF="$2"; shift ;;
    --ref=*) REF="${1#--ref=}" ;;
    --home) HOME_DIR="$2"; shift ;;
    --home=*) HOME_DIR="${1#--home=}" ;;
    *) INSTALL_ARGS="$INSTALL_ARGS $1" ;;
  esac
  shift
done

if command -v curl >/dev/null 2>&1; then
  fetch() { curl -fsSL "$1"; }
elif command -v wget >/dev/null 2>&1; then
  fetch() { wget -qO- "$1"; }
else
  echo "dsco-claude: need curl or wget on PATH" >&2
  exit 1
fi

url="https://github.com/$REPO/archive/$REF.tar.gz"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT INT TERM

echo "dsco-claude: fetching $REPO@$REF ..."
fetch "$url" | tar -xzf - -C "$tmp"

# the tarball's top dir is dsco-<ref>; locate the dsco-claude subtree inside it
src=""
for d in "$tmp"/*/dsco-claude; do
  if [ -d "$d" ]; then src="$d"; break; fi
done
if [ -z "$src" ]; then
  echo "dsco-claude: dsco-claude/ not found in $REPO@$REF" >&2
  echo "  (the bundle ships on master and on releases that include it)" >&2
  exit 1
fi

mkdir -p "$HOME_DIR"
# refresh the bundle contents in place
for item in skills install.sh install.ps1 bootstrap.sh bootstrap.ps1 VERSION README.md CHANGELOG.md; do
  rm -rf "$HOME_DIR/$item"
done
cp -R "$src"/. "$HOME_DIR"/
chmod +x "$HOME_DIR/install.sh" 2>/dev/null || true

echo "dsco-claude: bundle placed in $HOME_DIR"
# shellcheck disable=SC2086
sh "$HOME_DIR/install.sh" install $INSTALL_ARGS
