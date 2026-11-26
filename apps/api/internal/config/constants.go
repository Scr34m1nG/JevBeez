package config

import "time"

const (
	DefaultPort              = "8080"
	DefaultRedisAddr         = "localhost:6379"
	DefaultCacheTTL          = 1 * time.Hour
	UnverifiedCacheTTL       = 2 * time.Minute
	ShareCacheTTL            = 30 * 24 * time.Hour
	DefaultCooldownDuration  = 15 * time.Second
	DefaultBlockDuration     = 24 * time.Hour
	DefaultRequestBurstLimit = 10
	DefaultExpirationWindow  = 1 * time.Minute

	HTTPClientTimeout         = 8 * time.Second
	DialerTimeout             = 5 * time.Second
	ResponseHeaderTimeout     = 5 * time.Second
	MaxScrapeBodyBytes        = 2 * 1024 * 1024
	MinArticleWordCount       = 80
	MinSocialWordCount        = 3

	PlatformTwitter   = "twitter"
	PlatformTikTok    = "tiktok"
	PlatformYouTube   = "youtube"
	PlatformInstagram = "instagram"
	PlatformThreads   = "threads"
	PlatformFacebook  = "facebook"

	TwitterOEmbedEndpoint = "https://publish.twitter.com/oembed"
	FxTwitterEndpoint     = "https://api.fxtwitter.com"
	YouTubeOEmbedEndpoint = "https://www.youtube.com/oembed"

	JevEndpoint               = "https://api.typesafe.ai/v1/systemone"
	JevDefaultModel           = "jev-latest"
	TavilyEndpoint            = "https://api.tavily.com/search"
	FirecrawlScrapeEndpoint   = "https://api.firecrawl.dev/v1/scrape"
	FirecrawlSearchEndpoint   = "https://api.firecrawl.dev/v1/search"

	ThresholdContradictedHigh = 0.70
	ThresholdSupportedHigh    = 0.70
	ThresholdSensationalLow   = 0.40
	ThresholdSensationalHigh  = 0.70
	ThresholdContradictedLow  = 0.40

	VerdictHoax        = "HOAX"
	VerdictValid       = "VALID"
	VerdictMisleading  = "MISLEADING"
	VerdictUnverified  = "UNVERIFIED"
	VerdictPhishing    = "PHISHING"
	VerdictMalware     = "MALWARE"
	VerdictScam        = "SCAM"

	MaxShortenerHops       = 5
	MaxSandboxPayloadBytes = 100 * 1024
	SandboxTimeout         = 3 * time.Second

	VerdictAds            = "ADS"
	MaxUnwrapDepth        = 3
	MinHexSubdomainLength = 16
)

var DangerousExtensions = []string{
	".apk", ".exe", ".scr", ".bat", ".cmd", ".msi", ".vbs", ".iso", ".dmg", ".bin",
}

var HighRiskTLDs = []string{
	".xyz", ".top", ".tk", ".icu", ".buzz", ".cfd", ".fit", ".rest", ".cam", ".work",
	".gq", ".ml", ".ga", ".cf", ".click", ".link", ".sbs", ".quest",
}

var ShortenerDomains = []string{
	"s.id", "link.id", "tribun.link", "heylink.me", "lynk.id", "desty.page",
	"bit.ly", "tinyurl.com", "t.co", "is.gd", "cutt.ly", "rebrand.ly", "rb.gy",
	"shorturl.at", "ow.ly", "buff.ly", "qr.ae", "v.gd", "soo.gd", "rotf.lol",
}

var WrapperParams = []string{
	"adurl", "url", "u", "q", "to", "redirect", "redirect_uri", "redirect_url",
	"dest", "destination", "target", "continue", "goto", "next", "link", "out", "r",
}

var WrapperHostPatterns = []string{
	"google.com/pagead/", "googleadservices.com", "google.com/url",
	"doubleclick.net", "vk.com/away.php", "l.facebook.com/l.php",
	"lm.facebook.com/l.php", "l.instagram.com", "out.reddit.com",
	"linkedin.com/redir", "youtube.com/redirect",
}

var EcommerceDomains = []string{
	"ebay.com", "ebay.co.uk", "amazon.com", "amazon.co.id", "amazon.co.uk", "amazon.de", "amazon.sg",
	"shopee.co.id", "shopee.sg", "shopee.com.my", "shopee.ph", "shopee.co.th",
	"tokopedia.com", "blibli.com", "bukalapak.com", "lazada.co.id", "lazada.sg", "lazada.com.my",
	"aliexpress.com", "aliexpress.us", "alibaba.com", "temu.com", "shein.com",
	"zalora.co.id", "jd.id", "elevenia.co.id", "etsy.com", "walmart.com", "target.com",
	"rakuten.com", "olx.co.id", "carousell.com", "carousell.co.id", "bhphotovideo.com",
}

var BrandProfiles = map[string][]string{
	"bca":          {"bca.co.id", "klikbca.com", "mybca.bca.co.id"},
	"mandiri":      {"bankmandiri.co.id", "livinbymandiri.co.id", "mandirikartukredit.com"},
	"bri":          {"bri.co.id", "brimo.bri.co.id", "ib.bri.co.id"},
	"bni":          {"bni.co.id", "bni-life.co.id"},
	"cimb":         {"cimbniaga.co.id", "octoclicks.co.id"},
	"dana":         {"dana.id"},
	"ovo":          {"ovo.id"},
	"gopay":        {"gojek.com", "gopay.co.id"},
	"shopee":       {"shopee.co.id"},
	"tokopedia":    {"tokopedia.com"},
	"blibli":       {"blibli.com"},
	"pajak":        {"pajak.go.id"},
	"bpjs":         {"bpjs-kesehatan.go.id", "bpjsketenagakerjaan.go.id"},
	"pln":          {"pln.co.id"},
	"telkomsel":    {"telkomsel.com"},
	"indihome":     {"indihome.co.id", "telkom.co.id"},
	"kemenkes":     {"kemkes.go.id"},
	"polri":        {"polri.go.id"},
	"posindonesia": {"posindonesia.co.id"},
	"jne":          {"jne.co.id"},
	"jnt":          {"jet.co.id"},
}

