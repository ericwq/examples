# Mosh client coding analysis

![mosh-comm.svg](img/mosh-comm.svg)

## mosh-client.cc

`mosh-client` command options

```sh
Usage: mosh-client [-# 'ARGS'] IP PORT
Usage: mosh-client -c
Usage: mosh-client --help
Usage: mosh-client --version
```

- [STMClient constructor](#stmclient-constructor)
- [STMClient::init](#stmclientinit)
- [STMClient::main](#stmclientmain)
- [How to send keystroke to remote server](#how-to-send-keystroke-to-remote-server)
- [How to receive state from server](#how-to-receive-state-from-server)
- [How does the notification engine decide to show message?](#how-does-the-notification-engine-decide-to-show-message)
- [How does the prediction engine work?](#how-does-the-prediction-engine-work)

### main

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

In the `main` function, `STMClient` is the core to start `mosh` client.

- `main()` calls `Crypto::disable_dumping_core()` to make sure we don't dump core.
- `main()` parses command options.
- `mian()` gets ip argument and checks port argument range.
- `main()` reads key from environment "MOSH_KEY".
- `main()` reads prediction preference from environment "MOSH_PREDICTION_DISPLAY".
- `main()` queries locale.
- `main()` starts the `STMClient`.

### STMClient constructor

- `STMClient()` is called withe following parameters:`ip`,`port`,`key`,`predict_mode`,`verbose` parameter.
- `STMClient()` saves `key`, `ip`, `port` parameters as field member.
- `STMClient()` initializes `escape_key`, `escape_pass_key`, `escape_pass_key2`,
  - with "0x1E", "^", "^" corresponding value.
- `STMClient()` initializes `escape_requires_lf` with false value.
- `STMClient()` initializes empty `termios` structs: `raw_termios`, `saved_termios`.
- `STMClient()` initializes `new_state` frame buffer with 1\*1 size.
- `STMClient()` initializes `local_framebuffer` frame buffer with 1\*1 size.
- `STMClient()` initializes `overlay`, which is type of [`Overlay::OverlayManager`](#overlayoverlaymanager).
- `STMClient()` initializes a NULL `network`.
- `STMClient()` initializes `display`, which is type of [`Terminal::Display`](#terminaldisplay).
- `STMClient()` initializes `repaint_requested`, `lf_entered`, `quit_sequence_started`, `clean_shutdown` with false value.
- `STMClient()` initializes `display_preference` in prediction engine with the value from `predict_mode` parameter.

#### Overlay::OverlayManager

- `OverlayManager` has a `NotificationEngine`, which performs the notification work.
- `OverlayManager` has a `PredictionEngine`, which performs the prediction work.
- `OverlayManager` has a `TitleEngine`, which performs the title work.
- The default constructor initializes a `OverlayManager` without any parameters.

#### Terminal::Display

- `Display()` is initialized with true `use_environment`.
- `Display()` calls [`setupterm()`](https://linux.die.net/man/3/setupterm) to read in the `terminfo` database, initialize the `terminfo` structures.
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
- [Application Cursor Keys mode](https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-PC-Style-Function-Keys) is a way for the server to change the control sequences sent by the arrow keys.
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
- `output_new_frame()` gets `new_state` from the latest state saved in `received_states`.
- `new_state` is of type `Terminal::Framebuffer`.
- `output_new_frame()` calls `overlays.apply()` to [apply to local overlays](#overlaymanagerapply) with `new_state` as parameter.
- `output_new_frame()` calls [`display.new_frame()`](#how-to-calculate-frame-buffer-difference) to calculate minimal `diff` from where we are.
- `output_new_frame()` writes the `diff` to `STDOUT_FILENO`.
- `output_new_frame()` sets `repaint_requested` to false.
- `output_new_frame()` sets `local_framebuffer` to the new state.

#### OverlayManager.apply

- `overlays.apply()` aka `OverlayManager::apply()`
- `overlays.apply()` calls `predictions.cull()` to [prepare the engine for prediction](#predictionenginecull).
- `overlays.apply()` calls `predictions.apply()` to [manipulate frame buffer for prediction](#predictionengineapply).
- `overlays.apply()` calls `notifications.adjust_message()` to clear the `message` if time expires.
- `overlays.apply()` calls `notifications.apply()` to [draw notification string](#notificationengineapply) in frame buffer.
- `overlays.apply()` calls `title.apply()` to set up the `window_title` and `icon_name` in frame buffer.

#### NotificationEngine.apply

- `notifications.apply()` aka `NotificationEngine::apply`.
- Initialize `time_expired` by checking `last_word_from_server>6500` or `last_acked_state>10000`.
- Return early, if time does not expire and the `message` is empty.
- Hide cursor if cursor row in frame buffer is zero.
- Draw bar across top of screen. Adding a new row of `Cell` in farme buffer.
- Prepare notification string which compose of `explanation`, `time_elapsed`, `keystroke_str` and `message`.
- Draw notification string on top of screen.
  - Support normal character, wide character, combining character(up to 32)
  - [Combining Diacritical Marks](https://unicode.org/charts/PDF/U0300.pdf)
  - [Combining Diacritical Marks Supplement](https://unicode.org/charts/PDF/U1DC0.pdf),
  - [Combining Diacritical Marks for Symbols](https://unicode.org/charts/PDF/U20D0.pdf)
  - [Combining Half Marks](https://unicode.org/charts/PDF/UFE20.pdf)

#### PredictionEngine.apply

- `predictions.apply()` aka `PredictionEngine::apply()`.
- `apply()` has a `Framebuffer` parameter.
- Decide whether to show the prediction based on `display_preference`, `srtt_trigger`, `glitch_trigger`.
- Iterate through the `cursors` list, calls the second `apply()` method for each `ConditionalCursorMove` object.
  - `apply()` method aka `ConditionalCursorMove::apply()`.
  - If it's not `active`, returns.
  - If it's `tentative`: `tentative_until_epoch > confirmed_epoch`, returns.
  - Ensures `row < fb.ds.get_height()` and `col < fb.ds.get_width()`.
  - Ensures frame buffer's `DrawState.origin_mode` is false.
  - [Move the cursor to `row`](#drawstatemove_row) via calling `fb.ds.move_row()`.
  - [Move the cursor to `col`](#drawstatemove_col) via calling `fb.ds.move_col()`.
- Iterate through the `overlays` list, calls `apply()` method for each `ConditionalOverlayRow` object.
  - `apply()` method aka `ConditionalOverlayRow::apply`.
  - `apply()` iterates through each cell in the row.
  - `apply()` calls the third [`apply()` method for each `ConditionalOverlayCell`](#conditionaloverlaycellapply) object.

#### DrawState.move_row

- `move_row()` aka `DrawState::move_row()`
- `move_row()` has a `N` parameter, which is the row number, a `relative` parameter means relative or not.
- For relative number, add `N` to `cursor_row`.
- For absolute number, assign `N` to `cursor_row`.
- Check the `cursor_row` and `cursor_col` to make sure it comes within the window size.
- Update the `combining_char_col` and `combining_char_row` for combining character.
- Set `next_print_will_wrap` false.

#### DrawState.move_col

- `move_col()` aka `DrawState::move_col()`
- `move_row()` has a `N` parameter, which is the column number, a `relative` parameter means relative or not.
- If `implicit` is true, update the `combining_char_col` and `combining_char_row` for combining character.
- For relative number, add `N` to `cursor_col`.
- For absolute number, assign `N` to `cursor_col`.
- If `implicit` is true, `next_print_will_wrap` dependens on whether it reaches the window size.
- Check the `cursor_row` and `cursor_col` to make sure it comes within the window size.
- If `implicit` is false,
  - Update the `combining_char_col` and `combining_char_row` for combining character.
  - Set `next_print_will_wrap` false.

#### ConditionalOverlayCell.apply

- If the cell `row` and `col` exceeds the `fb.ds` size and not `active`, returns.
- If it's `tentative`: `tentative_until_epoch > confirmed_epoch`, returns.
- If the `replacement` and the cell from frame buffer is blank, set `flag` false.
- In case `unknown`,
  - If `flag` is true and `col` is not at the right most edge: `col != fb.ds.get_width() - 1`,
    - Set the cell in frame buffer with underline `Renditions` attribute.
  - Returns early.
- In case `replacement` is different from the cell in frame buffer.
  - Replace the cell in frame buffer with `replacement`.
  - If `flag` is true, set the cell in frame buffer with underline `Renditions` attribute.

#### PredictionEngine.cull

- `flagging` : whether we are underlining predictions.
- `srtt_trigger` : show predictions because of slow round trip time.
- `glitch_trigger` : show predictions temporarily because of long-pending prediction.

The above is the meaning of key variables. The result of `cull()` is to set up `srtt_trigger`, `glitch_trigger`, `flagging` and verify the validity of prediction.

- `cull()` aka `PredictionEngine::cull()`.
- `cull()` accepts frame buffer as parameter.
- Return early if `display_preference == Never`.
- [Reset engine](#predictionenginereset) if `last_height` or `last_width` is different from `fb.ds`.
- [Control `srtt_trigger`](#how-to-control-srtt_trigger) with hysteresis: `send_interval > SRTT_TRIGGER_HIGH`.
- Control underlining with hysteresis: `send_interval > FLAG_TRIGGER_HIGH`.
- Really big glitches also activate underlining: `glitch_trigger > GLITCH_REPAIR_COUNT`.
- Iterate through each row in `overlays`, which is type of `list<ConditionalOverlayRow>`.
  - Erase current row if `i->row_num < 0` or `i->row_num >= fb.ds.get_height()`, continue the next row.
  - Iterate through each cell in the row.
  - [Check cell validity](#conditionaloverlaycellget_validity) via calling `get_validity()`,
  - In case validity returns `IncorrectOrExpired`,
    - If cell tentative time is greater than engine's `confirmed_epoch`,
      - [Reset the cell](#conditionaloverlaycellreset), if `display_preference == Experimental`, otherwise [`kill_epoch()`](#predictionenginekill_epoch).
      - Continue the iteration.
    - If not, [Reset the cell](#conditionaloverlaycellreset), if `display_preference == Experimental`, otherwise [reset engine](#predictionenginereset) and return.
  - In case validity returns `Correct`,
    - If cell tentative time is greater than engine's `confirmed_epoch`, update `confirmed_epoch`.
    - When predictions come in quickly, slowly take away the glitch trigger.
    - Match rest of row to the actual renditions.
  - In case validity returns `CorrectNoCredit`, [reset the cell](#conditionaloverlaycellreset).
  - In case validity returns `Pending`, When a prediction takes a long time to be confirmed,
    - we activate the predictions even if `SRTT` is low.
  - Continue the next row.
- If `cursors` is not empty and the [cursor validity](#conditionalcursormoveget_validity) is `IncorrectOrExpired`,
  - Clear `cursors` list if `display_preference == Experimental`,
  - otherwise [reset engine](#predictionenginereset) and return.
- Iterate through the `cursors` list, if cursor validity is `Pending`, erases it from the `cursors` list.

#### How to control `srtt_trigger`

- If `send_interval > FLAG_TRIGGER_HIGH` set `srtt_trigger` is true.
- If `srtt_trigger` is true and `send_interval <= SRTT_TRIGGER_LOW` and `active()` returns false.
  - `active()` aka `PredictionEngine::active()`.
  - `active()` iterates through each cell in `overlays`.
  - `active()` true if any cell is active.
- If so, sets `srtt_trigger` false.

#### PredictionEngine::reset

- Clear the `cursors` list.
- Clear the `overlay` list.
- Increase the `prediction_epoch` to become tentative.

#### PredictionEngine::kill_epoch

- `kill_epoch` has a `epoch` parameter and a `fb` parameter.
- Remove item from `cursors` if the item's `tentative_until_epoch > epoch-1`.
- Add a new `ConditionalCursorMove` to `cursors`.
- Iterate through each cell in `overlays`,
  - If the cell's `tentative_until_epoch > epoch-1`, [reset the cell](#conditionaloverlaycellreset).
- Increase the `prediction_epoch` to become tentative.

#### ConditionalOverlayCell::reset

- Set the `unknown` false.
- Clear the `original_contents` list.
- Set both `expiration_frame` and `tentative_until_epoch` to -1.
- Set the `active` false.

#### ConditionalCursorMove.get_validity

- `get_validity()` aka `ConditionalCursorMove::get_validity()`
- If the cell is not active, returns `Inactive`.
- If the cell `row` and `col` exceeds the `fb.ds` size, returns `IncorrectOrExpired`.
- Here, `row` and `col` are the parameters.
- If `late_ack >= expiration_frame`
  - If cursor `row` and `col` is the same as `fb.ds`, return `Correct`.
  - If not, return `IncorrectOrExpired`.
- Return `Pending`.

#### ConditionalOverlayCell.get_validity

- `get_validity()` aka `ConditionalOverlayCell::get_validity`.
- If the cell is not active, returns `Inactive`.
- If the cell `row` and `col` exceed the `fb.ds` size, returns `IncorrectOrExpired`.
- Here, `row` is the parameter, `col` is a cell field.
- Gets the `current` cell from `fb.get_cell(row, col)`.
- If `late_ack >= expiration_frame`, we are going to see if it hasn't been updated yet.
  - `late_ack` is the parameter, `expiration_frame` is a cell field.
  - If `unknown` field is true, returns `CorrectNoCredit`.
  - If `replacement` field `is_blank()`, returns `CorrectNoCredit`.
  - If `current` cell `contents_match()` with the `replacement`,
    - Checks the `original_contents` contents matches with the `replacement`,
  - If so , returns `Correct`. Otherwise returns `CorrectNoCredit`.
  - If not, returns `IncorrectOrExpired`.
- If not, returns `Pending`.

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
- Builds difference string for output to terminal display.
- Has two `Framebuffer` parameters: `last` and `f`.
- Initializes a `FrameState`: `frame`, with the old `Framebuffer` `last`.
- Checks if the bell ring happened: if true, append escape sequence to `frame`.
- Checks if icon name or window title changed: if true, append escape sequence to `frame`.
- Checks if reverse video state changed: if true, append escape sequence to `frame`.
- Checks if window size changed:
  - If true, append escape sequence to `frame`. Reset cursor position in `frame` to (0,0).
  - If false, update the cursor position and rendition with `frame.last_frame.ds`
- Checks is cursor visibility initialized: if false, append escape sequence to `frame`.
- Copies the `framestate.last_frame.get_rows()` to new `rows`.
- Extends `rows` width if we've gotten a resize and new is wider than old.
- Adds `rows` hight if we've gotten a resize and new is taller than old.
- Checks if display moved up by a certain number of lines,
  - Calculate the scroll region.
  - Append escape sequence to `frame` fro scrolling.
  - Do the move in our local index: `rows`.
- For rest rows, [updates the display row by row](#displayput_row).
- Checks if cursor location changed, if true, append escape sequence to `frame`.
- Checks if cursor visibility changed, if true, append escape sequence to `frame`.
- Checks if renditions changed: if true, append escape sequence to `frame`.
- Checks if bracketed paste mode changed: if true, append escape sequence to `frame`.
- Checks if mouse reporting mode changed: if true, append escape sequence to `frame`.
- Checks if mouse focus mode changed: if true, append escape sequence to `frame`.
- Checks if mouse encoding mode changed: if true, append escape sequence to `frame`.
- Returns the final `frame` difference string for output.

#### Display::put_row

- Has the `FrameState`, new frame buffer, `frame_y` position, and the local index `rows` as parameters.
- If we're forced to write the first column because of wrap, go ahead and do so.
- If rows are the same object, we don't need to do anything at all.
- Iterate for every Cell,
  - Does cell need to be drawn? Skip all this.
  - Slurp up all the empty cells: just counting.
  - Clear or write cells within the row (not to end).
  - Now draw a character cell.
- Clear or write empty cells at EOL.

See [this post](https://www.lihaoyi.com/post/BuildyourownCommandLinewithANSIescapecodes.html) to understand more about escape sequence drawing.

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

- `process_user_input()` aka `STMClient::process_user_input()`.
- `process_user_input()` is called to get the user keystrokes from `STDIN_FILENO`.
- Calls `read()` system call to read the user keystrokes.
- Calls `set_local_frame_sent()` of prediction engine to save the last `send_states` number.
- Iterates through each input character:
  - Calls `new_user_byte()` to [predict the input character](#predictionenginenew_user_byte).
  - If `quit_sequence_started` is ready:
    - If current byte is ".", set message for notification engine and `start_shutdown()`, return true.
    - If current byte is "^Z", [suspend the mosh client](#how-does-it-suspend).
    - If current byte is `escape_pass_key`, [pushes it into current state](#how-to-save-the-user-input-to-state).
    - For other character, escape key followed by anything other than "." and "^" gets sent literally.
  - Checks whether current byte is `escape_key`: "Ctrl-^", set `quit_sequence_started` accordingly.
    - If true, set `lf_entered` to be false, set the message `escape_key_help` for notification engine.
  - If current byte is the `LF`, `CR` control character, set `lf_entered` accordingly.
  - If current byte is `FF` control character, set `repaint_requested` to be true.
  - For other character, [pushes it into current state](#how-to-save-the-user-input-to-state).
- The result of `process_user_input()` is that all the user keystrokes are saved in current state.

#### PredictionEngine::new_user_byte

- `new_user_byte()` has a `the_byte` parameter and a `local_framebuffer` parameter.
- Returns early if `display_preference == Never`.
- Sets `prediction_epoch = confirmed_epoch` if `display_preference == Experimental`,
- Calls `cull()` to [prepare for the prediction engine](#predictionenginecull).
- Translates application-mode cursor control function to ANSI cursor control sequence.
- Updates the `last_byte` with the value of `the_byte`.
- Initializes a new `actions`, which is type of `Parser::Actions`.
- Calls `parser.input()` to [parse character into actions](server.md#parse-unicode-character-to-action) with `the_byte` and `actions` as parameters.
- Iterates through each action in `actions`:
  - In case action is type of `Parser::Print`,
    - Calls [`init_cursor()`](#predictionengineinit_cursor) with frame buffer as parameter.
    - Extracts `ch` from `act->ch`.
    - If `ch` is backspace/delete: "ch == 0x7f", [process backspace character](#how-to-process-backspace-character).
    - If "ch<=0x20" or `wcwidth(ch) != 1`, the prediction becomes tentative.
    - For other character, [process printable character](#how-to-process-printable-character).
  - In case action is type of `Parser::Execute`,
    - If `the_byte` is `CR`, the prediction becomes tentative. [perform `CR` in frame buffer](#predictionenginenewline_carriage_return).
    - For other character, the prediction becomes tentative.
  - In case action is type of `Parser::Esc_Dispatch`,
    - The prediction becomes tentative.
  - In case action is type of `Parser::CSI_Dispatch`,
    - If `the_byte` is right arrow, move the cursor in frame buffer.
    - If `the_byte` is left arrow, move the cursor in frame buffer.
    - For other character, the prediction becomes tentative.
  - For other action type, deletes current action.
- `new_user_byte()` checks the validity of prediction, manipulates overlays to reflect the change.

#### How to process backspace character

- `new_user_byte()` continues to process backspace character (0x7f).
- [Finds or creates the specified row](#predictionengineget_or_make_row) from prediction engine.
- If the last cursor in `cursors` is `<=0`, ignores it.
- Decreases the last cursor by one.
- Sets expire frame number and expire time for `cursors`.
- Iterates through to the end of `the_row`, starting from current cursor column.
  - Replace current `Cell` with the next one.
- `unknown`: whether we have the replacement.
- `replacement`: means the original contents before the prediction.
- `active`: means this cell is a prediction.

![mosh-row.svg](img/mosh-row.svg)

#### How to process printable character

- `new_user_byte()` continues to process printable character.
- Makes sure `cursor().col` and `cursor().row` is in range.
- [Finds `the_row` from prediction engine](#predictionengineget_or_make_row) with specified row as parameter.
- If cursor position in the last column, the prediction becomes tentative.
- Reversely iterate through to the end of `cursor().col`, starting from the end of `the_row`.
  - Moves the current `Cell` to the right: do the insert.
- Clears the cursor position `Cell`, matches renditions of character to the left, inserts the character.
- Sets expire frame number and expire time for `cursors`.
- Do we need to wrap?
  - If don't need, just increases `cursor().col++`.
  - If do need, the prediction becomes tentative, [perform `CR` in frame buffer](#predictionenginenewline_carriage_return).

#### PredictionEngine::get_or_make_row

- `get_or_make_row()` has row number and the number of columns as paramters.
- `get_or_make_row()` finds the specified row from prediction engine or creates a new row.
- Tries to find the specified row in `overlays`. If found, returns it.
- Creates a new `ConditionalOverlayRow` with the row number as parameter.
- Initializes every cell in the above row.
- Pushes the new row into `overlays`.

#### PredictionEngine::newline_carriage_return

- `newline_carriage_return` has a frame buffer parameter.
- Calls [`init_cursor()`](#predictionengineinit_cursor) with frame buffer as parameter.
- Set the cursor column to 0.
- If the cursor row is at the bottom of screen: `cursor().row == fb.ds.get_height() - 1`.
  - Makes blank prediction for last row.
- If not, Add the cursor row.

#### PredictionEngine::init_cursor

- `init_cursor` has a frame buffer parameter.
- If `cursors` is empty,
  - Creates a `ConditionalCursorMove` with frame buffer cursor position as parameters.
  - Pushes the above object into `cursors` list.
  - Sets latest cursor in `cursors` list active.
- In latest cursor in `cursors` list, if the cursor's `tentative_until_epoch != prediction_epoch`,
  - Creates a `ConditionalCursorMove` with latest cursor position in `cursors` list as parameters.
  - Pushes the above object into `cursors` list.
  - Sets latest cursor in `cursors` list active.

#### How to save the user input to state

- `network->get_current_state()` is actually `TransportSender.get_current_state()`.
- The return value of `TransportSender.get_current_state()` is a `UserStream` object.
- `push_back()` creates `Parser::UserByte`.
- `push_back()` wraps `UserByte` with `UserEvent` and pushes it into `UserStream`.
- `UserStream` object contains two kinds of character: `Parser::UserByte` and `Parser::Resize`.

#### How does it suspend

- Close display via writing to `STDOUT_FILENO` with `display.close()`.
- Restores the `saved_termios` for `STDIN_FILENO`, via `tcsetattr()`.
- Prints "[mosh is suspended.]" message.
- Flushes the output,
- Send `SIGSTOP` signal.
- Waiting for resume.

#### How to process resize

- `process_resize()` gets the window size for `STDIN_FILENO` , via `ioctl()` and `TIOCGWINSZ` flag.
- `process_resize()` creates `Parser::Resize` with the window size.
- `process_resize()` pushes the above `Parser::Resize` into `network->get_current_state()`.
- `process_resize()` calls `overlays.get_prediction_engine().reset()` to tell prediction engine.
- `process_resize()` returns true.

#### How does the network tick

- `network->tick()` calls `sender.tick()` to send data or an ack if necessary.
- `sender.tick()` aka `TransportSender<MyState>::tick()`
- `sender.tick()` calls `calculate_timers()` to calculate next send and ack times.
  - `calculate_timers()` aka `TransportSender<MyState>::calculate_timers()`.
  - `calculate_timers()` calls [`update_assumed_receiver_state()`](#how-to-pick-the-reciver-state) to update assumed receiver state.
  - `calculate_timers()` calls [`rationalize_states()`](#how-to-rationalize-states) cut out common prefix of all states.
  - `calculate_timers()` calculate `next_send_time` and `next_ack_time`.
- `sender.tick()` makes sure it's time to send data to the receiver.
- `sender.tick()` compares `current_state` with `assumed_receiver_state` to calculate difference string
  - For `UserStream` state: see [How to calculate the diff for UserStream](#how-to-calculate-the-diff-for-userstream).
  - For `Complete` state: see [How to calculate the diff for Complete](#how-to-calculate-the-diff-for-complete).
- `sender.tick()` calls `attempt_prospective_resend_optimization()` to optimize diff.
- If `diff` is empty and if it's greater than the `next_ack_time`.
  - `sender.tick()` calls [`send_empty_ack()`](#how-to-send-empty-ack) to send empty ack.
- If `diff` is not empty and if it's greater than `next_send_time` or `next_ack_time`.
  - `sender.tick()` calls [`send_to_receiver()`](#how-to-send-to-receiver) to send diffs.

#### How to send empty ack

- `send_empty_ack()` aka `TransportSender<MyState>::send_empty_ack()`.
- `send_empty_ack()` gets the last state number from `sent_states.back()` and increases it one.
- `send_empty_ack()` calls `add_sent_state()` to push `current_state` into `sent_states`,
  - with the new state number and current time as parameters.
  - `add_sent_state()` limits the size of `sent_states` below 32.
  - `add_sent_state()` limits the size of `sent_states` by erasing state from middle of `send_states` queue.
- `send_empty_ack()` calls [`send_in_fragments()`](#how-to-send-data-in-fragments) to send the new state
  - with empty string as `diff` parameter.

#### How to send to receiver

- `send_to_receiver()` aka `TransportSender<MyState>::send_to_receiver()`.
- If `current_state` number is equal to `sent_states.back()` number,
- `send_to_receiver()` refreshes the `timestamp` field of the latest state in `sent_states`.
- If `current_state` number is not equal to `sent_states.back()` number, increase the state number.
- `send_to_receiver()` calls `add_sent_state()` to push `current_state` into `sent_states`,
  - with the new state number and current time as parameters.
  - `add_sent_state()` limits the size of `sent_states` below 32.
  - `add_sent_state()` limits the size of `sent_states` by erasing state from middle of `send_states` queue.
- Note `sent_states` is type of list `TimestampedState`, while `current_state` is type of `MyState`.
- `send_to_receiver()` calls [`send_in_fragments()`](#how-to-send-data-in-fragments) to send data.
- `send_to_receiver()` updates `assumed_receiver_state`, `next_ack_time` and `next_send_time`.

#### How to calculate the diff for UserStream

- `diff_from()` aka `UserStream::diff_from()`, which calculates diff based on user keystrokes.
- `diff_from()` has a existing `UserStream` as parameter.
- `diff_from()` compares current `UserStream` with existing `UserStream` to calculate the diff.
- `diff_from()` finds the position in the current `UserStream` which is different from existing `UserStream`.
- `diff_from()` iterates to the end of current `UserStream`, starting from the above position.
- `diff_from()` build `ClientBuffers::UserMessage`, with the `UserEvent` object in each iteration,
- `diff_from()` returns the serialized string representation of the `ClientBuffers::UserMessage` object.
- `ClientBuffers::UserMessage` is a proto2 message. See userinput.proto file.
- `ClientBuffers::UserMessage` contains several `ClientBuffers.Instruction`.
- `ClientBuffers.Instruction` is composed of `Keystroke` or `ResizeMessage`.
- Several `Keystroke` can be appended to one `ClientBuffers.Instruction`.
- `ResizeMessage` is added to one `ClientBuffers.Instruction`.

#### How to calculate the diff for Complete

- `diff_from()` aka `Complete::diff_from()`.
- `diff_from()` has a existing `Complete` as parameter.
- Compares current `echo_ack` with existing one.
  - If they are different, adds `EchoAck` intruction with current `echo_ack` as parameter.
- Compares current `Framebuffer` with existing one.
- If they are different, `diff_from()` compares the `Framebuffer` size.
  - If the `Framebuffer` size is different,
    - adds `ResizeMessage` instruction with current size as parameter.
  - Calls `display.new_frame()` to [calculate the `Framebuffer` diff](#how-to-calculate-frame-buffer-difference).
  - If `Framebuffer` diff is not empty, adds `HostBytes` instruction with the diff as parameter.
- `diff_from()` returns the serialized string representation of the `HostBuffers::HostMessage` obejct.
- `HostBuffers::HostMessage` is a proto2 message. See hostinput.proto file.
- `HostBuffers::HostMessage` contains several `HostBuffers::Instruction`.
- `HostBuffers::Instruction` is composed of `HostBytes`, `ResizeMessage`, `EchoAck`.

#### How to pick the reciver state

- `update_assumed_receiver_state()` chooses the most recent receiver state based on network traffic.
- `update_assumed_receiver_state()` picks the first item in `send_states`.
- `send_state` is type of `list<TimestampedState<MyState>>`.
- `assumed_receiver_state` point to the middle of `sent_states`.
- `send_state` skips the first item.
- For each item in `send_states`,
  - If the time gap for each state is lower than `connection->timeout() + ACK_DELAY`, `ACK_DELAY` is 100ms,
  - updates `assumed_receiver_state` to current item.
  - `connection->timeout()` aka `Connection::timeout()`.
  - `connection->timeout()` calcuates [RTO](https://datatracker.ietf.org/doc/html/rfc2988) based on `SRTT` and `RTTVAR`.
  - If the time gap for each state `now - i->timestamp` is greater than `connection->timeout() + ACK_DELAY`,
  - returns early.
- The result is saved in `assumed_receiver_state`.

#### How to rationalize states

- `rationalize_states()` aka `TransportSender<MyState>::rationalize_states()`.
- `rationalize_states()` picks the first state from `sent_states` as common prefix.
  - The comm prefix is the first state in `send_states`.
  - `sent_states` is type of `list<TimestampedState<MyStat>>`.
- `rationalize_states()` calls `current_state.subtract()` to cut out common prefix from `current_state`.
- `rationalize_states()` iterates through `send_states`.
  - For each state,
  - `rationalize_states()` calls `i->state.subtract()` to cut out common prefix from current state.
  - In case `MyState` is `UserStream`:
    - `subtract()` aka `UserStream::subtract()`.
    - `subtract()` cuts out any `UserEvent` from its `actions` deque, if it's the same `UserEvent` in `prefix`.
    - The result is the common prefix in `current_state` and early `sent_states` is cut out.
  - In case `MyState` is `Complete`:
    - `subtract()` aka `Complete:subtract()`.
    - `subtract()` does nothing.

#### How to send data in fragments

- `send_in_fragments()` aka `TransportSender<MyState>::send_in_fragments()`.
- `send_in_fragments()` creates `TransportBuffers.Instruction` with the `diff` created in previous step.
  - See [How to calculate the diff for UserStream](#how-to-calculate-the-diff-for-userstream).
  - See [How to calculate the diff for Complete](#how-to-calculate-the-diff-for-complete).
- `TransportBuffers.Instruction` contains the following fields.
  - `old_num` : is the source number. It's value is `assumed_receiver_state->num`.
  - `new_num` : is the target number. It's value is specified by `new_num` parameter.
  - `throwaway_num` : is the throwaway number. It's value is `sent_states.front().num`.
  - `diff` : contains the `diff`. It's value is specified by `diff` parameter.
  - `ack_num` : is the ack number. It's value is assigned by `ack_num`.
- `send_in_fragments()` calls `Fragmenter::make_fragments` to splits the `TransportBuffers.Instruction` into `Fragment`.
  - `make_fragments()` serializes `TransportBuffers.Instruction` into string and compresses it to string `payload`.
  - `make_fragments()` splits the `payload` string into fragments based on the size of `MTU`,
  - The default size of `MTU` is 1280.
- `Fragment` has the following fields:
  - `id` : which is the instruction id. It's the same id for all the fragment.
  - `fragment_num` : which starts from zero, and is increased one for each new fragment.
  - `final` : which is used to indicate the last fragment.
  - `contents` : which contains part of the instruction.
  - The fragments is saved in `Fragment` vector.
- `send_in_fragments()` calls [`connection->send()`](#how-to-send-a-packet) to send each `Fragment` to the server.

#### How to send a packet?

- `connection->send()` aka `Connection::send()`.
- `connection->send()` calls `new_packet()` to create a `Packet`.
  - `Packet` is type of `Network::Packet`.
  - Besides the `payload` field,
  - A `Packet` also contains a unique `seq` field, a `timestamp` field and a `timestamp_reply` field.
  - See [more about `Packet`](#networkfragment---networkpacket---cryptomessage)
- `connection->send()` calls `session.encrypt()` to encrypt the `Packet`.
- `connection->send()` calls `sendto()` system call to send the encrypted data to receiver.
  - `sendto()` use the last socket from socket list to send the encrypted data.
- `connection->send()` checks the time gap between now and `last_port_choice`, `last_roundtrip_success`.
- For client, [`hop_port()`](#how-does-the-client-roam) is called to roam the client,
  - if the time gap is greater than `PORT_HOP_INTERVAL`.
- For server, `last_heard` is checked to make sure after 40 seconds, `has_remote_addr` is false.
  - `has_remote_addr` is false meaning no more send.

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
- `network->recv()` calls `connection.recv()` to [receive payload string](#how-to-read-data-from-socket).
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
- `network->recv()` initializes a `RemoteState` and wraps it in `TimestampedState<RemoteState>`.
- If `Instruction` diff field is not empty, `network->recv()` calls `new_state.state.apply_string()`.
  - [`apply_string()`](server.md#apply_string) is called with `Instruction` diff field as parameter.
  - [`apply_string()`](server.md#apply_string) initializes `RemoteState` with remote data.
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
  - The `fragment_num` field, which is `uint16_t`, contains the fragment number and fragment `final` flag.
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
- `recv_one()` checks congestion flag, if so, set `congestion_experienced` to true.
<!-- congestion flag at Unix Network Programming: The Socket Networking API? P390, P588, check Connection::Socket::Socket()-->
- `recv_one()` calls `session.decrypt()` to decrypt the received data and transform it into `Message`.
- `recv_one()` creates a `Packet` object based on the `Message`.
- `recv_one()` checks `Packet`'s sequence number to make sure it is greater than the `expected_receiver_seq`.
  - if packet sequence number is greater than `expected_receiver_seq`,
  - `recv_one()` increases `expected_receiver_seq`.
  - `recv_one()` saves the `p.timestamp` in `saved_timestamp`, saves current time in `saved_timestamp_received_at`.
  - `recv_one()` signals counterparty to slow down via decrease `saved_timestamp`, if congestion is detected.
  - `recv_one()` calculates `SRTT` and `RTTVAR` based on each [RTT](https://datatracker.ietf.org/doc/html/rfc29880).
  - `recv_one()` updates `last_heard` with current time.
  - For server side, [client roaming](#how-does-the-server-support-client-roam) is supported here.
  - if packet sequence number is less than `expected_receiver_seq`,
  - `recv_one()` return out-of-order or duplicated packets to caller.
- `recv_one()` return the `payload` to caller.

#### How does the server support client roam

- `recv_one()` compares `packet_remote_addr` with `remote_addr`.
- If the packet remote address is different than remote address, update the `remote_addr` and `remote_addr_len`.
- `recv_one()` calls `getnameinfo()` to validate the new remote address.

### How to send keystroke to remote server

This section summarizes the data structures in the process.

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

### How to receive state from server

Upon network sockets is ready to read, `main()` calls `process_network_input()` to process it.

#### bytes -> `TransportBuffers:Instruction`

- `process_network_input()` calls `network->recv()` to do the job.
- `network->recv()` calls `connection.recv()` to read string representation from socket. See [here](server.md#cryptomessage---networkpacket---networkfragment) for detail.
- `network->recv()` calls `fragments.get_assembly()` to build the `TransportBuffers:Instruction` object.

#### `TransportBuffers:Instruction` -> `HostBuffers::HostMessage` -> terminal

- `network->recv()` finds the reference state (old state) from received states.
- `network->recv()` [apply `Instruciton` difference](server.md#apply_string) to the reference state to generate new state.
- `HostBuffers::HostMessage` contains [difference escape sequences](server.md#terminal---difference-escape-sequence---hostbuffershostmessage) for frame buffer.
- The escape sequences received from pty master is the instructions send from remote application
- The escape sequences received from server is the result of frame buffer after applying the above instructions.
- `network->recv()` add the new state to received states.

### How does the notification engine decide to show message?

- Every time `STMClient::output_new_frame()` is called to show the frame buffer.
- `overlays.apply()` is called to apply local overlay.
- `overlays.apply()` calls `notifications.adjust_message()` to clear expired message.
- `overlays.apply()` calls `notifications.apply()` to [show the message](#notificationengineapply).
- Before `notifications.apply()` decide to show the message,
- `notifications.apply()` checks the last word from server time and last ack state time, to decide whether to show the message.
- If the message is empty, `notifications.apply()` will not show the message either.

### How does the prediction engine work?

#### user keystroke -> overlays

- Upon user keystroke is ready to read, `process_user_input()` is called to process it.
- `process_user_input()` reads the user keystroke from `STDIN_FILENO`.
- `process_user_input()` calls `overlays.get_prediction_engine().new_user_byte()` for each byte.
- `PredictionEngine::new_user_byte()` [checks the validity of the prediction](#predictionenginenew_user_byte) and manipulates `overlays` to reflect the change.
- Now, all the change happens in `overlays` or `cursors`.
- At the same time, the user keystrokes is send to the server through network socket.

#### overlays -> show up

- When it's time to show the frame buffer for diasplay,`output_new_frame()` is called to perfrm the display.
- `output_new_frame()` calls `overlays.apply()` to apply local overlay.
- `overlays.apply()` calls `predictions.cull()` to [check the validity of the prediction](#predictionenginecull).
- `overlays.apply()` calls `predictions.apply()` to [manipulate the frame buffer](#predictionengineapply).
- Here if condition is ready, e.g. `flagging` is true, cell is `active`… underline prediction show up.
- If the `overlays` is incorrect, the contents of `overlays` will be reset by `predictions.cull()`.
