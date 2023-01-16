# cdmbar for i3wm

*i3bar compatible status bar*

---

This is a i3wm status bar that's built as a toy project. It started off pretty basic (I've since added click event support and support for partial updates), doesn't have many features and isn't very configurable unless you want to edit the source and recompile it.

This interacts with i3 using the [i3bar input protocol](https://i3wm.org/docs/i3bar-protocol.html). i3 versions earlier than v4.3 are not supported.

### Features

* Supports click events
* Supports partial refreshes
* SIGUSR1 forces a refresh
* It has colours
* Sometimes it breaks

### Included providers

* `AudioPlayer` - show the currently playing song
* `Battery` - show the current battery charge status and provide alerts if it leaves set boundaries
* `CPU` - show CPU load and provide alerts if it leaves set boundaries
* `DateTime` - show the current date and time
* `Disk` - show the current usage of a disk
* `IPAddress` - show the current local IPv4 address
* `Memory` - show the current memory usage and provide alerts it if leaves set boundaries
* `PlainText`
* `PulseaudioVolume` - show the current volume of a PulseAudio sink and control that using the scroll wheel
* `Timer` - provides a small timer that play/pauses with a left-click and resets with a right-click.
* `WiFi` - show the curent WiFi SSID, connection frequency and connection strength

### Compiling locally

These instructions assume you have a recent version of the Go 1.18 (or newer) toolchain installed on your machine and a copy of GNU Make.

`cdmbar` includes some stuff that works with the VCS stamping introduced in Go 1.18, however you can compile with `-buildvcs=false` and everything should still work fine without Git installed.

```
// Alternatively, download the source as a ZIP file
git clone https://github.com/codemicro/bar.git
cd bar

// Will install cdmbar to ~/.local/bin - to use it with i3, we need to put it on PATH
make clean install

// You can customise the install directory using `make install_dir=/usr/local/bin clean install`
// You can build without VCS stamping using `make build_args="-buildvcs=false" clean install`
```

### Using with i3wm

Update the `status_command` line in your i3wm config file to match the below.

```
bar {
        status_command cdmbar
}
```

### Changing options

Edit the arguments of the call to `b.RegisterBlockGenerator` inside of `cmd/bar/main.go`, then recompile.
