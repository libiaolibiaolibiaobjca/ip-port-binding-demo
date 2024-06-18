package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

func main() {
	log.Println("用于示例，同一个端口，可以在多个网卡上监听，同时提供服务")

	ipv4 := listIpv4()
	for _, ip := range ipv4 {
		go func(ip string) {
			address := ip + ":8080"
			listener, err := net.Listen("tcp", address)
			if err != nil {
				log.Println(err)
			}
			log.Println("Serving on", address)

			err = http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				log.Println("Received request %+v", r)
				w.WriteHeader(200)
				w.Write([]byte("Hello, World!" + address))
				return
			}))
			if err != nil {
				log.Fatal(err)
			}
		}(ip)
	}

	select {}
}

func listIpv4() []string {
	ipv4s := make([]string, 0)
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip != nil && ip.To4() != nil {
				ipv4s = append(ipv4s, ip.String())
			}
		}
	}
	return ipv4s
}
