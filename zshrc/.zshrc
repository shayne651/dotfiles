export PATH="/opt/homebrew/bin:/opt/homebrew/sbin:$HOME/.local/bin:$HOME/.local/bin:$PATH"

# Raise the per-process open-file limit; macOS default (256) is too low
# when running zsh plugins + tmux with multiple panes
ulimit -n 65536

[[ -z "$TMUX" ]] && command -v fastfetch &>/dev/null && fastfetch

eval "$(starship init zsh)"

# Arch Linux terminals send raw bracketed paste markers that zle doesn't strip
if [[ "$(uname)" != "Darwin" ]]; then
  unset zle_bracketed_paste
fi

bindkey '^[[Z' reverse-menu-complete
export TERM=xterm-256color

# Don't expand relative paths (e.g. ./foo) to absolute paths on tab completion
zstyle ':completion:*' expand prefix suffix

source ~/.config/zsh-autosuggestions/zsh-autosuggestions.zsh
source ~/.config/zsh-autocomplete/zsh-autocomplete.plugin.zsh
# zsh-syntax-highlighting must be sourced last — it hooks into ZLE at source time
# and will break/leak FDs if any plugin loads after it
source ~/.config/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh

# Add RVM to PATH for scripting. Make sure this is the last PATH variable change.
export PATH="$HOME/.rbenv/bin:$PATH"
eval "$(rbenv init -)"

# Load user specific zshrc
if [[ -f "$HOME/.zshrc-specific" ]]; then
  source "$HOME/.zshrc-specific"
fi

# Preview images in a directory with chafa (press any key to advance, q to quit)
# Linux/Arch only — wallpapers and chafa are not set up on macOS
if [[ "$(uname)" != "Darwin" ]]; then
  preview() {
    local dir="${1:-.}"
    local images=("${dir}"/**/*.{jpg,jpeg,png,gif,webp}(N))
    if [[ ${#images[@]} -eq 0 ]]; then
      echo "No images found in $dir"
      return 1
    fi
    for img in "${images[@]}"; do
      echo "\033[1m$img\033[0m"
      chafa --format kitty --size "${COLUMNS}x$((LINES - 2))" "$img"
      read -sk1 key
      [[ "$key" == "q" ]] && break
      clear
    done
  }
fi

# Save (rename) the current tmux session
tmux_save() {
  if [[ -z "$1" ]]; then
    echo "Usage: tmux_save <name>"
    return 1
  fi
  tmux rename-session "$1"
}

# Attach to a named tmux session, creating it if it doesn't exist
tmux_start() {
  if [[ -z "$1" ]]; then
    echo "Usage: tmux_start <name>"
    return 1
  fi
  tmux new-session -As "$1"
}

# History settings
export HISTFILE=~/.zsh_history
export HISTSIZE=1000
export SAVEHIST=1000
setopt appendhistory
setopt extendedhistory

alias zshrc="nvim ~/.zshrc-specific;source ~/.zshrc-specific"
