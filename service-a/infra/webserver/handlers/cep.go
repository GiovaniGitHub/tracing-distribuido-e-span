package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"time"

	"github.com/GiovaniGitHub/service-a/configs"
	"github.com/GiovaniGitHub/service-a/infra/monitoring"
	"github.com/GiovaniGitHub/service-a/internal/entity"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

// Post Address godoc
// @Summary      Handle CEP Request
// @Description  receive a CEP and forwards it if valid
// @Tags         postcep
// @Accept       json
// @Produce      json
// @Param        request     body      entity.Input  true  "CEP request"
// @Success      200  {object}  map[string]interface{}
// @Failure      422  {object}  map[string]string
// @Router       /cep	[post]
func PostTemperature(w http.ResponseWriter, r *http.Request) {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	// Trata o sinal de interrupção
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Inicia o monitoramento
	shutdown, err := monitoring.InitProvider(configs.OTEL_SERVICE_NAME, configs.OTEL_EXPORTER_OTLP_ENDPOINT)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	tracer := otel.Tracer(configs.MICROSERVICE_NAME)
	ctx = otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
	ctx, span := tracer.Start(ctx, "get-cep-temperature")
	defer span.End()

	if r.Method != "POST" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	var input entity.Input
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(r.Header))
	if match, _ := regexp.MatchString(`^\d{8}$`, input.CEP); !match {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}
	apiURL := configs.EXTERNAL_CALL_URL + ":" + configs.EXTERNAL_CALL_PORT + "/cep/%s"
	log.Printf("API esta disponível em: %s", apiURL)

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf(apiURL, input.CEP), nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	startTime := time.Now()
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Failed to request external service", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	endTime := time.Now()

	responseTime := endTime.Sub(startTime)
	span.SetAttributes(attribute.Float64("responseTime", responseTime.Seconds()))
	body, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "Failed to read response from external service", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)

}
