# Go Expert - LAB Desafio 

## Descrição
O sistema em Go que receba um CEP, identifica a cidade e retorna o clima atual (temperatura em graus celsius, fahrenheit e kelvin).
## Conteúdo

1. [Como Rodar o Projeto](#como-rodar-o-projeto)
2. [Testes Automatizados](#testes-automatizados)
3. [Docker](#docker)
4. [Deploy no Google Cloud Run](#deploy-no-google-cloud-run)

## Como Rodar o Projeto
### Ambiente de Desenvolvimento

1. Certifique-se de ter o Golang 1.19 instalado em sua máquina.
2. Clone o repositório: `git clone https://github.com/GiovaniGitHub/service-b.git`
3. Navegue até o diretório do projeto: `cd service-b`
4. Crie um .env a partir do .env.template e altere o campo
**Exemplo**
```bash
WEB_SERVER_PORT=8080
URL_BASE=http://localhost
```

### Rodar Sem Docker
 - Requisitos basicos:
   - Golang v1.19

```bash
    make run # Roda o projeto
```

```bash
    make test # Executa os testes
```

```bash
    make all # Executa os testes e o projeto
```

### Rodar Com Docker
 - Requisitos basicos:
   - Docker
- Altere o campo **CONTAINER_NAME** no arquivo makefile 

```bash
    make build-docker # Cria a imagem docker do projeto
```

```bash
    make run-docker # Roda o projeto
```

### Rodar com Docker Compose
 - Requisitos basicos:
   - Docker

```bash
    docker compose -f docker-compose.yml up -d # Roda o projeto # Cria a imagem docker do projeto
```


### Teste da API

#### Usando CURL

| Comando | Resultado                             |
|---------|-----------------------------------------|
| curl -X 'GET' 'http://localhost:8080/cep/70070080' -H 'accept: application/json' | {"temp_C":"36","temp_F":"96.80","temp_K":"309.00"} |
| curl -X 'GET' 'http://localhost:8080/cep/7007008A' -H 'accept: application/json' | invalid zipcode |
| curl -X 'GET' 'http://localhost:8080/cep/70070081' -H 'accept: application/json' | can not found zipcode |

#### Usando Swagger
- Acessar: http://localhost:8080/docs/index.html


### Aplicação em execução no Google Cloud Run

1. Em produção a aplicação esta rodando no Google Cloud Run.
2. Segue um teste possivel

```bash
    curl -H "Content-Type: application/json" https://service-b-prqp4ppyua-uc.a.run.app/cep/70070080
```
Onde a saída possível é:
```json
{"temp_C":"28","temp_F":"82.40","temp_K":"301.00"}
```

**Obs.**
- Caso queira subir a aplicação em usa conta no Google Cloud Run
é necessário modificar o Dockerfile e setar o valor da porta para 
expor a aplicação (que geralmente é a porta 8080)