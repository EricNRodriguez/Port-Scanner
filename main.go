package main

import (
	"golang.org/x/sync/semaphore"
	"net"
	"time"

	//"net"
	"os/exec"
	"strconv"
	"strings"
	//"time"
)

type PortNumbers []uint

type PortScanner struct {
	IP string
	sem *semaphore.Weighted

}

func (p *PortScanner) Scan(s, e int) (openPorts []int, err error) {
	if p.sem == nil {
		p.sem = semaphore.NewWeighted(1)
	}
	for i := s; i < e; i++ {
		p.sem.Acquire(nil, 1)
		go func(port int) {
			conn, _ := net.DialTimeout(p.IP, strconv.Itoa(port), 500*time.Millisecond)
			defer p.sem.Release(1)
			if conn != nil {
				defer conn.Close()
			}
			return
		}(i)
	}
	return
}

func (p *PortScanner) SetProcessLimit() {
	out, _ := exec.Command("ulimit", "-u").CombinedOutput()
	num, _ := strconv.Atoi(strings.TrimSpace(string(out)))
	p.sem = semaphore.NewWeighted(int64(num))
	return
}


func main() {
	sem := new(PortScanner)
	sem.SetProcessLimit()
	sem.Scan(0, 10)
	return
}

