package osc52

import (
	"encoding/base64"
	"io"
	"os"
	"strings"
)

var output = NewOutput(os.Stdout, os.Environ())

type envs map[string]string

func (e envs) Get(key string) string {
	v, ok := e[key]
	if !ok {
		return ""
	}
	return v
}

type Output struct {
	out  io.Writer
	envs envs
}

func NewOutput(out io.Writer, envs []string) *Output {
	e := make(map[string]string, 0)
	for _, env := range envs {
		s := strings.Split(env, "=")
		k := s[0]
		v := strings.Join(s[1:], "=")
		e[k] = v
	}
	o := &Output{
		out:  out,
		envs: e,
	}
	return o
}

func Copy(str string) {
	output.Copy(str)
}

func (o *Output) Copy(str string) {
	mode := "default"
	term := o.envs.Get("TERM")
	lcterm := o.envs.Get("LC_TERM")
	switch {
	case o.envs.Get("TMUX") != "":
		mode = "tmux"
	case strings.HasPrefix(term, "screen"), strings.HasPrefix(lcterm, "screen"):
		mode = "screen"
	case strings.HasPrefix(term, "kitty"), strings.HasPrefix(lcterm, "kitty"):
		mode = "kitty"
	}

	switch mode {
	case "default":
		o.copyDefault(str)
	case "tmux":
		o.copyTmux(str)
	case "screen":
		o.copyDCS(str)
	case "kitty":
		o.copyKitty(str)
	}
}

func (o *Output) copyDefault(str string) {
	b64 := base64.StdEncoding.EncodeToString([]byte(str))
	o.out.Write([]byte("\x1b]52;c;" + b64 + "\x07"))
}

func (o *Output) copyTmux(str string) {
	b64 := base64.StdEncoding.EncodeToString([]byte(str))
	o.out.Write([]byte("\x1bPtmux;\x1b\x1b]52;c;" + b64 + "\x07\x1b\\"))
}

func (o *Output) copyDCS(str string) {
	// " This function base64's the entire source, wraps it in a single OSC52, and then
	// " breaks the result into small chunks which are each wrapped in a DCS sequence.
	// " This is appropriate when running on `screen`. Screen doesn't support OSC52,
	// " but will pass the contents of a DCS sequence to the outer terminal unchanged.
	// " It imposes a small max length to DCS sequences, so we send in chunks.
	// let b64 = s:b64encode(a:str, 76)
	// " Remove the trailing newline.
	// let b64 = substitute(b64, '\n*$', '', '')
	// " Replace each newline with an <end-dcs><start-dcs> pair.
	// let b64 = substitute(b64, '\n', "\e/\eP", "g")
	// " (except end-of-dcs is "ESC \", begin is "ESC P", and I can't figure out
	// " how to express "ESC \ ESC P" in a single string. So the first substitute
	// " uses "ESC / ESC P" and the second one swaps out the "/". It seems like
	// " there should be a better way.)
	// let b64 = substitute(b64, '/', '\', 'g')
	// " Now wrap the whole thing in <start-dcs><start-osc52>...<end-osc52><end-dcs>.
	// return "\eP\e]52;c;" . b64 . "\x07\e\x5c"
	b64 := base64.StdEncoding.EncodeToString([]byte(str))
	s := strings.SplitN(b64, "\n", 76)
	q := "\x1bP\x1b]52;c;"
	for _, v := range s {
		q += "\x1b\\\x1bP" + v
	}
	q += "\x07\x1b\x5c"
	o.out.Write([]byte(q))
}

func (o *Output) copyKitty(str string) {
	o.out.Write([]byte("\x1b]52;c;!\x07"))
	o.copyDefault(str)
}
