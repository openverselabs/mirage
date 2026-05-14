package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type Result struct {
	Domain string
	CNAME  string
}

var (
	concurrency int
	outputFile  string
	listFile    string
	silent      bool
	resolver    string
	timeout     int
)

func init() {
	flag.IntVar(&concurrency, "c", 100, "Maximum concurrency / jumlah worker")
	flag.StringVar(&outputFile, "o", "", "File to write output to")
	flag.StringVar(&listFile, "l", "", "File containing list of domains to check")
	flag.BoolVar(&silent, "silent", false, "Silent mode (no terminal output)")
	flag.StringVar(&resolver, "r", "", "Custom DNS resolver (e.g., 8.8.8.8)")
	flag.IntVar(&timeout, "t", 5, "DNS Timeout in seconds")

	flag.Usage = func() {
		if !silent {
			showBanner()
		}
		fmt.Fprintf(os.Stderr, "Usage:\n  mirage [flags]\n\nFlags:\n")
		flag.PrintDefaults()
	}
}

func showBanner() {
	banner := "\n" +
		"'||\\   /||`                                   \n" +
		" ||\\\\.//||   ''                               \n" +
		" ||     ||   ||  '||''|  '''|.  .|''|| .|''|, \n" +
		" ||     ||   ||   ||    .|''||  ||  || ||..||\n" +
		".||     ||. .||. .||.   `|..||. `|..|| `|...  \n" +
		"                                    ||        \n" +
		"    openverselabs - v0.1.0      `...|'        \n"
	fmt.Fprintln(os.Stderr, banner)
}

func main() {
	flag.Parse()

	var outWriter *os.File
	if outputFile != "" {
		var err error
		outWriter, err = os.Create(outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[!] Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer outWriter.Close()
	}

	jobs := make(chan string)
	results := make(chan Result)

	var wg sync.WaitGroup
	var outWg sync.WaitGroup

	outWg.Add(1)
	go func() {
		defer outWg.Done()
		for res := range results {
			lineColor := fmt.Sprintf("\033[36m%s\033[0m -> \033[32m%s\033[0m", res.Domain, res.CNAME)
			lineRaw := fmt.Sprintf("%s -> %s", res.Domain, res.CNAME)

			if !silent {
				fmt.Println(lineColor)
			}
			if outWriter != nil {
				outWriter.WriteString(lineRaw + "\n")
			}
		}
	}()

	res := &net.Resolver{
		PreferGo: true,
	}

	if resolver != "" {
		res.Dial = func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Duration(timeout) * time.Second,
			}
			return d.DialContext(ctx, "udp", resolver+":53")
		}
	}

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for domain := range jobs {
				checkCNAME(domain, results, res, timeout)
			}
		}()
	}

	var scanner *bufio.Scanner

	if listFile != "" {
		if !silent {
			showBanner()
		}
		file, err := os.Open(listFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[!] Error opening input file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	} else {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			flag.Usage()
			os.Exit(0)
		}
		if !silent {
			showBanner()
		}
		scanner = bufio.NewScanner(os.Stdin)
	}

	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())
		if idx := strings.Index(domain, " "); idx != -1 {
			domain = domain[:idx]
		}

		if domain != "" {
			if strings.HasPrefix(domain, "http://") {
				domain = strings.TrimPrefix(domain, "http://")
			} else if strings.HasPrefix(domain, "https://") {
				domain = strings.TrimPrefix(domain, "https://")
			}
			if idx := strings.Index(domain, "/"); idx != -1 {
				domain = domain[:idx]
			}
			if domain != "" {
				jobs <- domain
			}
		}
	}
	close(jobs)

	wg.Wait()
	close(results)

	outWg.Wait()
}

func checkCNAME(domain string, results chan<- Result, res *net.Resolver, timeout int) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	cname, err := res.LookupCNAME(ctx, domain)
	if err != nil {
		return
	}

	cname = strings.TrimSuffix(cname, ".")
	domain = strings.TrimSuffix(domain, ".")

	if cname != domain && cname != "" {
		results <- Result{Domain: domain, CNAME: cname}
	}
}
