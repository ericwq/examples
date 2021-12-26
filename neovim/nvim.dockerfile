#------------------------------ NOTICE ------------------------------
# please perform the following command before the image build
# this command must be done in the dockerfile directory.
#
# git clone https://github.com/NvChad/NvChad.git
#
FROM alpine:edge
LABEL maintainer="ericwq057@qq.com"

# This is the base pacakges for neovim 
# https://github.com/NvChad/NvChad
#
# tree-sitter needs tree-sitter-cli, nodejs
# telscope needs ripgrep, fzf, fd
#
RUN apk add git neovim neovim-doc tree-sitter-cli nodejs ripgrep fzf fd ctags alpine-sdk --update

# additional pacakges for the IDE
# mainly go, ccls, tmux
#
# consider to add the following pacakges:
# py3-pip bash
#
RUN apk add tmux colordiff curl tzdata htop go ccls protoc --update

ENV HOME=/home/ide
ENV GOPATH /go
ENV PATH=$PATH:$GOPATH/bin
ENV ENV=$HOME/.profile

# Create user/group 
# ide/develop
#
RUN addgroup develop && adduser -D -h $HOME -s /bin/ash -G develop ide
RUN mkdir -p $GOPATH && chown -R ide:develop $GOPATH

USER ide:develop
WORKDIR $HOME

# Prepare for the nvim
RUN mkdir -p $HOME/.config/nvim/lua && mkdir -p $GOPATH

# Install go language server
RUN go install golang.org/x/tools/gopls@latest

# Install golangci-lint
# https://golangci-lint.run/usage/install/
# RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.43.0

# The source script
# https://hhoeflin.github.io/2020/08/19/bash-in-docker/
# https://unix.stackexchange.com/questions/176027/ash-profile-configuration-file
#
# ENV=$HOME/.profile
#
COPY --chown=ide:develop ./profile 		$HOME/.profile

# The clipboatd support for vim and tmux
# https://sunaku.github.io/tmux-yank-osc52.html
#
COPY --chown=ide:develop ./tmux.conf 		$HOME/.tmux.conf
COPY --chown=ide:develop ./vimrc 		$HOME/.config/nvim/vimrc
COPY --chown=ide:develop ./yank 		$GOPATH/bin/yank
RUN chmod +x $GOPATH/bin/yank

# Install packer.vim
# PackerSync command will install packer.vim automaticlly, while the
# installation  will stop to wait for user <Enter> input.
# So we install packer manually.
#
# we also move it to 'opt' directory instead of 'start' directory
# because NvChad install packer in 'opt' directory
# https://github.com/wbthomason/packer.nvim
#
RUN git clone --depth 1 https://github.com/wbthomason/packer.nvim \
	~/.local/share/nvim/site/pack/packer/opt/packer.nvim

# The neovim configuration
# based on https://github.com/NvChad/NvChad
#

COPY --chown=ide:develop ./NvChad/init.lua	$HOME/.config/nvim/
COPY --chown=ide:develop ./NvChad/lua		$HOME/.config/nvim/lua
COPY --chown=ide:develop ./custom		$HOME/.config/nvim/lua/custom


# Install the packer plugins
# https://github.com/wbthomason/packer.nvim/issues/502
#
# NvChad version
RUN nvim --headless -c 'autocmd User PackerComplete quitall' -c 'PackerSync'

# Install treesitter language parsers
# See :h packages
# https://github.com/wbthomason/packer.nvim/issues/237
#
RUN nvim --headless -c 'packadd nvim-treesitter' -c 'TSInstallSync go c cpp yaml lua json dockerfile markdown' +qall
CMD ["/bin/ash"]
