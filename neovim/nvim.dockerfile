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
# protoc py3-pip bash 
#
RUN apk add tmux colordiff curl tzdata htop fzf go ctags ccls --update

ENV HOME /home/ide
ENV GOPATH /go
ENV PATH=$PATH:$GOPATH/bin

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
COPY --chown=ide:develop ./v3/nvim/init.lua	$HOME/.config/nvim/
COPY --chown=ide:develop ./v3/nvim/lua		$HOME/.config/nvim/lua
COPY --chown=ide:develop ./custom		$HOME/.config/nvim/lua/custom

# Install the packer plugins
# https://github.com/wbthomason/packer.nvim/issues/502
#
# NvChad version
RUN nvim --headless -c 'autocmd User PackerComplete quitall' -c 'PackerSync'

CMD ["/bin/ash"]
