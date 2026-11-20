package security

import (
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const (
	ReasonSSRFAttempt      = "SSRF_ATTEMPT"
	ReasonPathTraversal    = "PATH_TRAVERSAL"
	ReasonInjectionAttempt = "INJECTION_ATTEMPT"
	ReasonScannerProbe     = "SCANNER_PROBE"
)

var (
	injectionRegex = regexp.MustCompile(`(?i)(<script|javascript:|onerror\s*=|onload\s*=|union\s+select|waitfor\s+delay|pg_sleep\(|sleep\(|1=1--|'(\s*or\s*)'1'='1|;\s*cat\s+|\|\s*cat\s+|;\s*id\b|\|\s*id\b|\$\(id\)|\$\(cat|\bwhoami\b)`)
	sensitiveFiles = []string{
		"/etc/passwd",
		"/etc/shadow",
		"/proc/self/",
		"/.env",
		"/.git",
		"/.svn",
		"/.aws",
		"/id_rsa",
		"wp-config",
		"wp-login",
		"phpmyadmin",
		"eval-stdin",
		".php",
	}
)

func DetectExploit(rawURL string) (bool, string, string) {
	clean := strings.TrimSpace(rawURL)
	if clean == "" {
		return false, "", ""
	}

	lower := strings.ToLower(clean)

	// Check non-http schemes
	for _, scheme := range []string{"file:", "gopher:", "dict:", "ftp:", "ldap:", "php:", "data:", "expect:"} {
		if strings.HasPrefix(lower, scheme) {
			return true, ReasonSSRFAttempt, "Protokol non-HTTP tidak diizinkan: " + scheme
		}
	}

	// Check unescaped and double-unescaped variants
	unescaped := lower
	for i := 0; i < 2; i++ {
		if u, err := url.QueryUnescape(unescaped); err == nil && u != unescaped {
			unescaped = strings.ToLower(u)
		}
	}

	// 1. Path traversal check
	if strings.Contains(unescaped, "../") || strings.Contains(unescaped, "..\\") || strings.Contains(rawURL, "%2e%2e") {
		return true, ReasonPathTraversal, "Percobaan path traversal terdeteksi"
	}

	// 2. Sensitive files / system files check
	for _, file := range sensitiveFiles {
		if strings.Contains(unescaped, file) {
			// Don't ban normal news URLs ending with .php if domain is public
			if file == ".php" {
				if strings.Contains(unescaped, "wp-login") || strings.Contains(unescaped, "wp-config") || strings.Contains(unescaped, "eval-stdin") || strings.Contains(unescaped, "c99.php") || strings.Contains(unescaped, "r57.php") || strings.Contains(unescaped, "wso.php") {
					return true, ReasonScannerProbe, "Percobaan probe file sensitif / webshell terdeteksi"
				}
				continue
			}
			return true, ReasonPathTraversal, "Percobaan akses file sensitif terdeteksi: " + file
		}
	}

	// 3. Injection checks (XSS, SQLi, Command Injection)
	if match := injectionRegex.FindString(unescaped); match != "" {
		return true, ReasonInjectionAttempt, "Pola injeksi berbahaya terdeteksi: " + match
	}

	// 4. Host-level SSRF and private IP checks
	parsed, err := url.Parse(clean)
	if err == nil && parsed.Host != "" {
		hostname := strings.ToLower(parsed.Hostname())

		// Loopback and Cloud Metadata hostnames
		if hostname == "localhost" || strings.HasSuffix(hostname, ".localhost") ||
			hostname == "metadata.google.internal" ||
			hostname == "instance-data" ||
			hostname == "metadata.azure.com" {
			return true, ReasonSSRFAttempt, "Target loopback / cloud metadata internal: " + hostname
		}

		// DNS wildcard services mapping to internal IP
		if strings.Contains(hostname, "nip.io") || strings.Contains(hostname, "sslip.io") || strings.Contains(hostname, "lvh.me") {
			return true, ReasonSSRFAttempt, "Penyedia DNS wildcard internal terdeteksi: " + hostname
		}

		// Direct IP check (IPv4 / IPv6)
		if ip := net.ParseIP(hostname); ip != nil {
			if IsPrivateIP(ip) {
				return true, ReasonSSRFAttempt, "Akses ke IP privat / internal diblokir: " + ip.String()
			}
		}

		// Decimal or Hex IP detection (e.g. http://2130706433 or http://0x7f000001)
		if isNumericOrHexIP(hostname) {
			return true, ReasonSSRFAttempt, "Format alamat IP alternatif (hex/decimal) terdeteksi"
		}
	}

	return false, "", ""
}

func isNumericOrHexIP(host string) bool {
	// Strip port if any
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}

	// Single decimal integer (e.g. 2130706433 for 127.0.0.1)
	if num, err := strconv.ParseUint(host, 10, 64); err == nil && num > 0 {
		return true
	}

	// Hex format (e.g. 0x7f000001)
	if strings.HasPrefix(strings.ToLower(host), "0x") {
		return true
	}

	return false
}
