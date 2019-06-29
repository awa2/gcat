package common

import (

	"os/exec"
)

func ExecuteWinCmd(cmd string) ([]byte,error){
	out, err := exec.Command("cmd.exe", "/k", "chcp 65001 && "+cmd).Output()
    if err != nil {
        return out,err
	}

	return out,err
}