# dsco skills

dsco-specific Claude Code skills live here. There are none yet — this folder
is the home for them when they exist.

Any skill added here must follow the bundle's version-targeting convention:

1. Declare the dsco version it targets in its `SKILL.md` frontmatter:
   ```yaml
   x-dsco-target: v1.4.0-rc.1
   x-bundle-version: 1.4.0-rc.1
   ```
2. Gate version-specific advice: before using an API, check the user's
   pinned dsco version (`go.mod`) and offer the upgrade if the feature needs
   a newer version. See the `dsco-expert` agent's *Version targeting* section
   for the exact behavior to mirror.
3. Bump alongside the bundle `VERSION` when the targeted dsco version changes.
