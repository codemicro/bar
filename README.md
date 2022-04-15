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

These instructions assume you have a recent version of the Go 1.18 (or newer) toolchain installed on your machine. `cdmbar` includes some stuff that works with the VCS stamping introduced in Go 1.18, however you can compile with `-buildvcs=false` and everything should still work fine.

```
git clone https://github.com/codemicro/bar.git
cd bar
go build -o cdmbar github.com/codemicro/bar/cmd/bar
// To use with i3, you probably want to put it somewhere that's on PATH
mv ./cdmbar ~/.local/bin
```

### Using with i3wm

Update the `status_command` line in your i3wm config file to match the below.

```
bar {
        status_command cdmbar
}
```

### Changing options

Edit the `blocks` variable inside of `cmd/bar/main.go`.
