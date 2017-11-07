/*
 * Copyright (c) 2017 Daniel Müller
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/fuchsi/bencode"
	"github.com/fuchsi/torrentfile"
	"github.com/pborman/getopt/v2"
)

var noColorFlag = getopt.BoolLong("nocolors", 'n', "No ANSI colour")
var t textFormatter

func main() {
	helpFlag := getopt.BoolLong("help", 'h', "Show this help message and exit")
	versionFlag := getopt.BoolLong("version", 'v', "Print version and quit")
	filesFlag := getopt.BoolLong("files", 'f', "Show files within the torrent")
	detailedFlag := getopt.BoolLong("detailed", 'd', "Show detailed information about the files")
	verboseFlag := getopt.BoolLong("everything", 'e', "Print everything about the torrent")

	getopt.SetParameters("<filename>")

	// Parse the program arguments
	getopt.Parse()

	if *detailedFlag && *filesFlag {
		getopt.Usage()
		fmt.Println("gotorrentinfo: error: argument -f/--files: not allowed with argument -d/--detailed")
		return
	}
	if *verboseFlag && *filesFlag {
		getopt.Usage()
		fmt.Println("gotorrentinfo: error: argument -f/--files: not allowed with argument -e/--everything")
		return
	}

	if *versionFlag {
		fmt.Println("gotorrentinfo v0.1.0")
		return
	}
	if *helpFlag || getopt.NArgs() == 0 {
		getopt.Usage()
		return
	}

	file, err := os.Open(getopt.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	t = textFormatter{}
	t.init()

	labelformat := t.bright | t.yellow
	indent := strings.Repeat(" ", 4)
	colWidth := 19

	if !*verboseFlag {
		tfile, err := torrentfile.DecodeTorrentFile(file)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(t.format(filepath.Base(file.Name()), t.bright))
		if !*detailedFlag {
			fmt.Println(indent + t.format("name", labelformat) + strings.Repeat(" ", colWidth-len(indent+"name")) + tfile.Name)
			fmt.Println(indent + t.format("comment", labelformat) + strings.Repeat(" ", colWidth-len(indent+"comment")) + tfile.Comment)
			fmt.Println(indent + t.format("announce url", labelformat) + strings.Repeat(" ", colWidth-len(indent+"announce url")) + tfile.AnnounceUrl)
			fmt.Println(indent + t.format("created by", labelformat) + strings.Repeat(" ", colWidth-len(indent+"created by")) + tfile.CreatedBy)
			fmt.Println(indent + t.format("created on", labelformat) + strings.Repeat(" ", colWidth-len(indent+"created on")) + t.format(tfile.CreationDate.String(), t.magenta))
			if tfile.Encoding != "" {
				fmt.Println(indent + t.format("enconding", labelformat) + strings.Repeat(" ", colWidth-len(indent+"enconding")) + tfile.Encoding)
			}
			fmt.Println(indent + t.format("num files", labelformat) + strings.Repeat(" ", colWidth-len(indent+"num files")) + strconv.FormatInt(int64(len(tfile.Files)), 10))
			fmt.Println(indent + t.format("total size", labelformat) + strings.Repeat(" ", colWidth-len(indent+"total size")) + t.format(formatBytes(tfile.TotalSize()), t.cyan))
			fmt.Println(indent + t.format("Info Hash", labelformat) + strings.Repeat(" ", colWidth-len(indent+"Info Hash")) + fmt.Sprintf("%x", tfile.InfoHash()))
		}

		if *filesFlag || *detailedFlag {
			fmt.Println(indent + t.format("files", t.bright|t.yellow))
			for key, value := range tfile.Files {
				fmt.Println(strings.Repeat(indent, 2) + t.format(fmt.Sprintf("%d", key), t.bright|t.yellow))
				fmt.Println(strings.Repeat(indent, 3) + t.format("path", t.bright|t.yellow))
				fmt.Println(strings.Repeat(indent, 4) + value.Path)
				fmt.Println(strings.Repeat(indent, 3) + t.format("length", t.bright|t.yellow))
				fmt.Println(strings.Repeat(indent, 4) + t.format(formatBytes(value.Length), t.cyan))

			}
		}
		if *detailedFlag {
			fmt.Println(t.format("    piece length", t.bright|t.yellow))
			fmt.Println("            " + t.format(formatBytes(tfile.PieceLength), t.cyan))
			fmt.Println(t.format("    pieces", t.bright|t.yellow))
			pieces := len(tfile.Pieces) * torrentfile.PIECE_SIZE
			fmt.Println(t.format(fmt.Sprintf("            [%d UTF-8 Bytes]", pieces), t.red|t.bright))
		}
	} else {
		dict, err := bencode.Decode(file)
		if err != nil {
			log.Fatal(err)
		}

		printDict(dict, 1)
	}

	fmt.Println()
}

type textFormatter struct {
	none, normal, bright, white, green, red, cyan, yellow, magenta, dull int

	escape  byte
	reset   string
	mapping map[int]string
}

func (t *textFormatter) init() {
	t.none = 0x000000
	t.normal = 0x000001
	t.bright = 0x000002
	t.white = 0x000004
	t.green = 0x000008
	t.red = 0x000010
	t.cyan = 0x000020
	t.yellow = 0x000040
	t.magenta = 0x000080
	t.dull = 0x000100

	t.escape = 0x1b
	t.reset = "[0m"

	t.mapping = make(map[int]string)

	t.mapping[t.normal] = "[0m"
	t.mapping[t.bright] = "[1m"
	t.mapping[t.dull] = "[22m"
	t.mapping[t.white] = "[37m"
	t.mapping[t.green] = "[32m"
	t.mapping[t.cyan] = "[36m"
	t.mapping[t.yellow] = "[33m"
	t.mapping[t.red] = "[31m"
	t.mapping[t.magenta] = "[35m"
}

func (t textFormatter) format(s string, format int) string {
	if !*noColorFlag {
		codestring := ""
		for name, code := range t.mapping {
			if name&format != 0 {
				codestring += fmt.Sprintf("%c%s", t.escape, code)
			}
		}
		if codestring != "" {
			return fmt.Sprintf("%s%s%c%s", codestring, s, t.escape, t.reset)
		}
	}

	return s
}

func formatBytes(size uint64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	}
	suffixArray := []byte{'K', 'M', 'G', 'T', 'P', 'E', 'Z', 'Y'}
	suffixCount := len(suffixArray)
	i := 0
	bytes := float64(size)

	for ; bytes > 1024 && i < suffixCount; i++ {
		bytes /= 1024
	}

	return fmt.Sprintf("%0.2f %cB", bytes, suffixArray[i-1])
}

func printDict(dict map[string]interface{}, indent int) {
	// Sort keys
	list := make(sort.StringSlice, len(dict))
	i := 0
	for key := range dict {
		list[i] = key
		i++
	}
	list.Sort()

	format := t.green
	if indent < 2 {
		format = t.bright | t.yellow
	}

	for _, key := range list {
		fmt.Print(strings.Repeat(" ", indent*4))
		fmt.Println(t.format(key, format))

		switch val := dict[key].(type) {
		case string:
			printString(val, key, indent)
		case []interface{}:
			printList(val, indent+1)
		case map[string]interface{}:
			printDict(val, indent+1)
		case int64:
			printInt(val, indent)
		}
	}
}

func printList(list []interface{}, indent int) {
	if len(list) == 1 {
		fmt.Print(strings.Repeat(" ", indent*4))
		fmt.Println(list[0])
		return
	}

	for key, val := range list {
		fmt.Print(strings.Repeat(" ", indent*4))
		fmt.Println(t.format(fmt.Sprintf("%d", key), t.bright|t.yellow))

		switch val := val.(type) {
		case string:
			printString(val, "", indent)
		case []interface{}:
			printList(val, indent+1)
		case map[string]interface{}:
			printDict(val, indent+1)
		case int64:
			printInt(val, indent)
		}
	}
}

func printInt(val int64, indent int) {
	fmt.Print(strings.Repeat(" ", (indent+1)*4))
	fmt.Println(t.format(fmt.Sprintf("%d", val), t.cyan))
}

func printString(val, key string, indent int) {
	fmt.Print(strings.Repeat(" ", (indent+1)*4))
	if key == "pieces" {
		fmt.Println(t.format(fmt.Sprintf("[%d UTF-8 Bytes]", len(val)), t.red|t.bright))
	} else {
		fmt.Println(val)
	}
}
