# Mosh server coding analysis

## mosh-server.cc

`mosh-server` command parameters.

```sh
Usage: mosh-server new [-s] [-v] [-i LOCALADDR] [-p PORT[:PORT2]] [-c COLORS] [-l NAME=VALUE] [-- COMMAND *]
Usage: mosh-server --help
Usage: mosh-server --version
```

In the `main` function: `run_server` is the core to start `mosh` server.

- `main()` calls `Crypto::disable_dumping_core()` to make sure we don't dump core.
- `main()` parses `new` command-line syntax.
- `main()` checks port range.
- `main()` gets shell name and prepares for shell command arguments.
- `main()` makes sure UTF8 locale is set.
- `main()` calls [`run_server()`](#run_server) with the ip, port and shell path, shell arguments,(etc.) as parameters.

### run_server

- `run_server()` gets network idle timeout.
- `run_server()` gets network signaled idle timeout.
- `run_server()` gets initial window size. They will be overwritten by client on first connection.
- `run_server()` calls [`Terminal::Complete()`](client.md#terminalcomplete) to open parser and terminal.
- `run_server()` creates blank [`Network::UserStream`](client.md#networkuserstream).
- `run_server()` creates ` network`, which is type of [`ServerConnection`](#serverconnection).
- `run_server()` sets the `verbose` mode via `network->set_verbose()`.
- `run_server()` sets the `verbose` mode via `Select::set_verbose()`.
- `run_server()` calls `network->port()` to [get the port string](#how-to-get-port-string) representation.
- `run_server()` calls `network->get_key()` to [get the session key string](#how-to-get-session-key-string) representation.
- `run_server()` prints port and session key to the standard output. The output starts with `"MOSH CONNECT "`.
- `run_server()` ignores signal `SIGHUP`, `SIGPIPE`.
- `run_server()` calls `fork()` to detach from terminal.
  - Parent process prints the license information and terminates.
  - Child process continues.
- `run_server()` redirects `STDIN_FILENO`, `STDOUT_FILENO`, `STDERR_FILENO` to `"/dev/null"` for non-verbose mode.
- `run_server()` calls [`forkpty()`](#forkpty) to create a new process operating in a pseudo terminal.

#### forkpty

- The child process, which will run a shell process:
  - Re-enable signals `SIGHUP`, `SIGPIPE` with default value.
  - Close server-related socket file descriptors, via calling `delete`.
  - Set terminal [UTF8 support](#how-to-set-the-terminal-utf8-support).
  - Set `"TERM"` envrionment variable to be `"xterm"` or `"xterm-256color"` based on `"-c color"` option.
  - Set `"NCURSES_NO_UTF8_ACS"` envrionment variable
    - to ask ncurses to send UTF-8 instead of ISO 2022 for line-drawing chars.
  - Clear `"STY"` envrionment variable so GNU screen regards us as top level.
  - Change to the home directory, via calling [`chdir_homedir()`](#chdir_homedir):
  - If `.hushlogin` file don't exist and `with_motd` is true,
    - print the motd from `"/run/motd.dynamic"`,
    - or print the motd from `"/var/run/motd.dynamic"` and `"/etc/motd"`.
  - Print warning message if there is unattached mosh session, via calling [`warn_unattached()`](#warn_unattached).
    - See the following parent process to understand mosh session.
  - Wait for parent to release us, via calling `fgets()` for `stdin`.
  - Enable core dump, via calling `Crypto::reenable_dumping_core()`.
  - Execute the shell command with arguments, via calling `execvp()`.
  - Terminate the child process, via calling `exit()`.
- The parent process, which will run the mosh server process:
  - Add utmp record via calling [`utempter_add_record()`](https://www.unix.com/man-page/suse/8/utempter),
    - with `master` as pseudo-terminal master file descriptor, `"mosh [%ld]"` as host name.
  - Serve the client, via calling [`serve()`](#serve).
    - with `master`, `terminal`, `network` as parameters.
  - Delete utmp record via calling [`utempter_remove_record()`](https://www.unix.com/man-page/suse/8/utempter),
    - with `master` as pseudo-terminal master file descriptor.
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
- Call `setenv()` to set the `"PWD"` envrionment variable.

#### warn_unattached

- `warn_unattached()` calls `getpwuid(getuid())` to get the current user.
- `warn_unattached()` checks the records in `utmp` file
  - If the `ut_user` field is the same user as the current user, via calling `getpwuid(getuid()`,
  - If the `ut_type` field is `USER_PROCESS`,
  - If the `ut_host` field does look like `"mosh [%ld]"`, where `%ld` is the process ID,
  - If the `ut_host` field doesn't euqal `ignore_entry`, which is the mosh session,
  - If pseudo-terminal device identified by the `ut_line` field exist,
  - Pushes the `ut_host` into `unattached_mosh_servers` vector.
- `warn_unattached()` returns if `unattached_mosh_servers` vector is empty.
- `warn_unattached()` prints warning message to `STDOUT`, if there exists unattached sessions.

### serve

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
- Upon network sockets is ready to read, process it with [`network.recv()`](client.md#how-to-receive-network-input).
  - After `network.recv()`, the remote state is saved in `received_states`,
  - and the opposite direction `ack_num` is saved.
  - If new user input available for the terminal,
  - TODO
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
- Upon pty master input is ready to read, process it with
