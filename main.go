package main

import (
	"os"
	"fmt"
	"net"
	"time"
	"bufio"
	"path/filepath"
)

const partitions = 256
const bitmapSize = 1 << 24

func partitionFile(filePath string) ([]string, error) {
	// temporary folder 
	tempDir := fmt.Sprintf("temp_%s", time.Now().Format("20060102150405"))
	err := os.Mkdir(tempDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary folder: %w", err)
	}

	// partition files
	partitionFiles := make([]*os.File, partitions)
	partitionPaths := make([]string, partitions)

	for i := 0; i < partitions; i++ {
		path := filepath.Join(tempDir, fmt.Sprintf("partition_%d.txt", i))
		partitionPaths[i] = path
		file, err := os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("failed to create partition file: %w", err)
		}
		partitionFiles[i] = file
	}

	
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ip := scanner.Text()
		ipUint32, err := ipToUint32(ip)
		if err != nil {
			continue
		}
	
		firstOctet := (ipUint32 >> 24) & 0xFF
	
		if firstOctet >= partitions {
			continue
		}
	
		partitionFiles[firstOctet].WriteString(ip + "\n")
	}
	
	for _, file := range partitionFiles {
		file.Close()
	}

	/*
	// for some reason cannot delete the empty folder (the partition files are already removed at this point)
	err = os.RemoveAll(tempDir)
	if err != nil {
		return nil, fmt.Errorf("failed to remove temporary folder: %w", err)
	}
	*/
	

	
	return partitionPaths, nil

}



func ipToUint32(ip string) (uint32, error) {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return 0, fmt.Errorf("invalid IP address: %s", ip)
	}
	ipv4 := parsedIP.To4()
	if ipv4 == nil {
		return 0, fmt.Errorf("not an IPv4 address: %s", ip)
	}
	return uint32(ipv4[0])<<24 | uint32(ipv4[1])<<16 | uint32(ipv4[2])<<8 | uint32(ipv4[3]), nil
}

func countUniqueIPsInPartition(filePath string) (int, error) {
	bitmap := make([]bool, bitmapSize)
	count := 0

	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to open partition file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ip, err := ipToUint32(scanner.Text())
		if err != nil {
			continue
		}

		index := ip % uint32(bitmapSize)
		if !bitmap[index] {
			bitmap[index] = true
			count++
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("error reading partition file: %w", err)
	}

	return count, nil

}

func countUniqueIPs(filePath string) (int, error) {
	partitionPaths, err := partitionFile(filePath)
	if err != nil {
		return 0, err
	}

	totalUnique := 0
	for _, path := range partitionPaths {
		count, err := countUniqueIPsInPartition(path) 
		if err != nil {
			return 0, err
		}
		totalUnique += count
		_ = os.Remove(path)
	}

	return totalUnique, nil
}

func main() {

	// if using command like  "go run main.go input.txt":
	filePath := os.Args[1]

	/* 
	filePath := "input.txt"
	*/

	uniqueCount, err := countUniqueIPs(filePath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Number of unique IP addresses: %d\n", uniqueCount)

}