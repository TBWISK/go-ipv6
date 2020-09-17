package ipv6

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

//TestInit 测试初始化
func TestInit(t *testing.T) {
	ipv6 := NewIPDBv6("./ipv6wry.db")
	addr := ipv6.GetIPAddr("2001:250:208:5809:8052:eb04:8087:fcc1")
	fmt.Println(addr)
}
func readipfile() []string {
	path := "./test_ipv6_10w.txt"
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	xs := string(b)
	xx := strings.Split(xs, "\n")
	var resp []string
	for idx := range xx {
		resp = append(resp, strings.Trim(xx[idx], " "))
	}
	return resp
}
func TestIPv6(t *testing.T) {
	items := readipfile()
	ipv6 := NewIPDBv6("./ipv6wry.db")
	results := make(map[string]int64, 100)
	for idx := range items {
		addrx := ipv6.GetIPAddr(items[idx])
		addr := addrx.String()
		cnt, ok := results[addr]
		if ok == false {
			results[addr] = 1
		} else {
			results[addr] = cnt + 1
		}
	}
	fmt.Println(results)
}
