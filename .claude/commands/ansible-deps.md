When adding any config, plugin, or tool that requires system dependencies, always update the Ansible playbooks.

## Rules

1. **Every dependency must be installed for all three platforms**: Arch Linux, macOS, and Debian/Ubuntu.

2. **File locations**:
   - `ansible/install_dependencies.yml` — core tools (git, stow, nvim, etc.)
   - `ansible/install_common.yml` — app packages; has bulk lists at the top (`arch_packages`, `mac_homebrew_packages`, `mac_cask_packages`) — prefer adding to these lists for simple packages
   - `ansible/install_mac_dependencies.yml` — Mac bootstrap (brew setup)
   - `ansible/configure_ubuntu_desktop.yml` — Debian GUI-only installs

3. **Package managers**:
   - Arch: prefer `pacman` module with `become: yes`; use `yay` (via `command` module) only for AUR packages not available in the official repos
   - macOS: `homebrew` (CLI tools) or `homebrew_cask` (GUI apps)
   - Debian: prefer `apt` module with `become: yes`; only use `snap` as a last resort when the package is unavailable or too outdated in apt

4. **GUI applications on Debian — headless guard required**:
   The `ubuntu_type` fact (`desktop` or `server`) is set in `configure_ubuntu_desktop.yml`.
   Any package that requires a display (GUI apps, window managers, display servers) must be guarded:
   ```yaml
   when: ansible_facts['os_family'] == "Debian" and ubuntu_type == "desktop"
   ```
   CLI tools, libraries, language servers, and terminal utilities do NOT need this guard.

5. **Do not install UI apps in `install_common.yml` or `install_dependencies.yml`** — put them in `configure_ubuntu_desktop.yml` for Debian, and in the `mac_cask_packages` list for macOS.

## Checklist before finishing any config change

- [ ] Identified all new system dependencies
- [ ] Added Arch install (pacman)
- [ ] Added macOS install (homebrew or cask)
- [ ] Added Debian install (apt/snap), with headless guard if UI
- [ ] Verified package names are correct for each platform
