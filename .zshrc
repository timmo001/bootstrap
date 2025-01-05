export PATH=$HOME/bin:/usr/local/bin:$PATH
export ZSH="$HOME/.oh-my-zsh"
ZSH_THEME="robbyrussell"
plugins=(git zsh-autosuggestions zsh-syntax-highlighting fast-syntax-highlighting zsh-autocomplete)
source $ZSH/oh-my-zsh.sh
export PATH=$PATH:/usr/local/go/bin
export PATH=/home/aidan/.local/share/fnm/bin:$PATH
eval "`fnm env`"
