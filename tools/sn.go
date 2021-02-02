package tools

import (
	"net"
	"fmt"
	"os/exec"
	"regexp"

	//"syscall"
)
/**
 * 获取电脑CPUId
 */
func GetCpuId() string {
	cmd := exec.Command("wmic", "cpu", "get", "ProcessorID")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(string(out))
	str := string(out)
	//匹配一个或多个空白符的正则表达式
	reg := regexp.MustCompile("\\s+")
	str = reg.ReplaceAllString(str, "")
	return str[11:]
}
//
////MAC地址:
//func GetMac() {
//	interfaces, err :=  net.Interfaces()
//	if err != nil {
//		panic("Poor soul, here is what you got: " + err.Error())
//	}
//
//	for _, inter := range interfaces {
//		fmt.Println(inter.Name, inter.HardwareAddr)
//	}
//}
//MAC地址,默认取第一条
func GetMac() string{
	// 获取本机的MAC地址
	interfaces, err := net.Interfaces()
	if err != nil {
		panic("Poor soul, here is what you got: " + err.Error())
	}
	//for _, inter := range interfaces {
	//fmt.Println(inter.Name)
	inter := interfaces[0]
	mac := inter.HardwareAddr.String() //获取本机MAC地址
	//fmt.Println("MAC = ", mac)

	return mac
	//}
}

////硬盘ID
//func GetDisk() {
//	var st syscall.Stat_t
//	err := syscall.Stat("/dev/disk0", &st)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Printf("%+v", st)
//}
