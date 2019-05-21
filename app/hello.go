package main

import (
	"fmt"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	//go get github.com/StackExchange/wmi	これもいった
)

func main() {
	v, _ := mem.VirtualMemory()
	// almost every return value is a struct
	fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)

	// convert to JSON. String() is also implemented
	fmt.Println(v)

	// cpuやってみた
	m, _ := cpu.Info()
	fmt.Printf("CPU: %v", m)

	// var data = v
	// fpw, err := os.Create("taka.json")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fpw.Write(data)
}
