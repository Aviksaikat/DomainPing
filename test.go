package main

import (
    "bufio"
    "flag"
    "fmt"
    "net"
    "net/http"
    "os"
)

func main() {
    // Define command-line flags
    domainFlag := flag.String("d", "", "single domain to check")
    fileFlag := flag.String("f", "", "file containing domains to check")
    ipFlag := flag.Bool("ip", false, "include resolved IP addresses in output")
    flag.Parse()

    // Read input domains
    domains := make([]string, 0)
    if *domainFlag != "" {
        domains = append(domains, *domainFlag)
    } else if *fileFlag != "" {
        file, err := os.Open(*fileFlag)
        if err != nil {
            fmt.Println("Error opening file:", err)
            os.Exit(1)
        }
        defer file.Close()

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            domains = append(domains, scanner.Text())
        }
    } else {
        scanner := bufio.NewScanner(os.Stdin)
        for scanner.Scan() {
            domains = append(domains, scanner.Text())
        }
    }

    // Check if each domain is alive
    for _, domain := range domains {
        httpURL := "http://" + domain
        httpsURL := "https://" + domain
        var ip string
        if *ipFlag {
            addrs, err := net.LookupIP(domain)
            if err != nil {
                ip = "unknown"
            } else {
                ip = addrs[0].String()
            }
        }
        fmt.Printf("Checking %s\n", domain)
        _, err := http.Get(httpURL)
        if err == nil {
            fmt.Printf("  HTTP: OK")
            if *ipFlag {
                fmt.Printf(" (%s)", ip)
            }
            fmt.Println()
        }
        _, err = http.Get(httpsURL)
        if err == nil {
            fmt.Printf("  HTTPS: OK")
            if *ipFlag {
                fmt.Printf(" (%s)", ip)
            }
            fmt.Println()
        }
    }
}
