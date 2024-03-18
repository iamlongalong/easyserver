package utils

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"
)

// GetPreferredOutboundIP 获取本机最合适的 ipv4 地址
func GetPreferredOutboundIP() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP, nil
			}
		}
	}

	return nil, fmt.Errorf("no suitable IPV4 address found")
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
