// Package pterm is a modern go module to beautify console output.
// It can be used without configuration, but if desired, everything can be customized down to the smallest detail.
//
// Official docs are available at: https://docs.pterm.sh
//
// View the animated examples here: https://github.com/pterm/pterm#-examples
package pterm

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"atomicgo.dev/cursor"
	"github.com/gookit/color"
)

var (
	// Output completely disables output from pterm if set to false. Can be used in CLI application quiet mode.
	Output = true

	// PrintDebugMessages sets if messages printed by the DebugPrinter should be printed.
	PrintDebugMessages = false

	// RawOutput is set to true if pterm.DisableStyling() was called.
	// The variable indicates that PTerm will not add additional styling to text.
	// Use pterm.DisableStyling() or pterm.EnableStyling() to change this variable.
	// Changing this variable directly, will disable or enable the output of colored text.
	RawOutput = false
)

func init() {
	color.ForceColor()

	// Make the cursor visible when the program stops
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		for s := range c {
			cursor.Show()

			// Re-raise the signal to trigger the default behavior
			signal.Stop(c)
			p, err := os.FindProcess(os.Getpid())
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "failed to find process: %v\n", err)
				return
			}
			err = p.Signal(s)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "failed to signal process: %v\n", err)
				return
			}
		}
	}()
}

// EnableOutput enables the output of PTerm.
func EnableOutput() {
	Output = true
}

// DisableOutput disables the output of PTerm.
func DisableOutput() {
	Output = false
}

// EnableDebugMessages enables the output of debug printers.
func EnableDebugMessages() {
	PrintDebugMessages = true
}

// DisableDebugMessages disables the output of debug printers.
func DisableDebugMessages() {
	PrintDebugMessages = false
}

// EnableStyling enables the default PTerm styling.
// This also calls EnableColor.
func EnableStyling() {
	RawOutput = false
	EnableColor()
}

// DisableStyling sets PTerm to RawOutput mode and disables all of PTerms styling.
// You can use this to print to text files etc.
// This also calls DisableColor.
func DisableStyling() {
	RawOutput = true
	DisableColor()
}

// RecalculateTerminalSize updates already initialized terminal dimensions. Has to be called after a terminal resize to guarantee proper rendering. Applies only to new instances.
func RecalculateTerminalSize() {
	// keep in sync with DefaultBarChart
	DefaultBarChart.Width = GetTerminalWidth() * 2 / 3
	DefaultBarChart.Height = GetTerminalHeight() * 2 / 3
	DefaultParagraph.MaxWidth = GetTerminalWidth()
}
