
# go-osc52

<p>
    <a href="https://github.com/aymanbagabas/go-osc52/releases"><img src="https://img.shields.io/github/release/aymanbagabas/go-osc52.svg" alt="Latest Release"></a>
    <a href="https://pkg.go.dev/github.com/aymanbagabas/go-osc52?tab=doc"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="GoDoc"></a>
</p>

A Go library to work with the [ANSI OSC52](https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-Operating-System-Commands) terminal sequence.

## Example

```go
str := "Hello World!"
osc52.Copy(str) // Copies str to system clipboard
osc52.CopyPrimary(str) // Copies str to primary clipboard (X11 only)
```

## SSH Example

You can use this over SSH using [gliderlabs/ssh](https://github.com/gliderlabs/ssh) for instance:

```go
envs := sshSession.Environ()
pty, _, _ := s.Pty()
envs = append(envs, "TERM="+pty.Term)
output := NewOutput(sshSession, envs)
// Copy text in your application
output.Copy("Hello awesome!")
```

### Tmux users

If you're using tmux, make sure you set `set -g default-terminal` in your
`~/.tmux.conf` file, to a value that starts with `tmux-`. `tmux-256color` for
instance. See [this](https://github.com/tmux/tmux/wiki/FAQ#why-do-you-use-the-screen-terminal-description-inside-tmux)
for more details.

You can instead pass the `TMUX` environment variable:

```sh
ssh -o SendEnv=TMUX <host>
```

## Credits

* [vim-oscyank](https://github.com/ojroques/vim-oscyank) this is heavily inspired by vim-oscyank.