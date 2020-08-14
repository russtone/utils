package dns

import "strings"

// Subdomains returns list of subdomains of base domain.
func Subdomains(sub, base string) []string {
	domains := make([]string, 0)

	ls := strings.Count(sub, ".")
	lb := strings.Count(base, ".")

	if ls <= lb {
		return nil
	}

	parts := strings.Split(sub, ".")

	for i := 0; i < ls-lb; i++ {
		d := strings.Join(parts[i:], ".")
		domains = append(domains, d)
	}

	return domains
}
