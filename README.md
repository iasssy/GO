# IP Address Unique Counter

This Go program reads a large file containing a list of IP addresses and calculates the number of unique IP addresses in the file. The program efficiently partitions the file and uses a bitmap to count unique IPs, making it suitable for processing large files.

## Features

- Reads an input file containing IP addresses.
- Partitions the file into smaller chunks to process efficiently.
- Uses a bitmap to track unique IPs.
- Outputs the total count of unique IP addresses.

## How to Use

To use the program, you need to run the Go command with the path to your IP addresses file as an argument. 

```bash
go run main.go fileName
```

## Further improvements
- Temporary folder removal bug.
- Progress reporting improvements.
