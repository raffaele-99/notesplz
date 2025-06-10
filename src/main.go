// main.go
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the YAML configuration structure
type Config struct {
	Project string            `yaml:"Project"`
	VpnIP   string            `yaml:"VPN IP"`
	BoxSets map[string]BoxSet `yaml:"Box Sets"`
}

// BoxSet represents a box set configuration
type BoxSet struct {
	Hosts               string `yaml:"hosts"`
	CredPagePerHostOrSet string `yaml:"cred page per host or set"`
	MakeGraphDirectory  string `yaml:"make graph directory"`
}

// Templates module functions
type Templates struct{}

func (t *Templates) MakeRootIndex(project string, boxSets []string) string {
	template := fmt.Sprintf(`
at %s index

---

`, project)
	for _, boxSet := range boxSets {
		template += fmt.Sprintf("-> [[%s index|%s]]\n\n", boxSet, boxSet)
	}
	return template
}

func (t *Templates) MakeSetIndex(project, boxSet, credsPerHostOrSet string) string {
	template := fmt.Sprintf(`
<- back to [[%s index]]

---

-> [[%s hosts]]

-> [[random notes]]

`, project, boxSet)
	if credsPerHostOrSet == "set" {
		template += fmt.Sprintf("-> [[%s/creds|creds]]\n\n", boxSet)
	}
	return template
}

func (t *Templates) MakeSetCredsPage(boxSet string) string {
	return fmt.Sprintf(`
<- back to [[%s index|%s index]]

---
store creds for %s here
`, boxSet, boxSet, boxSet)
}

func (t *Templates) MakeSetHostsIndex(boxSet string, hosts []string) string {
	template := fmt.Sprintf(`
<- back to [[%s index|%s index]]

---

`, boxSet, boxSet)
	for _, host := range hosts {
		template += fmt.Sprintf("-> [[%s/hosts/%s/%s|%s]]\n\n", boxSet, host, host, host)
	}
	return template
}

func (t *Templates) MakeHostIndex(boxSet, host, credsPerHostOrSet string) string {
	template := fmt.Sprintf(`
<- back to [[%s hosts|%s hosts]]

---

-> [[%s/hosts/%s/nmap|nmap]]

-> [[%s/hosts/%s/enum|enum]]
`, boxSet, boxSet, boxSet, host, boxSet, host)
	
	if credsPerHostOrSet == "host" {
		template += fmt.Sprintf("\n-> [[%s/hosts/%s/creds|creds]]", boxSet, host)
	}
	
	return template
}

func (t *Templates) MakeHostEnumPage(boxSet, host string) string {
	return fmt.Sprintf(`
<- back to [[%s/hosts/%s/%s|%s]]

---

enum notes for %s %s
---

`, boxSet, host, host, host, boxSet, host)
}

func (t *Templates) MakeHostNmapPage(boxSet, host string) string {
	return fmt.Sprintf(`
<- back to [[%s/hosts/%s/%s|%s]]

---

nmap notes for %s %s
---
`, boxSet, host, host, host, boxSet, host)
}

func (t *Templates) MakeHostCredsPage(boxSet, host string) string {
	return fmt.Sprintf(`
<- back to [[%s/hosts/%s/%s|%s]]

---

store creds for %s %s here
`, boxSet, host, host, host, boxSet, host)
}

// Main functions
func readYAML(filePath string) (*Config, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func makeGraphDirectory(hosts []string, project, boxSet, vpnIP string) error {
	graphDir := filepath.Join(project, boxSet, "graph")
	err := os.MkdirAll(graphDir, 0755)
	if err != nil {
		return err
	}

	// Create host graph files
	for _, host := range hosts {
		content := fmt.Sprintf(`---
can reach:
  - "[[%s]]"
  - "[[%s]]"
type:
  - computer
  - dc
  - delete/add as needed
---

`, host, vpnIP)
		
		err = ioutil.WriteFile(filepath.Join(graphDir, host+".md"), []byte(content), 0644)
		if err != nil {
			return err
		}
	}

	// Create VPN IP graph file
	vpnContent := fmt.Sprintf(`---
can reach:
  - "[[%s]]"
type:
  - attackbox
---

this is your machine`, vpnIP)
	
	err = ioutil.WriteFile(filepath.Join(graphDir, vpnIP+".md"), []byte(vpnContent), 0644)
	return err
}

func makeSetDirectory(hosts []string, project, boxSet, credsPerHostOrSet, graphStatus, vpnIP string) error {
	templates := &Templates{}
	
	// Create directories
	err := os.MkdirAll(filepath.Join(project, boxSet), 0755)
	if err != nil {
		return err
	}
	
	err = os.MkdirAll(filepath.Join(project, boxSet, "hosts"), 0755)
	if err != nil {
		return err
	}

	// Make the set index file
	setIndexContent := templates.MakeSetIndex(project, boxSet, credsPerHostOrSet)
	err = ioutil.WriteFile(filepath.Join(project, boxSet, boxSet+" index.md"), []byte(setIndexContent), 0644)
	if err != nil {
		return err
	}

	// Make random notes file
	err = ioutil.WriteFile(filepath.Join(project, boxSet, "random notes.md"), []byte("random notes for this set"), 0644)
	if err != nil {
		return err
	}

	// Make the set hosts index file
	hostsIndexContent := templates.MakeSetHostsIndex(boxSet, hosts)
	err = ioutil.WriteFile(filepath.Join(project, boxSet, "hosts", boxSet+" hosts.md"), []byte(hostsIndexContent), 0644)
	if err != nil {
		return err
	}

	// If creds are per set, make the creds file in the set directory
	if credsPerHostOrSet == "set" {
		credsContent := templates.MakeSetCredsPage(boxSet)
		err = ioutil.WriteFile(filepath.Join(project, boxSet, "creds.md"), []byte(credsContent), 0644)
		if err != nil {
			return err
		}
	}

	// Process each host
	for _, host := range hosts {
		// Create host directory
		hostDir := filepath.Join(project, boxSet, "hosts", host)
		err = os.MkdirAll(hostDir, 0755)
		if err != nil {
			return err
		}

		// Make the host index file
		hostIndexContent := templates.MakeHostIndex(boxSet, host, credsPerHostOrSet)
		err = ioutil.WriteFile(filepath.Join(hostDir, host+".md"), []byte(hostIndexContent), 0644)
		if err != nil {
			return err
		}

		// Make the host enum file
		enumContent := templates.MakeHostEnumPage(boxSet, host)
		err = ioutil.WriteFile(filepath.Join(hostDir, "enum.md"), []byte(enumContent), 0644)
		if err != nil {
			return err
		}

		// Make the host nmap file
		nmapContent := templates.MakeHostNmapPage(boxSet, host)
		err = ioutil.WriteFile(filepath.Join(hostDir, "nmap.md"), []byte(nmapContent), 0644)
		if err != nil {
			return err
		}

		// If creds are per host, make the creds file in the host directory
		if credsPerHostOrSet == "host" {
			hostCredsContent := templates.MakeHostCredsPage(boxSet, host)
			err = ioutil.WriteFile(filepath.Join(hostDir, "creds.md"), []byte(hostCredsContent), 0644)
			if err != nil {
				return err
			}
		}
	}

	// Set up the graph for this box set
	if graphStatus == "yes" {
		err = makeGraphDirectory(hosts, project, boxSet, vpnIP)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	// Parse command line arguments
	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatal("Error: Please provide path to YAML file")
	}
	
	filePath := flag.Arg(0)

	// Load data from YAML file
	config, err := readYAML(filePath)
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	if config.Project == "" {
		log.Fatal("Error: No project name found in the YAML file.")
	}

	templates := &Templates{}

	// Create root directory
	err = os.MkdirAll(config.Project, 0755)
	if err != nil {
		log.Fatalf("Error creating project directory: %v", err)
	}

	// Get box set names
	var boxSetNames []string
	for name := range config.BoxSets {
		boxSetNames = append(boxSetNames, name)
	}

	// Create root index
	rootIndex := templates.MakeRootIndex(config.Project, boxSetNames)
	err = ioutil.WriteFile(filepath.Join(config.Project, config.Project+" index.md"), []byte(rootIndex), 0644)
	if err != nil {
		log.Fatalf("Error creating root index: %v", err)
	}

	// Process each box set
	for boxSetName, boxSet := range config.BoxSets {
		// Split hosts string into slice
		hosts := strings.Fields(boxSet.Hosts)
		
		err = makeSetDirectory(hosts, config.Project, boxSetName, boxSet.CredPagePerHostOrSet, boxSet.MakeGraphDirectory, config.VpnIP)
		if err != nil {
			log.Fatalf("Error creating box set directory for %s: %v", boxSetName, err)
		}
	}

	fmt.Printf("Successfully created notes structure for project: %s\n", config.Project)
}
