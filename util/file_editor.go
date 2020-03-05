package util

import (
	"io/ioutil"
	"os"
	"os/exec"
)

const DefaultEditor = "vim"

func OpenInEditor(template string) ([]byte, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = DefaultEditor
	}

	file, err := ioutil.TempFile(os.TempDir(), "*")
	if err != nil {
		return []byte{}, err
	}

	filename := file.Name()

	// Defer removal of the temporary file in case any of the next steps fail.
	defer os.Remove(filename)

	if err = file.Close(); err != nil {
		return []byte{}, err
	}

	if template != "" {
		ioutil.WriteFile(filename, []byte(template), 0644)
	}

	execPath, err := exec.LookPath(editor)

	cmd := exec.Command(execPath, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}
