package cmdlib

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Input(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
