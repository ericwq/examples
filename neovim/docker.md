docker build -t nide:0.1 -f nvim-tmux.dockerfile .

docker run -it -d -h ggg --env TZ=Asia/Shanghai -u ide --name ggg -p 8654:22 nide:0.1 bash

ssh ide@localhost -p 8654 -t "tmux a -t golangide”

tmux new-session -s "IDE" -n "editor" -d “nvim"
tmux attach-session -t "IDE"
