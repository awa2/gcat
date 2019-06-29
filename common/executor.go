package common

func ExecuteWinCmd(cmd string) {
	out, err := exec.Command("cmd.exe", "/k", "chcp 65001 && "+cmd).Output()
    if err != nil {
        return nil, err
	}
	return out
}