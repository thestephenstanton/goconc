package main

// this is a very defensive implementation of this problem

import (
	"context"
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

type sig struct{}

var lightDirections = []string{"North", "East", "South", "West"}

func otherDirections(i int) []string {
	other := make([]string, 0, len(lightDirections)-1)

	for j, d := range lightDirections {
		if i == j {
			continue
		}

		other = append(other, d)
	}

	return other
}

func run() error {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pedestrian := make(chan sig)

	lights := make([]Light, 0, len(lightDirections))

	for i, d := range lightDirections {
		otherDirections := otherDirections(i)

		light := NewLight(d, otherDirections)
		light.Start()
		lights = append(lights, light)
	}

	go func() {
		d := time.Duration(rand.Intn(5)+1) * time.Second

		// d += time.Duration(200 * time.Millisecond)

		<-time.After(d)
		pedestrian <- sig{}
	}()

	time.Sleep(5 * time.Millisecond)

loop:
	for {
		for _, light := range lights {
			finished, ok := light.TurnGreen()
			if !ok {
				return fmt.Errorf("turning light %s green", light.direction)
			}

			finishCtx, finishCancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer finishCancel()

			select {
			case <-ctx.Done():
				break loop
			case <-finishCtx.Done():
				return fmt.Errorf("light %s took too long to change back", light.direction)
			case <-finished:
				finishCancel()
			case <-pedestrian:
				pedestrianStart := time.Now()

				ok := light.Interupt()
				if !ok {
					// this should be a very rare condition, this should only be able to happen
					// if between the time that we are signaled a pedestrian is cross and us calling
					// Interupt(), the light changes from green back to red. so basically after this
					// case statement starts but before the Interupt(). but we still continue as normal
					slog.Info("rare edge found!")
				}

				slog.Info("pedestrian crossing...")

				select {
				case <-time.After(1 * time.Second):
					slog.Info("...pedestrian finished crossing", "timeTaken", time.Since(pedestrianStart).String())
				case <-ctx.Done():
					break loop
				}
			}
		}
	}

	fmt.Printf("time since starting: %s\n", time.Since(start))

	return nil
}

type Light struct {
	direction       string
	otherDirections string

	done      chan sig
	turnGreen chan chan sig
	interupt  chan sig
}

func NewLight(direction string, otherDirections []string) Light {
	return Light{
		direction:       direction,
		otherDirections: strings.Join(otherDirections, "/"),
		done:            make(chan sig),
		turnGreen:       make(chan chan sig),
		interupt:        make(chan sig),
	}
}

func (l Light) Start() {
	go func() {
		for {
			select {
			case <-l.done:
				return
			case finished := <-l.turnGreen:
				currentTime := time.Now().Format("15:04:05")
				fmt.Printf("[Time: %s] %s is GREEN, %s are RED\n", currentTime, l.direction, l.otherDirections)

				select {
				case <-l.done:
					return
				case <-time.After(2 * time.Second):
					close(finished)
				case <-l.interupt:
					close(finished)
				}
			}
		}
	}()
}

// returns back a channel that will close once the light is no longer green
func (l Light) TurnGreen() (<-chan sig, bool) {
	finished := make(chan sig)

	select {
	case l.turnGreen <- finished:
		return finished, true
	default:
		// this should only happen if the light is already green
		return nil, false
	}
}

func (l Light) Interupt() bool {
	select {
	case l.interupt <- sig{}:
		return true
	default:
		return false
	}
}

func (l Light) Stop() {
	close(l.done)
}
