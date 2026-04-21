# Dotfiles Project — AI Instructions

## Ansible Dependency Rule

**Any time a config or plugin requires a system dependency, you MUST update the Ansible playbooks to install it.**

This applies to: packages, CLI tools, fonts, runtimes, language servers, luarocks, pip packages, etc.

### Required coverage for every dependency

Always add installs for all three platforms. Use the existing playbook files:

| Platform | Playbook | Module |
|----------|----------|--------|
| Arch Linux | `ansible/install_common.yml` or `ansible/install_dependencies.yml` | `pacman` preferred; `yay` (via `command`) for AUR-only packages |
| macOS | `ansible/install_common.yml` or `ansible/install_mac_dependencies.yml` | `homebrew` (CLI) or `homebrew_cask` (GUI apps) |
| Debian/Ubuntu | `ansible/install_common.yml` or `ansible/install_dependencies.yml` | `apt` preferred; `snap` only as last resort |

### GUI / Desktop applications on Debian

Debian may be headless (server). The `ubuntu_type` fact is set in `configure_ubuntu_desktop.yml`:
- `desktop` — has a GUI (gnome/plasma/ubuntu-desktop detected)
- `server` — headless, no display

**If a dependency is for a UI application** (anything that opens a window, requires a display, or is a GUI tool), guard it with:

```yaml
when: ansible_facts['os_family'] == "Debian" and ubuntu_type == "desktop"
```

CLI tools, language servers, libraries, and terminal utilities do NOT need this guard.

### Bulk package lists

`install_common.yml` has bulk package list variables at the top (`arch_packages`, `mac_homebrew_packages`, `mac_cask_packages`). Prefer adding to these lists for simple packages rather than adding individual tasks.

### Package name differences

Package names often differ across platforms. Always verify the correct name for each package manager. For example:
- imagemagick is `imagemagick` on apt/pacman, `imagemagick` on brew
- luarocks is `luarocks` on all three

---

Use `/ansible-deps` to be reminded of this workflow during a session.

## After Making Changes

**Always tell the user what command(s) to run to pick up the changes**, with a brief explanation of what each command does. Do not assume the user knows what to do next. Examples:

- Config file changes that require a process restart → tell them to kill and relaunch that process
- Shell config changes → tell them to run `exec zsh` (replaces the current shell with a fresh one, picking up new config)
- Tmux config changes → tell them to run `tmux source ~/.tmux.conf` (reloads config in the running session without restarting)
- Hyprland config changes → tell them to run `hyprctl reload` (reloads hyprland.conf in the running session; note exec-once lines do NOT re-run on reload — a full restart is needed for those)
- Neovim plugin changes → tell them to run `:Lazy sync` inside neovim

Never give a command without explaining what it does and why it's needed.

## Troubleshooting Commands Reference

Useful commands for diagnosing issues in this repo:

| Command | What it does |
|---------|--------------|
| `rg "pattern"` | Search file contents recursively from the current directory (ripgrep — faster alternative to grep). Use to find config values, function names, or package references across many files. |
| `pgrep -af <name>` | List running processes matching a name, with full command line |
| `tail -n <N> <file>` | Show last N lines of a file (useful for logs) |
| `find -L <dir> -type f` | Find files recursively; `-L` follows symlinks |
| `ls -la` | List files with permissions, ownership, and symlink targets |
| `stat <file>` | Show detailed file metadata including ownership and permissions |
| `readlink <path>` | Show where a symlink points |
| `pacman -Qq <pkg>` | Check if a package is installed on Arch |
