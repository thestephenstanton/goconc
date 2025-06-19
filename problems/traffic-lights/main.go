package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"log/slog"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	err := run()
	if err != nil {
		slog.Error("starting application", "error", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

var lightDirections = []string{"North", "East", "South", "West"}

type lightControl struct {
	turnGreen chan<- any
	backToRed <-chan any
}

func run() error {
	start := time.Now()

	done := make(chan any)

	lightControls := make([]lightControl, 0, len(lightDirections))

	for _, d := range lightDirections {
		turnGreen := make(chan any)
		backToRed := light(done, turnGreen, d)

		lightControls = append(lightControls, lightControl{
			turnGreen: turnGreen,
			backToRed: backToRed,
		})
	}

	time.AfterFunc(10*time.Second, func() { close(done) })

	pedestrian := make(chan any)

	go func() {
		time.After(time.Duration(rand.Intn(5)+1) * time.Second)
		pedestrian <- nil
	}()

loop:
	for {
		for _, light := range lightControls {
			light.turnGreen <- nil

			select {
			case <-light.backToRed:
			case <-done:
				break loop
			}
		}
	}

	fmt.Printf("time since starting: %s\n", time.Since(start))

	return nil
}

func light(done <-chan any, turnGreen <-chan any, direction string) <-chan any {
	backToRed := make(chan any)

	var otherDirections string

	for _, d := range lightDirections {
		if d == direction {
			continue
		}

		otherDirections += d + "/"
	}

	otherDirections = strings.TrimRight(otherDirections, "/")

	go func() {
		defer close(backToRed)

		for {
			select {
			case <-turnGreen:
				currentTime := time.Now().Format("15:04:05")
				fmt.Printf("[Time: %s] %s is GREEN, %s are RED\n", currentTime, direction, otherDirections)

				select {
				case <-time.After(2 * time.Second):
					backToRed <- nil
				case <-done:
					return
				}
			case <-done:
				return
			}
		}
	}()

	return backToRed
}
