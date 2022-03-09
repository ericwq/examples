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

### STMClient::init

`Overlay::OverlayManager overlays` is initialized in `STMClient` construction function.

`Terminal::Display display` is initialized in `STMClient` construction function.

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
- Get the window size for `stdin` , via `ioctl()` and `TIOCGWINSZ`.

## reference

How the terminal works? Who is responsible for terminal rendering? Does GPU-rendering in terminal matter?

- [Linux terminals, tty, pty and shell](https://dev.to/napicella/linux-terminals-tty-pty-and-shell-192e)
- [Linux terminals, tty, pty and shell - part 2](https://dev.to/napicella/linux-terminals-tty-pty-and-shell-part-2-2cb2)
- [How does a Linux terminal work?](https://unix.stackexchange.com/questions/79334/how-does-a-linux-terminal-work)
- [How Zutty works: Rendering a terminal with an OpenGL Compute Shader](https://tomscii.sig7.se/2020/11/How-Zutty-works)
- [A totally biased comparison of Zutty (to some better-known X terminal emulators)](https://tomscii.sig7.se/2020/12/A-totally-biased-comparison-of-Zutty)
- [A look at terminal emulators, part 1](https://lwn.net/Articles/749992/)
- [A look at terminal emulators, part 2](https://lwn.net/Articles/751763/)
- [High performant 2D renderer in a terminal](https://blog.ghaiklor.com/2020/07/27/high-performant-2d-renderer-in-a-terminal/)
- [The TTY demystified](http://www.linusakesson.net/programming/tty/)

### typing

- [Typing with pleasure](https://pavelfatin.com/typing-with-pleasure/)
- [Measured: Typing latency of Zutty (compared to others)](https://tomscii.sig7.se/2021/01/Typing-latency-of-Zutty)

## clangd format

- [Clang-Format Style Options](https://clang.llvm.org/docs/ClangFormatStyleOptions.html)
- [clangd format generator](https://zed0.co.uk/clang-format-configurator/)

### reference links

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
