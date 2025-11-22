package routes

import (
	"bigdata-api/internal/api/controllers"
	"bigdata-api/internal/api/middlewares"
	"bigdata-api/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func InitRoutes(cfg *config.Config) *fiber.App {
	app := fiber.New()

	// -----------------------------------------------------
	// GLOBAL CORS
	// -----------------------------------------------------
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept, x-api-key",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
	}))

	api := app.Group("/v1")

	// -----------------------------------------------------
	// PUBLIC ROUTES (No Auth)
	// -----------------------------------------------------
	api.Get("/health", controllers.HealthCheck)
	api.Get("/metrics", controllers.MetricsController)

	// -----------------------------------------------------
	// USER ROUTES (API KEY REQUIRED + RATE LIMIT)
	// -----------------------------------------------------
	user := api.Group("/",
		middlewares.ApiKeyMiddleware(),  // must have valid API key
		middlewares.LoggingMiddleware(), // save logs
	)

	// user query (rate-limited)
	user.Post("query",
		middlewares.RateLimitMiddleware(), // ONLY query is rate limited
		middlewares.CreditsMiddleware(),
		controllers.QueryController,
	)

	// -----------------------------------------------------
	// ADMIN ROUTES (Admin Key Only — NO RATE LIMIT)
	// -----------------------------------------------------
	admin := api.Group("/admin",
		middlewares.AdminMiddleware(), // only admin API key allowed
	)

	// moved stats → no rate-limit now
	admin.Get("/stats", controllers.StatsController)

	// user management
	admin.Post("/create-user", controllers.CreateUserController)
	admin.Post("/add-credits", controllers.AddCreditsController)
	admin.Get("/users", controllers.GetUsersController)
	admin.Post("/regenerate-key", controllers.RegenerateAPIKeyController)

	// logs
	admin.Get("/logs", controllers.GetLogsController)
	admin.Get("/api-key", controllers.GetAdminAPIKey)

	return app
}
