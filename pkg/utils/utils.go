package utils

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"
)

// GetPreferredOutboundIP 获取本机最合适的 ipv4 地址
func GetPreferredOutboundIP() (net.IP, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	// 网卡名 => []ip
	tars := map[string][]net.IP{}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					tars[i.Name] = append(tars[i.Name], ipNet.IP)
				}
			}
		}
	}

	// 主网卡优先级顺序
	priority := []string{"eth0", "eth1", "wlan0", "wlan1"}

	for _, name := range priority {
		if ips, ok := tars[name]; ok && len(ips) > 0 {
			return ips[0], nil
		}
	}

	// enp 和 wlp 的, eg: "enp0s3", "enp0s8", "wlp2s0", "wlp3s0b1"
	for name, ips := range tars {
		if strings.HasPrefix(name, "enp") || strings.HasPrefix(name, "wlp") {
			if len(ips) > 0 {
				return ips[0], nil
			}
		}
	}

	// 其他的随意了
	for _, ips := range tars {
		if len(ips) > 0 {
			return ips[0], nil
		}
	}

	// 什么都没有就是本地回环
	return net.IPv4(127, 0, 0, 1), nil
}

func GetAvailablePort(start int, end int) (int, error) {
	for port := start; port <= end; port++ {
		ln, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
		if err != nil {
			continue // 如果此端口不可用，尝试下一个端口
		}
		ln.Close()
		//如果能够成功监听端口，我们需要立即关闭它，
		//以免在我们再次使用它时还处于连接状态
		return port, nil
	}
	return 0, fmt.Errorf("no available ports in the range %d-%d", start, end)
}

func GetHttpAddrString(useSSL bool, ip string, port int, path string) string {
	schema := "http"

	if useSSL {
		schema = "https"
	}

	base := fmt.Sprintf("%s://%s:%d", schema, ip, port)
	p, err := url.JoinPath(base, path)
	if err != nil {
		log.Printf("[GetHttpAddrString] error: %s. params: %v, %s, %d, %s", err, useSSL, ip, port, path)
		return base
	}

	return p
}
