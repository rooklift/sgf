package sgf

// Basic NGF parser (i.e. for WBaduk files)

import (
	"fmt"
	"strconv"
	"strings"
)

func load_ngf(ngf string) (*Node, error) {

	ngf = strings.TrimSpace(ngf)

	lines := strings.Split(ngf, "\n")

	var boardsize, handicap int
	var pw, pb, rawdate, re string
	var komi float64

	if len(lines) < 12 {
		return nil, fmt.Errorf("load_ngf(): file too short")
	}

	// ------------------------------------

	boardsize, _ = strconv.Atoi(strings.TrimSpace(lines[1]))

	// ------------------------------------

	pw_fields := strings.Fields(lines[2])
	pb_fields := strings.Fields(lines[3])

	if len(pw_fields) > 0 {
		pw = pw_fields[0]
	}

	if len(pb_fields) > 0 {
		pb = pb_fields[0]
	}

	// ------------------------------------

	handicap, _ = strconv.Atoi(strings.TrimSpace(lines[5]))

	if handicap < 0 || handicap > 9 {
		return nil, fmt.Errorf("load_ngf(): got bad handicap")
	}

	// ------------------------------------

	komi, err := strconv.ParseFloat(strings.TrimSpace(lines[7]), 64)
	if err == nil {
		if float64(int(komi)) == komi {
			komi += 0.5
		}
	}

	// ------------------------------------

	if len(lines[8]) >= 8 {
		rawdate = lines[8][0:8]
	}

	// ------------------------------------

	if strings.Contains(strings.ToLower(lines[10]), "white win") || strings.Contains(strings.ToLower(lines[10]), "black lose") {
		re = "W+"
	}
	if strings.Contains(strings.ToLower(lines[10]), "black win") || strings.Contains(strings.ToLower(lines[10]), "white lose") {
		re = "B+"
	}

	margin := ""

	if strings.Contains(strings.ToLower(lines[10]), "resign") {
		margin = "R"
	}
	if strings.Contains(strings.ToLower(lines[10]), "time") {
		margin = "T"
	}

	// Finding a numeric margin in the file is... tricky. Not trying.

	if re != "" {
		re += margin
	}

	// ------------------------------------

	root := NewTree(boardsize)
	node := root

	if handicap > 1 {
		root.SetValue("HA", strconv.Itoa(handicap))
		root.SetValues("AB", HandicapPoints(boardsize, handicap, true))			// Uses Tygem layout
	}

	if komi != 0 {
		root.SetValue("KM", fmt.Sprintf("%.1f", komi))
	}

	if len(rawdate) == 8 {
		ok := true
		for n := 0; n < 8; n++ {
			if rawdate[n] < '0' || rawdate[n] > '9' {
				ok = false
			}
		}
		if ok {
			date := rawdate[0:4] + "-" + rawdate[4:6] + "-" + rawdate[6:8]
			root.SetValue("DT", date)
		}
	}

	if pw != "" { root.SetValue("PW", pw) }
	if pb != "" { root.SetValue("PB", pb) }
	if re != "" { root.SetValue("RE", re) }

	for _, line := range lines {

		line = strings.TrimSpace(line)
		line = strings.ToUpper(line)

		if len(line) < 7 {
			continue
		}

		if line[0:2] == "PM" {

			if line[4] == 'B' || line[4] == 'W' {

				key := byte_to_string(line[4])

				// Coordinates are from 1-19, but with "B" representing
				// the digit 1. (Presumably "A" would represent 0.)

				x := int(line[5]) - 66		// Thus, 66 is correct, to map B to 0, etc
				y := int(line[6]) - 66

				if x >= 0 && x < boardsize && y >= 0 && y < boardsize {
					node = NewNode(node)
					node.SetValue(key, Point(x, y))
				}
			}
		}
	}

	if root.MainChild() == nil {		// We'll assume we failed in this case
		return nil, fmt.Errorf("load_ngf(): root ended up with zero children")
	}

	return root, nil
}
