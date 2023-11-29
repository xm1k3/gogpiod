/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/warthog618/gpiod"
)

type Song struct {
	Frequency float64
	Duration  time.Duration
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gogpiod",
	Short: "Golang test with gpiod",
	Long:  `Golang test with gpiod`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		chipFLag, _ := cmd.Flags().GetString("chip")
		chip, err := gpiod.NewChip(chipFLag)
		if err != nil {
			log.Fatal(err)
		}
		defer chip.Close()

		fmt.Println("CHIP:", chip)

		line, err := chip.RequestLine(4, gpiod.AsOutput(0))
		if err != nil {
			log.Fatal(err)
		}
		defer line.Reconfigure(gpiod.AsInput)

		fmt.Println("LINE:", line)

		line.SetValue(1)
		time.Sleep(time.Second * 4)
		line.SetValue(0)

		chip.RequestLine(18, gpiod.WithEventHandler(handler), gpiod.WithBothEdges)

		// for _, note := range song {
		// 	playTone(line, note.Frequency, note.Duration)
		// }
	},
}

func handler(evt gpiod.LineEvent) {
	fmt.Println("EVENT:", evt)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gogpiod.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringP("chip", "c", "gpiochip0", "Chip name")
}

func playTone(line *gpiod.Line, freq float64, duration time.Duration) {
	high := time.Duration(float64(time.Second) / freq / 2)
	low := high

	start := time.Now()
	for time.Since(start) < duration {
		line.SetValue(1)
		time.Sleep(high)
		line.SetValue(0)
		time.Sleep(low)
	}
}
