package usecase

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/linde12/gowol"
	"golang.org/x/crypto/ssh"
)

type Commands struct {
	PCAddress        string
	UserName         string
	UserPassword     string
	MacAddress       string
	BroadcastAddress string
}

func NewCommandProvider(pcAddress, userName, userPassword, macAddress, broadcastAddress string) Commands {
	return Commands{
		PCAddress:        pcAddress,
		UserName:         userName,
		UserPassword:     userPassword,
		MacAddress:       macAddress,
		BroadcastAddress: broadcastAddress,
	}
}

func (c Commands) PcReboot() error {
	return c.runCommandOverSSH("sudo reboot")
}

func (c Commands) PcOff() error {
	return c.runCommandOverSSH("sudo poweroff")
}

func (c Commands) PcOn() error {
	pkt, err := gowol.NewMagicPacket(c.MacAddress)
	if err != nil {
		return fmt.Errorf("cannot create magic packet: %w", err)
	}

	if err := pkt.Send(c.BroadcastAddress); err != nil {
		return fmt.Errorf("cannot send packet: %w", err)
	}

	return nil
}

func (c Commands) PcStatus() error {
	cmd := exec.Command("ping", "-c", "1", "-W", "2", c.PCAddress)

	return cmd.Run()
}

func (c Commands) runCommandOverSSH(cmd string) error {
	log.Printf(c.UserName)
	log.Printf(c.UserPassword)

	config := &ssh.ClientConfig{
		User: c.UserName,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.UserPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	conn, err := ssh.Dial("tcp", c.PCAddress+":22", config)
	if err != nil {
		return fmt.Errorf("cannot connect to %s: %w", c.PCAddress, err)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return fmt.Errorf("cannot create session: %w", err)
	}
	defer session.Close()

	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	return nil
}
