.PHONY: build deploy clean test run

# Configurações
FUNCTION_NAME=climia-api
REGION=us-east-1

# Build para Lambda
build:
	@echo "📦 Build para Lambda..."
	GOOS=linux GOARCH=amd64 go build -o bootstrap cmd/api/main.go
	zip function.zip bootstrap
	@echo "✅ Build concluído!"

# Deploy para Lambda
deploy: build
	@echo "🚀 Deploy para Lambda..."
	aws lambda update-function-code \
		--function-name $(FUNCTION_NAME) \
		--zip-file fileb://function.zip \
		--region $(REGION)
	@echo "✅ Deploy concluído!"

# Criar função Lambda (primeira vez)
create:
	@echo "🆕 Criando função Lambda..."
	aws lambda create-function \
		--function-name $(FUNCTION_NAME) \
		--runtime provided.al2 \
		--role arn:aws:iam::YOUR_ACCOUNT:role/lambda-execution-role \
		--handler bootstrap \
		--zip-file fileb://function.zip \
		--region $(REGION) \
		--timeout 30 \
		--memory-size 512
	@echo "✅ Função criada!"

# Testar local
test:
	@echo "🧪 Testando local..."
	go test ./...

# Rodar local
run:
	@echo "🏃 Rodando local..."
	go run cmd/api/main.go

# Limpar arquivos
clean:
	@echo "🧹 Limpando..."
	rm -f bootstrap function.zip
	@echo "✅ Limpeza concluída!"

# Ver logs da Lambda
logs:
	@echo "📋 Logs da Lambda..."
	aws logs tail /aws/lambda/$(FUNCTION_NAME) --follow

# Ver status da Lambda
status:
	@echo "📊 Status da Lambda..."
	aws lambda get-function --function-name $(FUNCTION_NAME) --region $(REGION)

# Help
help:
	@echo "📋 Comandos disponíveis:"
	@echo "  make build   - Build para Lambda"
	@echo "  make deploy  - Deploy para Lambda"
	@echo "  make create  - Criar função Lambda (primeira vez)"
	@echo "  make test    - Testar local"
	@echo "  make run     - Rodar local"
	@echo "  make clean   - Limpar arquivos"
	@echo "  make logs    - Ver logs da Lambda"
	@echo "  make status  - Ver status da Lambda" 