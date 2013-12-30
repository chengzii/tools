package main
import(
	"fmt"
	"syscall"
	"runtime"
	"os"
	"mahonia"
)

type DiskStatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

type MemStatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
	Self uint64 `json:"self"`
}

func Mem() MemStatus {
	//自身占用
	memStat := new(runtime.MemStats)
	runtime.ReadMemStats(memStat)
	mem := MemStatus{}
	mem.Self = memStat.Alloc

	//系统占用,仅linux/mac下有效
	//system memory usage
	sysInfo := new(syscall.Sysinfo_t)
	err := syscall.Sysinfo(sysInfo)
	if err == nil {
		mem.All = sysInfo.Totalram/1024/1024
		//mem.All = sysInfo.Totalram * uint64(syscall.Getpagesize())
 		mem.Free = sysInfo.Freeram/1024/1024
 		//mem.Free = sysInfo.Freeram * uint64(syscall.Getpagesize())
		mem.Used = mem.All - mem.Free
	}
	return mem
}

func main() {
	os_lang:=os.Getenv("LANG")
	os_type:=runtime.GOOS
	os_cpunum:=runtime.NumCPU()
	os_df:=Df("/")
	os_mem:=Mem()
	encoder := mahonia.NewEncoder("GBK")
	if os_lang == "en_US.UTF-8"{
		fmt.Println("")
		fmt.Println("------------------系统信息---------------- ")
		fmt.Println("操作系统类型: ",os_type)
		fmt.Println("CPU逻辑核数 : ",os_cpunum)
		fmt.Println("系统语言环境: ",os_lang)
		fmt.Println("")
		fmt.Println("====>硬盘信息")
		fmt.Println("硬盘总大小  :",os_df.All," (单位:M)")
		fmt.Println("硬盘已用大小:",os_df.Used," (单位:M)")
		fmt.Println("硬盘剩余大小:",os_df.Free," (单位:M)")
		fmt.Println("***注意！该数值统计的是根目录('/')下挂载硬盘的使用情况***")
		fmt.Println("")
		fmt.Println("====>内存信息")
		fmt.Println("内存总大小  :",os_mem.All," (单位:M)")
		fmt.Println("内存已用大小:",os_mem.Used," (单位:M)")
		fmt.Println("内存剩余大小:",os_mem.Free," (单位:M)")
		fmt.Println("***注意！该数值是系统可内存分配的情况，不反映程序可使用内存的情况，程序可分配内存请使用  free -m 命令查询***")
		fmt.Println("--------------------------------------------------------------------------------------------------------------")
		fmt.Println("")
	}else{
		fmt.Println("")
		fmt.Println(encoder.ConvertString("------------------系统信息---------------- "))
		fmt.Println(encoder.ConvertString("操作系统类型: "),os_type)
		fmt.Println(encoder.ConvertString("CPU逻辑核数 : "),os_cpunum)
		fmt.Println(encoder.ConvertString("系统语言环境: "),os_lang)
		fmt.Println("")
		fmt.Println(encoder.ConvertString("====>硬盘信息"))
		fmt.Println(encoder.ConvertString("硬盘总大小  :"),os_df.All,encoder.ConvertString(" (单位:M)"))
		fmt.Println(encoder.ConvertString("硬盘已用大小:"),os_df.Used,encoder.ConvertString(" (单位:M)"))
		fmt.Println(encoder.ConvertString("硬盘剩余大小:"),os_df.Free,encoder.ConvertString(" (单位:M)"))
		fmt.Println(encoder.ConvertString("***注意！该数值统计的是根目录('/')下挂载硬盘的使用情况***"))
		fmt.Println("")
		fmt.Println(encoder.ConvertString("====>内存信息"))
		fmt.Println(encoder.ConvertString("内存总大小  :"),os_mem.All,encoder.ConvertString(" (单位:M)"))
		fmt.Println(encoder.ConvertString("内存已用大小:"),os_mem.Used,encoder.ConvertString(" (单位:M)"))
		fmt.Println(encoder.ConvertString("内存剩余大小:"),os_mem.Free,encoder.ConvertString(" (单位:M)"))
		fmt.Println(encoder.ConvertString("***注意！该数值是系统可内存分配的情况，不反映程序可使用内存的情况，程序可分配内存请使用  free -m 命令查询***"))
		fmt.Println("--------------------------------------------------------------------------------------------------------------")
		fmt.Println("")
		
	}
}

// disk usage of path/disk
func Df(path string) (disk DiskStatus) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)/1024/1024
	disk.Free = fs.Bfree * uint64(fs.Bsize)/1024/1024
	disk.Used = disk.All - disk.Free
	return
}
