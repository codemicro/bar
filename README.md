# cdmbar for i3wm

*i3bar compatible status bar*

---

This is a i3wm status bar that's built as a toy project. It's pretty basic, doesn't have many features and isn't very configurable unless you want to edit the source and recompile it.

This interacts with i3 using the [i3bar input protocol](https://i3wm.org/docs/i3bar-protocol.html).

### Features

* SIGUSR1 forces a refresh
* It has colours
* Sometimes it breaks

### Compiling locally

These instructions assume you have a recent version of the Go 1.18 (or newer) toolchain installed on your machine and a copy of GNU Make.

`cdmbar` includes some stuff that works with the VCS stamping introduced in Go 1.18, however you can compile with `-buildvcs=false` and everything should still work fine.

```
git clone https://github.com/codemicro/bar.git
cd bar
// Will install cdmbar to ~/.local/bin - to use it with i3, we need to put it on PATH
make clean install

// You can customise the install directory using `make install_dir=/usr/local/bin clean install`
// You can build without VCS stamping using `make go_args="-buildvcs=false" clean install`
```

### Using with i3wm

Update the `status_command` line in your i3wm config file to match the below.

```
bar {
        status_command cdmbar
}
```

### Changing options

Edit the `blocks` variable inside of `cmd/bar/main.go`, then recompile.
