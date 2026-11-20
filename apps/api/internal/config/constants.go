package config

import "time"

const (
	DefaultPort              = "8080"
	DefaultRedisAddr         = "localhost:6379"
	DefaultCacheTTL          = 1 * time.Hour
	UnverifiedCacheTTL       = 2 * time.Minute
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
)
