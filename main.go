package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strings"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Version: "1.0.0",
		Name:    "comelint",
		Usage:   "Linter for commit messages",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-merge",
				Value: false,
				Usage: "Prohibit MERGE messages",
			},
			&cli.BoolFlag{
				Name:  "no-revert",
				Value: false,
				Usage: "Prohibit REVERT messages",
			},
			&cli.BoolFlag{
				Name:  "no-wip",
				Value: false,
				Usage: "Prohibit WIP messages",
			},
			&cli.BoolFlag{
				Name:  "no-multiline",
				Value: false,
				Usage: "Prohibit multiline messages",
			},
			&cli.IntFlag{
				Name:        "min-length",
				Usage:       "Minimum length of messages",
				DefaultText: "1",
			},
			&cli.IntFlag{
				Name:        "max-length",
				Usage:       "Maximum length of messages",
				Value:       math.MaxInt64,
				DefaultText: "Unlimited",
			},
			&cli.StringFlag{
				Name:        "regexp",
				Usage:       "RegExp rule for messages",
				Value:       ".*",
				DefaultText: ".*",
			},
		},
		Action: func(context *cli.Context) error {
			message := context.Args().First()

			err := ValidateMessage(message, context)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			os.Exit(0)

			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func ValidateMessage(message string, context *cli.Context) error {

	if len(message) == 0 {
		return returnError("Please pass commit message")
	}

	if minLength := context.Int("min-length"); len(message) < minLength {
		return returnError("Message length is less than recommended")
	}

	if maxLength := context.Int("max-length"); len(message) > maxLength {
		return returnError("Message length is longer than recommended")
	}

	if context.Bool("no-multiline") {
		if len(strings.Split(strings.TrimSuffix(message, "\n"), "\n")) > 1 {
			return returnError("Multiline commits are prohibited")
		}
	}

	if context.Bool("no-merge") {
		match, _ := regexp.MatchString("(?i)^merge.*", message)
		if match {
			return returnError("Merge commits are prohibited")
		}
	}

	if context.Bool("no-revert") {
		match, _ := regexp.MatchString("(?i)^revert.*", message)
		if match {
			return returnError("Revert commits are prohibited")
		}
	}

	if context.Bool("no-wip") {
		match, _ := regexp.MatchString("(?i)^wip.*", message)
		if match {
			return returnError("WIP commits are prohibited")
		}
	}

	match, err := regexp.MatchString("(?i)"+context.String("regexp"), message)
	if err != nil {
		return returnError("Regexp rule is invalid")
	}
	if !match {
		return returnError("Message does not pass RegExp matching")
	}

	return nil
}

type ValidationError struct {
	message string
}

func (e *ValidationError) Error() string {
	return e.message
}

func returnError(errorMessage string) error {
	return &ValidationError{message: errorMessage}
}
