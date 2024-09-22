package main

import (
	"okieoth/schemaguesser/cmd/schemaguesser/cmd"
	"okieoth/schemaguesser/internal/pkg/logger"
)

func main() {
	logger.Init("schemaguesser.log")
	cmd.Execute()
}
