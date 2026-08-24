## Setup

### 1. Clone (with submodules)
```sh
git clone --recurse-submodules git@github.com:shayne651/dotfiles.git
```

If you already cloned without `--recurse-submodules`, initialize them after the fact:
```sh
git submodule update --init --recursive
```

### 2. Run Ansible
```sh
ansible-playbook ansible/main.yml -u $USER --ask-become-pass

Run as your normal login user, not with `sudo`; the playbook uses sudo only
for system changes so dotfiles are installed into your user home directory.
```

Or via Make:
```sh
make install
```

## Submodules

| Path | Description |
|------|-------------|
| `tmux/.tmux/plugins/tpm` | Tmux Plugin Manager |
| `tmux/.tmux/plugins/catppuccin/tmux` | Catppuccin theme for tmux |
| `zshrc/.config/zsh-autosuggestions` | Zsh autosuggestions |
| `zshrc/.config/zsh-syntax-highlighting` | Zsh syntax highlighting |
| `zshrc/.config/zsh-autocomplete` | Zsh autocomplete |

To pull the latest upstream changes for all submodules:
```sh
git submodule update --remote --recursive
```

Or via Make:
```sh
make update
```

## Make targets

| Target | What it does |
|--------|--------------|
| `make install` | Run the Ansible playbook |
| `make update` | Pull latest repo changes and update all submodules |
| `make submodules` | Initialize/fetch submodules (if cloned without `--recurse-submodules`) |

## What's included

- **Zsh** — config, plugins (autosuggestions, syntax highlighting, autocomplete), Starship prompt
- **Neovim** — full config via stow
- **Tmux** — config + TPM plugins + Catppuccin theme
- **Ghostty** — terminal config (platform-specific opacity via Ansible)
- **Hyprland** — window manager config (Arch only)

## Dependencies

- `ansible`
- `git`

All other dependencies are installed by the Ansible playbook.
