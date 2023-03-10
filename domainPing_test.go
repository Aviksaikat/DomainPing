package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
    "sync"
    "strings"
    "bufio"
    "encoding/json"
	"math/rand"
	"github.com/fatih/color"
)

type domainResult struct {
	Host string `json:"host"`
	IP   string `json:"ip,omitempty"`
}



func getRandomColour() (*color.Color){

	rand.Seed(time.Now().UnixNano())
	colors := []*color.Color{
		color.New(color.FgRed),
		color.New(color.FgGreen),
		color.New(color.FgYellow),
		color.New(color.FgBlue),
		color.New(color.FgMagenta),
		color.New(color.FgCyan),
		color.New(color.FgWhite),
		color.New(color.FgHiRed),
		color.New(color.FgHiGreen),
		color.New(color.FgHiYellow),
		color.New(color.FgHiBlue),
		color.New(color.FgHiMagenta),
		color.New(color.FgHiCyan),
		color.New(color.FgHiWhite),
		color.New(color.Bold),
		color.New(color.Faint),
		color.New(color.Italic),
		color.New(color.Underline),
		color.New(color.BlinkSlow),
		color.New(color.BlinkRapid),
		color.New(color.ReverseVideo),
		color.New(color.Concealed),
		color.New(color.CrossedOut),
		color.New(color.FgHiBlack),
	}
	randomColor := colors[rand.Intn(len(colors))]
	return randomColor
}

func printBanner() {
	const banner1 = `
██████╗  ██████╗ ███╗   ███╗ █████╗ ██╗███╗   ██╗██████╗ ██╗███╗   ██╗ ██████╗ 
██╔══██╗██╔═══██╗████╗ ████║██╔══██╗██║████╗  ██║██╔══██╗██║████╗  ██║██╔════╝ 
██║  ██║██║   ██║██╔████╔██║███████║██║██╔██╗ ██║██████╔╝██║██╔██╗ ██║██║  ███╗
██║  ██║██║   ██║██║╚██╔╝██║██╔══██║██║██║╚██╗██║██╔═══╝ ██║██║╚██╗██║██║   ██║
██████╔╝╚██████╔╝██║ ╚═╝ ██║██║  ██║██║██║ ╚████║██║     ██║██║ ╚████║╚██████╔╝
╚═════╝  ╚═════╝ ╚═╝     ╚═╝╚═╝  ╚═╝╚═╝╚═╝  ╚═══╝╚═╝     ╚═╝╚═╝  ╚═══╝ ╚═════╝
`
	const banner2 = `
██████   ██████  ███    ███  █████  ██ ███    ██ ██████  ██ ███    ██  ██████  
██   ██ ██    ██ ████  ████ ██   ██ ██ ████   ██ ██   ██ ██ ████   ██ ██       
██   ██ ██    ██ ██ ████ ██ ███████ ██ ██ ██  ██ ██████  ██ ██ ██  ██ ██   ███ 
██   ██ ██    ██ ██  ██  ██ ██   ██ ██ ██  ██ ██ ██      ██ ██  ██ ██ ██    ██ 
██████   ██████  ██      ██ ██   ██ ██ ██   ████ ██      ██ ██   ████  ██████ 
`
	col := getRandomColour()	
	col.Println(banner1)

	col = getRandomColour()
	col.Println("\nAuthor: avik_saikat")
	col.Println("Github: https://github.com/aviksaikat")
	col.Println("Gitlab: https://gitlab.com/aviksaikat")
	fmt.Println("\n----------------------------------------------------------------\n")
}

func checkProtocol(domain, protocol, ip string, results *[]string, wg *sync.WaitGroup) {
	defer wg.Done()

	client := http.Client{
		Timeout: time.Second * 5,
	}

	url := fmt.Sprintf("%s://%s", protocol, domain)
	resp, err := client.Get(url)
	if err == nil && resp.StatusCode == 200 {
		var result string
		if ip != "" {
			result = fmt.Sprintf("%s [%s]", url, ip)
		} else {
			result = url
		}
		*results = append(*results, result)
	}
}


func checkAlive(domain string, resolveIP bool) []string {
	fmt.Printf("Checking if %s is alive...\n", domain)

	var results []string
	var wg sync.WaitGroup

	if resolveIP {
		addrs, err := net.LookupIP(domain)
		if err != nil {
			fmt.Println("Error looking up IP address:", err)
			os.Exit(1)
		}
		for _, addr := range addrs {
			ip := addr.String()
			wg.Add(1)
			go checkProtocol(domain, "http", ip, &results, &wg)
			go checkProtocol(domain, "https", ip, &results, &wg)
		}
	} else {
		wg.Add(1)
		go checkProtocol(domain, "http", "", &results, &wg)
		go checkProtocol(domain, "https", "", &results, &wg)
	}

	wg.Wait()

	return results
}


func saveResultsToFile(outputFile string, results []string, includeIP bool) {
    printResults(results)
	
    file, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		os.Exit(1)
	}
	defer file.Close()

	var formattedResults []domainResult

	for _, result := range results {
		host := strings.Split(result, " ")[0]
		ip := ""

		if includeIP {
			ip = strings.Split(result, " ")[1][1 : len(strings.Split(result, " ")[1])-1]
		}

		formattedResults = append(formattedResults, domainResult{host, ip})
	}

	jsonResults, err := json.MarshalIndent(formattedResults, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling results as JSON:", err)
		os.Exit(1)
	}

	_, err = file.Write(jsonResults)
	if err != nil {
		fmt.Println("Error writing to output file:", err)
		os.Exit(1)
	}

	fmt.Printf("Results saved to %s\n", outputFile)
    
}


func printResults(results []string) {
	for _, result := range results {
		fmt.Println(result)
	}
}

func main() {
    var domain string
    var domains []string
    var domains_file string
    var resolveIP bool
    var outputFile string
    var results []string

    flag.StringVar(&domain, "d", "", "Specify the domain to check")
    flag.StringVar(&domains_file, "l", "", "List of domains to check")
    flag.BoolVar(&resolveIP, "ip", false, "Resolve IP address")
    flag.StringVar(&outputFile, "o", "", "Output file path")

    flag.Parse()

    if domain == "" && domains_file == "" {
		fmt.Println("Please provide either -d or -l option")
		os.Exit(1)
	}

    // read data from the file
    if domains_file != "" { 
        file, err := os.Open(domains_file)
        if err != nil {
            fmt.Println("Error opening input file:", err)
            os.Exit(1)
        }
        defer file.Close()


        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            domains = append(domains, scanner.Text())
        }
        for _, domain := range domains {
            results = append(results, checkAlive(domain, resolveIP)...)
        }
    } else {
        results = append(results, checkAlive(domain, resolveIP)...)
    }

    if outputFile != "" {
        saveResultsToFile(outputFile, results, resolveIP)
    } else {
        printResults(results)
    }
}