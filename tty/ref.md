# Terminal emulator

## SSP design in [mosh research paper](https://mosh.org/mosh-paper.pdf)

SSP is organized into two layers. A datagram layer sends UDP packets over the network, and a transport layer is responsible for conveying the current object state to the remote host.

### Datagram Layer

The datagram layer maintains the “roaming” connection. It accepts opaque payloads from the transport layer, prepends an incrementing sequence number, encrypts the packet, and sends the resulting ciphertext in a UDP datagram. It is responsible for estimating the timing characteristics of the link and keeping track of the client’s current public IP address.

- Client roaming.
  - Every time the server receives an authentic datagram from the client with a sequence number greater than any before, it sets the packet’s source IP address and UDP port number as its new “target.” See [the client implementation](client.md#how-does-the-client-roam) and [the server implementation](client.md#how-does-the-server-support-client-roam).
- Estimating round-trip time and RTT variation. See [RTT and RTTVAR calculation](client.md#how-to-receive-datagram-from-socket).
  - Every outgoing datagram contains a millisecond timestamp and an optional “timestamp reply,” containing the most recently received timestamp from the remote host. See [the implementation](client.md#how-to-send-a-packet)
  - SSP adjusts the “timestamp reply” by the amount of time since it received the corresponding timestamp.

### Transport Layer

The transport layer synchronizes the contents of the local state to the remote host, and is agnostic to the type of objects sent and received.

- Transport sender behavior
  - The transport sender updates the receiver to the current state of the object by sending an Instruction: a self-contained message listing the source and target states and the binary “diff” between them.
  - This “diff” is a logical one, calculated by the object implementation.
  - The ultimate semantics of the protocol depend on the type of object, and are not dictated by SSP.
  - For user inputs, the diff contains every intervening keystroke.
  - For screen states, it is only the minimal message that transforms the client’s frame to the current one.
- Transport sender timing
  - It is not required to send every octet it receives from the host and can modulate the “frame rate” based on network conditions.
  - The minimum interval between frames is set at half the smoothed RTT estimate, so there is about one Instruction in flight to the receiver at any time.
  - The transport sender uses delayed acks, similar to TCP, to cut down on excess packets.
  - The server also pauses from the first time its object has changed before sending off an Instruction, because updates to the screen tend to clump together, and it would be wasteful to send off a new frame with a partial update and then have to wait the full “frame rate” interval before sending another.
  - SSP sends an occasional heartbeat to allow the server to learn when the client has roamed to a new IP address, and to allow the client to warn the user when it hasn’t recently heard from the server.

## SSP design in [github.com](https://github.com/mobile-shell/mosh/issues/1087#issuecomment-641801909)

- The sender always sends diffs. There is no "full update" instruction.
- The diff has three important fields: the source, target, and throwaway number. See [the `Instruction` implementation](client.md#how-to-send-data-in-fragments).
  - The target of the diff is always the current sender-side state.
  - The throwaway number of the diff is always the most recent state that has been explicitly acknowledged by the receiver.
  - The source of the diff is allowed to be:
    - The most recent state that was explicitly acknowledged by the receiver
    - Any more-recent state that the sender thinks the receiver probably will have by the time the current diff arrives
  - The sender gets to make this choice based on what is most efficient and how likely it thinks the receiver actually has the more-recent state.
- Upon receiving a diff, the receiver throws away anything older than the throwaway number and attempts to apply the diff.
  - If it has the source state, and if the target state is newer than the receiver's current state, it succeeds and then acknowledges the new target state.
  - Otherwise it fails to apply the diff and just acks its current state number.

## reference project

### st

[st - simple terminal project](https://st.suckless.org/), a c project, adding dependency packages with root privilege.

```sh
ssh root@localhost
# apk add ncurses-terminfo-base fontconfig-dev freetype-dev libx11-dev libxext-dev libxft-dev
```

switch to the ide user of `nvide`.

```sh
ssh ide@localhost
git clone https://git.suckless.org/st
cd st
make clean
bear -- make st
```

now you can check the source code of `st` via [nvide](https://github.com/ericwq/nvide).

### mosh

[mosh 1.3.2 - mobile shell project](https://mosh.org/), a C++ project, adding dependency packages with root privilege.

```sh
ssh root@localhost
apk add ncurses-dev zlib-dev openssl-dev perl-dev perl-io-tty protobuf-dev automake autoconf libtool gzip
```

switch to the ide user of `nvide`. Download mosh from [mosh-1.3.2.tar.gz](https://mosh.org/mosh-1.3.2.tar.gz)

```sh
ssh ide@localhost
curl -O https://mosh.org/mosh-1.3.2.tar.gz
tar xvzf mosh-1.3.2.tar.gz
cd mosh-1.3.2
./configure
bear -- make
```

### mosh 1.4

[mosh 1.4 - mobile shell project](https://mosh.org/), a C++ project, adding dependency packages with root privilege. Note the following script only works on `alpine:3.18`. `alpine:edge` will cause protobuf compile problem.

```sh
ssh root@localhost
apk add build-base autoconf automake gzip libtool ncurses-dev openssl-dev perl-dev perl-io-tty protobuf-dev zlib-dev perl-doc bear
```

switch to the ide user of `nvide`.

```sh
ssh ide@localhost
git clone https://github.com/mobile-shell/mosh.git
cd mosh
./autogen.sh
./configure
bear -- make
```

### zutty

[zutty project](https://github.com/tomszilagyi/zutty), terminal emulator rendering through OpenGL ES Compute Shaders]

```sh
ssh root@localhost
apk add libxmu-dev mesa-dev freetype-dev
```

switch to the ide user of `nvide`.

```sh
ssh ide@localhost
git clone https://github.com/tomszilagyi/zutty.git
cd zutty
curl -O https://raw.githubusercontent.com/socrocket/core/master/core/waf/clang_compilation_database.py
```

Please skip this step. It doesn't work as expected.

- [Compilation database](https://docs.embold.io/compilation-database/)
- modify the `wscript` file, add the following line to it according to [Compilation database](https://sarcasm.github.io/notes/dev/compilation-database.html#waf)

```python
def configure(conf):
    conf.load('compiler_cxx')
    …
    conf.load('clang_compilation_database')
```

run the following commands. Note the place of `build/compile_commands.json`.

```sh
./waf distclean
./waf configure
bear -- ./waf
```

### utmps

[utmps](https://skarnet.org/software/utmps/), a C project, adding dependency packages with root privilege.

```sh
ssh root@localhost
apk add skalibs-dev utmps-dev
```
switch to the ide user of `nvide`.

```sh
ssh ide@localhost
git clone https://github.com/tomszilagyi/zutty.git
git clone git://git.skarnet.org/utmps
cd utmps
./configure
bear -- make
```

## reference doc

How the terminal works? Who is responsible for terminal rendering? Does GPU-rendering in terminal matter?

- [Documentation for State Synchronization Protocol](https://github.com/mobile-shell/mosh/issues/1087)
- [Text-Terminal-HOWTO](https://tldp.org/HOWTO/Text-Terminal-HOWTO.html)
- [Linux terminals, tty, pty and shell](https://dev.to/napicella/linux-terminals-tty-pty-and-shell-192e)
- [Linux terminals, tty, pty and shell - part 2](https://dev.to/napicella/linux-terminals-tty-pty-and-shell-part-2-2cb2)
- [How does a Linux terminal work?](https://unix.stackexchange.com/questions/79334/how-does-a-linux-terminal-work)
- [How Zutty works: Rendering a terminal with an OpenGL Compute Shader](https://tomscii.sig7.se/2020/11/How-Zutty-works)
- [A totally biased comparison of Zutty (to some better-known X terminal emulators)](https://tomscii.sig7.se/2020/12/A-totally-biased-comparison-of-Zutty)
- [A look at terminal emulators, part 1](https://lwn.net/Articles/749992/)
- [A look at terminal emulators, part 2](https://lwn.net/Articles/751763/)
- [High performant 2D renderer in a terminal](https://blog.ghaiklor.com/2020/07/27/high-performant-2d-renderer-in-a-terminal/)
- [The TTY demystified](http://www.linusakesson.net/programming/tty/)
- [Control sequence](https://ttssh2.osdn.jp/manual/4/en/about/ctrlseq.html#ESC)
- [The ASCII Character Set](https://www.w3schools.com/charsets/ref_html_ascii.asp#:~:text=The%20ASCII%20Character%20Set&text=ASCII%20is%20a%207%2Dbit,are%20all%20based%20on%20ASCII.)

### C++ reference

- [c++ reference](https://www.cplusplus.com/reference/)
- [c++ grammar](https://www.runoob.com/cplusplus/cpp-modifier-types.html)

### typing

- [Typing with pleasure](https://pavelfatin.com/typing-with-pleasure/)
- [Measured: Typing latency of Zutty (compared to others)](https://tomscii.sig7.se/2021/01/Typing-latency-of-Zutty)

### clangd format

- [Clang-Format Style Options](https://clang.llvm.org/docs/ClangFormatStyleOptions.html)
- [clangd format generator](https://zed0.co.uk/clang-format-configurator/)

### reference links

- [Using (neo)vim for C++ development](https://idie.ru/posts/vim-modern-cpp)
- [Getting Started with Mosh (Mobile Shell)](https://bitlaunch.io/blog/getting-started-with-mosh/)
- [example language server](https://github.com/ChrisAmelia/dotfiles/blob/master/nvim/lua/lsp.lua#L108-L120)
- [nvim-lua/kickstart.nvim](https://github.com/nvim-lua/kickstart.nvim)
- [Add, Delete And Grant Sudo Privileges To Users In Alpine Linux](https://ostechnix.com/add-delete-and-grant-sudo-privileges-to-users-in-alpine-linux/)
- [Why is GO111MODULE everywhere, and everything about Go Modules](https://maelvls.dev/go111module-everywhere/#go111module-with-go-117)
- [Understanding go.mod and go.sum](https://faun.pub/understanding-go-mod-and-go-sum-5fd7ec9bcc34)
- [spellsitter.nvim](https://github.com/lewis6991/spellsitter.nvim)
- [Neovim Tips for a Better Coding Experience](https://alpha2phi.medium.com/neovim-tips-for-a-better-coding-experience-3d0f782f034e)
- [Neovim - Treesitter Syntax Highlighting](https://www.youtube.com/watch?v=hkxPa5w3bZ0)
  now you can check the source code of `mosh` via [nvide](https://github.com/ericwq/nvide).
