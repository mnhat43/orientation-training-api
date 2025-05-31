package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	cf "orientation-training-api/configs"

	r "orientation-training-api/cmd/frontapi/router"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func initializeDatabase() {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	pwd, _ := os.Getwd()
	sourceURL := "file:///" + pwd + "/internal/platform/db/migrations"
	databaseURL := "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbName + "?sslmode=disable"

	if os.Getenv("ENV") != "dev" {
		// ref: https://github.com/golang-migrate/migrate/issues/275#issuecomment-523469298
		databaseURL = "postgres://" + user + ":" + password + "@/" + dbName + "?host=/cloudsql/" + host
	}

	log.Info("databaseURL: ", databaseURL)

	log.Info("==> MIGRATION: RUN")
	m, err := migrate.New(sourceURL, databaseURL)

	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	log.Info("==> MIGRATION: DONE!!!")
}

func main() {

	e := echo.New()
	log.Info("==> initializeDatabase: Start init !!!")
	initializeDatabase()
	log.Info("==> initializeDatabase: End init !!!")

	if os.Getenv("ENV") != "prod" {
		e.Debug = true
	} else {
		e.Logger.SetLevel(log.INFO)
	}

	e.Use(middleware.Logger())

	// custom http error, change error message of jwt
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}

		// customer error message of jwt token
		if strings.Contains(err.Error(), "jwt") {
			c.JSON(code, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Login invalid. Please login again",
			})
		} else {
			// customer another error message of jwt token
			c.JSON(code, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: err.Error(),
			})
		}
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodOptions, http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello , you are in routes")
	})

	router := r.NewAppRouter(e.Logger)
	router.UserRoute(e.Group("/user"))
	router.AuthRoute(e.Group("/auth"))
	router.CourseRoute(e.Group("/course"))
	router.ModuleItemRoute(e.Group("/module-item"))
	router.ModuleRoute(e.Group("/module"))
	router.LectureRoute(e.Group("/lecture"))
	router.UserProgressRoute(e.Group("/user-progress"))
	router.TemplatePathRoute(e.Group("/template-path"))
	router.QuizRoute(e.Group("/quiz"))

	go func() {
		if err := e.Start(":8080"); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
