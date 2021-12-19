FROM alpine:edge
LABEL maintainer="ericwq057@qq.com"

# This is the base pacakges for neovim 
# https://github.com/NvChad/NvChad
#
RUN apk add git nodejs neovim ripgrep alpine-sdk --update

# additional pacakges for golang IDE
# mainly go, ccls, tmux, fzf
#
# consider add the following pacakges:
# protoc py3-pip bash ctags ccls 
#
RUN apk add tmux colordiff curl tzdata htop fzf go --update

ENV HOME /home/ide
ENV GOPATH $HOME/go

# Create user/group 
# ide/develop
#
RUN addgroup develop && adduser -D -h $HOME -s /bin/ash -G develop ide

USER ide:develop
WORKDIR $HOME
ENV PATH=$PATH:$GOPATH/bin

# Prepare for the nvim
RUN mkdir -p $HOME/.config/nvim/lua && mkdir -p $GOPATH

# Install go language server
RUN go install golang.org/x/tools/gopls@latest

# TODO: The source script
# https://hhoeflin.github.io/2020/08/19/bash-in-docker/
#
COPY --chown=ide:develop ./profile 		$HOME/.profile

# The clipboatd support for vim and tmux
# https://sunaku.github.io/tmux-yank-osc52.html
#
COPY --chown=ide:develop ./tmux.conf 		$HOME/.tmux.conf
COPY --chown=ide:develop ./vimrc 		$HOME/.config/nvim/vimrc
COPY --chown=ide:develop ./yank 		$GOPATH/bin/yank
RUN chmod +x $GOPATH/bin/yank

# The neovim configuration
# based on https://github.com/NvChad/NvChad
#
COPY --chown=ide:develop ./v3/nvim/init.lua	$HOME/.config/nvim/
COPY --chown=ide:develop ./v3/nvim/lua		$HOME/.config/nvim/lua

# Install the packer plugins
# 
RUN  nvim --headless -c 'autocmd User PackerComplete quitall' -c 'PackerSync'

CMD ["/bin/ash"]
