package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/valyala/fasthttp"
)

type Config struct {
	APIKey string `json:"apiKey"`
}

type SubdomainResponse struct {
	Domain     string   `json:"domain"`
	Subdomains []string `json:"subdomains"`
}

func loadConfig() (Config, error) {
	var config Config
	file, err := os.ReadFile("config.json")
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(file, &config)
	return config, err
}

func fetchSubdomains(client *fasthttp.Client, domain, apiKey string) ([]string, error) {
	url := fmt.Sprintf("https://api.xreverselabs.my.id/subdomain?apiKey=%s&url=%s", apiKey, domain)

	// Create a new request and response
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	// Set the request details
	req.SetRequestURI(url)

	// Perform the GET request
	if err := client.Do(req, resp); err != nil {
		return nil, err
	}

	// Check for a 200 OK status code
	if resp.StatusCode() != fasthttp.StatusOK {
		return nil, fmt.Errorf("non-200 response code: %d", resp.StatusCode())
	}

	// Parse the response body
	var result SubdomainResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}

	// Release the request and response objects
	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)

	return result.Subdomains, nil
}

func scanDomain(client *fasthttp.Client, domain string, apiKey string, wg *sync.WaitGroup, output chan<- string) {
	defer wg.Done()
	fmt.Printf("%s Scanning domain: %s\n", aurora.Cyan("[*]").String(), aurora.BrightBlue(domain).String())

	subdomains, err := fetchSubdomains(client, domain, apiKey)
	if err != nil {
		fmt.Printf("Error fetching subdomains for %s: %v\n", domain, err)
		return
	}

	if len(subdomains) > 0 {
		output <- strings.Join(subdomains, "\n") + "\n"
		fmt.Printf("%s Found %d subdomains for %s\n", aurora.Green("[+]"), len(subdomains), aurora.BrightBlue(domain))
	} else {
		fmt.Printf("%s No subdomains found for %s\n", aurora.Red("[-]"), aurora.BrightBlue(domain))
	}
}

func writeOutputToFile(filename string, output <-chan string, done chan<- bool) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		done <- false
		return
	}
	defer file.Close()

	for line := range output {
		_, err := file.WriteString(line)
		if err != nil {
			fmt.Printf("Error writing to output file: %v\n", err)
			done <- false
			return
		}
	}
	done <- true
}

func displayHelp() {
	helpText := `
Usage of Subdomain Scanner:

-f  [file]      File containing list of domains to scan
-d  [domain]    Single domain to scan
-o  [file]      Output file (default: output.txt)
-t  [threads]   Number of concurrent threads (default: 5)

Examples:
go run main.go -f list_domain.txt -t 10 -o output.txt
go run main.go -d xreverselabs.my.id -o output.txt

Credits : https://github.com/xReverselabs/SubRecon`
	fmt.Println(helpText)
}

func main() {
	// Load API key from config file
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	// Command-line flags
	fileFlag := flag.String("f", "", "File containing list of domains")
	domainFlag := flag.String("d", "", "Single domain to scan")
	outputFlag := flag.String("o", "output.txt", "Output file")
	threadsFlag := flag.Int("t", 5, "Number of threads")

	helpFlag := flag.Bool("help", false, "Show help")

	flag.Parse()

	if *helpFlag {
		displayHelp()
		return
	}

	if *fileFlag == "" && *domainFlag == "" {
		fmt.Println("Please specify either a file with -f or a single domain with -d\nOr use --help for more information")
		return
	}

	var domains []string

	// Handle domain list from file
	if *fileFlag != "" {
		file, err := os.Open(*fileFlag)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			domains = append(domains, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return
		}
	}

	// Handle single domain
	if *domainFlag != "" {
		domains = append(domains, *domainFlag)
	}

	// Initialize fasthttp client
	client := &fasthttp.Client{}

	// Create channels and wait group
	output := make(chan string, len(domains))
	done := make(chan bool)
	var wg sync.WaitGroup

	// Write output to file asynchronously
	go writeOutputToFile(*outputFlag, output, done)

	// Start scanning with concurrency
	semaphore := make(chan struct{}, *threadsFlag)
	start := time.Now()
	for _, domain := range domains {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(domain string) {
			defer func() { <-semaphore }()
			scanDomain(client, domain, config.APIKey, &wg, output)
		}(domain)
	}

	wg.Wait()
	close(output)
	success := <-done

	elapsed := time.Since(start)
	if success {
		fmt.Printf("\n%s Results saved to %s in %s\n", aurora.Green("[+]"), aurora.BrightGreen(*outputFlag), aurora.BrightYellow(elapsed.String()))
	} else {
		fmt.Println("Failed to save results.")
	}
}
