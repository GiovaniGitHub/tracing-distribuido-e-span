// server_test.go
package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GiovaniGitHub/service-b/infra/webserver/handlers"
	"github.com/GiovaniGitHub/service-b/internal/entity"
)

func TestGetTemperatureEndpoint(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/cep/59067400", nil)
	w := httptest.NewRecorder()
	handlers.GetTemperature(w, req)
	res := w.Result()
	defer res.Body.Close()

	if status := res.StatusCode; status != http.StatusOK {
		t.Errorf("Handler retornou um código de status incorreto: esperado %v, obtido %v", http.StatusOK, status)
	}

	expectedContentType := "text/plain; charset=utf-8"
	if contentType := res.Header.Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("Tipo de conteúdo incorreto: esperado %v, obtido %v", expectedContentType, contentType)
	}

	var actualTemperature entity.Temperature
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Erro ao ler o corpo da resposta: %v", err)
	}

	err = json.Unmarshal(body, &actualTemperature)
	if err != nil {
		t.Errorf("Erro ao decodificar o corpo da resposta: %v", err)
	}
}
