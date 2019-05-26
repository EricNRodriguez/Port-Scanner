package main

import (
	"bytes"
	"flag"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

type PortScanner struct {
	ip 			string
	timeOut 	time.Duration
	sem chan 	int
}

func (p *PortScanner) Start(s, e int) (openPorts []int, err error) {
	openPortsChan := make(chan int)
	wg := &sync.WaitGroup{}
	go func() {
		for {
			openPorts = append(openPorts, <- openPortsChan)
		}
	}()
	for i := s; i < e; i++ {
		p.sem <- 1
		wg.Add(1)
		go p.scanPort(i, openPortsChan)
		<- p.sem
		wg.Done()
	}
	wg.Wait()
	p.printFormattedData(openPorts)
	return
}

func (p *PortScanner) scanPort(port int, open chan <- int) {
	if conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", p.ip, strconv.Itoa(port)), p.timeOut*time.Millisecond); err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(p.timeOut * time.Millisecond)
			p.scanPort(port, open)
		}
	} else {
		conn.Close()
		open <- port
	}
	return
}

func (p *PortScanner) printFormattedData(openPorts []int) {
	output := bytes.Buffer{}
	output.WriteString(fmt.Sprintf("OPEN PORTS - %s \n-----------------------\n", p.ip))
	for _, p := range openPorts {
		output.WriteString(fmt.Sprintf(" + %d \n", p))
	}
	fmt.Println(output.String())
	return
}

func ProcessLimit() int64 {
	out, _ := exec.Command("ulimit", "-u").CombinedOutput()
	num, _ := strconv.Atoi(strings.TrimSpace(string(out)))
	return int64(num)
}


func main() {
	var (
		ip 		= flag.String("ip", "127.0.0.1", "ip address")
		timeOut = flag.Int("t", 1000, "tcp dial timeout (milliseconds)")
		start 	= flag.Int("s", 0, "port scan range lower bound (inclusive)")
		end 	= flag.Int("e", 65535, "port scan range upper bound (exclusive)")
	)
	flag.Parse()
	scanner := &PortScanner{
		*ip,
		time.Duration(*timeOut),
		make(chan int, ProcessLimit()),
	}
	scanner.Start(*start, *end)
}

