package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bmatcuk/doublestar/v3"
	"github.com/maeda6uiui/go-counter"
	"github.com/urfave/cli/v2"
)

func addCounts(counts map[string]int, inputFilepath string, threshold int, toLower bool) error {
	fp, err := os.Open(inputFilepath)
	if err != nil {
		return err
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		splits := strings.SplitN(line, " ", 2)
		if len(splits) != 2 {
			continue
		}

		var key string
		if toLower {
			key = strings.ToLower(splits[1])
		} else {
			key = splits[1]
		}

		count, err := strconv.Atoi(splits[0])
		if err != nil {
			return err
		}

		if count < threshold {
			continue
		}

		if curCount, exists := counts[key]; exists {
			counts[key] = curCount + count
		} else {
			counts[key] = count
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func appAction(c *cli.Context) error {
	inputDir := c.String("inputDir")
	outputFilepath := c.String("outputFilepath")
	threshold := c.Int("threshold")
	toLower := c.Bool("toLower")

	globPattern := filepath.Join(inputDir, "*.txt")
	inputFilepaths, err := doublestar.Glob(globPattern)
	if err != nil {
		return err
	}

	counts := make(map[string]int)

	for _, inputFilepath := range inputFilepaths {
		fmt.Println(inputFilepath)

		if err := addCounts(counts, inputFilepath, threshold, toLower); err != nil {
			return err
		}
	}

	outputFile, err := os.Create(outputFilepath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)

	cnter := counter.NewCounterFromMap(counts)
	keys, freqs := cnter.MostCommon()

	for i := 0; i < len(keys); i++ {
		writer.WriteString(fmt.Sprintf("%v\t%v\n", keys[i], freqs[i]))
	}

	writer.Flush()

	return nil
}

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "inputDir",
				Aliases: []string{"i"},
			},
			&cli.StringFlag{
				Name:    "outputFilepath",
				Aliases: []string{"o"},
			},
			&cli.IntFlag{
				Name:    "threshold",
				Aliases: []string{"t"},
				Value:   1,
			},
			&cli.BoolFlag{
				Name:    "toLower",
				Aliases: []string{"l"},
			},
		},

		Action: appAction,
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
