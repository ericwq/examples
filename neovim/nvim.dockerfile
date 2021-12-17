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
	fzf \
	g++ \
#	bash \
	tmux 

ENV HOME /home/ide
ENV GOPATH $HOME/go

# Create user/group : ide/develop
RUN addgroup develop && adduser -D -h $HOME -s /bin/ash -G develop ide
#RUN chown -R ide:develop $GOPATH

USER ide:develop
WORKDIR $HOME
ENV PAHT=$PATH:$GOPATH/bin

# Prepare for the nvim
RUN mkdir -p $HOME/.config/nvim/lua/plugins && mkdir -p $GOPATH

## Copy the tmux.conf 
COPY --chown=ide:develop ./tmux.conf 		$HOME/.tmux.conf
COPY --chown=ide:develop ./yank 		$GOPATH/bin/yank
COPY --chown=ide:develop ./vimrc 		$HOME/.config/nvim/vimrc

RUN chmod +x $GOPATH/bin/yank

CMD ["/bin/ash"]
