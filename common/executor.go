package common

import (

	"os/exec"
)

func ExecuteWinCmd(cmd string) ([]byte,error){
	return exec.Command("cmd.exe", "/k", "chcp 65001 && "+cmd).Output()
}