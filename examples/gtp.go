package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	sgf ".."
)

var gtp_names = map[sgf.Colour]string {
	sgf.EMPTY: "??",
	sgf.BLACK: "b",
	sgf.WHITE: "w",
}

type Engine struct {
	stdin	io.WriteCloser
	stdout	*bufio.Scanner
	stderr	*bufio.Scanner
}

func (self *Engine) Start(path string, args ...string) {

	var cmd exec.Cmd

	cmd.Dir = filepath.Dir(path)
	cmd.Path = filepath.Base(path)
	cmd.Args = append([]string{cmd.Path}, args...)

	var err1 error
	self.stdin, err1 = cmd.StdinPipe()

	stdout_pipe, err2 := cmd.StdoutPipe()
	self.stdout = bufio.NewScanner(stdout_pipe)

	stderr_pipe, err3 := cmd.StderrPipe()
	self.stderr = bufio.NewScanner(stderr_pipe)

	err4 := cmd.Start()

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		panic(fmt.Sprintf("%v\n%v\n%v\n%v\n", err1, err2, err3, err4))
	}

	go self.ConsumeStderr()
}

func (self *Engine) ConsumeStderr() {
	for self.stderr.Scan() {
		// fmt.Printf("%s\n", self.stderr.Text())
	}
}

func (self *Engine) SendAndReceive(msg string) string {

	// FIXME: detect dead engine

	msg = strings.TrimSpace(msg)
	fmt.Fprintf(self.stdin, "%s\n", msg)

	var response bytes.Buffer
	for self.stdout.Scan() {
		response.WriteString(self.stdout.Text())
		response.WriteString("\n")
		if self.stdout.Text() == "" {
			return response.String()
		}
	}

	// If we get to here, Scan() returned false, likely meaning the engine is dead.
	// We should do something.

	return ""
}

func main() {

	engine := new(Engine)
	engine.Start("../../../../../Programs (self-installed)/Leela Zero/leelaz.exe",
					"--gtp", "--noponder", "-p", "25", "-w", "networks/better_192_163e407b")

	root := sgf.NewTree(19)
	root.SetValue("KM", "7.5")

	outfilename := "foo.sgf"

	engine.SendAndReceive("boardsize 19")
	engine.SendAndReceive("komi 7.5")
	engine.SendAndReceive("clear_board")

	last_save_time := time.Now()
	node := root
	colour := sgf.WHITE

	passes_in_a_row := 0

	for {
		colour = colour.Opposite()

		if time.Now().Sub(last_save_time) > 5 * time.Second {
			node.Save(outfilename)
			last_save_time = time.Now()
		}

		response := engine.SendAndReceive(fmt.Sprintf("genmove %s", gtp_names[colour]))

		var move string
		fmt.Sscanf(response, "= %s", &move)

		if move == "pass" {
			node = node.PassColour(colour)
			passes_in_a_row++
			if passes_in_a_row >= 3 {
				break					// We have to tell the score somehow.
			}
			continue
		} else {
			passes_in_a_row = 0
		}

		if move == "resign" {
			s := fmt.Sprintf("%s+R", colour.Opposite().Upper())
			root.SetValue("RE", s)
			break
		}

		sgf := move_to_sgf(move, 19)

		var err error
		node, err = node.PlayMoveColour(sgf, colour)
		if err != nil {
			fmt.Printf("%v\n", err)
			break
		} else {
			node.Board().Dump()
		}
	}

	node.Save(outfilename)
}

func move_to_sgf(s string, size int) string {

	if len(s) < 2 || len(s) > 3 {
		return ""
	}

	if s[0] < 'A' || s[0] > 'Z' {
		return ""
	}

	if s[1] < '0' || s[1] > '9' {
		return ""
	}

	if len(s) == 3 && (s[2] < '0' || s[2] > '9') {
		return ""
	}

	x := int(s[0]) - 65
	if x >= 8 {
		x--
	}

	y_int, _ := strconv.Atoi(s[1:])
	y := size - int(y_int)

	return sgf.Point(x, y)
}
