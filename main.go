package main

import (
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	gosxnotifier "github.com/deckarep/gosx-notifier"
	"github.com/urfave/cli/v2"
)

func parsePosition(positionStr string) (uint8, time.Duration) {
	positionSlice := strings.Split(positionStr, ":")

	targetPosition, _ := strconv.Atoi(positionSlice[0])
	positionDuration, _ := strconv.Atoi(positionSlice[1])

	return uint8(targetPosition), time.Duration(positionDuration) * time.Minute
}

func notifyUser(targetHeight, height uint8) {
	var note *gosxnotifier.Notification
	if targetHeight > height {
		note = gosxnotifier.NewNotification("Time to STAND UP!")
	} else {
		note = gosxnotifier.NewNotification("Time to SIT")
	}

	note.Title = "Desk Controller"
	note.Sound = gosxnotifier.Funk
	note.Push()
}

func main() {
	app := &cli.App{
		Name:  "desk-controller",
		Usage: "TMotion desk BLE controller",
		Commands: []*cli.Command{
			{
				Name:  "scan",
				Usage: "scan bluetooth devices, CTRL+c to stop",
				Action: func(cCtx *cli.Context) error {
					slog.Info("TODO: Scanning")
					return nil
				},
			},
			{
				Name:  "cycle",
				Usage: "Cycles through specified positions.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "desk-address",
						Usage:    "UUID desk address, can be obtained with `scan` subcommand",
						Required: true,
					},
					&cli.StringSliceFlag{
						Name:  "position",
						Usage: "Positions to cycle through formatted as <centimeters:minutes>",
						Action: func(cCtx *cli.Context, positionStrings []string) error {
							positionRegex, _ := regexp.Compile(`^\d+:\d+$`)
							for _, positionStr := range positionStrings {
								if !positionRegex.MatchString(positionStr) {
									return fmt.Errorf("incorrect format for position: %s", positionStr)
								}
							}

							return nil
						},
					},
					&cli.IntFlag{
						Name:    "repeat",
						Aliases: []string{"r"},
						Usage:   "Repeat the cycle for a specified number of times",
						Value:   1,
					},
					&cli.IntFlag{
						Name:    "delay",
						Aliases: []string{"d"},
						Usage:   "Wait after notification for a specified number of seconds",
						Value:   3,
						Action: func(cCtx *cli.Context, delay int) error {
							if delay < 0 {
								return fmt.Errorf("delay must be a non-negative integer: %d", delay)
							}

							return nil
						},
					},
				},
				Action: func(cCtx *cli.Context) error {
					programLevel := new(slog.LevelVar)
					h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})
					slog.SetDefault(slog.New(h))

					if cCtx.Bool("verbose") {
						programLevel.Set(slog.LevelDebug)
					}

					desk := NewDesk(cCtx.String("desk-address"))
					err := desk.Connect()
					if err != nil {
						slog.Error(err.Error())
					}

					for i := 0; i < cCtx.Int("repeat"); i++ {
						for _, positionStr := range cCtx.StringSlice("position") {
							targetPosition, positionDurationMin := parsePosition(positionStr)

							notifyUser(targetPosition, desk.GetHeight())
							time.Sleep(time.Duration(cCtx.Int("delay")) * time.Second)

							slog.Info("Moving to position for specified duration", "position", targetPosition, "duration", positionDurationMin)
							desk.MoveTo(targetPosition)
							time.Sleep(positionDurationMin)
						}
					}

					return nil
				},
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Set log level to DEBUG",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error(err.Error())
	}
}
