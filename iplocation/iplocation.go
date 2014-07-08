package iplocation

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	path = "qqwry.dat"
)

type Location struct {
	Country string
	Area    string
}

var (
	f           *os.File
	ip          int64
	test        string = "202.171.253.103"
	count       int64  = 0
	first_index int64  = 0
	last_index  int64  = 0
)

func init() {

}

func main() {
	f, _ = os.OpenFile(path, os.O_RDONLY, 666)
	defer f.Close()
	buf := make([]byte, 8)
	f.Read(buf)

	first_index = byte2int64(buf[:4])
	last_index = byte2int64(buf[4:])
	count = (last_index-first_index)/7 + 1
	s := GetIpLocation(test)
	fmt.Println(s)

}
func GetIpLocation(ipstr string) string {
	ip = sting2ip(ipstr)
	index := find(ip, 0, count-1)
	ioffset := first_index + index*7
	aoffset := getLong3(ioffset + 4)
	address := getAddr(aoffset)
	return address
}

func getAddr(offset int64) string {
	f.Seek(offset+4, 0)
	buf := make([]byte, 1)
	countryAddr := ""
	areaAddr := ""
	f.Read(buf)
	if buf[0] == 0x01 {
		countryOffset := getLong3(0)
		f.Seek(offset, 0)
		b := make([]byte, 1)
		if b[0] == 0x02 {
			countryAddr = getString(getLong3(0))
			f.Seek(countryOffset+4, 0)
		} else {
			countryAddr = getString(countryOffset)
			areaAddr = getAreaAddr(0)
		}

	} else if buf[0] == 0x02 {
		countryAddr = getString(getLong3(0))
		areaAddr = getAreaAddr(offset + 8)
	} else {
		countryAddr = getString(offset + 4)
		areaAddr = getAreaAddr(0)

	}
	return countryAddr + "/" + areaAddr
}
func getAreaAddr(offset int64) string {
	var result string
	if offset > 0 {
		f.Seek(offset, 0)
	}
	buf := make([]byte, 1)
	f.Read(buf)

	if buf[0] == 0x01 || buf[0] == 0x02 {
		p := getLong3(0)
		if p > 0 {
			result = getString(p)
		} else {
			result = ""
		}
	} else {
		result = getString(offset)
	}
	return result
}
func getString(offset int64) string {
	if offset > 0 {
		f.Seek(offset, 0)
	}
	buf := make([]byte, 1)
	f.Read(buf)
	str := ""
	for buf[0] != 0 {
		str = str + string(buf)
		buf = make([]byte, 1)
		f.Read(buf)
	}
	return str
}
func find(ip, left, right int64) int64 {
	var result int64
	if right-left == 1 {
		result = left
	} else {
		middle := (left + right) / 2
		offset := first_index + middle*7
		f.Seek(offset, 0)
		buf := make([]byte, 4)
		f.Read(buf)
		new_ip := byte2int64(buf)
		if ip <= new_ip {
			result = find(ip, left, middle)
		} else {
			result = find(ip, middle, right)
		}
	}
	return result
}
func getLong3(offset int64) int64 {
	if offset > 0 {
		f.Seek(offset, 0)
	}
	buf := make([]byte, 3)
	f.Read(buf)
	return byte2int64(buf)
}

func sting2ip(str string) int64 {
	ss := strings.Split(str, ".")
	var ip int64 = 0
	for i := 0; i < len(ss); i++ {
		ip = ip << 8
		s, _ := strconv.ParseInt(ss[i], 10, 0)
		ip = ip + s
	}
	return int64(ip)
}

func byte2int64(ss []byte) int64 {
	var tem int64 = 0
	//文件读取 逆序
	for i := len(ss) - 1; i >= 0; i-- {
		tem = tem << 8
		s := int64(ss[i])
		tem = tem + s
	}
	return tem
}
