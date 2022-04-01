# Mosh server coding analysis

![mosh-comm.svg](img/mosh-comm.svg)

## mosh-server.cc

`mosh-server` command parameters.

```sh
Usage: mosh-server new [-s] [-v] [-i LOCALADDR] [-p PORT[:PORT2]] [-c COLORS] [-l NAME=VALUE] [-- COMMAND *]
Usage: mosh-server --help
Usage: mosh-server --version
```

- [Run the server](#run_server)
- [Serve the client](#serve)
- [How to read from client connection](#how-to-read-from-client-connection)

### main

```cpp
int main(int argc, char* argv[]) {
	//+--184 lines: folding ----------
	try {
		return run_server(desired_ip, desired_port, command_path, command_argv, colors, verbose, with_motd);
	} catch (const Network::NetworkException& e) {
		fprintf(stderr, "Network exception: %s\n", e.what());
		return 1;
	} catch (const Crypto::CryptoException& e) {
		fprintf(stderr, "Crypto exception: %s\n", e.what());
		return 1;
	}
}
```

In the `main` function: `run_server` is the core to start `mosh` server.

- `main()` calls `Crypto::disable_dumping_core()` to make sure we don't dump core.
- `main()` parses command parameters
- `main()` checks port range.
- `main()` prepares shell name and shell command arguments.
- `main()` makes sure UTF8 locale is set.
- `main()` calls [`run_server()`](#run_server) with the ip, port and shell path, shell arguments,(etc.) as parameters.

### run_server

- `run_server()` gets network idle timeout.
- `run_server()` gets network signaled idle timeout.
- `run_server()` gets initial window size. They will be overwritten by client on first connection.
- `run_server()` calls [`Terminal::Complete()`](client.md#terminalcomplete) to open parser and terminal.
- `run_server()` creates blank [`Network::UserStream`](client.md#networkuserstream) for newtork.
- `run_server()` initializes network, which is [`ServerConnection`](#serverconnection).
- `run_server()` sets the `verbose` mode via `network->set_verbose()`.
- `run_server()` sets the `verbose` mode via `Select::set_verbose()`.
- `run_server()` calls `network->port()` to [get the port string](#how-to-get-port-string) representation.
- `run_server()` calls `network->get_key()` to [get the session key string](#how-to-get-session-key-string) representation.
- `run_server()` prints port and session key to the standard output. The output starts with "MOSH CONNECT ".
- `run_server()` ignores signal `SIGHUP`, `SIGPIPE`.
- `run_server()` calls `fork()` to detach from terminal.
  - Parent process prints the license information and exits.
  - Child process continues.
- `run_server()` redirects `STDIN_FILENO`, `STDOUT_FILENO`, `STDERR_FILENO` to `"/dev/null"` for non-verbose mode.
- `run_server()` calls [`forkpty()`](#forkpty) to create a new process operating in a pseudo terminal.

#### forkpty

- The child process, which will run a shell process:
  - Re-enable signals `SIGHUP`, `SIGPIPE` with default value.
  - Close server-related socket file descriptors, via calling `delete`.
  - Set terminal [UTF8 support](#how-to-set-the-terminal-utf8-support).
  - Set "TERM" environment variable to be "xterm" or "xterm-256color" based on "-c color" option.
  - Set `"NCURSES_NO_UTF8_ACS"` environment variable
    - to ask ncurses to send UTF-8 instead of ISO 2022 for line-drawing chars.
  - Clear `"STY"` environment variable so GNU screen regards us as top level.
  - Change to the home directory, via calling [`chdir_homedir()`](#chdir_homedir):
  - If `.hushlogin` file don't exist and `with_motd` is true,
    - print the motd from `"/run/motd.dynamic"`,
    - or print the motd from `"/var/run/motd.dynamic"` and `"/etc/motd"`.
    - Print warning message if there is unattached mosh session, via calling [`warn_unattached()`](#warn_unattached).
    - See the following parent process to understand mosh session.
  - Wait for parent to release us, via calling `fgets()` for `stdin`.
  - Enable core dump, via calling `Crypto::reenable_dumping_core()`.
  - Execute the shell command with arguments, via calling `execvp()`.
  - If error happens during `execvp()`, Terminate the child process, via calling `exit()`.
- The parent process, which will run the mosh server process:
  - Add utmp record via calling [`utempter_add_record()`](https://www.unix.com/man-page/suse/8/utempter),
    - with `master` as pty master parameter, `"mosh [%ld]"` as host name parameter.
    - as login service update utmp record is required.
  - Serve the client, via calling [`serve()`](#serve).
    - with `master`, `terminal`, `network` as parameters.
  - Delete utmp record via calling [`utempter_remove_record()`](https://www.unix.com/man-page/suse/8/utempter),
    - with `master` as pty master parameter.
    - as login service update utmp record is required.
  - Close the master pseudo-terminal.
  - Close server-related socket file descriptors, via calling `delete`.
  - Print exiting message.

#### `ServerConnection`

- `ServerConnection` aka `Network::Transport<Terminal::Complete, Network::UserStream>`.
- `Network::Transport` is constructed with:
  - the `terminal` which is type of `Terminal::Complete`,
  - the `blank` which is type of `UserStream`,
  - `ip`, `port` as parameters.
- `Network::Transport` has a `Connection`, which represents the underlying, encrypted network connection.
- `Network::Transport` has a `TransportSender<Terminal::Complete>`, which represents the sender.
- `Network::Transport` has a `list<TimestampedState<Network::UserStream>>`, which represents receiver.
- `Network::Transport` calls `connection(desired_ip, desired_port)` to [initialize the connection](#how-to-initialize-connection).
- `Network::Transport` calls `sender(connection, initial_state)` to [initialize sender](#how-to-initialize-sender).
- In the constructor of `Network::Transport`,
  - `received_states` is a list type of `TimestampedState<Network::UserStream>`.
  - `received_states` is initialized with the `blank`as parameter.
  - `received_states` adds the `blank` to its list.
- `Network::Transport()` set `receiver_quench_timer` to zero.
- `Network::Transport()` set `last_receiver_state` to be `terminal`.
- `Network::Transport()` creates `fragments`, which is type of `FragmentAssembly`.

#### How to initialize connection

- `connection(desired_ip, desired_port)` is called to create the connection with server.
- `connection()` is the constructor of `Network::Connection`
- `connection()` initializes a empty deque of `Socket`: `socks`.
- `connection()` initializes `has_remote_addr` to true.
- `connection()` initializes the `key`, which is type of `Base64Key`
  - `Base64Key` reads 16 bytes from `/dev/urandom` as the `key`.
- `connection()` initializes `session` with `key` as parameter,
  - `session` object is used to encrypt/decrypt message.
- `connection()` calls `setup()` to set the `last_port_choice` to current time.
- `connection()` calls `parse_portrange()` to parse port range from `desired_port` parameter.
- `connection()` calls [`try_bind()`](#try_bind) to bind the port to network interface.
  - If `desired_ip` is given, use `desired_ip` as parameter to call `try_bind()`.
  - `try_bind()` is called with port range parameters.
- `connection()` returns if `try_bind()` returns true.

#### How to initialize sender

- `sender(connection, initial_state)` is called to initialize the sender.
- `sender()` is the constructor of `TransportSender<Terminal::Complete>`.
- `sender()` initializes `connection` pointer with the `connection` as parameter.
- `sender()` initializes `current_state` with the `initial_state` as parameter.
- `sender()` initializes `sent_states` list with the `initial_state` as the first state.

#### try_bind

- `try_bind()` initializes a `AddrInfo` object with `desired_ip` as parameter.
  - `AddrInfo` calls `getaddrinfo()` to get the `addrinfo` object.
- `try_bind()` creates a `Socket` and pushes it into `socks` deque.
  - `Socket` uses `setsockopt()` to set socket options: `IP_MTU_DISCOVER`, `IP_TOS`, `IP_RECVTOS`.
- `try_bind()` searches the port range, calls `bind()` bind the server socket to that port.
  - If `bind()` returns successfully. `try_bind()` calls `set_MTU()` to set the MTU and return true.
  - If `bind()` fails, `try_bind()` throw exceptions and return false.

#### How to get port string

- `network->port()` calls `connection.port()` to get the port.
- `connection.port()` calls `getsockname()` to get the `sockaddr` of server socket.
- `connection.port()` calls `getnameinfo()` to get the port string representation.

#### How to get session key string

- `network->get_key()` calls `connection.get_key()` to get the session key.
- `connection.get_key()` calls `key.printable_key()` instead.
- `key.printable_key()` aka `Base64Key::printable_key()`.
- `key.printable_key()` calls `base64_encode()` to show the `key` base64 representation.

#### How to set the terminal UTF8 support

- Get the `termios` struct for `STDIN_FILENO`, via calling `tcgetattr()`.
- Set `IUTF8` flag for `termios`.
- Set the `termios` struct for `STDIN_FILENO`, via calling `tcsetattr()`.

#### chdir_homedir

- Call `getenv()` or `getpwuid(getuid())` to get the `home` path.
- Call `chdir()` to change to the `home` path.
- Call `setenv()` to set the `"PWD"` environment variable.

#### warn_unattached

- `warn_unattached()` calls `getpwuid(getuid())` to get the current user.
- `warn_unattached()` checks the records in `utmp` file
  - If the `ut_user` field is the same user as the current user, via calling `getpwuid(getuid()`,
  - If the `ut_type` field is `USER_PROCESS`,
  - If the `ut_host` field does look like `"mosh [%ld]"`, where `%ld` is the process ID,
  - If the `ut_host` field isn't equal to `ignore_entry`, which is the mosh session,
  - If pseudo-terminal device identified by the `ut_line` field exist,
  - Pushes the `ut_host` into `unattached_mosh_servers` vector.
- `warn_unattached()` returns if `unattached_mosh_servers` vector is empty.
- `warn_unattached()` prints warning message to `STDOUT`, if there exists unattached sessions.

#### serve

- `serve()` initializes the singleton `Select` object: `sel`.
- `serve()` [registers signal handler](client.md#how-to-register-signal-handler) in `sel`, for `SIGTERM`, `SIGINT`, `SIGUSR1`.
- `serve()` gets the latest remote state number, via calling `network.get_remote_state_num()`.
- `serve()` sets `child_released` to false.
- `serve()` sets `connected_utmp` to false, initializes `saved_addr` and `saved_addr_len` to zero and false.

In the main loop(while loop), It performs the following steps:

- Calculate `timeout` based on `network.wait_time()` and `terminal.wait_time()`.
- Clear file descriptor, via calling `sel.clear_fds()`.
- Get network socket, via calling `network->fds()`, `fd_list.back()`.
- Add network socket to `Select` object.
- Add pty master to `Select` object.
- Wait for socket input, signal, pty master input via calling `sel.select()`, with the `timeout` as parameter.
- Upon receive signals, the corresponding item in `Select.got_signal` array is set.
- Upon network sockets is ready to read, read it with [`network.recv()`](client.md#how-to-receive-network-input).
  - After `network.recv()`, the remote state is saved in `received_states`,
  - and the opposite direction `ack_num` is saved.
  - [forward input to terminal](#how-to-forward-input-to-terminal) if remote state number is not equal to `last_remote_num`.
- Upon pty master input is ready to read, read it via calling `read()` system call.
  - If it does read some input, call [`terminal.act()`](#terminalactstring) to process the input data.
  - append the return value of `terminal.act()` to `terminal_to_host`.
  - set currrent state via calling `network.set_current_state()` with the `terminal` as parameter.
- `swrite` TODO

#### How to forward input to terminal

- Update `last_remote_num` with the latest one.
- Initialize a empty `UserStream` object: `us`.
- Find the difference between new state and `last_receiver_state`, via calling [`network.get_remote_diff()`](#get_remote_diff).
- Initialize the `UserStream` from the above difference string via calling [`apply_string()`](#apply_string).
- Iterate through the above `us`,
  - get the `action` object of type `Parser::Action` via calling `us.get_action()`.
  - `UserEvent` object contains `Parser::UserByte` object or `Parser::Resize` object.
  - `Parser::UserByte` and `Parser::Resize` are sub-class of `Parser::Action`.
  - For `Resize` action:
    - skip the consecutive Resize action,
    - convert the action into `Parser::Resize`,
    - get the window size for `STDIN_FILENO` , via `ioctl()` and `TIOCGWINSZ` flag,
    - set the window size for `STDIN_FILENO` , via `ioctl()` and `TIOCSWINSZ` flag.
  - For other action:
    - call [`terminal.act()`](#terminalact) to get the transcript character.
    - append the transcript character to `terminal_to_host`.
- If `us` is not empty,
  - register input frame number for future echo ack via calling `terminal.register_input_frame()`.
- Set the current state via calling `network.set_current_state()`, if `network` is not shutdown.
- If `connected_utmp` is false and `saved_addr_len` is different from `network.get_remote_addr_len()`,
  - delete `utmp` record via calling [`utempter_remove_record()`](https://www.unix.com/man-page/suse/8/utempter), with pty master as parameter.
  - store the value from `network.get_remote_addr()` into `saved_addr_len`,
  - store the value from `network.get_remote_addr_len()` into `saved_addr_len`,
  - get the `host` name via calling `getnameinfo()`,
  - add `utmp` record via calling [`utempter_add_record`](https://www.unix.com/man-page/suse/8/utempter), with pty master and `host` as parameters.
  - set `connected_utmp` to true.
- If `child_released` is false, release the child process via writing `\n` to pty master,
  - upon receive the empty line input, the child process will start the shell.
  - set `child_released` to true.

#### get_remote_diff

- `get_remote_diff()` aka `Transport<MyState, RemoteState>::get_remote_diff()`.
- Here `RemoteState` is `UserStream`.
- `get_remote_diff()` calls `diff_from()` to [calculate diff](client.md#how-to-calculate-the-diff-for-userstream).
  - `diff_from()` to compare the newest `received_states` with `last_receiver_state` .
  - `diff_from()` returns the difference string representation of the `ClientBuffers::UserMessage` object.
- Next `get_remote_diff()` rationalizes `received_states`.
  - `get_remote_diff()` sets the `oldest_receiver_state` with the value of the oldest `received_states`.
  - `get_remote_diff()` iterates through the `received_states` list in reverse order (newest to oldest).
  - `oldest_receiver_state` is the target to be evluated for each iteration.
  - For eache iterating state, calls `UserStream::subtract()` to subtract shared `UserEvent`.
  - `get_remote_diff()` stores the newest `received_states` in `last_receiver_state`.
- `get_remote_diff()` returns the difference string representation of `ClientBuffers::UserMessage`.

#### apply_string

- `apply_string()` aka `UserStream::apply_string()`.
- `apply_string()` creates a `ClientBuffers::UserMessage` object.
- `apply_string()` parses string representation of `ClientBuffers::UserMessage`.
- `apply_string()` iterates through `ClientBuffers::UserMessage` object.
- `apply_string()` extracts `UserByte` or `Resize` from `ClientBuffers::UserMessage`.
- `apply_string()` wraps `UserByte` or `Resize` in `UserEvent` and pushes `UserEvent` into `UserStream`.
- For each iteration, `apply_string()` builds `UserEvent` and pushes it into `UserStream.actions`.
- That means `apply_string()` initializes `UserStream` with `ClientBuffers::UserMessage`.

#### terminal.act

- `terminal.act()` aka `Complete::act()`.
- `terminal.act()` apply action to terminal via calling `act->act_on_terminal()` with the `terminal` as parameter.
  - For `UserByte` action, [`UserByte::act_on_terminal`](#userbyteact_on_terminal) will store the transcript character to `dispatch.terminal_to_host`.
  - For `Resize` action, [`Resize::act_on_terminal`](#resizeact_on_terminal) will change the terminal frame buffer size.
- `terminal.act()` returns `terminal.read_octets_to_host()`.
  - `terminal.read_octets_to_host()` aka `Emulator::read_octets_to_host()`.
  - `terminal.read_octets_to_host()` reads `dispatch.terminal_to_host` to `ret`.
  - `terminal.read_octets_to_host()` clears `dispatch.terminal_to_host`.
  - `terminal.read_octets_to_host()` returns `ret`.
- The above implementation means `terminal.act()` return the transcript character to caller.

#### UserByte::act_on_terminal

- `act_on_terminal()` has parameter `Terminal::Emulator* emu`.
- `act_on_terminal()` calls `emu->user.input()` to convert user's cursor control sequence to ANSI cursor control sequence.
  - `emu->user.input()` has `UserByte` parameter and `application_mode_cursor_keys` parameter.
  - `emu->user.input()` checks the `UserByte` parameter.
  - For `Ground` character, `emu->user.input()` returns the raw character string.
  - For "0x1b" character, `emu->user.input()` sets state to `ESC` and returns raw character string.
  - If state is `ESC` and character is `O`,
    - `emu->user.input()` sets state to `SS3` and return empty string.
  - If state is `ESC` and character isn't `O`,
    - `emu->user.input()` sets state to `Ground` and return raw character string.
  - If state is `SS3` and character isn't `A-D` and `application_mode_cursor_keys` is false,
    - `emu->user.input()` sets state to `Ground` and return `ESC [ [A-D]` string.
  - If state is `SS3` and character isn't `A-D` and `application_mode_cursor_keys` is true,
    - `emu->user.input()` sets state to `Ground` and return `ESC O [A-D]` string.
- `act_on_terminal()` calls `emu->dispatch.terminal_to_host.append()` to append the above string to terminal `dispatch.terminal_to_host`.
- `act_on_terminal()` returns void.

#### Resize::act_on_terminal

- `act_on_terminal()` has parameter `Terminal::Emulator* emu`.
- `act_on_terminal()` calls `emu->resize()` to adjust terminal frame buffer size.
  - `emu->resize()` aka `Emulator::resize()`.
  - `emu->resize()` has `s_width` and `s_height` parameter.
  - `emu->resize()` calls `fb.resize()` to finish the job.
  - `fb.resize()` aka `Framebuffer::resize()`.
  - `fb.resize()` has `s_width` and `s_height` as parameters.
  - `fb.resize()` adjust `Framebuffer.ds` size according to the width and height parameters.
  - `fb.resize()` adjust `Framebuffer.row` size according to the width and height parameters.
  - The above implementation means the `Framebuffer.ds` and `Framebuffer.row` is changed according to the parameters.
- `act_on_terminal()` returns void.

#### terminal.act(string)

- `terminal.act(string)` aka `Complete::act(string)`.
- Iterate the string parameter,
- For each character, call `parser.input()` to parse octet into up to three actions.
  - Iterate the `actions`, for each action,
  - apply action on terminal via calling `act->act_on_terminal()` with the `terminal` as parameter.
- Clear the `actions` via calling `actions.clear()`.
- Return `terminal.read_octets_to_host()`.

![mosh-parse.svg](img/mosh-parse.svg)

#### Parse unicode character to action

- The first `parser.input()` is actually `Parser::UTF8Parser::input()`.
- If ASCII code ( less than "0x7f") and `buf_len` is 0,
  - call the [second `parser.input()`](#parse-wide-character-according-to-transition) with the `actions` as parameter.
  - `actions` aka `Complete.actions`, a vector of type `Parser::Action`
  - return early.
- Assign the `c` to `buf[buf_len++]`.
- Parse the `buf` fields of `Parser::UTF8Parser` in a loop.
  - According to Unicode 6.0, section 3.9 [Best Practices for using U+FFFD](https://www.unicode.org/versions/Unicode6.0.0/ch03.pdf).
  - Convert multi-byte sequence to wide character via calling `mbrtowc()`.
  - Call the [second `parser.input()`](#parse-wide-character-according-to-transition) for the wide character with the `actions` as parameter.
  - Continue the loop until all byte is parsed.

#### Parse wide character according to `Transition`

- The second `parser.input()` is actually `Parser::Parser::input()`.
- `parser.input()` calls `state->input()` to [parse the wide character to `Transition`](#parse-wide-character-to-transition).
- `parser.input()` calls `append_or_delete()` if `tx.next_state` is not NULL.
- `append_or_delete()` decides whether to push the `Action`:`state->exit()` into `actions`.
- `parser.input()` calls `append_or_delete()` to push the `Action`:`tx.action` into `actions`.
- `parser.input()` clears `tx.action`.
- `parser.input()` calls `append_or_delete()` if `tx.next_state` is not NULL.
- `append_or_delete()` decides whether to push the `Action`:`tx.next_state->enter()` into `actions`.
- `parser.input()` updates `state` with `tx.next_state`.

#### Parse wide character to `Transition`

- `state->input()` parses the character into `Transition` via calling `anywhere_rule()`.
  - `anywhere_rule()` creates `Transition` based on character coding rule.
- If the created `Transition.next_state` is not empty,
  - `state->input()` assigns value to `char_present` and `ch` fields of `Transition.action`.
  - returns early with created `Transition`.
- `state->input()` parses the character into `Transition` via calling `this->input_state_rule()`.
  - `this->input_state_rule()` parses high Unicode code-points.
  - The behaviour of `this->input_state_rule()` depends on the implementation of `State` sub-class.
  - The default `State` is `Paser:Ground`.
  - `Ground::input_state_rule()` parses character according to `C0_prime` rule and `GLGR` rule.
  - `C0_prime` rule returns a `Transition` whose `action` field is `Parser::Execute`.
  - `GLGR` rule returns a `Transition` whose `action` field is `Parser::Print`.
  - `Ground::input_state_rule()` returns the second `Transition` whose `action` field is `Parser::Ignore`.
- `state->input()` assigns value to `char_present` and `ch` fields of the second `Transition.action`.
- Return with the second created `Transition`.

### How to read from client connection

Upon server network socket is ready to read, `connection->recv()` starts to read it.

#### `Crypto::Message` -> `Network::Packet` -> `Network::Fragment`

- [`Connection::recv_one()`](client.md#how-to-receive-datagram-from-socket) reads UDP datagram from server socket.
- [`Connection::recv_one()`](client.md#how-to-receive-datagram-from-socket) decrypts the UDP datagram into `Crypto::Message`.
- `Crypto::Message` is a utility class for crypto.
- `Crypto::Message` contains the following fields:
  - `Nonce`: contains `direction` and `seq` fields in `Network::Packet`.
  - `text`: contains `timestamp`, `timestamp_reply` and `payload` fields in `Network::Packet`.
- [`Connection::recv_one()`](client.md#how-to-receive-datagram-from-socket) transforms `Crypto::Message` into `Network::Packet`.
- `Network::Packet` belongs to in [datagram layter](ref.md#datagram-layer).
- `Network::Packet` contains the following fields:
  - `seq`,
  - `timestamp`,
  - `timestamp_reply`,
  - `payload`,
  - `direction`.
- [`Connection::recv_one()`](client.md#how-to-receive-datagram-from-socket) returns string representation of `Network::Packet` paylod field.
- [`Connection::recv()`](client.md#how-to-read-data-from-socket) returns string representation of `Network::Packet` payload field.
- [`Fragment()`](client.md#how-to-create-the-frament-from-string) transform `Network::Packet` into `Network::Fragment`
- `Network::Fragment` contains the following fields:
  - `id`,
  - `fragment_num`,
  - `final`,
  - `contents`.

#### `Network::Fragment` -> `TransportBuffers.Instruction` -> `Network::UserStream`

- [`fragments.get_assembly()`](client.md#how-to-build-instruction-from-fragments) aka `FragmentAssembly::get_assembly()`.
- [`fragments.get_assembly()`](client.md#how-to-build-instruction-from-fragments) concatenates the `contents` field of each `Network::Fragment` into one string.
- [`fragments.get_assembly()`](client.md#how-to-build-instruction-from-fragments) decompress the string and transforms it into `TransportBuffers.Instruction`.
- `TransportBuffers.Instruction` is the "state" in [transport layter](ref.md#transport-layer).
- `TransportBuffers.Instruction` contains the following fields:
  - `old_num`,
  - `new_num`,
  - `ack_num`,
  - `throwaway_num`,
  - `diff`.
- [`connection->recv()`](client.md#how-to-receive-network-input) creates an empty `TimestampedState<Network::UserStream>`.
- `Network::UserStream` is wrapped in `TimestampedState<Network::UserStream>`.
- `TimestampedState<Network::UserStream>` contains the following fields:
  - `timestamp`,
  - `num`,
  - `state`.

#### `Parser::UserByte` -> `Network::UserEvent` -> `Network::UserStream`

- [`network.get_remote_diff()`](#get_remote_diff) compares the newest `Network::UserStream` and existing `Network::UserStream`.
- [`network.get_remote_diff()`](#get_remote_diff) returns the difference string representation of `ClientBuffers::UserMessage`.
- [`UserStream.apply_string()`](#apply_string) initializes `UserStream` with `ClientBuffers::UserMessage`.
- `Network::UserStream` contains a deque of type `Network::UserEvent`.
- `Network::UserEvent` contains the following fields:
  - `type`,
  - `userbyte`,
  - `resize`.
- [`connection->recv()`](client.md#how-to-receive-network-input) stores `TimestampedState<Network::UserStream>` in `received_states`.
