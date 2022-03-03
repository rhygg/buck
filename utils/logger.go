package utils

import (
	"fmt"

	"github.com/fatih/color"
)

//simple logger with colorful outputs.

func Log(msg string) string {
	return color.CyanString(msg)
}

func LogError(msg string) string {
	return color.RedString(msg)
}

func LogInfo(msg string) string {
	return color.GreenString(msg)
}

func LogWarning(msg string) {
	fmt.Println(color.YellowString(msg))
}

func LogDebug(msg string) {
	fmt.Println(color.BlueString(msg))
}
