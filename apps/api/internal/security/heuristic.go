package security

import (
	"fmt"
	"net"
	"net/url"
	"path"
	"strings"

	"jevbeez/api/internal/config"
)

type ThreatCheckResult struct {
	IsThreat    bool     `json:"is_threat"`
	ThreatType  string   `json:"threat_type,omitempty"`
	Confidence  float64  `json:"confidence"`
	Reasons     []string `json:"reasons"`
	TargetBrand string   `json:"target_brand,omitempty"`
	ApexDomain  string   `json:"apex_domain,omitempty"`
}

var legitimateBrandDomains = map[string][]string{
	"bca":       {"bca.co.id", "klikbca.com", "mybca.bca.co.id"},
	"mandiri":   {"bankmandiri.co.id", "livinbymandiri.co.id", "mandirikartukredit.com"},
	"bri":       {"bri.co.id", "brimo.bri.co.id", "ib.bri.co.id"},
	"bni":       {"bni.co.id", "bni-life.co.id"},
	"cimb":      {"cimbniaga.co.id", "octoclicks.co.id"},
	"dana":      {"dana.id"},
	"ovo":       {"ovo.id"},
	"gopay":     {"gojek.com", "gopay.co.id"},
	"shopee":    {"shopee.co.id"},
	"tokopedia": {"tokopedia.com"},
	"blibli":    {"blibli.com"},
	"pajak":     {"pajak.go.id"},
	"bpjs":      {"bpjs-kesehatan.go.id", "bpjsketenagakerjaan.go.id"},
	"pln":       {"pln.co.id"},
	"telkomsel": {"telkomsel.com"},
	"indihome":  {"indihome.co.id", "telkom.co.id"},
	"kemenkes":  {"kemkes.go.id"},
	"polri":     {"polri.go.id"},
	"posindonesia": {"posindonesia.co.id"},
	"jne":       {"jne.co.id"},
	"jnt":       {"jet.co.id"},
}

func CheckURLHeuristics(rawURL string) *ThreatCheckResult {
	res := &ThreatCheckResult{
		IsThreat:   false,
		Confidence: 0.0,
		Reasons:    []string{},
	}

	u, err := url.Parse(rawURL)
	if err != nil || u.Host == "" {
		return res
	}

	rawLower := strings.ToLower(rawURL)
	host := strings.ToLower(u.Hostname())
	pathLower := strings.ToLower(u.Path)

	// 1. Check for @ symbol trick (URL userinfo deception)
	if u.User != nil || strings.Contains(u.RawPath, "@") || strings.Contains(rawURL, "@") {
		userInfoStr := ""
		if u.User != nil {
			userInfoStr = strings.ToLower(u.User.String())
		}
		// If userInfo has brand or domain pattern, user was tricked to open evil host
		for brand := range config.BrandProfiles {
			if strings.Contains(userInfoStr, brand) || strings.Contains(rawLower, brand+"@") {
				res.IsThreat = true
				res.ThreatType = config.VerdictPhishing
				res.TargetBrand = brand
				res.Confidence = 95.0
				res.Reasons = append(res.Reasons, "Menggunakan manipulasi karakter '@' untuk menyamarkan domain tujuan asli di belakang nama brand terpercaya.")
				return res
			}
		}
	}

	// 2. Check for dangerous file extension (APK dropper / executable malware)
	ext := strings.ToLower(path.Ext(pathLower))
	for _, dangerousExt := range config.DangerousExtensions {
		if ext == dangerousExt || strings.Contains(pathLower, dangerousExt) {
			res.IsThreat = true
			res.ThreatType = config.VerdictMalware
			res.Confidence = 98.0
			res.Reasons = append(res.Reasons, fmt.Sprintf("Tautan mengarah langsung ke unduhan file berbahaya berformat '%s' (indikasi malware/APK pencuri data seperti surat undangan, resi kurir, atau tilang palsu).", dangerousExt))
			return res
		}
	}

	// 3. Check for Punycode / Homoglyph deception
	if strings.HasPrefix(host, "xn--") || strings.Contains(host, ".xn--") {
		res.IsThreat = true
		res.ThreatType = config.VerdictPhishing
		res.Confidence = 90.0
		res.Reasons = append(res.Reasons, "Domain menggunakan karakter penyamar Punycode/Homoglyph untuk meniru ejaan situs terpercaya.")
		return res
	}

	// 4. Check for Raw IP as Host
	isIP := net.ParseIP(host) != nil
	if isIP {
		// If raw IP has suspicious paths (e.g. bca, login, admin, apk)
		for brand := range config.BrandProfiles {
			if strings.Contains(pathLower, brand) {
				res.IsThreat = true
				res.ThreatType = config.VerdictPhishing
				res.TargetBrand = brand
				res.Confidence = 95.0
				res.Reasons = append(res.Reasons, "Menggunakan alamat IP mentah tanpa domain resmi dengan path yang mencatut nama instansi resmi.")
				return res
			}
		}
		if strings.Contains(pathLower, "login") || strings.Contains(pathLower, "auth") || strings.Contains(pathLower, "verify") {
			res.IsThreat = true
			res.ThreatType = config.VerdictPhishing
			res.Confidence = 85.0
			res.Reasons = append(res.Reasons, "Halaman login/verifikasi berada di alamat IP publik langsung tanpa sertifikat nama domain.")
			return res
		}
	}

	// 5. Check Brand Impersonation & Subdomain Trap
	for brand := range config.BrandProfiles {
		brandPresentInHost := strings.Contains(host, brand)
		brandPresentInPath := strings.Contains(pathLower, brand)

		if !brandPresentInHost && !brandPresentInPath {
			continue
		}

		// Check if the host belongs to legitimate domains for this brand
		isLegit := false
		if legitList, ok := legitimateBrandDomains[brand]; ok {
			for _, legitDomain := range legitList {
				if host == legitDomain || strings.HasSuffix(host, "."+legitDomain) {
					isLegit = true
					break
				}
			}
		}

		if isLegit {
			continue
		}

		// If brand keyword is present in host or path, but host is not legitimate:
		// Check if it's on a high-risk TLD or subdomain trap
		isHighRiskTLD := false
		for _, tld := range config.HighRiskTLDs {
			if strings.HasSuffix(host, tld) {
				isHighRiskTLD = true
				break
			}
		}

		isSubdomainTrap := false
		// Subdomain trap: e.g. bca.co.id.login.xyz (contains .co.id. or official domain in non-apex position)
		for _, legit := range legitimateBrandDomains[brand] {
			if strings.Contains(host, legit+".") {
				isSubdomainTrap = true
				break
			}
		}

		if brandPresentInHost || isSubdomainTrap || (brandPresentInPath && isHighRiskTLD) {
			res.IsThreat = true
			res.ThreatType = config.VerdictPhishing
			res.TargetBrand = brand
			res.Confidence = 92.0

			if isSubdomainTrap {
				res.Reasons = append(res.Reasons, fmt.Sprintf("Menyamarkan subdomain menyerupai domain resmi %s untuk mengelabui pengguna (subdomain trap).", strings.ToUpper(brand)))
			} else {
				res.Reasons = append(res.Reasons, fmt.Sprintf("Mencatut nama merek resmi '%s' pada domain tidak terafiliasi.", strings.ToUpper(brand)))
			}

			if isHighRiskTLD {
				res.Reasons = append(res.Reasons, "Menggunakan ekstensi domain murah berisiko tinggi yang kerap disalahgunakan untuk penipuan digital.")
			}
			return res
		}
	}

	return res
}
