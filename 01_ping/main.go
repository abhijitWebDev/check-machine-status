package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type Host struct {
    IP  string
    URL string
}

func main() {
    var hosts []Host

    file, err := os.Open("hosts.txt")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        parts := strings.Split(line, " ")
        host := Host{}
        if len(parts) >= 1 {
            host.IP = parts[0]
        }
        if len(parts) >= 2 {
            host.URL = parts[1]
        }
        hosts = append(hosts, host)
    }

    if err := scanner.Err(); err != nil {
        fmt.Println(err)
        return
    }

    f := excelize.NewFile()
    f.SetCellValue("Sheet1", "A1", "IP")
    f.SetCellValue("Sheet1", "B1", "URL")
    f.SetCellValue("Sheet1", "C1", "Status")

    // Set column widths
    f.SetColWidth("Sheet1", "A", "A", 30)
    f.SetColWidth("Sheet1", "B", "B", 30)
    f.SetColWidth("Sheet1", "C", "C", 30)

    for i, host := range hosts {
        upCount := 0
        for j := 0; j < 5; j++ {
            if host.URL != "" {
                resp, err := http.Get(host.URL)
                if err == nil && resp.StatusCode == 200 {
                    upCount++
                }
            } else if host.IP != "" {
                _, err := net.DialTimeout("tcp", net.JoinHostPort(host.IP, "80"), time.Second)
                if err == nil {
                    upCount++
                }
            }
            time.Sleep(1 * time.Second) // wait for a second before the next request
        }
        if upCount > 0 {
            f.SetCellValue("Sheet1", fmt.Sprintf("C%d", i+2), "Up")
        } else {
            f.SetCellValue("Sheet1", fmt.Sprintf("C%d", i+2), "Down")
        }
        f.SetCellValue("Sheet1", fmt.Sprintf("A%d", i+2), host.IP)
        f.SetCellValue("Sheet1", fmt.Sprintf("B%d", i+2), host.URL)
    }

    timestamp := time.Now().Format("20060102_150405")
    filename := fmt.Sprintf("output_%s.xlsx", timestamp)
    if err := f.SaveAs(filename); err != nil {
        fmt.Println(err)
    }
}