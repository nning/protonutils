package utils

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// AskYesOrNo writes text to w and reads a line from r. It returns true if the
// answer starts with "y" or "Y". End of input is treated as "no".
func AskYesOrNo(r io.Reader, w io.Writer, text string) (bool, error) {
	fmt.Fprint(w, text, " [y/N] ")

	line, err := bufio.NewReader(r).ReadString('\n')
	if err == io.EOF {
		fmt.Fprintln(w)
	} else if err != nil {
		return false, err
	}

	answer := strings.TrimSpace(line)
	return strings.HasPrefix(answer, "y") || strings.HasPrefix(answer, "Y"), nil
}
