#!/bin/sh
# dsco-claude — cross-platform installer/updater for the dsco-expert agent
# (and any dsco-specific skills) into Claude Code.
#
# POSIX sh: Linux, macOS, WSL, Git Bash. On native Windows PowerShell use
# install.ps1 instead.
#
# Usage:
#   ./install.sh [install|update|uninstall|status] [--project [DIR]] [--copy]
#
#   install    symlink the agent (and skills) into Claude Code (default)
#   update     refresh the links/copies to the current bundle (same as install)
#   uninstall  remove what this bundle installed
#   status     show the bundle version and what is currently installed
#
#   --project [DIR]  target DIR/.claude instead of ~/.claude (DIR defaults to .)
#   --copy           copy files instead of symlinking (filesystems w/o symlinks)

set -eu

usage() {
  sed -n '2,20p' "$0" | sed 's/^# \{0,1\}//'
}

# --- resolve this script's directory, following symlinks (the bundle root) ---
script="$0"
while [ -h "$script" ]; do
  d=$(cd -P "$(dirname "$script")" && pwd)
  script=$(readlink "$script")
  case "$script" in
    /*) ;;
    *) script="$d/$script" ;;
  esac
done
BUNDLE_DIR=$(cd -P "$(dirname "$script")" && pwd)

AGENT_SRC="$BUNDLE_DIR/agents/dsco-expert.md"
SKILLS_SRC="$BUNDLE_DIR/skills"
VERSION=$(cat "$BUNDLE_DIR/VERSION" 2>/dev/null || echo "unknown")

# --- parse arguments ---
CMD=""
CLAUDE_DIR="$HOME/.claude"
USE_COPY=0
while [ $# -gt 0 ]; do
  case "$1" in
    install|update|uninstall|status) CMD="$1" ;;
    --project)
      if [ $# -ge 2 ] && [ "${2#-}" = "$2" ] && [ "$2" != "install" ] \
         && [ "$2" != "update" ] && [ "$2" != "uninstall" ] \
         && [ "$2" != "status" ]; then
        proj="$2"; shift
      else
        proj="."
      fi
      CLAUDE_DIR="$proj/.claude"
      ;;
    --copy) USE_COPY=1 ;;
    -h|--help) usage; exit 0 ;;
    *) echo "dsco-claude: unknown argument '$1'" >&2; usage; exit 2 ;;
  esac
  shift
done
CMD="${CMD:-install}"

link_one() { # src dst
  _src="$1"; _dst="$2"
  mkdir -p "$(dirname "$_dst")"
  rm -rf "$_dst"
  if [ "$USE_COPY" -eq 1 ]; then
    cp -R "$_src" "$_dst"
    echo "  copied   $_dst"
  else
    ln -s "$_src" "$_dst"
    echo "  linked   $_dst -> $_src"
  fi
}

each_skill() { # calls: $1 <name> <srcdir>
  for _d in "$SKILLS_SRC"/*/; do
    [ -f "$_d/SKILL.md" ] || continue
    "$1" "$(basename "$_d")" "${_d%/}"
  done
}

do_install() {
  echo "dsco-claude v$VERSION -> $CLAUDE_DIR"
  if [ -f "$AGENT_SRC" ]; then
    link_one "$AGENT_SRC" "$CLAUDE_DIR/agents/dsco-expert.md"
  fi
  each_skill _install_skill
  echo "done."
}
_install_skill() { link_one "$2" "$CLAUDE_DIR/skills/$1"; }

do_uninstall() {
  echo "dsco-claude: removing from $CLAUDE_DIR"
  _rm "$CLAUDE_DIR/agents/dsco-expert.md"
  each_skill _uninstall_skill
  echo "done."
}
_uninstall_skill() { _rm "$CLAUDE_DIR/skills/$1"; }
_rm() {
  if [ -e "$1" ] || [ -h "$1" ]; then rm -rf "$1"; echo "  removed  $1"; fi
}

do_status() {
  echo "dsco-claude v$VERSION"
  echo "bundle:  $BUNDLE_DIR"
  echo "target:  $CLAUDE_DIR"
  if [ -f "$AGENT_SRC" ]; then _status_one "$CLAUDE_DIR/agents/dsco-expert.md"; fi
  each_skill _status_skill
}
_status_skill() { _status_one "$CLAUDE_DIR/skills/$1"; }
_status_one() {
  _t="$1"
  if [ -h "$_t" ]; then
    echo "  linked   $_t -> $(readlink "$_t")"
  elif [ -e "$_t" ]; then
    echo "  present  $_t (copy)"
  else
    echo "  missing  $_t"
  fi
}

case "$CMD" in
  install|update) do_install ;;
  uninstall)      do_uninstall ;;
  status)         do_status ;;
esac
