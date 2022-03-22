# Mosh client coding analysis

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

### STMClient constructor

- Set the `ip`,`port`,`key`,`predict_mode`,`verbose` parameter. Here, `network` is `NULL`.
- `Overlay::OverlayManager overlays` is initialized in construction function.
  - `overlays` contains `NotificationEngine`, `PredictionEngine`, `TitleEngine`. The design of these engine is unclear.
- `Terminal::Display` is initialized in construction function.
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

#### Transport<MyState, RemoteState>::Transport

- Initialize the `connection`, initialize the sender with `connection` and `initial_state`.
  - `connection` is the underlying, encrypted network connection.

### STMClient::main

`STMClient::main()` calls [`main_init()`](#stmclientmain_init) to initialize signal handling and structures. In the main loop(while loop), It performs the following steps:

- Output terminal content to the `STDOUT_FILENO` via [`output_new_frame()`](#how-to-output-content).
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

#### How to output content

- `output_new_frame()` aka `STMClient::output_new_frame()`.
- `output_new_frame()` gets the `Framebuffer` from the latest state in `received_states`.
- `output_new_frame()` calls `overlays.apply()` to apply local overlays.
- `output_new_frame()` calls `display.new_frame()` to calculate minimal `diff` from where we are.
- `output_new_frame()` writes the `diff` to `STDOUT_FILENO`.
- `output_new_frame()` sets `repaint_requested` to true.
- `output_new_frame()` sets `local_framebuffer` to the new state.

#### STMClient::main_init

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
  - Note `sent_states` is type of list `TimestampedState`, while `current_state` is type of `MyState`.
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

<!--
TODO What's the behavior of the serverside.
TODO what the purpose of `overlay`.
TODO what the meaning of `display`.
TODO In fragment, if f.contents size is smaller than MTU, how to know the content size?
-->

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
- `network->recv()` applies `diff` to reference state, if `diff` is not empty.
  - It creates a new state, which applies the `diff` to the reference state.
- `network->recv()` inserts new state if out-of-order state is received, `network->recv()` returns directly.
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
