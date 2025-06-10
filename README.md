# notesplz
ai go port of notes-template

## setup
```
go mod init notes-template
go get gopkg.in/yaml.v3
```

## usage
```bash
go run src/main.go config.yaml
```
or
```bash
cd src && go build -o notesplz
./notesplz config.yaml
```

### config file

yaml file should follow this structure:

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

### result

the generated structure would be:
```
OSCP A/
├── Project Name index.md
├── ad set/
│   ├── ad set index.md
│   ├── random notes.md
│   ├── creds.md
│   ├── hosts/
│   │   ├── ad set hosts.md
│   │   ├── 192.168.226.141/
│   │   │   ├── 192.168.226.141.md
│   │   │   ├── enum.md
│   │   │   └── nmap.md
│   │   └── ...
│   └── graph/ 
│       ├── 192.168.226.141.md
│       ├── 10.10.186.140.md
│       ├── 10.10.186.142.md
│       └── 192.168.45.166.md
└── ...
```
