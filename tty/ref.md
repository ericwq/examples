# Terminal emulator

## SSP design in [mosh research paper](https://mosh.org/mosh-paper.pdf)

SSP is organized into two layers. A datagram layer sends UDP packets over the network, and a transport layer is responsible for conveying the current object state to the remote host.

### Datagram Layer

The datagram layer maintains the “roaming” connection. It accepts opaque payloads from the transport layer, prepends an incrementing sequence number, encrypts the packet, and sends the resulting ciphertext in a UDP datagram. It is responsible for estimating the timing characteristics of the link and keeping track of the client’s current public IP address.

- Client roaming.
  - Every time the server receives an authentic datagram from the client with a sequence number greater than any before, it sets the packet’s source IP address and UDP port number as its new “target.” See [the client implementation](#how-does-the-client-roam) and [the server implementation](#how-does-the-server-support-client-roam).
- Estimating round-trip time and RTT variation. See [RTT and RTTVAR calculation](#how-to-receive-datagram-from-socket).
  - Every outgoing datagram contains a millisecond timestamp and an optional “timestamp reply,” containing the most recently received timestamp from the remote host. See [the implementation](#how-to-send-a-packet)
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
- The diff has three important fields: the source, target, and throwaway number. See [the `Instruction` implementation](#how-to-send-data-in-fragments).
  - The target of the diff is always the current sender-side state.
  - The throwaway number of the diff is always the most recent state that has been explicitly acknowledged by the receiver.
  - The source of the diff is allowed to be:
    - The most recent state that was explicitly acknowledged by the receiver
    - Any more-recent state that the sender thinks the receiver probably will have by the time the current diff arrives
  - The sender gets to make this choice based on what is most efficient and how likely it thinks the receiver actually has the more-recent state.
- Upon receiving a diff, the receiver throws away anything older than the throwaway number and attempts to apply the diff.
  - If it has the source state, and if the target state is newer than the receiver's current state, it succeeds and then acknowledges the new target state.
  - Otherwise it fails to apply the diff and just acks its current state number.

## Mosh coding analysis

### mosh-client.cc

In the `main` function, `STMClient` is the core to start `mosh` client.

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
  - `network->get_current_state().push_back()` adds `Parser::Resize` to `UserStream`.
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
- Wait for socket input or user keystroke or signal via `sel.select()`, within the `waittime` timeout.
- Upon receive signals, the corresponding item in `Select.got_signal` array is set.
- Upon network sockets is ready to read, process it with [`process_network_input()`](#how-to-process-the-network-input).
- Upon user keystroke is ready to read, process it with [`process_user_input()`](#how-to-process-the-user-input)
- Upon receive `SIGWINCH` signal, resize the terminal with `process_resize()`.
- Upon receive `SIGCONT` signal, process it with `resume()`.
- Upon receive `SIGTERM, SIGINT, SIGHUP, SIGPIPE` signals, showdown the process via `network->start_shutdown()`.
- Perform [`network->tick()`](#how-does-the-network-tick) to synchronizes the data to the server.

#### How to process the user input

- `STMClient::main` calls `process_user_input()` if the main loop got the user keystrokes from `STDIN_FILENO`.
- `process_user_input()` aka `STMClient::process_user_input()`.
- `process_user_input()` calls `read()` system call to read the user keystrokes.
- `process_user_input()` check the input character,
- If it get the `LF`, `CR` character, set `repaint_requested` to be true.
- For each character, `process_user_input()` calls `network->get_current_state().push_back()` to save it in `UserStream` object.
  - `network->get_current_state()` is actually `TransportSender.get_current_state()`.
  - `UserStream` object contains two kinds of character: `Parser::UserByte` and `Parser::Resize`.
  - Here the keystroke is wrapped in `Parser::UserByte`.
  - The return value of `TransportSender.get_current_state()` is a `UserStream` object.
  - `network->get_current_state().push_back()` adds `Parser::UserByte` to `UserStream`.
- The result of `process_user_input()` is that all the user keystrokes are saved in current state.

#### How does the network tick

- `STMClient::main` calls `network->tick()` in the main loop to procee the data in current state.
- `network->tick()` calls `sender.tick()` to send data or an ack if necessary.
- `sender.tick()` aka `TransportSender<MyState>::tick()`
- `sender.tick()` calls `calculate_timers()` to calculate next send and ack times.
  - `calculate_timers()` aka `TransportSender<MyState>::calculate_timers()`.
  - `calculate_timers()` calls [`update_assumed_receiver_state()`](#how-to-pick-the-reciver-state) to update assumed receiver state.
  - `calculate_timers()` calls [`rationalize_states()`](#how-to-rationalize-states) cut out common prefix of all states.
  - `calculate_timers()` calculate `next_send_time` and `next_ack_time`.
- `sender.tick()` calls `current_state.diff_from()` to [calculate diff](#how-to-calculate-the-diff-client-side).
- `sender.tick()` calls `attempt_prospective_resend_optimization()` to optimize diff.
- If `diff` is empty and if it's greater than the `next_ack_time`.
  - `sender.tick()` calls `send_empty_ack()` to send ack.
  - `send_empty_ack()` aka `TransportSender<MyState>::send_empty_ack()`.
  - `send_empty_ack()` calls [`send_in_fragments()`](#how-to-send-data-in-fragments) to send data.
- If `diff` is not empty and if it's greater than `next_send_time` or `next_ack_time`.
  - `sender.tick()` calls `send_to_receiver()` to send diffs.
  - `send_to_receiver()` aka `TransportSender<MyState>::send_to_receiver()`.
  - `send_to_receiver()` calls `add_sent_state()` to send a new state.
  - `add_sent_state()` adds the new state to `sent_states` and limits the size of `send_states` list.
  - Or `send_to_receiver()` refreshes the `timestamp` field of the latest state in `sent_states`.
  - Note `sent_states` is list of type `TimestampedState`, while `current_state` is of type `MyState`.
  - `send_to_receiver()` calls [`send_in_fragments()`](#how-to-send-data-in-fragments) to send data.
  - `send_to_receiver()` updates `assumed_receiver_state`, `next_ack_time` and `next_send_time`.

#### How to calculate the diff (client side)

- `current_state.diff_from()` aka `UserStream::diff_from()`, who calculate diff based on user keystrokes.
- `diff_from()` compares `current_state` with `assumed_receiver_state` to calculate the diff.
- For client side:
  - `diff_from()` compares two `UserStream` object.
  - `diff_from()` finds the different position and build `ClientBuffers::UserMessage`, which is a proto2 message.
  - `diff_from()` returns the serialized string for the `ClientBuffers::UserMessage` object.
  - `UserMessage` contains several `ClientBuffers.Instruction`.
  - `ClientBuffers.Instruction` is composed of `Keystroke` or `ResizeMessage` (see userinput.proto file)
  - Several `Keystroke` can be appended to one `ClientBuffers.Instruction`.
  - `ResizeMessage` is added to one `ClientBuffers.Instruction`.

#### How to pick the reciver state

- `update_assumed_receiver_state()` chooses a most recent receiver state based on network traffic.
- `update_assumed_receiver_state()` picks the first item in `send_state`.
- `send_state` is of type `list<TimestampedState<MyState>>`.
- `send_state` skips the first item.
- For each item in `send_state`, if the time gap is lower than `connection->timeout()`. Update `assumed_receiver_state`.
  - `connection->timeout()` aka `Connection::timeout()`.
  - `connection->timeout()` calcuates [RTO](https://datatracker.ietf.org/doc/html/rfc2988) based on `SRTT` and `RTTVAR`.
- The result is saved in `assumed_receiver_state`.
- `assumed_receiver_state` point to the middle of `sent_states`.

#### How to rationalize states

- `rationalize_states()` aka `TransportSender<MyState>::rationalize_states()`.
- `rationalize_states()` picks the first state from `sent_states` as common prefix.
  - `sent_states` is of type `list<TimestampedState<MyStat>>`.
- The comm prefix is the first state in `send_state`.
- `rationalize_states()` calls `current_state.subtract()` to cut out common prefix from `current_state`.
- `rationalize_states()` calls `i->state.subtract()` to cut out common prefix for all states in `sent_states`.
  - For client side:
  - `subtract()` aka `UserStream::subtract()`.
  - `subtract()` cuts out any `UserEvent` from 's `actions` deque, if it's the same `UserEvent` in `prefix`.
  - The result is the caller of `subtract()` cut out common prefix.
- The result is that the common prefix in `current_state` and `sent_states` is cut out.

#### How to send data in fragments

- `send_in_fragments()` aka `TransportSender<MyState>::send_in_fragments()`.
- `send_in_fragments()` creates `TransportBuffers.Instruction` with the `diff` created in [previous](#how-to-calculate-the-diff-client-side) step.
- `TransportBuffers.Instruction` contains the following fields.
  - `old_num` field is the source number. It's value is `assumed_receiver_state->num`.
  - `new_num` field is the target number. It's value is specified by `new_num` parameter.
  - `throwaway_num` field is the throwaway number. It's value is `sent_states.front().num`.
  - `diff` field contains the `diff`. It's value is specified by `diff` parameter.
  - `ack_num` field is the ack number. It's value is assigned by `ack_num`.
- `send_in_fragments()` calls `Fragmenter::make_fragments` to splits the `TransportBuffers.Instruction` into `Fragment`.
  - `make_fragments()` serializes `TransportBuffers.Instruction` into string and compresses it to string `payload`.
  - `make_fragments()` splits the `payload` string into fragments based on the size of `MTU`,
  - The default size of `MTU` is 1280.
  - Fragment has a `id` field, which is the instruction id. It's the same id for all the fragment.
  - Fragment has a `fragment_num` field, which starts from zero, and is increased one for each new fragment.
  - Fragment has a `final` field, which is used to indicate the last fragment.
  - Fragment has a `contents` field, which contains part of the instruction.
  - The fragments is saved in `Fragment` vector.
- `send_in_fragments()` calls [`connection->send()`](#how-to-send-a-packet) to send each `Fragment` to the server.

### How to send a packet?

- `connection->send()` aka `Connection::send()`.
- `connection->send()` calls `new_packet()` to create a `Packet`.
  - `timestamp_reply` means?
  - `Packet` is of type `Network::Packet`.
  - Besides the `payload` field,
  - A `Packet` also contains a unique `seq` field, a `timestamp` field and a `timestamp_reply` field.
- `connection->send()` calls `session.encrypt()` to encrypt the `Packet`.
- `connection->send()` calls `sendto()` system call to send the encrypted data to receiver.
  - `sendto()` use the last socket from socket list to send the encrypted data.
- `connection->send()` checks the time gap between now and `last_port_choice`, `last_roundtrip_success`.
- `connection->send()` calls [`hop_port()`](#how-does-the-client-roam), if the time gap is greater than `PORT_HOP_INTERVAL`.

#### How does the client roam.

- `hop_port()` aka `Connection::hop_port()`. `hop_port()` only works for client.
- `hop_port()` calls `setup()` to update `last_port_choice`.
- `hop_port()` creates a new `Socket` object and calls `socks.push_back()` to save it in `socks` list.
- `hop_port()` calls [`prune_sockets()`](#how-to-prune-the-sockets) to prune the old sockets.
- `last_port_choice` is changed, when a new `Socket` is created.
- `last_roundtrip_success` is changed, when a new datagram is received.
- `PORT_HOP_INTERVAL` is 10s. Which means every 10 seconds a new socket is added to the socket list.

#### How to process the network input

<!--
TODO What's the behavior of the serverside.
TODO what the purpose of `overlay`.
TODO what the meaning of `display`.
TODO how to receive network input.
-->

- `STMClient::main` calls `process_network_input()` if network is ready to read.
- `process_network_input()` aka `STMClient::process_network_input()`
- `process_network_input()` calls [`network->recv()`](#how-to-receive-network-input) to receive network input.

#### How to receive network input

- `network->recv()` aka `Transport<MyState, RemoteState>::recv()`
- `network->recv()` calls [`connection.recv()`](#how-to-read-data-from-socket) to receive the data.
- `network->recv()` calls `fragments.add_fragment()` get the complete packet.
- `network->recv()` calls `fragments.get_assembly()` to build the `Instruction`.
- `network->recv()` calls `sender.process_acknowledgment_through()` to update `last_roundtrip_success`.
- The above implementation means that last send timestamp is saved as `last_roundtrip_success`.
- `network->recv()` makes sure we don't already have the new state?
- `network->recv()` makes sure we do have the old state.
- `network->recv()` throws away the unnecessary state via `process_throwaway_until()`.
- `network->recv()` limit on state queue.
- `network->recv()` apply diff to reference state.
- `network->recv()` Insert new state in sorted place.
- `network->recv()` calls `received_states.push_back()` to store the received state.
- `network->recv()` calls `sender.set_ack_num()` to set `ack_num`.
- `network->recv()` calls `sender.remote_heard()` to set last time received new state.
- `network->recv()` calls `sender.set_data_ack()` to accelerate reply ack.

#### How to read data from socket

- `connection.recv()` aka `Connection::recv()`
- `connection.recv()` calls [`recv_one()`](#how-to-receive-datagram-from-socket) on the first `Socket` in `socks`.
- If [`recv_one()`](#how-to-receive-datagram-from-socket) returns `EAGAIN` or `EWOULDBLOCK`, try the next `Socket` in `socks` until the last one.
- `connection.recv()` calls [`prune_sockets()`](#how-to-prune-the-sockets) to prune the old sockets.
- `connection.recv()` returns the `payload` got from `recv_one()`.

#### How to prune the sockets.

- `prune_sockets()` aka `Connection::prune_sockets()`
- `prune_sockets()` removes old sockets if the new socket has been working for long enough.
- `prune_sockets()` makes sure we don't have too many receive sockets open.

#### How to receive datagram from socket

- `recv_one()` aka `Connection::recv_one()`
- `recv_one()` calls `recvmsg()` system call to receive data from socket.
- `recv_one()` calls `session.decrypt()` to decrypt the received message.
- `recv_one()` creates a `Packet` object based on the decrypted data.
- `recv_one()` checks `Packet`'s sequence number to make sure it is greater than the `expected_receiver_seq`.
  - if packet sequence number is greater than `expected_receiver_seq`,
    - `recv_one()` increases `expected_receiver_seq`.
    - `recv_one()` saves the `p.timestamp` in `saved_timestamp`, saves the time in `saved_timestamp_received_at`.
    - `recv_one()` signals counterparty to slow down via decrease `saved_timestamp`, if congestion is detected.
    - `recv_one()` calculates `SRTT` and `RTTVAR` based on each [RTT](https://datatracker.ietf.org/doc/html/rfc29880).
    - `recv_one()` updates `last_heard` with current time.
  - For server side, [client roaming](#how-does-the-server-support-client-roam) is supported here.
  - if packet sequence number is less than `expected_receiver_seq`
    - `recv_one()` return out-of-order or duplicated packets to caller, .
- `recv_one()` return the `payload` to caller.

#### How does the server support client roam

- `recv_one()` compares `packet_remote_addr` with `remote_addr`.
- If the packet remote address is different than remote address, update the `remote_addr` and `remote_addr_len`.
- `recv_one()` calls `getnameinfo()` to validate the new remote address.

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
