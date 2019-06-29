package common

import(
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
)

func GetOSInfo(){
	kernelVersion, _ := host.KernelVersion()

}

func GetCPUs(){
	InfoStats, err := cpu.Info()
}