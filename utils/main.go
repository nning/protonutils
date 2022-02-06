package utils

import (
	"bufio"
	"fmt"
	"os"
)

// AskYesOrNo shows text and waits for users to confirm
func AskYesOrNo(text string) (bool, error) {
	fmt.Print(text, " [y/N] ")

	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		return false, err
	}

	if char != 121 && char != 89 {
		return false, nil
	}

	return true, nil
}
