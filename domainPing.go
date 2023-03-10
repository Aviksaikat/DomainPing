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
	banners := []string {`
██████╗  ██████╗ ███╗   ███╗ █████╗ ██╗███╗   ██╗██████╗ ██╗███╗   ██╗ ██████╗ 
██╔══██╗██╔═══██╗████╗ ████║██╔══██╗██║████╗  ██║██╔══██╗██║████╗  ██║██╔════╝ 
██║  ██║██║   ██║██╔████╔██║███████║██║██╔██╗ ██║██████╔╝██║██╔██╗ ██║██║  ███╗
██║  ██║██║   ██║██║╚██╔╝██║██╔══██║██║██║╚██╗██║██╔═══╝ ██║██║╚██╗██║██║   ██║
██████╔╝╚██████╔╝██║ ╚═╝ ██║██║  ██║██║██║ ╚████║██║     ██║██║ ╚████║╚██████╔╝
╚═════╝  ╚═════╝ ╚═╝     ╚═╝╚═╝  ╚═╝╚═╝╚═╝  ╚═══╝╚═╝     ╚═╝╚═╝  ╚═══╝ ╚═════╝
`,
`
██████   ██████  ███    ███  █████  ██ ███    ██ ██████  ██ ███    ██  ██████  
██   ██ ██    ██ ████  ████ ██   ██ ██ ████   ██ ██   ██ ██ ████   ██ ██       
██   ██ ██    ██ ██ ████ ██ ███████ ██ ██ ██  ██ ██████  ██ ██ ██  ██ ██   ███ 
██   ██ ██    ██ ██  ██  ██ ██   ██ ██ ██  ██ ██ ██      ██ ██  ██ ██ ██    ██ 
██████   ██████  ██      ██ ██   ██ ██ ██   ████ ██      ██ ██   ████  ██████ 
`,
	
	}
	col := getRandomColour()
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(banners))
	col.Println(banners[randomIndex])

	col = getRandomColour()
	col.Println("\nAuthor: avik_saikat")
	col.Println("Github: https://github.com/aviksaikat")
	col.Println("Gitlab: https://gitlab.com/aviksaikat")
	fmt.Println("\n----------------------------------------------------------------\n")
}

func showHelp() {
	printBanner()
    fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] INPUT_FILE\n", os.Args[0])
    fmt.Fprintln(os.Stderr, "Ping a list of domains and save the results to a file.")
    fmt.Fprintln(os.Stderr, "\nOptions:")
	flag.PrintDefaults()
	os.Exit(0)
}



func checkAlive(domain string, resolveIP bool) []string {
    // fmt.Printf("Checking if %s is alive...\n", domain)

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
            go func() {
                defer wg.Done()
                if err := checkProtocol(domain, "http", ip, &results); err != nil {
                    fmt.Println(err)
                }
                if err := checkProtocol(domain, "https", ip, &results); err != nil {
                    fmt.Println(err)
                }
            }()
        }
    } else {
        wg.Add(2)
        go func() {
            defer wg.Done()
            if err := checkProtocol(domain, "http", "", &results); err != nil {
                fmt.Println(err)
            }
        }()
        go func() {
            defer wg.Done()
            if err := checkProtocol(domain, "https", "", &results); err != nil {
                fmt.Println(err)
            }
        }()
    }

    wg.Wait()

    return results
}

func checkProtocol(domain, protocol, ip string, results *[]string) error {
    client := http.Client{
        Timeout: time.Second * 5,
    }

    url := fmt.Sprintf("%s://%s", protocol, domain)
    resp, err := client.Get(url)
    if err != nil {
        return fmt.Errorf("Error checking %s: %s", url, err)
    }
    if resp.StatusCode == 200 {
        var result string
        if ip != "" {
            result = fmt.Sprintf("%s [%s]", url, ip)
        } else {
            result = url
        }
        *results = append(*results, result)
    }
    return nil
}



func saveResultsToFile(outputFile string, results []string, includeIP bool) string {
    printResults(results)
	
	var formattedResults []domainResult
	var outputBytes []byte
    var outputFormat string

    file, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        fmt.Println("Error opening output file:", err)
        os.Exit(1)
    }
    defer file.Close()

    
    for _, result := range results {
        host := strings.Split(result, " ")[0]
        ip := ""

        if includeIP {
            ip = strings.Split(result, " ")[1][1 : len(strings.Split(result, " ")[1])-1]
        }

        formattedResults = append(formattedResults, domainResult{host, ip})
    }

    if strings.HasSuffix(outputFile, ".json") {
        jsonResults, err := json.MarshalIndent(formattedResults, "", "    ")
        if err != nil {
            fmt.Println("Error marshaling results as JSON:", err)
            os.Exit(1)
        }
        outputBytes = jsonResults
        outputFormat = "JSON"
    } else {
        textResults := make([]string, 0, len(formattedResults))
        for _, result := range formattedResults {
            if includeIP {
                textResults = append(textResults, fmt.Sprintf("%s\t%s\n", result.Host, result.IP))
            } else {
				// fmt.Println(result.Host)
                textResults = append(textResults, fmt.Sprintf("%s\n", result.Host))
            }
        }
        outputBytes = []byte(strings.Join(textResults, "\n"))
        outputFormat = "plain text"
    }

    _, err = file.Write(outputBytes)
    if err != nil {
        fmt.Println("Error writing to output file:", err)
        os.Exit(1)
    }

    // fmt.Printf("Results saved to %s (%s format)\n", outputFile, outputFormat)
	return outputFormat
}


func printResults(results []string) {
	for _, result := range results {
		fmt.Println(result)
	}
}

func main() {
	var domain string
	var outputFile string
	var inputFilePath string
	var outputFormat string
	var resolveIP bool
	var banner bool
	var help bool

	flag.StringVar(&domain, "d", "", "Specify the domain to check")
	flag.BoolVar(&resolveIP, "ip", false, "Resolve IP address")
	flag.StringVar(&outputFile, "o", "", "Output file path")
	flag.StringVar(&inputFilePath, "f", "", "Input file path")
	flag.BoolVar(&banner, "banner", false, "Print banner")
	flag.BoolVar(&help, "h", false, "show help")

	flag.Parse()

	// check if there is any input from stdin
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		if (domain == "" && inputFilePath == "") {
			showHelp()
		}
	}

	if help {
		showHelp()
	}

	if banner != false {
		printBanner()
		os.Exit(0)
	}

	printBanner()

	if (domain == "" && inputFilePath == "") {
		scanner := bufio.NewScanner(os.Stdin)
        for scanner.Scan() {
            domain := scanner.Text()
			results := checkAlive(domain, resolveIP)

			if outputFile != "" {
				outputFormat = saveResultsToFile(outputFile, results, resolveIP)
				// fmt.Printf("Results saved to %s (%s format)\n", outputFile, outputFormat)
			} else {
				printResults(results)
			}
        }	
	} else if inputFilePath != "" {
		inputFile, err := os.Open(inputFilePath)
		if err != nil {
			fmt.Println("Error opening input file:", err)
			os.Exit(1)
		}
		defer inputFile.Close()

		scanner := bufio.NewScanner(inputFile)
		for scanner.Scan() {
			domain := scanner.Text()
			results := checkAlive(domain, resolveIP)

			if outputFile != "" {
				outputFormat = saveResultsToFile(outputFile, results, resolveIP)
				// fmt.Printf("Results saved to %s (%s format)\n", outputFile, outputFormat)
			} else {
				printResults(results)
			}
		}
		if outputFile != "" {
			fmt.Printf("Results saved to %s (%s format)\n", outputFile, outputFormat)
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading input file:", err)
			os.Exit(1)
		}
	} else {
		if domain == "" {
			fmt.Println("Please specify a domain to check with the -d flag.")
			os.Exit(1)
		}

		results := checkAlive(domain, resolveIP)

		if outputFile != "" {
			outputFormat := saveResultsToFile(outputFile, results, resolveIP)
			fmt.Printf("Results saved to %s (%s format)\n", outputFile, outputFormat)
		} else {
			printResults(results)
		}
	}
}
