package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"jevbeez/api/internal/cache"
	"jevbeez/api/internal/config"
	"jevbeez/api/internal/decision"
	"jevbeez/api/internal/handler"
	"jevbeez/api/internal/jev"
	"jevbeez/api/internal/scraper"
	"jevbeez/api/internal/search"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Config initialization failed: %v", err)
	}

	redisCache, err := cache.NewRedisService(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}
	defer redisCache.Close()

	scraperService := scraper.NewScraperService(cfg.FirecrawlAPIKey)

	searchManager, err := search.NewSearchManager(cfg.TavilyAPIKey, cfg.FirecrawlAPIKey)
	if err != nil {
		log.Fatalf("Search manager initialization failed: %v", err)
	}

	jevClient := jev.NewClient(cfg.TypesafeAPIKey)
	evaluator := decision.NewEvaluator()

	analyzeHandler := handler.NewAnalyzeHandler(
		redisCache,
		scraperService,
		searchManager,
		jevClient,
		evaluator,
	)

	app := fiber.New(fiber.Config{
		EnableTrustedProxyCheck: true,
		ProxyHeader:             fiber.HeaderXForwardedFor,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSOrigin,
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, OPTIONS",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "jevbeez-api",
		})
	})

	app.Post("/api/v1/analyze", analyzeHandler.Handle)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("JevBeez API server listening on port %s", cfg.Port)
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Printf("Server listen error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down JevBeez API server gracefully...")
	_ = app.Shutdown()
}
