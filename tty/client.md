# Mosh client coding analysis

![mosh-comm.svg](img/mosh-comm.svg)

## mosh-client.cc

In the `main` function, `STMClient` is the core to start `mosh` client.

```cpp
  try {
    STMClient client( ip, desired_port, key, predict_mode, verbose );
    client.init();

    try {
      success = client.main();
    } catch ( Exception e) {
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

- [STMClient constructor](#stmclient-constructor)
- [STMClient::init](#stmclientinit)
- [STMClient::main](#stmclientmain)
- [How to send keystroke to remote server](#how-to-send-keystroke-to-remote-server)

<!-- TODO what the purpose of `overlay`. -->

### STMClient constructor

- `STMClient()` is called withe following parameters:`ip`,`port`,`key`,`predict_mode`,`verbose` parameter.
- `STMClient()` saves `key`, `ip`, `port` parameters as field member.
- `STMClient()` initializes `escape_key`, `escape_pass_key`, `escape_pass_key2`,
  - with `0x1E`, `^`, `^` corresponding value.
- `STMClient()` initializes `escape_requires_lf` with false value.
- `STMClient()` initializes empty `termios` structs: `raw_termios`, `saved_termios`.
- `STMClient()` initializes `new_state` frame buffer with 1\*1 size.
- `STMClient()` initializes `local_framebuffer` frame buffer with 1\*1 size.
- `STMClient()` initializes `overlay`, which is type of [`Overlay::OverlayManager`](#overlayoverlaymanager).
- `STMClient()` initializes a NULL `network`.
- `STMClient()` initializes `display`, which is type of [`Terminal::Display`](#terminaldisplay).
- `STMClient()` initializes `repaint_requested`, `lf_entered`, `quit_sequence_started`, `clean_shutdown` with false value.

#### Overlay::OverlayManager

- `OverlayManager` has a `NotificationEngine`, which performs the notification work.
- `OverlayManager` has a `PredictionEngine`, which performs the prediction work.
- `OverlayManager` has a `TitleEngine`, which performs the title work.
- The default constructor initializes a `OverlayManager` without any parameters.

#### Terminal::Display

- `Display()` is initialized with true `use_environment`.
- `Display()` calls [`setupterm()`](https://linux.die.net/man/3/setupterm) to reads in the `terminfo` database, initialize the `terminfo` structures.
- `Display()` calls [`tigetstr()`](https://linux.die.net/man/3/setupterm) to check (erase character) [ech](https://pubs.opengroup.org/onlinepubs/7908799/xcurses/terminfo.html) support.
- `Display()` calls [`tigetflag()`](https://linux.die.net/man/3/setupterm) to check (back color erase) [bce](https://pubs.opengroup.org/onlinepubs/7908799/xcurses/terminfo.html) support.
- `Display()` get the `TERM` environment variable and compare it with the following value.
  - `xterm`, `rxvt`, `kterm`, `Eterm`, `screen`
- If `TERM` environment variable contains the above string, `has_title` is true.
- `Display()` get the `MOSH_NO_TERM_INI` environment variable.
- If `MOSH_NO_TERM_INI` environment variable is set,
  - `Display()` calls [`tigetstr()`](https://linux.die.net/man/3/setupterm) to get the (enter ca mode) [smcup](https://pubs.opengroup.org/onlinepubs/7908799/xcurses/terminfo.html) string.
  - `Display()` calls [`tigetstr()`](https://linux.die.net/man/3/setupterm) to get the (exit ca mode) [rmcup](https://pubs.opengroup.org/onlinepubs/7908799/xcurses/terminfo.html) string.

### STMClient::init

- If the client terminal doesn't [support UTF8](#how-to-check-the-utf8-support), exit the application.
- Get the `termios` struct for `STDIN_FILENO`, via `tcgetattr()`.
- Set `IUTF8` flag for `termios`.
- Set the terminal to raw mode, via `cfmakeraw()`.
- Set the `termios` struct for `STDIN_FILENO`, via `tcsetattr()`.
- Put terminal in [application-cursor-key mode](#the-application-cursor-key-mode), via `swrite()` write `display.open()` to `STDOUT_FILENO`.
- Set terminal window title, via `overlays.set_title_prefix()`.
- Set variable, `escape_key`, `escape_pass_key`, `escape_pass_key2`.
- Set variable, `escape_key_help`.
- Set `escape_key_string` string, via `overlays.get_notification_engine().set_escape_key_string()`.
  - `overlays.get_notification_engine()` returns `NotificationEngine`.
  - `set_escape_key_string()` aka `NotificationEngine::set_escape_key_string()`.
  - `set_escape_key_string()` sets the `escape_key_string` field in `NotificationEngine` object.
- Set variable `connecting_notification`.

#### How to check the UTF8 support

- `is_utf8_locale()` checks `locale_charset()` to compare the locale with UTF-8.
  - `locale_charset()` calls `nl_langinfo()` to return a string with the name of the character encoding.
- `is_utf8_locale()` return true if the terminal support UTF8, otherwise return false.

#### The application-cursor-key mode

- `display.open()` aka `Display::open()`.
- `display.open()` returns a control sequence to set the application-cursor-key mode.
- Application Cursor Keys mode is a way for the server to change the control sequences sent by the arrow keys.
- In normal mode, the arrow keys send `ESC [A` through to `ESC [D`.
- In application mode, they send `ESC OA` through to `ESC OD`.

### STMClient::main

`STMClient::main()` calls [`main_init()`](#stmclientmain_init) to initialize signal handling and structures. In the main loop(while loop), It performs the following steps:

- Output terminal content to the `STDOUT_FILENO` via [`output_new_frame()`](#how-to-output-content).
- Get the network sockets from `network->fds()`.
- Add network sockets and `STDIN_FILENO` to the singleton `Select` object.
- Wait for socket input or user keystroke or signal via `sel.select()`, within the `waittime` timeout.
- Upon receive signals, the corresponding item in `Select.got_signal` array is set.
- Upon network sockets is ready to read, process it with [`process_network_input()`](#how-to-process-the-network-input).
- Upon user keystroke is ready to read, process it with [`process_user_input()`](#how-to-process-the-user-input)
- Upon receive `SIGWINCH` signal, resize the terminal with [`process_resize()`](#how-to-process-resize).
- Upon receive `SIGCONT` signal, process it with `resume()`.
- Upon receive `SIGTERM, SIGINT, SIGHUP, SIGPIPE` signals, showdown the process via `network->start_shutdown()`.
- Perform [`network->tick()`](#how-does-the-network-tick) to synchronizes the data to the server.

#### STMClient::main_init

In `client.main()`, `main_init()` is called to init the `mosh` client.

- `main_init()` [registers signal handler](#how-to-register-signal-handler) for `SIGWINCH`, `SIGTERM`, `SIGINT`, `SIGHUP`, `SIGPIPE`, `SIGCONT`.
- `main_init()` gets the window size for `STDIN_FILENO` , via `ioctl()` and `TIOCGWINSZ` flag.
- `main_init()` [initializes `local_framebuffer`](#how-to-initialize-frame-buffer) frame buffer with the above window size.
- `main_init()` [initializes `new_state`](#how-to-initialize-frame-buffer) frame buffer with `1*1` size.
- `main_init()` calls [`display.new_frame()`](#how-to-calculate-frame-buffer-difference) to get the initial screen.
- `main_init()` calls `swrite()` to write initial screen to `STDOUT_FILENO`.
- `main_init()` creates blank [`Network::UserStream`](#networkuserstream).
- `main_init()` creates `local_terminal` of type [`Terminal::Complete`](#terminalcomplete).
- `main_init()` creates `network` of type [`Network::Transport<Network::UserStream, Terminal::Complete>`](#networktransportnetworkuserstream-terminalcomplete).
- `main_init()` calls `network->set_send_delay()` to set minimal delay on outgoing keystrokes to 1 ms.
  - `set_send_delay()` calls `sender.set_send_delay()` to set the minimal delay.
- `main_init()` [tells server the terminal size](#how-to-tell-server-the-terminal-size).
- `main_init()` sets the `verbose` mode via `network->set_verbose()`.
- `main_init()` sets the `verbose` mode via `Select::set_verbose()`.

#### How to output content

- `output_new_frame()` aka `STMClient::output_new_frame()`.
- `output_new_frame()` gets `new_state` (the `Framebuffer`) from the state saved in `received_states`.
- `new_state` is of type `Terminal::Framebuffer`.
- `output_new_frame()` calls `overlays.apply()` to apply `new_state` to local overlays.
- `output_new_frame()` calls [`display.new_frame()`](#how-to-calculate-frame-buffer-difference) to calculate minimal `diff` from where we are.
- `output_new_frame()` writes the `diff` to `STDOUT_FILENO`.
- `output_new_frame()` sets `repaint_requested` to true.
- `output_new_frame()` sets `local_framebuffer` to the new state.

#### How to tell server the terminal size

- `main_init()` creates a `Parser::Resize` object and pushes it into `network->get_current_state()`.
- Here `network->get_current_state()` is actually `TransportSender.get_current_state()`.
- The return value of `TransportSender.get_current_state()` is a `UserStream` object.
- `network->get_current_state().push_back()` adds `Parser::Resize` to `UserStream`.
- The `Parser::Resize` object is initialized with the current terminal window size.
- The current state will ben send to server later.

#### Network::Transport<Network::UserStream, Terminal::Complete>

- `Network::Transport` is constructed with:
  - the blank `UserStream`,
  - local terminal which is type of `Terminal::Complete`,
  - `key`, `ip`, `port` as parameters.
- `Network::Transport` has a `Connection`, which represents the underlying, encrypted network connection.
- `Network::Transport` has a `TransportSender<Network::UserStream>`, which represents the sender.
- `Network::Transport` has a `list<TimestampedState<Terminal::Complete>>`, which represents receiver.
- `Network::Transport` calls `connection(key_str, ip, port)` to [initialize the connection](#how-to-initialize-connection).
- `Network::Transport` calls `sender(connection, initial_state)` to [initialize sender](#how-to-initialize-sender).
- In the constructor of `Network::Transport`,
  - `received_states` is a list type of `TimestampedState<Terminal::Complete>`.
  - `received_states` is initialized with the `local_terminal` as parameter.
  - `received_states` adds the `local_terminal` to its list.
- `Network::Transport()` set `receiver_quench_timer` to zero.
- `Network::Transport()` set `last_receiver_state` to be `local_terminal`.
- `Network::Transport()` creates `fragments`, which is type of `FragmentAssembly`.

#### How to initialize sender

- `sender(connection, initial_state)` is called to initialize the sender.
- `sender()` is the constructor of `TransportSender<Network::UserStream>`.
- `sender()` initializes `current_state` with the `initial_state` as parameter.
- `sender()` initializes `connection` pointer with the `connection` as parameter.
- `sender()` initializes `sent_states` list with the `initial_state` as the first state.

#### How to initialize connection

- `connection(key_str, ip, port)` is called to create the connection with server.
- `connection()` is the constructor of `Network::Connection`
- `connection()` calls `setup()` to set the `last_port_choice` to current time.
- `connection()` initializes a empty deque of `Socket`: `socks`.
- `connection()` initializes `remote_addr` with `ip`, `port` as parameters,
  - `remote_addr` represents server address.
- `connection()` initializes `session` with `key` as parameter,
  - `session` object is used to encrypt/decrypt message.
- `connection()` creates a `Socket` and pushes it into `socks` deque.
- `connection()` calls `set_MTU()` to set the MTU.
- `connection()` sets `has_remote_addr` to true.

#### Terminal::Complete

- `Terminal::Complete` represents the complete terminal, a `UTF8Parser` feeding `Actions` to an `Emulator`.
- `Complete()` creates a `Parser::UTF8Parser` object.
- `Complete()` creates a `Terminal::Emulator` object with `width` and `hight` as parameters.
- `Complete()` initializes a [`Terminal::Display`](#terminaldisplay) object with `false` as parameter.
- `Complete()` creates a `Parser::Actions` object.

#### Network::UserStream

- `Network::UserStream` has a deque of type `UserEvent`.
- `UserEvent` can store `Parser::UserByte` or `Parser::Resize` object.
- `Parser::UserByte` is used to store user keystroke.
- `Parser::Resize` is used to store resize event.
- The default constructor of `Network::UserStream` builds a empty `Network::UserStream` object.

#### How to calculate frame buffer difference

- `new_frame()` aka `Display::new_frame()`
- `new_frame()` receives two `Framebuffer`: `last` and `f`, and builds the difference for output to terminal display.
- `new_frame()` initializes a `FrameState`: `frame`, with the old `Framebuffer` `last`.
- `new_frame()` checks if the bell ring happened: if true, append escape sequence to `frame`.
- `new_frame()` checks if icon name or window title changed: if true, append escape sequence to `frame`.
- `new_frame()` checks if reverse video state changed: if true, append escape sequence `frame`.
- `new_frame()` checks if window size changed: if true, append escape sequence to `frame`.
- `new_frame()` checks is cursor visibility initialized: if false, append escape sequence to `frame`.
- `new_frame()` extends rows if we've gotten a resize and new is wider than old.
- `new_frame()` adds rows if we've gotten a resize and new is taller than old.
- `new_frame()` checks if display moved up by a certain number of lines
- `new_frame()` updates the display, row by row, via calling `put_row()` for each row. TODO detail of `put_row()`.
- `new_frame()` checks if cursor location changed, append escape sequence to `frame`.
- `new_frame()` checks if cursor visibility changed, append escape sequence to `frame`.
- `new_frame()` checks if renditions changed: if true, append escape sequence to `frame`.
- `new_frame()` checks if bracketed paste mode changed: if true, append escape sequence to `frame`.
- `new_frame()` checks if mouse reporting mode changed: if true, append escape sequence to `frame`.
- `new_frame()` checks if mouse focus mode changed: if true, append escape sequence to `frame`.
- `new_frame()` checks if mouse encoding mode changed: if true, append escape sequence to `frame`.
- `new_frame()` returns the final `frame` escape sequence string.

#### How to initialize frame buffer

- `new_state` and `local_framebuffer` is type of `Terminal::Framebuffer`.
- `Framebuffer` has a vector of `Row`, the rows number is determined by terminal hight.
- Each `Row` in `Framebuffer` has a vector of `Cell`, the `Cell` number is determined by terminal width.
- The `Cell` has the content string and content attributes: `Renditions`.
- `Renditions` determines the foreground color, background color, bold, faint, italic, underlined, etc.

#### How to register signal handler

- Disposition the previous signals via `sel.add_signal()` .
- `sel.add_signal()` aka `Select::add_signal()`.
- `add_signal()` calls `sigprocmask()` system call to add the specified signal mask.
- `add_signal()` calls `sigaction()` to register signal handler.
  - Here all Signals is blocked during handler invocation
- It blocks the signal outside of `pselect()`.
- In `pselect()`, the signal mask is replaced by a empty mask set.

#### How to process the user input

- `STMClient::main` calls `process_user_input()` if the main loop got the user keystrokes from `STDIN_FILENO`.
- `process_user_input()` aka `STMClient::process_user_input()`.
- `process_user_input()` calls `read()` system call to read the user keystrokes.
- `process_user_input()` calls `overlays.get_prediction_engine().set_local_frame_sent()` to save the last `send_states` number.
- `process_user_input()` check each input character:
- `process_user_input()` calls `overlays.get_prediction_engine().new_user_byte()` to TODO.
- If it get the `LF`, `CR` character, set `repaint_requested` to be true.
- For each character, `process_user_input()` calls `network->get_current_state().push_back()` to save it in `UserStream` object.
  - `network->get_current_state()` is actually `TransportSender.get_current_state()`.
  - `UserStream` object contains two kinds of character: `Parser::UserByte` and `Parser::Resize`.
  - Here the keystroke is wrapped in `Parser::UserByte`.
  - The return value of `TransportSender.get_current_state()` is a `UserStream` object.
  - `network->get_current_state().push_back()` adds `Parser::UserByte` to `UserStream`.
- The result of `process_user_input()` is that all the user keystrokes are saved in current state.

#### How to process resize

- `process_resize()` gets the window size for `STDIN_FILENO` , via `ioctl()` and `TIOCGWINSZ` flag.
- `process_resize()` creates `Parser::Resize` with the window size.
- `process_resize()` pushes the above `Parser::Resize` into `network->get_current_state()`.
- `process_resize()` calls `overlays.get_prediction_engine().reset()` to tell prediction engine.
- `process_resize()` returns true.

#### How does the network tick

- `STMClient::main` calls `network->tick()` in the main loop to procee the data in current state.
- `network->tick()` calls `sender.tick()` to send data or an ack if necessary.
- `sender.tick()` aka `TransportSender<MyState>::tick()`
- `sender.tick()` calls `calculate_timers()` to calculate next send and ack times.
  - `calculate_timers()` aka `TransportSender<MyState>::calculate_timers()`.
  - `calculate_timers()` calls [`update_assumed_receiver_state()`](#how-to-pick-the-reciver-state) to update assumed receiver state.
  - `calculate_timers()` calls [`rationalize_states()`](#how-to-rationalize-states) cut out common prefix of all states.
  - `calculate_timers()` calculate `next_send_time` and `next_ack_time`.
- `sender.tick()` calls `diff_from()` to compare `current_state` with `assumed_receiver_state` to [calculate diff](#how-to-calculate-the-diff-for-userstream).
- `sender.tick()` calls `attempt_prospective_resend_optimization()` to optimize diff.
- If `diff` is empty and if it's greater than the `next_ack_time`.
  - `sender.tick()` calls [`send_empty_ack()`](#how-to-send-empty-ack) to send ack.
- If `diff` is not empty and if it's greater than `next_send_time` or `next_ack_time`.
  - `sender.tick()` calls [`send_to_receiver()`](#how-to-send-to-receiver) to send diffs.

#### How to send empty ack

- `send_empty_ack()` aka `TransportSender<MyState>::send_empty_ack()`.
- `send_empty_ack()` gets the last state number via `sent_states.back()` and increases one.
- `send_empty_ack()` calls `add_sent_state()` to push `current_state` into `sent_states`,
  - with the new state number and current time as parameters.
  - limit the size of `sent_states` below 32.
- `send_empty_ack()` calls [`send_in_fragments()`](#how-to-send-data-in-fragments) to send the new state
  - with empty string as `diff` parameter.

#### How to send to receiver

- `send_to_receiver()` aka `TransportSender<MyState>::send_to_receiver()`.
- If `current_state` number is equal to `sent_states.back()` number,
- `send_to_receiver()` refreshes the `timestamp` field of the latest state in `sent_states`.
- If `current_state` number is not equal to `sent_states.back()` number, increase the state number.
- `send_to_receiver()` calls `add_sent_state()` to push `current_state` into `sent_states`,
  - with the new state number and current time as parameters.
  - limit the size of `sent_states` below 32.
- Note `sent_states` is type of list `TimestampedState`, while `current_state` is type of `MyState`.
- `send_to_receiver()` calls [`send_in_fragments()`](#how-to-send-data-in-fragments) to send data.
- `send_to_receiver()` updates `assumed_receiver_state`, `next_ack_time` and `next_send_time`.

#### How to calculate the diff for UserStream

- `current_state.diff_from()` aka `UserStream::diff_from()`, who calculate diff based on user keystrokes.
- `diff_from()` compares current `UserStream` with `existing` `UserStream` to calculate the diff.
- `diff_from()` finds the position in the current `UserStream` which is different from `existing` `UserStream`.
- `diff_from()` iterates to the end of current `UserStream` starting from the above position.
- `diff_from()` build `ClientBuffers::UserMessage`, with the `UserEvent` object in each iteration,
- `diff_from()` returns the serialized string of the `ClientBuffers::UserMessage` object.
- `ClientBuffers::UserMessage` is a proto2 message. See userinput.proto file.
- `ClientBuffers::UserMessage` contains several `ClientBuffers.Instruction`.
- `ClientBuffers.Instruction` is composed of `Keystroke` or `ResizeMessage`.
- Several `Keystroke` can be appended to one `ClientBuffers.Instruction`.
- `ResizeMessage` is added to one `ClientBuffers.Instruction`.

#### How to pick the reciver state

- `update_assumed_receiver_state()` chooses a most recent receiver state based on network traffic.
- `update_assumed_receiver_state()` picks the first item in `send_state`.
- `send_state` is type of `list<TimestampedState<MyState>>`.
- `send_state` skips the first item.
- For each item in `send_state`, if the time gap is lower than `connection->timeout()`. Update `assumed_receiver_state`.
  - `connection->timeout()` aka `Connection::timeout()`.
  - `connection->timeout()` calcuates [RTO](https://datatracker.ietf.org/doc/html/rfc2988) based on `SRTT` and `RTTVAR`.
- The result is saved in `assumed_receiver_state`.
- `assumed_receiver_state` point to the middle of `sent_states`.

#### How to rationalize states

- `rationalize_states()` aka `TransportSender<MyState>::rationalize_states()`.
- `rationalize_states()` picks the first state from `sent_states` as common prefix.
  - `sent_states` is type of `list<TimestampedState<MyStat>>`.
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
- `send_in_fragments()` creates `TransportBuffers.Instruction` with the `diff` created in [previous](#how-to-calculate-the-diff-for-userstream) step.
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

#### How to send a packet?

- `connection->send()` aka `Connection::send()`.
- `connection->send()` calls `new_packet()` to create a `Packet`.
  - `timestamp_reply` means?
  - `Packet` is type of `Network::Packet`.
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

- `STMClient::main` calls `process_network_input()` if network is ready to read.
- `process_network_input()` aka `STMClient::process_network_input()`
- `process_network_input()` calls [`network->recv()`](#how-to-receive-network-input) to receive network input.

#### How to receive network input

- `network->recv()` aka `Transport<MyState, RemoteState>::recv()`
- `network->recv()` calls `connection.recv()` to [receive payload](#how-to-read-data-from-socket) string.
- `network->recv()` calls `Fragment(const string& x)` to [build a `Fragment` object](#how-to-create-the-frament-from-string) from the payload string.
- `network->recv()` calls `fragments.add_fragment()` to [get the complete packet](#how-to-get-the-complete-packet).
- `network->recv()` calls `fragments.get_assembly()` to [build the `Instruction` object](#how-to-build-instruction-from-fragments).
- `network->recv()` calls `sender.process_acknowledgment_through()` to remove states from `send_states`.
  - It removes any `sent_states` whose `num` field is less than `ack_num`.
- `network->recv()` calls `connection.set_last_roundtrip_success()` to update `last_roundtrip_success`.
  - It means that last send timestamp is saved as `last_roundtrip_success`.
- `network->recv()` checks the `Instruction.new_num` does not exist in `received_states`.
  - It makes sure we don't already have the new state.
- `network->recv()` checks the `Instruction.old_num` does exist in `received_states`.
  - It makes sure we do have the old state.
- `network->recv()` throws away the unnecessary state via `process_throwaway_until()`.
  - Any state whose `num` field less than `throwaway_num` is thrown away.
- `network->recv()` limits the `received_states` queue size via drop the received state:
  - If `received_states.size() < 1024` and current time is less than `receiver_quench_timer`.
  - The value of `receiver_quench_timer` is `now` plus 15000ms.
- If `Instruction` diff field is not empty, `network->recv()` calls [`new_state.state.apply_string()`](server.md#apply_string).
  - `apply_string()` is called with `Instruction` diff field as parameter.
  - `apply_string()` initializes `RemoteState` with `ClientBuffers::UserMessage`.
- `network->recv()` initializes a `RemoteState` and wraps it in `TimestampedState<RemoteState>`.
- If out-of-order state is received, `network->recv()` inserts new state and returns directly,
- `network->recv()` calls `received_states.push_back()` to store the new state.
- `network->recv()` calls `sender.set_ack_num()` to set `ack_num`.
  - It means the `sender` set `ack_num` got from the `received_states` number.
- `network->recv()` calls `sender.remote_heard()` to set `last_heard`: last time received new state.
- `network->recv()` calls `sender.set_data_ack()` to set `pending_data_ack`: accelerate reply ack.

#### How to build instruction from fragments

- `fragments.get_assembly()` aka `FragmentAssembly::get_assembly()`
- `get_assembly()` concatenates the `contents` field of each `Fragment` into one piece.
- `get_assembly()` calls `get_compressor().uncompress_str()` to decompress the string.
- `get_assembly()` calls `ret.ParseFromString()` to build the `Instruction` object.
- `get_assembly()` clears the fragments, reset `fragments_arrived` and `fragments_total`.
- `get_assembly()` returns the `Instruction` object.

#### How to get the complete packet

- `fragments.add_fragment()` adds a frament into the `fragments`
- `fragments` is type of `FragmentAssembly`.
- `fragments.add_fragment()` checks fragment id and fragment final flag, adds new fragment to vector.
- `fragments.add_fragment()` returns true if the final fragment is received.
- Otherwise, `fragments.add_fragment()` returns false.

#### How to create the frament from string

- From the `Fragment::tostring()`, the format of the network fragment data is:
  - The `id` field, which is `uint64_t`, contains the fragment id.
  - The `fragment_num` field, which is `uint16_t`, contains the fragment number and fragment final flag.
  - The `contents` field, which is `string`, contains the fragment payload.
- `Fragment(const string& x)` constructs the `Fragment` using the above format.
- `Fragment(const string& x)` returns one `Fragment`.

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

### How to send keystroke to remote server

#### user keystroke -> `Parser::UserByte` -> `Network::UserEvent` -> `Network::UserStream`

- Upon receiving user keystroke:
  - [`STMClient::process_user_input()`](#how-to-process-the-user-input) reads user keystroke from `STDIN_FILENO`.
  - [`STMClient::process_user_input()`](#how-to-process-the-user-input) wraps user keystroke with `Parser::UserByte`,
  - `Parser::UserByte` contains `c` field.
  - `Parser::UserByte` is wrapped in `Network::UserEvent` and pushed into `Network::UserStream` object.
- Upon receiving signal `SIGWINCH`,
  - [`STMClient::process_resize()`](#how-to-process-resize) gets the window size.
  - [`STMClient::process_resize()`](#how-to-process-resize) creates `Parser::Resize` object with the above window size.
  - `Parser::Resize` contains `width` and `height` fields.
  - `Parser::Resize` is wrapped in `Network::UserEvent` and pushed into `Network::UserStream` object.
- `Network::UserStream` contains a deque of type `Network::UserEvent`.
- `Network::UserEvent` contains the following fields:
  - `type`,
  - `userbyte`,
  - `resize`.

#### `Network::UserStream` -> `TransportBuffers.Instruction` -> `Network::Fragment`

When it's time to send the `Network::UserStream` to remote server:

- [`sender.tick()`](#how-does-the-network-tick) calculates the difference between two `Network::UserStream` objects.
- The difference is transformed into string representation of `ClientBuffers::UserMessage`.
- [`send_in_fragments()`](#how-to-send-data-in-fragments) constructs the `TransportBuffers.Instruction` object.
- The string representation of `ClientBuffers::UserMessage` is assigned to the `diff` field of `TransportBuffers.Instruction`.
- `TransportBuffers.Instruction` is the "state" in [transport layter](ref.md#transport-layer).
- `TransportBuffers.Instruction` contains the following fields:
  - `old_num`,
  - `new_num`,
  - `ack_num`,
  - `throwaway_num`,
  - `diff`.
- [`send_in_fragments()`](#how-to-send-data-in-fragments) splits `TransportBuffers.Instruction` into one or several `Network::Fragment` based on `MTU` size.
- `Network::Fragment` is a utility class because of `MTU`.
- `Network::Fragment` contains the following fields:
  - `id`,
  - `fragment_num`,
  - `final`,
  - `contents`.
- [`send_in_fragments()`](#how-to-send-data-in-fragments) transforms `Network::Fragment` into network order string.

#### `Network::Fragment` -> `Network::Packet` -> `Crypto::Message`

- `Connection::send()` transforms the above network order string into `Network::Packet`.
- `Network::Packet` belongs to in [datagram layter](ref.md#datagram-layer).
- `Network::Packet` contains the following fields:
  - `seq`,
  - `timestamp`,
  - `timestamp_reply`,
  - `payload`,
  - `direction`.
- [`Connection::send()`](#how-to-send-a-packet) transfroms `Network::Packet` into `Crypto::Message`.
- `Crypto::Message` is a utility class for crypto.
- `Crypto::Message` contains the following fields:
  - `Nonce`: contains `direction` and `seq` fields in `Network::Packet`.
  - `text`: contains `timestamp`, `timestamp_reply` and `payload` fields in `Network::Packet`.
- [`Connection::send()`](#how-to-send-a-packet) encrypts `Crypto::Message`.
- [`Connection::send()`](#how-to-send-a-packet) sents `Crypto::Message` to remote server in UDP datagram.
