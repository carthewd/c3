package gitconfig

import (
	"path/filepath"
	"strings"

	"gopkg.in/go-ini/ini.v1"
)

// GetOpts will parse local .git config for options
func GetOpts(optionSect string, optionKey string) (string, error) {
	cfg, err := ini.Load(".git/config")
	if err != nil {
		return "NoSuchFile", err
	}

	val := cfg.Section(optionSect).Key(optionKey).String()

	return val, err
}

// GetOrigin returns the remote origin of a repo in the CWD
func GetOrigin() (string, error) {
	args := []string{"config", "remote.origin.url"}

	origin, err := GitCmd(args...)
	if err != nil || !strings.Contains(origin, "codecommit") {
		return "", err
	}

	origin = strings.Replace(origin, "\n", "", -1)
	return filepath.Base(origin), err
}
