#!/bin/bash
SERVICE_NAME="sumologic-server"

# Executa os testes do Go
echo "Executando testes do Go..."
go test ./...
if [ $? -ne 0 ]; then
    echo "Um ou mais testes falharam. Saindo..."
    exit 1
fi

echo "Todos os testes do Go passaram."

echo "Iniciando o servidor..."
docker-compose up -d $SERVICE_NAME

sleep 5

echo "Realizando testes de integração..."

JSON_PAYLOAD='{"command":["/scripts/await.sh","10"],"timeout":1000}'

# Realiza a chamada TCP
echo "Enviando requisição TCP..."
go run main.go client -p 3000 --script "./scripts/task_await_linux.sh" --script "60" -t 20000

sleep 30
# Interrompe o servidor
echo "Parando o servidor..."
docker-compose down

echo "Todos os testes e integração concluídos."