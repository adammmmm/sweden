package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func downloadAndFilterFile(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fileName := "delegated-ripencc-latest"
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	file.Seek(0, 0) // reset file pointer to beginning of file

	swedenFile, err := os.Create("sweden.txt")
	if err != nil {
		return err
	}
	defer swedenFile.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "SE") && strings.Contains(line, "ipv4") {
			parts := strings.Split(line, "|")
			if len(parts) >= 5 {
				ip := parts[3]
				numAddresses, err := strconv.Atoi(parts[4])
				if err != nil {
					continue
				}
				bits := int(math.Log2(float64(numAddresses)))
				cidr := fmt.Sprintf("%s/%d", ip, 32-bits)
				fmt.Fprintln(swedenFile, cidr)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func main() {
	url := "https://ftp.ripe.net/ripe/stats/delegated-ripencc-latest"
	if err := downloadAndFilterFile(url); err != nil {
		fmt.Println("Failed to download or filter file:", err)
	}
}
