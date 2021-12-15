FROM alpine:edge
LABEL maintainer="ericwq057@qq.com"

# This is the runtime package it contains
# base: bash colordiff git curl neovim tzdata htop
RUN apk add --no-cache \
	colordiff \
        neovim \
        git \
        curl \
        tzdata \
        htop \
	go \
	ripgrep \
	tmux \
	g++ \
	fzf 

ENV HOME /home/ide
ENV GOPATH /go

# Create user/group : ide/develop
RUN addgroup develop && adduser -D -h $HOME -s /bin/bash -G develop ide
RUN mkdir -p $GOPATH
RUN chown -R ide:develop $GOPATH

USER ide:develop
WORKDIR $HOME

# Prepare for the nvim
RUN mkdir -p $HOME/.config/nvim/ 

## Copy the .init.vim : init0.vim contains only the plugin part
COPY --chown=ide:develop ./tmux.conf 		$HOME/.tmux.conf
COPY --chown=ide:develop ./entrypoint.sh 	$HOME/entrypoint.sh

RUN chmod +x $HOME/entrypoint.sh
ENTRYPOINT ["./entrypoint.sh"]
EXPOSE 22/tcp
