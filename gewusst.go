package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type Gewusst [2]string

const (
	defaultFolder = "gewusst"
)

var (
	gewFolder   = flag.String("gewusst", path.Join(os.Getenv("pwd"), defaultFolder), "folder to use for 'gewusst.html' template")
	messages    []Gewusst
	gewusstWait sync.Mutex
)

func NewGewusst(title, text string) Gewusst {
	return Gewusst{title, text}
}

// update returns the amount of messages in the gewusst folder
func updateGewusst() int {
	gewusstWait.Lock()
	defer gewusstWait.Unlock()
	quart := "Q"
	switch time.Now().Month() {
	case 12, 1, 2:
		quart += "1"
		break
	case 3, 4, 5:
		quart += "2"
		break
	case 6, 7, 8:
		quart += "3"
		break
	case 9, 10, 11:
		quart += "4"
		break
	}
	files, err := ioutil.ReadDir(path.Join(*gewFolder, quart))
	if err != nil {
		fmt.Errorf("can't find 'gewusst' folder: %s\n", *gewFolder)
		messages = make([]Gewusst, 0)
		return 0
	}
	messages = make([]Gewusst, 0, len(files))
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".txt") {
			content, err := ioutil.ReadFile(path.Join(*gewFolder, quart, file.Name()))
			if err != nil {
				continue
			}
			messages = append(messages, NewGewusst(file.Name()[:len(file.Name())-4], string(content)))
		}
	}

	return len(messages)
}

func testOutput(w io.Writer) {
	for i := range messages {
		fmt.Fprintf(w, "%s:\n\n%s\n\n\n", messages[i].Title(), messages[i].Message())
	}
}

func (g *Gewusst) Title() string {
	return g[0]
}

func (g *Gewusst) Message() string {
	return g[1]
}

func (g *Gewusst) String() string {
	return fmt.Sprintf("%s:\n%s\n", g[0], g[1])
}
