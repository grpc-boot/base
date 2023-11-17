package utils

import (
	"errors"
	"math"
	"net"
	"strconv"
	"strings"
)

var (
	ErrIpV4Address = errors.New(`invalid ip v4 address`)
	ErrNotFoundIp  = errors.New("not found ip address")
)

// Long2Ip uint32转换为Ip
func Long2Ip(ipVal uint32) string {
	var (
		buffer     strings.Builder
		a, b, c, d = strconv.Itoa(int(ipVal >> 24)), strconv.Itoa(int(ipVal << 8 >> 24)), strconv.Itoa(int(ipVal << 16 >> 24)), strconv.Itoa(int(ipVal << 24 >> 24))
	)

	buffer.Grow(len(a) + len(b) + len(c) + len(d) + 3)

	buffer.WriteString(a)
	buffer.WriteByte('.')
	buffer.WriteString(b)
	buffer.WriteByte('.')
	buffer.WriteString(c)
	buffer.WriteByte('.')
	buffer.WriteString(d)

	return buffer.String()
}

// Ip2Long Ip转换为uint32
func Ip2Long(ip string) (ipVal uint32, err error) {
	var (
		val     uint32
		start   = 0
		leftMax = 4 * 8
	)

	for index, ch := range ip {
		if ch >= '0' && ch <= '9' {
			continue
		}

		if ch == '.' && start != index && (index-start < 4) {
			for i := index - start; i > 0; i-- {
				val += uint32(ip[index-i]-'0') * uint32(math.Pow(10, float64(i-1)))
			}
			if val > 0xff {
				return 0, ErrIpV4Address
			}

			leftMax -= 8
			if leftMax < 8 {
				return 0, ErrIpV4Address
			}

			ipVal += val << leftMax
			start = index + 1
			val = 0
			continue
		}

		return 0, ErrIpV4Address
	}

	//长度过长或过短判断
	if leftMax != 8 {
		return 0, ErrIpV4Address
	}

	for i := len(ip) - start; i > 0; i-- {
		val += uint32(ip[len(ip)-i]-'0') * uint32(math.Pow(10, float64(i-1)))
	}
	if val > 0xff {
		return 0, ErrIpV4Address
	}
	ipVal += val

	return ipVal, nil
}

// LocalIp 获取本机Ip
func LocalIp() (ip string, err error) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String(), nil
					}
				}
			}
		}
	}

	return "", ErrNotFoundIp
}
