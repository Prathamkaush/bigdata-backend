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
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, x-api-key",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
	}))

	api := app.Group("/v1")

	// -----------------------------------------------------
	// PUBLIC ROUTES (No Auth)
	// -----------------------------------------------------
	api.Get("/health", controllers.HealthCheck)
	api.Get("/metrics", controllers.MetricsController)
	api.Post("/admin/verify-key", controllers.VerifyKeyController)

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

	// feedback
	user.Post("/feedback", controllers.SubmitFeedback)

	// -----------------------------------------------------
	// ADMIN ROUTES (Admin Key Only — NO RATE LIMIT)
	// -----------------------------------------------------
	admin := api.Group("/admin",
		middlewares.AdminMiddleware(), // only admin API key allowed
		middlewares.LoggingMiddleware(),
	)

	// moved stats → no rate-limit now
	admin.Get("/stats", controllers.StatsController)

	// user management
	admin.Post("/create-user", controllers.CreateUserController)
	admin.Post("/add-credits", controllers.AddCreditsController)
	admin.Get("/users", controllers.GetUsersController)
	admin.Post("/user/:id/regenerate-key", controllers.RegenerateAPIKeyController)
	admin.Get("/user/:id", controllers.GetUserDetails)
	admin.Get("/user/:id/logs", controllers.GetUserLogs)
	admin.Get("/user/:id/usage", controllers.GetUserUsage)
	admin.Delete("/user/:id", controllers.DeleteUserController)

	//feedback
	admin.Get("/feedback", controllers.AdminGetFeedback)

	// logs
	admin.Get("/logs", controllers.GetLogsController)
	admin.Get("/api-key", controllers.GetAdminAPIKey)

	return app
}
