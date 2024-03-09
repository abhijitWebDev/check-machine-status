package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
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

    // Create a new progress bar
    p := mpb.New(mpb.WithWidth(60))
    bar := p.AddBar(int64(len(hosts)),
        mpb.PrependDecorators(
            decor.CountersNoUnit("%d / %d", decor.WCSyncSpace),
        ),
        mpb.AppendDecorators(
            decor.Percentage(decor.WCSyncSpace),
        ),
    )

    for i, host := range hosts {
        upCount := 0
        for j := 0; j < 5; j++ {
            if host.URL != "" {
                resp, err := http.Get(host.URL)
                if err == nil && resp.StatusCode == 200 {
                    upCount++
                }
            } else if host.IP != "" {
                var cmd *exec.Cmd
                if runtime.GOOS == "windows" {
                    cmd = exec.Command("ping", "-n", "1", "-w", "1000", host.IP)
                } else {
                    cmd = exec.Command("ping", "-c", "1", "-W", "1", host.IP)
                }
                err := cmd.Run()
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

        // Increment the progress bar
        bar.Increment()
    }

    // Wait for the progress bar to finish
    p.Wait()

    timestamp := time.Now().Format("20060102_150405")
    filename := fmt.Sprintf("output_%s.xlsx", timestamp)
    if err := f.SaveAs(filename); err != nil {
        fmt.Println(err)
    }
}