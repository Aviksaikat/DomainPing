# 🌐 DomainPing

DomainPing is a command line tool written in Go that checks the availability of a list of domains and saves the results to a file.

> Saikat Karmakar | 10 Mar : 2023

![](media/banner.gif)

## 🚀 Features
- 🌐 Bulk domain checking: Check the availability and response time of multiple domains at once.

- 🕵️‍♀️ IP address resolution: Optionally resolve the IP address of each domain and include it in the output.

- 📄 Flexible output format: Save the results to a file in either JSON or plain text format.

- 🔎 HTTP GET requests: Send an HTTP GET request to each domain to check if it's alive.

- 🔌 Easy integration: Use as a standalone tool or integrate into your own Go program as a package.

# 💾 Requirements 
```bash
- Go 1.19.5 or higher
```

# 🛠️ Installation 
- Using `go install`
```bash
go install github.com/aviksaikat/DomainPing@latest
```

- Build from source
```bash
go build domainPing.go
./domainPing
```


# 🤖 Usage
```bash
go run domainPing.go [OPTIONS] INPUT_FILE

Options:
  -banner
    	Print banner
  -d string
    	Specify the domain to check
  -f string
    	Input file path
  -h	show help
  -ip
    	Resolve IP address
  -o string
    	Output file path

```

![](media/help.gif)

- `bin/domainPing -h`


```bash
bin/domainPing -h                                            

██████╗  ██████╗ ███╗   ███╗ █████╗ ██╗███╗   ██╗██████╗ ██╗███╗   ██╗ ██████╗
██╔══██╗██╔═══██╗████╗ ████║██╔══██╗██║████╗  ██║██╔══██╗██║████╗  ██║██╔════╝
██║  ██║██║   ██║██╔████╔██║███████║██║██╔██╗ ██║██████╔╝██║██╔██╗ ██║██║  ███╗
██║  ██║██║   ██║██║╚██╔╝██║██╔══██║██║██║╚██╗██║██╔═══╝ ██║██║╚██╗██║██║   ██║
██████╔╝╚██████╔╝██║ ╚═╝ ██║██║  ██║██║██║ ╚████║██║     ██║██║ ╚████║╚██████╔╝
╚═════╝  ╚═════╝ ╚═╝     ╚═╝╚═╝  ╚═╝╚═╝╚═╝  ╚═══╝╚═╝     ╚═╝╚═╝  ╚═══╝ ╚═════╝


Author: avik_saikat
Github: https://github.com/aviksaikat
Gitlab: https://gitlab.com/aviksaikat

----------------------------------------------------------------

Usage: bin/domainPing [OPTIONS] INPUT_FILE
Ping a list of domains and save the results to a file.

Options:
  -banner
    	Print banner
  -d string
    	Specify the domain to check
  -f string
    	Input file path
  -h	show help
  -ip
    	Resolve IP address
  -o string
    	Output file path
```

## Check the availability of a single domains
![](media/solitary.gif)
![](media/single_domain.gif)

## Check the availability of domains in a file and print the results to the console
![](media/files.gif)

## Save the results to a JSON file
![](media/json.gif)

## Check the availability of domains in a file, include the resolved IP address
![](media/ip.gif)

## Using pipes
![](media/pipe.gif)
![](media/pipe_file.gif)


# 🤝 Contributing 
Contributions, issues and feature requests are welcome. Feel free to check the [issues page](https://github.com/Aviksaikat/DomainPing/issues) if you want to contribute.


# 💖 Show your support 
Give a ⭐️ if this project helped you!