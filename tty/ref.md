# Terminal emulator

## SSP design in [mosh research paper](https://mosh.org/mosh-paper.pdf)

- SSP is organized into two layers. A datagram layer sends UDP packets over the network, and a transport layer is responsible for conveying the current object state to the remote host.

### Datagram Layer

The datagram layer maintains the “roaming” connection. It accepts opaque payloads from the transport layer, prepends an incrementing sequence number, encrypts the packet, and sends the resulting ciphertext in a UDP datagram. It is responsible for estimating the timing characteristics of the link and keeping track of the client’s current public IP address.

- Client roaming.
  - Every time the server receives an authentic datagram from the client with a sequence number greater than any before, it sets the packet’s source IP address and UDP port number as its new “target.”
- Estimating round-trip time and RTT variation.
  - Every outgoing datagram contains a millisecond timestamp and an optional “timestamp reply,” containing the most recently received timestamp from the remote host.
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
- The diff has three important fields: the source, target, and throwaway number.
  - The target of the diff is always the current sender-side state.
  - The throwaway number of the diff is always the most recent state that has been explicitly acknowledged by the receiver.
  - The source of the diff is allowed to be:
    - The most recent state that was explicitly acknowledged by the receiver
    - Any more-recent state that the sender thinks the receiver probably will have by the time the current diff arrives
  - The sender gets to make this choice based on what is most efficient and how likely it thinks the receiver actually has the more-recent state.
- Upon receiving a diff, the receiver throws away anything older than the throwaway number and attempts to apply the diff.
  - If it has the source state, and if the target state is newer than the receiver's current state, it succeeds and then acknowledges the new target state.
  - Otherwise it fails to apply the diff and just acks its current state number.

## mosh coding analysis

### mosh-client.cc

In `main` function, `STMClient` is the core to start the `mosh` client.

```cpp
  try {
    STMClient client( ip, desired_port, key, predict_mode, verbose );
    client.init();

    try {
      success = client.main();
    } catch ( ... ) {
      client.shutdown();
      throw;
    }

    client.shutdown();
  } catch ( const Network::NetworkException &e ) {
    fprintf( stderr, "Network exception: %s\r\n",
	     e.what() );
    success = false;
  }
```

### STMClient constructor

- Set the `ip`,`port`,`key`,`predict_mode`,`verbose` parameter. Here, `network` is `NULL`.
- `Overlay::OverlayManager overlays` is initialized in construction function.
  - `overlays` contains `NotificationEngine`, `PredictionEngine`, `TitleEngine`. The design of these engine is unclear.
- `Terminal::Display display` is initialized in construction function.
  - `display` uses `Ncurses` libtool to setup the terminal.

### STMClient::init

- Check whether the client terminal support utf8 locale, via `is_utf8_locale()`, `locale_charset()`, `nl_langinfo()`.
- Get the `termios` struct for `STDIN`, via `tcgetattr()`.
- Set `IUTF8` flag for `termios`.
- Set the terminal for `STDIN` to raw mode, via `cfmakeraw()`.
- Set the terminal for `STDOUT` in application cursor key mode. `Display::open()`.
  - In Application mode, the cursor keys generate escape sequences that the application uses for its own purpose.
  - Application Cursor Keys mode is a way for the server to change the control sequences sent by the arrow keys. In normal mode, the arrow keys send `ESC [A` through to `ESC [D`. In application mode, they send `ESC OA` through to `ESC OD`.
- Set terminal window title, via `overlays.set_title_prefix()`. ?
- Set escape key string, via `overlays.get_notification_engine().set_escape_key_string()`. ?

### STMClient::main_init

In `client.main()`, `main_init()` is called to init the `mosh` client.

- Register signal handler for `SIGWINCH`, `SIGTERM`, `SIGINT`, `SIGHUP`, `SIGPIPE`, `SIGCONT`.
  - `sel.add_signal()` disposition the above signals. It blocks the signal outside of `pselect()`. In `pselect()`, the signal mask is replaced by a empty mask set.
- Get the window size for `STDIN` , via `ioctl()` and `TIOCGWINSZ`.
- Create `local_framebuffer` and set the `window_size`.?
- Create `new_state` frame buffer and set the size to `1*1`.
- initialize screen via `display.new_frame()`. Write screen to `STDOUT` via `swrite()`?
- Create the `Network::UserStream`, create the `Terminal::Complete local_terminal` with window size.?
- Open the network via `Network::Transport<Network::UserStream, Terminal::Complete>`.?
  - In the constructor function (client side), `connection(key_str, ip, port)` is called to create the socket with the server.
  - In `Connection::Connection()`, the `Addr remote_addr` is created and saved in `Connection.remote_addr`.
  - In `Connection::Connection()`, the `Connection::Socket` is created and saved in `Connection.socks`.
  - In `Connection::Connection()`, `Connection::set_MTU()` is called to set the MTU.
- Set minimal delay on outgoing keystrokes to 1 ms, via `network->set_send_delay()`.
- Tell server the size of the terminal via `network->get_current_state().push_back()`.
  - Here `network->get_current_state()` is actually `TransportSender.get_current_state()`.
  - The return value of `TransportSender.get_current_state()` is a `UserStream` object.
  - The `Parser::Resize` object is set with the current terminal window size.
- Set the `verbose` mode via `network->set_verbose()`.

### Transport<MyState, RemoteState>::Transport

- Initialize the `connection`, initialize the sender with `connection` and `initial_state`.
  - `connection` is the underlying, encrypted network connection.

### STMClient::main

In the main loop(while loop), It performs the following steps:

- Output terminal content to the `STDOUT_FILENO` via `output_new_frame()`.
- Get the network sockets from `network->fds()`.
- Add network sockets and `STDIN_FILENO` to the singleton `Select` object.
- Wait for socket input or user keystroke or signal via `sel.select()`, which specify `waittime`.
- Upon receive signals, the corresponding item in `Select.got_signal` array is set.
- Upon network sockets is ready to read, process it with `process_network_input()`.
- Upon user keystroke is ready to read, process it with `process_user_input`.
- Upon receive `SIGWINCH` signal, resize the terminal with `process_resize()`.
- Upon receive `SIGCONT` signal, process it with `resume()`.
- Upon receive `SIGTERM, SIGINT, SIGHUP, SIGPIPE` signals, showdown the process via `network->start_shutdown()`.
- Perform `network->tick()` to synchronizes the data to the server.

How the mosh client send the keystrokes to the server.

- `STMClient::main` calls `process_user_input()` if the main loop got the user keystrokes from `STDIN_FILENO`.
  - `process_user_input()` aka `STMClient::process_user_input()` calls `read()` system call to read the user keystrokes.
  - `process_user_input()` check the input character, for `LF`, `CR` etc. special character, they should be treated accordingly.
  - For each character, `process_user_input()` calls `network->get_current_state().push_back()` to save it in `UserStream` object.
  - `UserStream` object contains two kinds of character: `Parser::UserByte` and `Parser::Resize`.
  - Here `network->get_current_state()` is actually `TransportSender.get_current_state()`.
  - The return value of `TransportSender.get_current_state()` is a `UserStream` object.
- `STMClient::main` calls `network->tick()` in the main loop.
  - `network->tick()` calls `sender.tick()` to send data or an ack if necessary.
  - `sender.tick()` aka `TransportSender<MyState>::tick()`
  - `sender.tick()` calls `calculate_timers()` to calculate next send and ack times.
  - `sender.tick()` calls `current_state.diff_from()` to calculate diff.
  - Here `current_state.diff_from()` is actually `UserStream::diff_from()`, who calculate diff based on user keystrokes.
    - `UserStream::diff_from()` compares two `UserStream` object.
    - `UserStream::diff_from()` finds the different position and build `ClientBuffers::UserMessage`, which is a proto2 message.
    - `UserStream::diff_from()` returns the serialized string for the `ClientBuffers::UserMessage` object.
  - If `diff` is empty and if it's the ack time,
    - `sender.tick()` calls `send_empty_ack()` to send ack.
    - `send_empty_ack()` aka `TransportSender<MyState>::send_empty_ack()`.
    - `send_empty_ack()` calls `send_in_fragments()` to send data.
  - If `diff` is not empty and if it's the send or ack time,
    - `sender.tick()` calls `send_to_receiver()` to send diffs.
    - `send_to_receiver()` aka `TransportSender<MyState>::send_to_receiver()`.
  - `send_to_receiver()` calls `send_in_fragments()` to send data.

Send data process.

- `send_in_fragments()` aka `TransportSender<MyState>::send_in_fragments()`.
- `send_in_fragments()` creates `Instruction` with the `diff` created in previous step.
- `send_in_fragments()` splits the `Instruction` into `Fragment`.
  - Here `Fragmenter::make_fragments()` is called to serialize the `Instruction` into string,
  - compress it and splits it into `Fragment` based on the size of `MTU`,
  - The default size of `MTU` is 1280.
- `send_in_fragments()` calls `connection->send()` to send the `Fragment` to the receiver.
- `connection->send()` aka `Connection::send()` calls `sendto()` system call to send the real datagram to receiver.

How the mosh client receive the screen from the server.

- `STMClient::main` calls `process_network_input()` if network is ready to read.
  - `process_network_input()` aka `STMClient::process_network_input()` calls `network->recv()` to receive the data from server.
  - `network->recv()` aka `Transport<MyState, RemoteState>::recv()` calls `connection.recv()` to receive the data.
  - `connection.recv()` aka `Connection::recv()` calls `recv_one()` to read.
  - `recv_one()` aka `Connection::recv_one()` calls `recvmsg()` system call to receive data from socket.

## reference

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

## reference project

### st

[st - simple terminal project](https://st.suckless.org/), a c project, adding dependency packages with root privilege.

```sh
% ssh root@localhost
# apk add ncurses-terminfo-base fontconfig-dev freetype-dev libx11-dev libxext-dev libxft-dev
```

switch to the ide user of `nvide`.

```sh
% ssh ide@localhost
% git clone https://git.suckless.org/st
% cd st
% make clean
% bear -- make st
```

now you can check the source code of `st` via [nvide](https://github.com/ericwq/nvide).

### mosh

[mosh - mobile shell project](https://mosh.org/), a C++ project, adding dependency packages with root privilege.

```sh
% ssh root@localhost
# apk add ncurses-dev zlib-dev openssl1.1-compat-dev perl-dev perl-io-tty protobuf-dev automake autoconf libtool gzip
```

switch to the ide user of `nvide`. Download mosh from [mosh-1.3.2.tar.gz](https://mosh.org/mosh-1.3.2.tar.gz)

```sh
% ssh ide@localhost
% curl -O https://mosh.org/mosh-1.3.2.tar.gz
% tar xvzf mosh-1.3.2.tar.gz
% cd mosh-1.3.2
% ./configure
% bear -- make
```

now you can check the source code of `mosh` via [nvide](https://github.com/ericwq/nvide).
