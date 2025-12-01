eval "$(starship init zsh)"

bindkey '^[[Z' reverse-menu-complete
export TERM=xterm-256color

source ~/.config/zsh-autosuggestions/zsh-autosuggestions.zsh
source ~/.config/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh
source ~/.config/zsh-autocomplete/zsh-autocomplete.plugin.zsh

# Add RVM to PATH for scripting. Make sure this is the last PATH variable change.
export PATH="$PATH:$HOME/.rbenv/bin:$PATH"
eval "$(rbenv init -)"

# Load user specific zshrc
if [[ -f "$HOME/.zshrc-specific" ]]; then
  source "$HOME/.zshrc-specific"
fi
