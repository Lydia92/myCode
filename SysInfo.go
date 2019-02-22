package main

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"time"
	"strconv"
	"net"
	"os"
	"os/exec"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/load"
	"strings"

)

type HostInfo struct {
	Id          uint64    `orm:"auto" description:"主键"`
	Host        string    `orm:"size(15);unique" description:"表示主机地址"`
	Port        uint16    `orm:"default(22)" description:"表示主机ssh端口"`
	Disk        string    `orm:"size(255);null" description:"磁盘信息,{目录:{总大小，空闲量，使用比}}"`
	Cpu         string    `orm:"size(32);null" description:"cpu信息,{物理id:[核心数,频率],}"`
	Mem         string    `orm:"null" description:"内存信息,单位是MB"`
	Environment uint8     `orm:"index;default(0)" description:"0:unknown,1:dev,2:sit,3:uat,4:vis,5:prd"`
	Created     time.Time `orm:"auto_now_add;type(datetime);auto_now_add"`
	Updated     time.Time `orm:"auto_now;type(datetime);auto_now"`
}
type OsDisk struct {
	Id  uint64             `orm:"auto" description:"主键"`
	NodeAddr string             `orm:"size(15)" description:"表示主机地址"`
	Mounted  string             `orm:"size(25)" description:"磁盘挂载"`
	TotalSize float64
	UsedSize float64
	AvailSize float64
	UsedRate float64
	CreateTime  time.Time          `orm:"auto_now_add;type(datetime)"`
}

type OsDiskHistory struct {
	Id uint64             `orm:"auto" description:"主键"`
	NodeAddr string             `orm:"size(15)" description:"表示主机地址"`
	Mounted  string             `orm:"size(25)" description:"磁盘挂载"`
	TotalSize float64
	UsedSize float64
	AvailSize float64
	UsedRate float64
	CreateTime  time.Time `orm:"type(datetime)"`
}

type OsMem struct {
	Id  uint64             `orm:"auto" description:"主键"`
	NodeAddr string             `orm:"size(15)" description:"表示主机地址"`
	MemoryTotal  float64
	MemoryUsed float64
	MemoryUsedPercent float64
	MemoryCached  float64
	SwapTotal  float64
	SwapFree  float64
	CpuCount float64             `orm:"null" description:"cpu核数"`
	Load1 float64  `orm:"null" description:"1分钟负载"`
	Load5 float64  `orm:"null" description:"5分钟负载"`
	Load15  float64  `orm:"null" description:"15分钟负载"`
	CreateTime  time.Time          `orm:"auto_now_add;type(datetime);auto_now_add"`
}

type OsMemHistory struct {
	Id  uint64             `orm:"auto" description:"主键"`
	NodeAddr string             `orm:"size(15)" description:"表示主机地址"`
	MemoryTotal  float64
	MemoryUsed float64
	MemoryUsedPercent float64
	MemoryCached  float64
	SwapTotal  float64
	SwapFree  float64
	CpuCount float64             `orm:"null" description:"cpu核数"`
	Load1 float64  `orm:"null" description:"1分钟负载"`
	Load5 float64  `orm:"null" description:"5分钟负载"`
	Load15  float64  `orm:"null" description:"15分钟负载"`
	CreateTime  time.Time
}


func init(){
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", "root", "123456", "192.168.160.134", 3306, "aa")
	//fmt.Println(dns)
	//注册数据库，连接对应的生产环境
	orm.RegisterDataBase("default", "mysql", dns)
	orm.RegisterModel(new(HostInfo),new(OsDisk),new(OsDiskHistory),new(OsMem),new(OsMemHistory),)
	orm.RunSyncdb("default", false, true)
}
func main() {
	DB := orm.NewOrm()
	DB.Using("aa")
	var host []orm.Params
	//获取磁盘信息表中的ip
	sql := "select distinct host  from host_info "
	DB.Raw(sql).Values(&host)
	ip := GetIntranetIp()
	for _, v := range host {
		if ip == v["host"] {
			InsertDiskAllInfo(DB, ip)
			InertMemAllInfo(DB, ip)
		}
	}
}


//磁盘信息表及磁盘历史信息表插入数据
func InsertDiskAllInfo(DB orm.Ormer,ip string){
	var osDisk []OsDisk
	//获取磁盘信息表中的ip
	DB.Raw("select id,node_addr,mounted,total_size,used_size,avail_size,used_rate," +
		"create_time from os_disk where node_addr=?",ip).QueryRows(&osDisk)
	//将信息表中的详情插入到历史详情表中，并清空详情表对应的记录
	for _,v:=range osDisk{
		InsertDiskHistoryInfo(DB,ip,v.Mounted,v.TotalSize,v.UsedSize,v.AvailSize,v.UsedRate,v.CreateTime)
		DB.Delete(&OsDisk{Id:v.Id})
	}
	InserDisk(DB,ip)
}

//内存及内存历史表插入数据
func InertMemAllInfo(DB orm.Ormer,ip string){
	var osMem []OsMem
	//如果匹配到本机内网ip，则查取信息表中的详情。
	DB.Raw("select id,node_addr,memory_total,memory_used,memory_used_percent,memory_cached," +
		"	swap_total,swap_free,create_time, cpu_count,load1,load5,load15 from os_mem where node_addr=?",ip).QueryRows(&osMem)
	//将信息表中的详情插入到历史详情表中，并清空详情表对应的记录
	for _,v:=range osMem{
		InsertMemHistoryInfo(DB,ip,v.MemoryTotal,v.MemoryUsed,v.MemoryUsedPercent,
			v.MemoryCached,v.SwapTotal,v.SwapFree,v.CreateTime,v.CpuCount,v.Load1,v.Load5,v.Load15)
		DB.Delete(&OsMem{Id:v.Id})
	}
	InserMem(DB,ip)
}
//磁盘信息表插入新数据
func InserDisk(DB orm.Ormer,ip string) {
	//获取磁盘详情
	diskinfo := GetDiskInfo()
	osdisk := OsDisk{}
	for _, v := range diskinfo {
		diskTotal, _ := strconv.ParseFloat(v["diskTotal"], 64)
		mount := v["path"]
		diskUsed, _ := strconv.ParseFloat(v["diskUsed"], 64)
		diskFree, _ := strconv.ParseFloat(v["diskFree"], 64)
		diskUsedPercent, _ := strconv.ParseFloat(v["diskUsedPercent"], 64)
		osdisk = OsDisk{
			NodeAddr:  ip,
			Mounted:   mount,
			TotalSize: diskTotal,
			UsedSize:  diskUsed,
			AvailSize: diskFree,
			UsedRate:  diskUsedPercent,
		}
		DB.Insert(&osdisk)

	}
}

//内存信息表插入新数据
func InserMem(DB orm.Ormer,ip string) {

	CpuInfo := GetCpuInfo()
	meminfo := GetMemInfo()
	cpucount := CpuInfo["CpuCount"]
	Load1 := CpuInfo["Load1"]
	Load5 := CpuInfo["Load5"]
	Load15 := CpuInfo["Load15"]
	memoryTotal := meminfo["MemoryTotal"]
	memoryUsed := meminfo["MemoryUsed"]
	memoryUsedPercent := meminfo["MemoryUsedPercent"]
	memoryCached := meminfo["MemoryCached"]
	swapTotal := meminfo["SwapTotal"]
	swapFree := meminfo["SwapFree"]
	var osmem = OsMem{
		NodeAddr:          ip,
		MemoryTotal:       memoryTotal,
		MemoryUsed:        memoryUsed,
		MemoryUsedPercent: memoryUsedPercent,
		MemoryCached:      memoryCached,
		SwapTotal:         swapTotal,
		SwapFree:          swapFree,
		CpuCount:          cpucount,
		Load1:             Load1,
		Load5:             Load5,
		Load15:            Load15,
	}
	DB.Insert(&osmem)
}

//磁盘历史信息表插入新数据
func InsertDiskHistoryInfo(DB orm.Ormer,Host ,Mounted  string,TotalSize ,UsedSize ,AvailSize ,UsedRate float64 ,createTime time.Time){
	var osdiskhistory=OsDiskHistory{
		NodeAddr:Host,
		Mounted:Mounted,
		TotalSize:TotalSize,
		UsedSize:UsedSize,
		AvailSize:AvailSize,
		UsedRate:UsedRate,
		CreateTime:createTime,
	}
	DB.Insert(&osdiskhistory)
}


//内存历史信息表插入新数据
func InsertMemHistoryInfo(DB orm.Ormer,Host string,MemoryTotal,MemoryUsed,MemoryUsedPercent,
MemoryCached,SwapTotal,SwapFree float64 ,createTime time.Time,CpuCount ,Load1,Load5,Load15 float64 ){
	var osmemhistory=OsMemHistory{
		NodeAddr:Host,
		MemoryTotal:MemoryTotal,
		MemoryUsed:MemoryUsed,
		MemoryUsedPercent:MemoryUsedPercent,
		MemoryCached:MemoryCached,
		SwapTotal:SwapTotal,
		SwapFree:SwapFree,
		CreateTime:createTime,
		CpuCount:CpuCount,
		Load1:Load1,
		Load5:Load5,
		Load15:Load15,
	}
	DB.Insert(&osmemhistory)
}
func GetDiskInfo()[]map[string]string{
	cmd :=exec.Command("df","-hP")
	out,_:=cmd.CombinedOutput()
	str:=string(out)
	var info []map[string]string
	ds:=make(map[string]string)
	st:=strings.Fields(str)
	for i,v :=range st{
		if i%6==0 && i !=0 && i!=6{
			d,_:=disk.Usage(v)
			ds = map[string]string{"path":d.Path,"diskTotal":strconv.Itoa(int(d.Total/1024/1024)),
				"diskUsed":strconv.Itoa(int(d.Used/1024/1024)),"diskFree":strconv.Itoa(int(d.Free/1024/1024)),
				"diskUsedPercent":strconv.Itoa(int(d.UsedPercent))}
			info=append(info,ds)
		}
	}
	return info
}
func GetMemInfo()map[string]float64 {
	memory := make(map[string]float64)
	m, _ := mem.VirtualMemory()
	MemoryTotal := float64(m.Total)
	Used := float64(m.Total - m.Free - m.Buffers - m.Cached)
	UsedPercent := float64((100 * Used / MemoryTotal))
	memory["memoryTotal"] = float64(m.Total / 1024 / 1024)
	memory["memoryUsed"] = float64((m.Total - m.Free - m.Buffers - m.Cached) / 1024 / 1024)
	memory["memoryUsedPercent"] = float64(UsedPercent * 100)
	memory["memoryCached"] = float64(m.Cached / 1024 / 1024)
	memory["SwapTotal"] = float64(m.SwapTotal / 1024 / 1024)
	memory["SwapFree"] = float64(m.SwapFree / 1024 / 1024)
	//info["memory"]=memory
	return memory
}

func GetCpuInfo() map[string]float64{

	var cp map[string]float64
	cp=make(map[string]float64)
	i,_:=cpu.Counts(true)
	cp["CpuCount"]=float64(i)
	loadinfo,_:=load.Avg()
	cp["Load1"]=loadinfo.Load1
	cp["Load5"]=loadinfo.Load5
	cp["Load15"]=loadinfo.Load15
	return  cp
}
func GetIntranetIp() string {
	addrs, err := net.InterfaceAddrs()
	var ip string
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip= ipnet.IP.String()
			}
		}
	}
	return ip
}

