package transmission

import "os/exec"

type Transmission struct{}

func TranmissionClient() Transmission {
	return Transmission{}
}

func (t Transmission) AddTorrent(filename string) error {
	cmd := exec.Command("transmission-remote", "-a", filename)

	return cmd.Run()
}
