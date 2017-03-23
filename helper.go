package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// GetInput - simplifies getting and formatting user input
func GetInput(statement ...string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(statement)
	answer, _ := reader.ReadString('\n')

	return strings.TrimSpace(answer)
}
