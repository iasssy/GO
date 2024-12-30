package main

import (
	"fmt"
	"os"
	"bufio"
	"net"
	"hash/fnv"
)

const numHashTables = 16  

type HashTable struct {
	data map[string]bool
}

// initializes a new HashTable
func NewHashTable() *HashTable {
	return &HashTable{
		data: make(map[string]bool),
	}
}

func (ht *HashTable) Add(ip string) {
	ht.data[ip] = true
}

func (ht *HashTable) Contains(ip string) bool {
	_, exists := ht.data[ip]
	return exists
}

// for partitioning the IPs into different hash tables
func hash(ip string) int {
	h := fnv.New32a()
	h.Write([]byte(ip))
	return int(h.Sum32()) % numHashTables
}

func countUniqueIPs(filePath string) (int, error) {
	// slice of hash tables
	hashTables := make([]*HashTable, numHashTables)
	for i := 0; i < numHashTables; i++ {
		hashTables[i] = NewHashTable()
	}

	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	uniqueCount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ip := scanner.Text()
		if !isValidIP(ip) {
			continue
		}

		hashValue := hash(ip)
		hashTable := hashTables[hashValue]

		// if unique
		if !hashTable.Contains(ip) {
			hashTable.Add(ip)
			uniqueCount++
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("error reading file: %w", err)
	}

	return uniqueCount, nil
}

func isValidIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && parsedIP.To4() != nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <filePath>")
		return
	}
	filePath := os.Args[1]

	uniqueCount, err := countUniqueIPs(filePath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Number of unique IP addresses: %d\n", uniqueCount)
}
