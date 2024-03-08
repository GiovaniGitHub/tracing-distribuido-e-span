package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GiovaniGitHub/service-a/configs"
	_ "github.com/GiovaniGitHub/service-a/docs"
	"github.com/GiovaniGitHub/service-a/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Desafio 2.0 - service-a
// @version         1.0
// @description     Fullcycle Pós Go Expert Go Expert

// @termsOfService  http://swagger.io/terms/

// @contact.name   Giovani Angelo
// @contact.email  giovani.angelo@gmail.com

// @host      localhost:8081
// @BasePath  /
// @in header
// @name Authorization
func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/cep", func(r chi.Router) {
		r.Post("/", handlers.PostTemperature)
	})
	r.Handle("/metrics", promhttp.Handler())
	// Inicia o servidor
	apiURL := configs.URL_BASE + ":" + configs.WEB_SERVER_PORT + "/cep"
	log.Printf("API está disponível em: %s", apiURL)

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL(configs.URL_BASE+":"+configs.WEB_SERVER_PORT+"/docs/doc.json")))
	log.Printf("API Swagger está disponível em: %s", configs.URL_BASE+":"+configs.WEB_SERVER_PORT+"/docs/index.html")

	http.ListenAndServe(fmt.Sprintf(":%s", configs.WEB_SERVER_PORT), r)
}
