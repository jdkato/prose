package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/jdkato/aptag"
	"github.com/jdkato/aptag/tokenize"
	"github.com/urfave/cli"
)

// Version is the semantic version number
var Version string

func main() {
	var file string
	var text []byte
	var tags []tokenize.Token
	var err error
	var isTokenized bool

	app := cli.NewApp()
	app.Name = "aptag"
	app.Usage = "A command-line POS tagger."
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "path",
			Usage:       "read `path` as source input instead of stdin",
			Destination: &file,
		},
		cli.BoolFlag{
			Name:        "tokens",
			Usage:       "assume input has already been tokenized",
			Destination: &isTokenized,
		},
	}

	app.Action = func(c *cli.Context) error {
		if file != "" {
			text, err = ioutil.ReadFile(file)
			if err != nil {
				panic(err)
			}
		} else {
			stat, serr := os.Stdin.Stat()
			if serr != nil {
				panic(err)
			} else if (stat.Mode() & os.ModeCharDevice) == 0 {
				reader := bufio.NewReader(os.Stdin)
				text, err = ioutil.ReadAll(reader)
				if err != nil {
					panic(err)
				}
			}
		}
		if len(text) > 0 {
			tagger := aptag.NewPerceptronTagger()
			if isTokenized {
				tags = tagger.Tag(strings.Split(string(text), " "))
			} else {
				tags = tagger.TokenizeAndTag(string(text))
			}
			b, jerr := json.Marshal(tags)
			if jerr != nil {
				return jerr
			}
			fmt.Println(string(b))
		}
		return err
	}

	if app.Run(os.Args) != nil {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
