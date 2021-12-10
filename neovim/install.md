## Neovim 0.6 setup

alpine:edge

color: curl -s https://raw.githubusercontent.com/JohnMorales/dotfiles/master/colors/24-bit-color.sh | ash

1. link: https://sourabhbajaj.com/mac-setup/iTerm/
2. link: https://www.getman.io/posts/programming-go-in-neovim/
3. link: https://toroid.org/modern-neovim#update
4. link: https://github.com/hrsh7th/nvim-cmp
5. link: https://oroques.dev/notes/neovim-init/
6. link: https://github.com/golang/tools/tree/master/gopls

```
$more .profile 
export GOPATH=/go
export PATH=$PATH:$GOPATH/bin
export PS1='\u@\h:\w $ '
alias vi=nvim
```

- apk add neovim neovim-doc go git curl tzdata htop python3 fzf
- export HOME=/home/ide
- export GOPATH=/go
- mkdir /go
- addgroup develop && adduser -D -h $HOME -s /bin/ash -G develop ide
- chown -R ide:develop $GOPATH
- su - ide
- git clone https://github.com/savq/paq-nvim.git \
    "${XDG_DATA_HOME:-$HOME/.local/share}"/nvim/site/pack/paqs/opt/paq-nvim
- 
- git clone --depth 1 https://github.com/wbthomason/packer.nvim ~/.local/share/nvim/site/pack/packer/start/packer.nvim
- mkdir -p ~/.config/nvim/lua

* copy init.lua only keep the basics.lua

* git clone https://github.com/optimizacija/neovim-config.git
- cd
- makedir .config
- cd .config
- ln -s ../neovim-config/ nvim

