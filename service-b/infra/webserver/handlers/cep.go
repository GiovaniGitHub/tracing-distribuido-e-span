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
	"strconv"
	"strings"
	"time"

	"github.com/GiovaniGitHub/service-b/configs"
	"github.com/GiovaniGitHub/service-b/infra/monitoring"
	"github.com/GiovaniGitHub/service-b/internal/entity"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

func handleError(w http.ResponseWriter, errorMessage string, statusCode int) {
	http.Error(w, errorMessage, statusCode)
}

func getRequestWithContext(ctx context.Context, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	return req, nil
}

func makeRequest(req *http.Request) (*http.Response, error) {
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func parseResponseBody(response *http.Response, data interface{}) error {
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, data)
}

func getTemperatureFromWeatherData(weatherData entity.WeatherData, city string) (entity.Temperature, error) {
	var temperature entity.Temperature
	temperature.TempC = weatherData.CurrentCondition[0].TempC
	value, err := strconv.ParseFloat(weatherData.CurrentCondition[0].TempC, 32)
	if err != nil {
		return temperature, err
	}
	temperature.TempF = strconv.FormatFloat(value*1.8+32, 'f', 2, 32)
	temperature.TempK = strconv.FormatFloat(value+273, 'f', 2, 32)
	temperature.City = city
	return temperature, nil
}

// Get Address godoc
// @Summary      Get Address 2
// @Description  Get Address by Post Code 2
// @Tags         address
// @Accept       json
// @Produce      json
// @Param        cep     path       string    true    "Cep"
// @Success      200     {object}   entity.Temperature
// @Failure      400     string     "Bad Request"
// @Failure      404     string     "Not Found"
// @Failure      500     string     "Internal Server Error"
// @Router       /cep/{cep}	[get]
func GetTemperature(w http.ResponseWriter, r *http.Request) {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

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
	ctx, span := tracer.Start(ctx, "get-city-name")

	cep := r.URL.Path[len("/cep/"):]
	viacepURL := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)

	req, err := getRequestWithContext(ctx, viacepURL)
	if err != nil {
		handleError(w, "Erro ao criar requisição HTTP", http.StatusInternalServerError)
		return
	}

	startTime := time.Now()
	response, err := makeRequest(req)
	endTime := time.Now()
	responseTime := endTime.Sub(startTime)
	span.SetAttributes(attribute.Float64("responseTime", responseTime.Seconds()))
	span.End()
	if err != nil {
		handleError(w, "Erro ao fazer requisição HTTP", http.StatusInternalServerError)
		return
	}

	if response.StatusCode == 400 {
		handleError(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	var address entity.Address
	err = parseResponseBody(response, &address)
	if err != nil {
		handleError(w, "Erro ao processar resposta da API", http.StatusInternalServerError)
		return
	}

	if address.Erro {
		handleError(w, "can not found zipcode", http.StatusNotFound)
		return
	}

	city := strings.ReplaceAll(address.Localidade, " ", "+")
	tracer = otel.Tracer(configs.MICROSERVICE_NAME)
	ctx_temperature, span_temperature := tracer.Start(ctx, "get-city-temperature")
	defer span_temperature.End()

	urlWeather := fmt.Sprintf("https://wttr.in/%s?format=j1", city)

	req, err = getRequestWithContext(ctx_temperature, urlWeather)
	if err != nil {
		handleError(w, "Erro ao criar requisição HTTP", http.StatusInternalServerError)
		return
	}
	startTime = time.Now()
	response, err = makeRequest(req)
	endTime = time.Now()

	span_temperature.SetAttributes(attribute.Float64("responseTime", endTime.Sub(startTime).Seconds()))
	if err != nil {
		handleError(w, "Erro ao fazer requisição HTTP", http.StatusInternalServerError)
		return
	}

	var weatherData entity.WeatherData
	err = parseResponseBody(response, &weatherData)
	if err != nil {
		handleError(w, "Erro ao processar resposta da API", http.StatusInternalServerError)
		return
	}

	temperature, err := getTemperatureFromWeatherData(weatherData, city)
	if err != nil {
		handleError(w, "Erro ao calcular temperatura", http.StatusInternalServerError)
		return
	}

	// endTime := time.Now()
	// responseTime := endTime.Sub(startTime)
	// span.SetAttributes(attribute.Float64("responseTime", responseTime.Seconds()))
	json.NewEncoder(w).Encode(temperature)
}
