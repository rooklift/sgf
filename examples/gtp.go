package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	sgf ".."
)

const config_filename = "gtp_config.json"

type ConfigStruct struct {
	Engine1Path			string		`json:"engine_1_path"`
	Engine1Args			[]string	`json:"engine_1_args"`
	Engine2Path			string		`json:"engine_2_path"`
	Engine2Args			[]string	`json:"engine_2_args"`
}

var Config ConfigStruct

func init() {
	file, err := ioutil.ReadFile(config_filename)
	if err != nil {
		panic("Couldn't load config file " + config_filename)
	}

	err = json.Unmarshal(file, &Config)
	if err != nil {
		panic(err)
	}
}

type Engine struct {
	stdin	io.WriteCloser
	stdout	*bufio.Scanner
	stderr	*bufio.Scanner

	dir		string
	base	string
	args	[]string		// Not including base
}

func (self *Engine) Start(path string, args []string) {

	self.dir = filepath.Dir(path)
	self.base = filepath.Base(path)

	for _, a := range args {
		self.args = append(self.args, a)
	}

	var cmd exec.Cmd

	cmd.Dir = self.dir
	cmd.Path = self.base
	cmd.Args = append([]string{self.base}, self.args...)

	var err1 error
	self.stdin, err1 = cmd.StdinPipe()

	stdout_pipe, err2 := cmd.StdoutPipe()
	self.stdout = bufio.NewScanner(stdout_pipe)

	stderr_pipe, err3 := cmd.StderrPipe()
	self.stderr = bufio.NewScanner(stderr_pipe)

	err4 := cmd.Start()

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		panic(fmt.Sprintf("\nerr1: %v\nerr2: %v\nerr3: %v\nerr4: %v\n", err1, err2, err3, err4))
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

	engine1 := new(Engine)
	engine2 := new(Engine)
	engine1.Start(Config.Engine1Path, Config.Engine1Args)
	engine2.Start(Config.Engine2Path, Config.Engine2Args)

	player_map := map[sgf.Colour]*Engine{sgf.BLACK: engine1, sgf.WHITE: engine2}

	// -------------------------------------------------------------

	root := sgf.NewTree(19)
	root.SetValue("KM", "7.5")

	for _, engine := range player_map {
		engine.SendAndReceive("boardsize 19")
		engine.SendAndReceive("komi 7.5")
		engine.SendAndReceive("clear_board")
	}

	last_save_time := time.Now()
	node := root
	colour := sgf.WHITE

	passes_in_a_row := 0

	outfilename := time.Now().Format("2006-01-02-15-04-05") + ".sgf"

	for {
		colour = colour.Opposite()

		if time.Now().Sub(last_save_time) > 5 * time.Second {
			node.Save(outfilename)
			last_save_time = time.Now()
		}

		engine := player_map[colour]
		response := engine.SendAndReceive(fmt.Sprintf("genmove %s", colour.Lower()))

		var move string
		fmt.Sscanf(response, "= %s", &move)

		var err error

		if move == "resign" {
			root.SetValue("RE", fmt.Sprintf("%s+R", colour.Opposite().Upper()))
			break
		} else if move == "pass" {
			passes_in_a_row++
			node = node.PassColour(colour)
			if passes_in_a_row >= 3 {
				break
			}
		} else {
			passes_in_a_row = 0
			node, err = node.PlayMoveColour(move_to_sgf(move, 19), colour)
			if err != nil {
				fmt.Printf("%v\n", err)
				break
			}
		}

		// Must only get here with a valid move...

		other_engine := player_map[colour.Opposite()]
		other_engine.SendAndReceive(fmt.Sprintf("play %s %s", colour.Lower(), move))

		node.Board().Dump()
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
