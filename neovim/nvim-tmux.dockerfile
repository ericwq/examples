FROM alpine:edge
LABEL maintainer="ericwq057@qq.com"

# This is the runtime package it contains
# base: bash colordiff git curl neovim tzdata htop
RUN apk add --no-cache \
#	colordiff \
        neovim \
#        git \
#        curl \
        tzdata \
        htop \
#	go \
#	ripgrep \
#	fzf \
#	g++ \
	bash \
	tmux 

ENV HOME /home/ide
ENV GOPATH $HOME/go

# Create user/group : ide/develop
RUN addgroup develop && adduser -D -h $HOME -s /bin/bash -G develop ide
#RUN chown -R ide:develop $GOPATH

USER ide:develop
WORKDIR $HOME
ENV PAHT=$PATH:$GOPATH/bin:$HOME

# Prepare for the nvim
RUN mkdir -p $HOME/.config/nvim/ 
RUN mkdir -p $GOPATH

## Copy the .init.vim : init0.vim contains only the plugin part
COPY --chown=ide:develop ./tmux.conf 		$HOME/.tmux.conf
COPY --chown=ide:develop ./entrypoint.sh 	$HOME/entrypoint.sh

RUN chmod +x $HOME/entrypoint.sh
RUN echo $HOME/entrypoint.sh
RUN echo | ls -al
RUN echo $PATH
EXPOSE 22/tcp

CMD ["./entrypoint.sh"]
