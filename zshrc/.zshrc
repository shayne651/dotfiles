eval "$(starship init zsh)"

bindkey '^[[Z' reverse-menu-complete
export TERM=xterm-256color

source ~/.config/zsh-autosuggestions/zsh-autosuggestions.zsh
source ~/.config/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh
source ~/.config/zsh-autocomplete/zsh-autocomplete.plugin.zsh

# Add RVM to PATH for scripting. Make sure this is the last PATH variable change.
export PATH="$PATH:$HOME/.rbenv/bin:$PATH"
eval "$(rbenv init -)"

# Reefmonitor env
export LOCAL_DEV=true
export GOOGLE_APPLICATION_CREDENTIALS="$HOME/.config/reefmonitor.json"

# ESPTOOL auto complete
export PATH="$PATH:/Users/shaynetaylor/Library/Python/3.9/bin"


alias python="python3"
alias py="python3"
export PATH="/opt/homebrew/opt/openjdk/bin:$PATH"
