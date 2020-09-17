package ipv6

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"strings"
)

//ipv6数据来源 http://ip.zxinc.org/

//IPv6 ipv6
type IPv6 struct {
	Country  string
	Province string
	City     string
	Info     string
}

func (c *IPv6) String() string {
	if c.Country+c.Province+c.City == "" {
		return ""
	}
	return c.Country + "," + c.Province + "," + c.City
}

//NewIPv6 ipv6的类
func NewIPv6(Country, Province, City, Info string) *IPv6 {
	return &IPv6{Country: Country, Province: Province, City: City, Info: Info}
}

//IP6toInt ipv6转int
func IP6toInt(IPv6Address net.IP) *big.Int {
	IPv6Int := big.NewInt(0)
	IPv6Int.SetBytes(IPv6Address.To16())
	return IPv6Int
}

//IPDBv6 ipv6
type IPDBv6 struct {
	firstIndex int64
	indexCount int64
	offlen     int64
	iplen      int64
	img        []byte
}

func (c *IPDBv6) getimg(offset int64, size int64) []byte {
	img := c.img[offset : offset+size]
	xx := make([]byte, size, size)
	for x := range img {
		xx[x] = img[x]
	}
	return xx
}
func (c *IPDBv6) init() {
	c.firstIndex = c.getLong8(16, 8)
	c.indexCount = c.getLong8(8, 8)
	c.offlen = c.getLong8(6, 1)
	c.iplen = c.getLong8(7, 1)
}
func (c *IPDBv6) find(ip int64, l int64, r int64) int64 {
	if r-l <= 1 {
		return l
	}
	m := (l + r) / 2
	o := c.firstIndex + m*(8+c.offlen)
	newip := c.getLong8(o, 8)
	if ip < newip {
		return c.find(ip, l, m)
	}
	return c.find(ip, m, r)
}

func (c *IPDBv6) getLong8(offset int64, size int64) int64 {
	s := c.getimg(offset, size)
	x := bytes.NewBuffer(s)
	var y int64
	var i int64
	for i = 0; i < 8-size; i++ {
		var xx byte
		x.WriteByte(xx)
	}
	binary.Read(x, binary.LittleEndian, &y)
	return y
}

//NewIPDBv6 ip初始化
func NewIPDBv6(path string) *IPDBv6 {
	// path := "/Users/tbwisk/Downloads/ip/ipv6wry.db"
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	if string(b[0:4]) != "IPDB" {
		panic("数据包格式错误")
	}
	item := &IPDBv6{img: b}
	item.init()
	return item
}
func (c *IPDBv6) getAddr(ipRecOff int64) *IPv6 {
	ixmg := c.img
	b := ixmg[ipRecOff]
	if b == 1 {
		return c.getAddr(c.getLong8(ipRecOff+1, c.offlen))
	}
	cArea := c.getAreaAddr(ipRecOff)
	ipv6 := c.formatType(cArea)
	if b == 2 {
		ipRecOff = ipRecOff + 1 + c.offlen
	} else {
		img := c.img[ipRecOff : ipRecOff+99]
		ipRecOff = int64(bytes.Index(img, []byte{0}) + 1)
	}
	aArea := c.getAreaAddr(ipRecOff) //这部分理论可移除处理
	ipv6.Info = aArea
	return ipv6
	// return c.getAreaAddr(ipRecOff)
}
func (c *IPDBv6) getString(offset int64) string {
	img := c.img[offset : offset+99]
	var temp []byte
	last := bytes.Index(img, []byte{0})
	if last != -1 {
		temp = img[0:last]
	} else {
		fmt.Println("error for ipv6 file")
	}
	value := string(temp)
	return value
}
func (c *IPDBv6) formatType(value string) *IPv6 {
	first := strings.Index(value, "国")
	var country string
	var province string
	var city string
	if first >= 0 {
		country = value[0 : first+3]
	}
	isCity := false
	second := strings.Index(value, "省")
	if second >= 0 {
		second = second + 3
		province = value[first+3 : second]
	} else {
		provinceSlice := []string{"内蒙古", "广西", "西藏", "宁夏", "新疆"}
		goverCity := []string{"北京市", "天津市", "上海市", "重庆市"}
		for i := 0; i < len(goverCity); i++ {
			second = strings.Index(value, goverCity[i])
			if second >= 0 {
				second = second + len(goverCity[i])
				province = value[first+3 : second]
				isCity = true
				break
			}
		}
		if isCity != true {
			for i := 0; i < len(provinceSlice); i++ {
				second = strings.Index(value, provinceSlice[i])
				if second >= 0 {
					second = second + len(provinceSlice[i])
					province = value[first+3 : second]
					break
				}
			}
		}
	}
	if isCity != true {
		_city1 := strings.Index(value, "市")
		_city2 := strings.Index(value, "州")
		fmt.Println("second", second, _city1, _city2, string(value))
		if _city1 >= 0 && _city2 >= 0 && _city1 < _city2 && second < _city1 {
			city = value[second : _city1+3]
		} else if _city1 >= 0 && _city2 >= 0 && _city1 > _city2 && second < _city2 {
			if _city1-_city2 == 3 {
				city = value[second : _city2+6]
			} else {
				city = value[second : _city2+3]
			}
		} else if _city1 >= 0 && second < _city1 {
			city = value[second : _city1+3]
		} else if _city2 >= 0 && second < _city2 {
			city = value[second : _city2+3]
		}
	}
	if country+province+city == "" {
		return NewIPv6(country, province, city, "")
	}
	// return country + "," + province + "," + city
	return NewIPv6(country, province, city, "")
}

func (c *IPDBv6) getAreaAddr(offset int64) string {
	b := c.img[offset]
	if b == 1 || b == 2 {
		p := c.getLong8(offset+1, c.offlen)
		return c.getAreaAddr(p)
	}
	return c.getString(offset)
}

//GetIPAddr 获取ip地址
func (c *IPDBv6) GetIPAddr(ip string) *IPv6 {
	ipv6Decimal := IP6toInt(net.ParseIP(ip))
	x := ipv6Decimal.Rsh(ipv6Decimal, 64)
	i := c.find(x.Int64(), 0, c.indexCount)
	ipOffset := c.firstIndex + i*(8+c.offlen)
	ipRecOff := c.getLong8(ipOffset+8, c.offlen)
	resp := c.getAddr(ipRecOff)
	return resp
}
