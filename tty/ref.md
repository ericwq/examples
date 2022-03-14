# Terminal emulator

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
- Pet the terminal for `STDOUT` in application cursor key mode. `Display::open()`.
  - In Application mode, the cursor keys generate escape sequences that the application uses for its own purpose.
  - Application Cursor Keys mode is a way for the server to change the control sequences sent by the arrow keys. In normal mode, the arrow keys send `ESC [A` through to `ESC [D`. In application mode, they send `ESC OA` through to `ESC OD`.
- Set terminal window title, via `overlays.set_title_prefix()`. ?
- Set escape key string, via `overlays.get_notification_engine().set_escape_key_string()`. ?

### STMClient::main_init

In `client.main()`, `main_init()` is called to init the `mosh` client.

- Register signal handler for `SIGWINCH`, `SIGTERM`, `SIGINT`, `SIGHUP`, `SIGPIPE`, `SIGCONT`.
- Get the window size for `STDIN` , via `ioctl()` and `TIOCGWINSZ`.
- Create `local_framebuffer` and set the `window_size`.?
- Create `new_state` frame buffer and set the size to `1*1`.
- initialize screen via `display.new_frame()`. Write screen to `STDOUT` via `swrite()`?
- Create the `Network::UserStream`, create the `Terminal::Complete local_terminal` with window size.?
- Open the network via `Network::Transport<Network::UserStream, Terminal::Complete>`.?
- Set minial delay on outgoing keystrokes via `network->set_send_delay(1)`.?
- Tell server the size of the terminal via `network->get_current_state().push_back()`.?
- Set the `verbose` mode via `network->set_verbose()`.

### Transport<MyState, RemoteState>::Transport

- Initialize the `connection`, initialize the sender with `connection` and `initial_state`.
  - `connection` is the underlying, encrypted network connection.

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
