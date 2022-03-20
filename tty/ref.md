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
