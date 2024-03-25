package cash

import "os/exec"

var CommandProcess map[int]*exec.Cmd

func NewCommandProcess() {
	CommandProcess = make(map[int]*exec.Cmd)
}
