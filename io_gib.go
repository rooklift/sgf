package sgf

// Basic GIB parser (i.e. for Tygem files)

import (
	"fmt"
	"strconv"
	"strings"
)

func load_gib(gib string) (*Node, error) {

	root := NewTree(19)
	node := root

	lines := strings.Split(gib, "\n")

	for _, line := range lines {

		line = strings.TrimSpace(line)

		// Names...

		if strings.HasPrefix(line, "\\[GAMEBLACKNAME=") && strings.HasSuffix(line, "\\]") {
			root.SetValue("PB", line[16: len(line) - 2])
		}

		if strings.HasPrefix(line, "\\[GAMEWHITENAME=") && strings.HasSuffix(line, "\\]") {
			root.SetValue("PW", line[16: len(line) - 2])
		}

		// Game info...

		if strings.HasPrefix(line, "\\[GAMETAG=") {
			dt, re, km := parse_gib_gametag(line)
			if dt != "" { root.SetValue("DT", dt) }
			if re != "" { root.SetValue("RE", re) }
			if km != "" { root.SetValue("KM", km) }
		}

		// Split the line into tokens for the handicap and move parsing...

		fields := strings.Fields(line)

		// Handicap...

		if len(fields) >= 4 && fields[0] == "INI" {

			if node != root {
				return nil, fmt.Errorf("load_gib(): got INI field after moves were made")
			}

			handicap, _ := strconv.Atoi(fields[3])

			if handicap > 1 {
				root.SetValue("HA", strconv.Itoa(handicap))
				root.SetValues("AB", HandicapPoints19(handicap, true))
			}
		}

		// Moves...

		if len(fields) >= 6 && fields[0] == "STO" {
			x, err1 := strconv.Atoi(fields[4])
			y, err2 := strconv.Atoi(fields[5])
			if err1 == nil && err2 == nil {
				key := "B"; if fields[3] == "2" { key = "W" }
				node = NewNode(node)
				node.SetValue(key, Point(x, y))
			}
		}
	}

	return root, nil
}

func parse_gib_gametag(line string) (dt, re, km string) {

	fields := strings.Split(line, ",")

	var zipsu int

	for _, s := range fields {

		if len(s) < 2 {
			continue
		}

		if s[0] == 'C' {
			dt = s[1:]
			if len(dt) > 10 {
				dt = dt[:10]
			}
			dt = strings.Replace(dt, ":", "-", -1)
		}

		if s[0] == 'W' {
			grlt, err := strconv.Atoi(s[1:])
			if err == nil {
				switch grlt {
				case 0:
					re = "B+"
				case 1:
					re = "W+"
				case 3:
					re = "B+R"
				case 4:
					re = "W+R"
				case 7:
					re = "B+T"
				case 8:
					re = "W+T"
				}
			}
		}

		if s[0] == 'G' {
			gongje, err := strconv.Atoi(s[1:])
			if err == nil {
				km = fmt.Sprintf("%.1f", float64(gongje) / 10.0)
				if strings.HasSuffix(km, ".0") {
					km = km[:len(km) - 2]
				}
			}
		}

		if s[0] == 'Z' {
			zipsu, _ = strconv.Atoi(s[1:])
		}
	}

	if re == "B+" || re == "W+" {
		if zipsu > 0 {
			re += fmt.Sprintf("%.1f", float64(zipsu) / 10.0)
		}
		if strings.HasSuffix(re, ".0") {
			re = re[:len(re) - 2]
		}
	}

	return dt, re, km
}
