
# Subdomain Scanner Tool

A high-performance subdomain scanner written in Go, utilizing the [xReverseLabs Subdomain API](https://api.xreverselabs.my.id) for retrieving subdomains. This tool is optimized for speed using the `fasthttp` library and supports concurrent scanning, making it ideal for large-scale subdomain enumeration.

## Features

![SubRecon](https://raw.githubusercontent.com/xReverseLabs/SubRecon/main/Subrecon.png)

- **Fast and Efficient:** Leverages the `fasthttp` library for high-performance HTTP requests.
- **Concurrent Scanning:** Supports multi-threaded scanning to process multiple domains simultaneously.
- **Customizable Output:** Outputs the results to a specified file.
- **Flexible Usage:** Can scan from a single domain or a list of domains provided in a file.

## Installation

### Prerequisites

- **Go 1.16+** is required to run this tool.

### Clone the Repository

```bash
git clone https://github.com/xReverseLabs/SubRecon.git
cd SubRecon
```

### Install Dependencies

```bash
go mod tidy
```

## Configuration

Before running the tool, you need to set up the `config.json` file with your API key.

### Create `config.json`

```json
{
    "apiKey": "YOUR_API_KEY_HERE"
}
```

Replace `"YOUR_API_KEY_HERE"` with your actual API key obtained from [xReverseLabs](https://xreverselabs.my.id/clientarea/register).

## Usage

### Command-Line Options

- `-f [file]` : Path to the file containing a list of domains to scan.
- `-d [domain]` : A single domain to scan.
- `-o [file]` : Output file for saving the results (default: `output.txt`).
- `-t [threads]` : Number of concurrent threads to use (default: `5`).
- `--help` : Display the help message with usage instructions.

### Examples

#### Scanning a Single Domain

```bash
Linux :
go run main.go -d example.com -o output.txt

Windows :
SubRecon-x64.exe -d example.com -o output.txt
```

This will scan `example.com` for subdomains and save the results in `output.txt`.

#### Scanning Multiple Domains from a File

```bash
Linux :
go run main.go -f domains.txt -t 10 -o output.txt

Windows :
SubRecon-x64.exe -f domains.txt -t 10 -o output.txt
```

This will scan all domains listed in `domains.txt` with `10` concurrent threads and save the results in `output.txt`.

#### Displaying Help

```bash
Linux :
go run main.go --help

Windows :
SubRecon-x64.exe --help
```

This will display a help message with information on how to use the tool.

## Output

The output file will contain the discovered subdomains, one per line, like this:

```plaintext
subdomain1.example.com
subdomain2.example.com
subdomain3.example.com
...
```

Only domains with subdomains will be included in the output file.

## Performance

The tool is designed to be fast and efficient by utilizing the `fasthttp` library for HTTP requests and supporting concurrency. This makes it well-suited for scanning large numbers of domains quickly.

## Contributing

Contributions are welcome! If you have any ideas, suggestions, or bug reports, feel free to open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [fasthttp](https://github.com/valyala/fasthttp) for the high-performance HTTP library.
- [aurora](https://github.com/logrusorgru/aurora) for the beautiful terminal output colors.
- [Yon3zu](https://github.com/yon3zu) Developer

---

Feel free to reach out with any questions or suggestions!
