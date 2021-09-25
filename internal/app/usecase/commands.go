package usecase

import (
	"os/exec"
)

type Commands struct{}

func NewCommandProvider() Commands {
	return Commands{}
}

func (c Commands) Reboot() error {
	cmd := exec.Command("reboot")

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (c Commands) Poweroff() error {
	cmd := exec.Command("poweroff")

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
