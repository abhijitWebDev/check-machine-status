package main

import (
	"bufio"
	"fmt"
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
        if len(parts) != 2 {
            fmt.Println("Invalid line:", line)
            continue
        }
        hosts = append(hosts, Host{IP: parts[0], URL: parts[1]})
    }

    if err := scanner.Err(); err != nil {
        fmt.Println(err)
        return
    }

    f := excelize.NewFile()
    f.SetCellValue("Sheet1", "A1", "IP")
    f.SetCellValue("Sheet1", "B1", "URL")
    f.SetCellValue("Sheet1", "C1", "Status")

    for i, host := range hosts {
        upCount := 0
        for j := 0; j < 5; j++ {
            resp, err := http.Get(host.URL)
            if err == nil && resp.StatusCode == 200 {
                upCount++
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

    if err := f.SaveAs("output.xlsx"); err != nil {
        fmt.Println(err)
    }
}