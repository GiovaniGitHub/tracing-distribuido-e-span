# Go Expert - LAB Desafio 

## Descrição
O sistema em Go que receba um CEP, identifica a cidade e retorna o clima atual (temperatura em graus celsius, fahrenheit e kelvin) com monitoramento usando OpenTelemetry, jaeger, e zipkin.

## Conteúdo

1. [Como Rodar o Projeto](#como-rodar-o-projeto)
2. [Testando API](#testando-api)
2. [Monitoramento](#monitoramento)

## Como Rodar o Projeto
1. Certifique-se de ter o Docker instalado em sua máquina.
2. Clone o projeto do GitHub: [link](https://github.com/GiovaniGitHub/tracing-distribuido-e-span.git)
3. Entre no diretorio do projeto e execute o comando:
4. Crie um `.env` na pasta `service-a` e outro na pasta `service-b` a partir do `.env.template` e altere o campos:
    
    Exemplo:
 - **service-a**:
    ```bash
    WEB_SERVER_PORT=8081
    URL_BASE=http://localhost
    EXTERNAL_CALL_URL=http://serviceb
    EXTERNAL_CALL_PORT=8080
    OTEL_SERVICE_NAME=service-a
    OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:4317
    MICROSERVICE_NAME="microservice-tracer-service-a"
    ```

 - **service-b**:
    ```bash
    WEB_SERVER_PORT=8080
    URL_BASE=http://localhost
    MICROSERVICE_NAME="microservice-tracer-service-b"
    OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:4317
    OTEL_SERVICE_NAME=service-b
    ```

5. Execute o comando:
    ```sh
    docker-compose up -d --build
    ```

## Testando API
### Usando o Swagger:
Existe dois links para acesso ao Swagger:
 - **service-a**:
    Exemplo: http://localhost:8081/docs/index.html
 - **service-b**:
    Exemplo: http://localhost:8080/docs/index.html

### Usando Curl:
#### service-a
| Comando | Resultado                             |
|---------|-----------------------------------------|
|curl -X 'POST' 'http://localhost:8081/cep' -H 'accept: application/json' -H 'Content-Type: application/json' -d '{"cep": "59067400"}'| {"city": Natal, "temp_C":"26","temp_F":"78.80","temp_K":"293.00"} |
|curl -X 'POST' 'http://localhost:8081/cep' -H 'accept: application/json' -H 'Content-Type: application/json' -d '{"cep": "5906740A"}'| invalid zipcode |

#### service-b
| Comando | Resultado                             |
|---------|-----------------------------------------|
| curl -X 'GET' 'http://localhost:8080/cep/70070080' -H 'accept: application/json' | {"city":"Brasília", "temp_C":"36","temp_F":"96.80","temp_K":"309.00"} |
| curl -X 'GET' 'http://localhost:8080/cep/7007008A' -H 'accept: application/json' | invalid zipcode |
| curl -X 'GET' 'http://localhost:8080/cep/70070081' -H 'accept: application/json' | can not found zipcode |

## Monitoramento
- Acessar o zipkin: http://localhost:9411/zipkin/
- Acessar o Jaeger: http://localhost:16686
- Acessar o Prometheus: http://localhost:9090