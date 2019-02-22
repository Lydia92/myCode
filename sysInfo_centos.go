package main

import (
	"fmt"

	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/cpu"

	"os/exec"
	//"log"
	"strings"
	"strconv"
	//"encoding/json"
)

func main() {
	//Info :=make(map[string]string)
	Info:=sysInfo()
	fmt.Println(Info)


}

func sysInfo()map[string]interface{}{
	info:=make(map[string]interface{})
	cp:=make(map[string]string)
	memory:=make(map[string]string)
	cmd :=exec.Command("df","-hP")
	out,_:=cmd.CombinedOutput()
	str:=string(out)
    var diskinfo []interface{}
	st:=strings.Fields(str)
	for i,v :=range st{
		if i%6==0 && i !=0 && i!=6{
			d,err:=disk.Usage(v)
			if err!=nil{
				fmt.Println("err")
			}
			diskinfo=append(diskinfo,d.Path)
			//dis["Path"]=d.Path
			//dis["diskTotal"]=strconv.Itoa(int(d.Total/1024/1024/1024))
			//dis["diskFree"]=strconv.Itoa(int(d.Free/1024/1024/1024))
			//dis["diskUsed"]=strconv.Itoa(int(d.Used/1024/1024/1024))
			//dis["diskUsedPercent"]=strconv.Itoa(int(d.UsedPercent))
			ds := map[string]string{"diskTotal":strconv.Itoa(int(d.Total/1024/1024)),
				"diskUsed":strconv.Itoa(int(d.Used/1024/1024)),"diskFree":strconv.Itoa(int(d.Free/1024/1024)),
				"diskUsedPercent":strconv.Itoa(int(d.UsedPercent))}
			diskinfo=append(diskinfo,ds)
			info["diskinfo"]=diskinfo
			//aa:=make(map[string]string)

			//info["disk"]=dis
			//fmt.Println("Path:",info["Path"],"total:",info["total"],"G",
			// "free:",info["Free"],"G","used:",info["used"],"G","usedPercent:",info["UsedPercent"])

		}
	}
	m,_:=mem.VirtualMemory()
	MemoryTotal:=float64(m.Total)
	Used:=float64(m.Total-m.Free-m.Buffers-m.Cached)

	UsedPercent:=float64((100*Used/MemoryTotal))

	//fmt.Printf("UsedPercent:%.2f",UsedPercent)
	//fmt.Println("MemoryTotal:",m.Total/1024/1024,"M","Used:",(m.Total-m.Free-m.Buffers-m.Cached)/1024/1024,"M","UsedPercent:",UsedPercent*100,"%","Cached:",m.Cached/1024/1024,"M","SwapTotal:",m.SwapTotal/1024/1024,"M","SwapFree",m.SwapFree/1024/1024,"M")
	memory["memoryTotal"]=strconv.Itoa(int(m.Total/1024/1024))
	memory["memoryUsed"]=strconv.Itoa(int((m.Total-m.Free-m.Buffers-m.Cached)/1024/1024))
	memory["memoryUsedPercent"]=strconv.Itoa(int(UsedPercent*100))
	memory["memoryCached"]=strconv.Itoa(int(m.Cached/1024/1024))
	memory["SwapTotal"]=strconv.Itoa(int(m.SwapTotal/1024/1024))
	memory["SwapFree"]=strconv.Itoa(int(m.SwapFree/1024/1024))
	info["memory"]=memory
	i,_:=cpu.Counts(true)
	//fmt.Println("CPU:",i)
	cp["cpuCounts"]=strconv.Itoa(int(i))
	info["cpu"]=cp
/*	s,err:=json.Marshal(info)
	if err!=nil{
		fmt.Println("error")
	}
	ss:=string(s)*/
	return info

	//return info

}
