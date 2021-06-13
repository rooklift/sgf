package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rooklift/sgf"
)

type Record struct {
	Path		string		`json:"path"`
	Filename	string		`json:"filename"`
	Dyer		string		`json:"dyer"`
	SZ			int			`json:"SZ"`
	HA			int			`json:"HA"`
	PB			string		`json:"PB"`
	PW			string		`json:"PW"`
	BR			string		`json:"BR"`
	WR			string		`json:"WR"`
	RE			string		`json:"RE"`
	DT			string		`json:"DT"`
	EV			string		`json:"EV"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <dir>\n", filepath.Base(os.Args[0]))
		return
	}
	filepath.Walk(os.Args[1], handle_file)
}

func handle_file(path string, _ os.FileInfo, err error) error {

	// Returning an error halts the whole walk. So don't.

	if err != nil {
		return nil
	}

	root, err := sgf.LoadMainLine(path)
	if err != nil {
		return nil
	}

	record := new(Record)

	record.Path = strings.ReplaceAll(filepath.Dir(path), "\\", "/")
	record.Filename = filepath.Base(path)
	record.Dyer = root.Dyer()

	record.SZ = root.RootBoardSize()
	record.HA = root.RootHandicap()

	record.PB, _ = root.GetValue("PB")
	record.PW, _ = root.GetValue("PW")
	record.BR, _ = root.GetValue("BR")
	record.WR, _ = root.GetValue("WR")
	record.RE, _ = root.GetValue("RE")
	record.DT, _ = root.GetValue("DT")
	record.EV, _ = root.GetValue("EV")

	foo, _ := json.Marshal(record)

	fmt.Printf("%v\n", string(foo))

	return nil
}
