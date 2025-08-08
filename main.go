package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
)

func worker(ports, results chan int) { // 01
	for port := range ports {
		address := fmt.Sprintf("scanme.nmap.org:%d", port)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- port
	}
}

type ScanFlags struct {
	hosts string
	ports []int
}

func parseFlags() ScanFlags {
	var flags ScanFlags

	flag.StringVar(&flags.hosts, "host", "127.0.0.1", "Host to scan") // ! I need to fix these. I need the ability to pass multiple hosts and ports.
	flag.IntVar(flags.ports, "port", 80, "Port to scan")

	flag.Parse()
	return flags
}

// parsePortList parses a comma-separated list of ports and returns a slice of valid port numbers.
// It returns an error if no valid ports are provided or if any port is out of range
func parsePortList(portsString string) ([]int, error) {
	var portNums []int
	portSet := make(map[int]struct{})
	portList := strings.Split(portsString, ",")
	for _, port := range portList {
		var portNum int
		// port range functionality
		if strings.Contains(port, "-") {
			portsRange := strings.SplitN(port, "-", 2)
			if len(portsRange) != 2 {
				continue
			}
			startPort, err := strconv.Atoi(strings.TrimSpace((portsRange[0])))
			if err != nil || startPort < 1 || startPort > 65535 {
				continue
			}
			endPort, err := strconv.Atoi(strings.TrimSpace((portsRange[1])))
			if err != nil || endPort < 1 || endPort > 65535 {
				continue
			}
			if startPort > endPort {
				tmp := startPort
				startPort = endPort
				endPort = tmp
			}
			for i := startPort; i <= endPort; i++ {
				if _, exists := portSet[i]; !exists {
					portSet[i] = struct{}{}
					portNums = append(portNums, i)
				}
			}
		} else {
			port = strings.TrimSpace(port)
			portNum, err := strconv.Atoi(port)
			if err != nil || port < 1 || port > 65535 {
				continue
			}
			if _, exists := portSet[portNum]; !exists {
				portSet[portNum] = struct{}{}
				portNums = append(portNums, portNum)
			}
		}
	}
	if len(portNums) == 0 {
		return nil, errors.New("no valid ports provided")
	}
	return portNums, nil
}

// parseHostList parses a comma-separated list of hosts and returns a slice of unique hostnames.
// It returns an error if no hosts are provided or if any host is empty.
func parseHostList(hostList string) ([]string, error) {
	hostSlice := []string{}
	hostSet := make(map[string]struct{})

	if hostList == "" {
		return nil, errors.New("no hosts provided")
	}
	hosts := strings.Split(hostList, ",")
	for _, host := range hosts {
		host = strings.TrimSpace(strings.ToLower(host))
		if host == "" {
			continue
		}
		if _, exists := hostSet[host]; !exists {
			hostSet[host] = struct{}{}
			hostSlice = append(hostSlice, host)
		}
	}
	return hostSlice, nil
}

/*



 */

func main() {
	flags := parseFlags()

	hosts := flags.hosts
	ports := flags.ports

	// address := fmt.Sprintf("%v:%v", host, port)
	// conn, err := net.Dial("tcp", address)
	// if err != nil {
	// 	conn.Close()
	// }
	// conn.Close()
	// fmt.Printf("Port %v open on %v\n", port, host)

	ports := make(chan int, 100)
	results := make(chan int)
	var openPorts []int

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}

	go func() {
		for i := 1; i <= 1024; i++ {
			ports <- i
		}
	}()

	for i := 0; i < 1024; i++ {
		port := <-results
		if port != 0 {
			openPorts = append(openPorts, port)
		}
	}

	close(ports)
	close(results)
	sort.Ints(openPorts)
	for _, port := range openPorts {
		fmt.Printf("%d open\n", port)
	}
}
