# Notes Template Generator - Go Port

This is a Go port of the Python script for generating Obsidian vault templates for OSCP-like box sets.

## Requirements

- Go 1.20 or higher
- gopkg.in/yaml.v3 package

## Setup

1. Initialize the Go module:
```bash
go mod init notes-template
```

2. Install dependencies:
```bash
go get gopkg.in/yaml.v3
```

## Usage

Run the program with the path to your YAML configuration file:

```bash
go run main.go example.yaml
```

Or build and run:

```bash
go build -o notes-template
./notes-template example.yaml
```

## YAML Configuration Format

The YAML file should follow this structure:

```yaml
Project:
  "OSCP A"
VPN IP:
  192.168.45.166
Box Sets:
  ad set:
    hosts:
      10.10.186.140
      10.10.186.142
      192.168.226.141
    cred page per host or set:
      "set"
    make graph directory:
      "yes"
  standalone:
    hosts:
      192.168.226.143
      192.168.226.144
      192.168.226.145
    cred page per host or set:
      "host"
    make graph directory:
      "no"
```

## Features

- Creates organized directory structure for note-taking
- Supports both per-host and per-set credential pages
- Optional graph directory generation for network visualization
- Generates markdown files with proper linking for Obsidian

## Differences from Python Version

- Uses Go's standard library and `gopkg.in/yaml.v3` for YAML parsing
- Error handling follows Go conventions
- File I/O uses `ioutil` package
- Command-line parsing uses `flag` package instead of `argparse`

## Project Structure

The generated structure will be:
```
Project Name/
├── Project Name index.md
├── Box Set 1/
│   ├── Box Set 1 index.md
│   ├── random notes.md
│   ├── creds.md (if per-set)
│   ├── hosts/
│   │   ├── Box Set 1 hosts.md
│   │   ├── Host 1/
│   │   │   ├── Host 1.md
│   │   │   ├── enum.md
│   │   │   ├── nmap.md
│   │   │   └── creds.md (if per-host)
│   │   └── ...
│   └── graph/ (optional)
│       ├── Host 1.md
│       ├── Host 2.md
│       └── VPN IP.md
└── ...
```
