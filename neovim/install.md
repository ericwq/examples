# Neovim 0.6 setup
Try to find the better solution for neovim IDE (golang, c, c++, java, lua, html, css, vim script, makrdown, javascript)

## Base image
```
docker pull alpine:edge
```

## true color test:
```
curl -s https://raw.githubusercontent.com/JohnMorales/dotfiles/master/colors/24-bit-color.sh | ash
```
## Guide for neovim and lua
- [Color scheme - sonokai ](https://github.com/sainnhe/sonokai)
- [Neovim 0.5 features and the switch to init.lua](https://oroques.dev/notes/neovim-init/)
- [LSP- go language server](https://github.com/golang/tools/tree/master/gopls)
- [学习Neovim全配置，逃离VSCode](https://zhuanlan.zhihu.com/p/434727338)
- [iTerm2](https://sourabhbajaj.com/mac-setup/iTerm/)
- [Go neovim configuration](https://www.getman.io/posts/programming-go-in-neovim/)
- [Base neovim configuration](https://github.com/brainfucksec/neovim-lua)

## [Moving to modern Neovim](https://toroid.org/modern-neovim#update)
- [Package management - packer](https://github.com/wbthomason/packer.nvim)
- [Telescope](https://github.com/nvim-telescope/telescope.nvim)
- [Status line? - lualine](https://github.com/hoob3rt/lualine.nvim)
- [Key mappings? - which-key.nvim](https://github.com/folke/which-key.nvim)
- [LuaSnip](https://github.com/L3MON4D3/LuaSnip)
- [NvimTree](https://github.com/kyazdani42/nvim-tree.lua)
- [Treesitter](https://tree-sitter.github.io/tree-sitter/)
- [Treesitter - romgrk/nvim-treesitter-context](https://github.com/romgrk/nvim-treesitter-context)
- [LSP - nvim-lspconfig](https://github.com/neovim/nvim-lspconfig)
- [LSP - nvim-cmp](https://github.com/hrsh7th/nvim-cmp)
- [LSP - symbols-outline.nvim](https://github.com/simrat39/symbols-outline.nvim)
- [LSP - lsp-signature](https://github.com/ray-x/lsp_signature.nvim)
- [Debug? - nvim-dap](https://github.com/mfussenegger/nvim-dap)
- [Debug? - nvim-dap-ui](https://github.com/rcarriga/nvim-dap-ui)

## .profile
```
$more .profile 
export GOPATH=/go
export PATH=$PATH:$GOPATH/bin
export PS1='\u@\h:\w $ '
alias vi=nvim
```

## [Nert font support](https://github.com/ryanoasis/nerd-fonts#glyph-sets)
```
% brew tap homebrew/cask-fonts
% brew install --cask font-hack-nerd-font
```

## docker file
- apk add neovim neovim-doc go git curl tzdata htop python3 fzf
- apk add tree-sitter nodejs
- apk add make musl-dev g++
- export HOME=/home/ide
- export GOPATH=/go
- mkdir /go
- addgroup develop && adduser -D -h $HOME -s /bin/ash -G develop ide
- chown -R ide:develop $GOPATH
- su - ide
- git clone --depth 1 https://github.com/wbthomason/packer.nvim ~/.local/share/nvim/site/pack/packer/start/packer.nvim


## others
- git clone https://github.com/savq/paq-nvim.git "${XDG_DATA_HOME:-$HOME/.local/share}"/nvim/site/pack/paqs/opt/paq-nvim
- go install golang.org/x/tools/gopls@latest
- mkdir -p ~/.config/nvim/lua
- copy init.lua only keep the basics.lua
- git clone https://github.com/optimizacija/neovim-config.git
- cd
- makedir .config
- cd .config
- ln -s ../neovim-config/ nvim
