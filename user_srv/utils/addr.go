package utils

import (
	"fmt"
	"net"
)

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, nil
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer func(l *net.TCPListener) {
		err_1 := l.Close()
		if err_1 != nil {

		}
	}(l)
	return l.Addr().(*net.TCPAddr).Port, nil
}

func main() {
	port, _ := GetFreePort()
	fmt.Println(port)
}
