package gitconfig

import (
	"os/exec"
)

// GitCmd runs the git command with a list of strings as arguments
func GitCmd(c ...string) (string, error) {
	rawOutput := exec.Command("git", c...)

	o, err := rawOutput.Output()

	output := string(o[:])

	return output, err
}
