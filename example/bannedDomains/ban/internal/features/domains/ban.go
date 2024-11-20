package domains

import (
	"bufio"
	"os"
)

// Domains is a feature that checks if a domain is banned.
type Domains struct {
	bannedDomains map[string]interface{}
}

// NewDomains creates a new Domains feature.
func NewDomains() *Domains {
	d := &Domains{
		bannedDomains: make(map[string]interface{}),
	}

	d.readDomains()

	return d
}

// IsBanned checks if a domain is banned.
func (d Domains) IsBanned(domain string) bool {
	_, ok := d.bannedDomains[domain]

	return ok
}

// readDomains reads a list of banned domains from a file.
func (d Domains) readDomains() {
	bannedDomainsSource, err := os.Open("banned_domains.txt")
	if err != nil {
		panic(err)
	}
	defer bannedDomainsSource.Close()

	bannedDomains := bufio.NewScanner(bannedDomainsSource)
	for bannedDomains.Scan() {
		d.bannedDomains[bannedDomains.Text()] = nil
	}
}
