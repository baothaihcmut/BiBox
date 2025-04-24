ENV_FILE=.env
CONFIG_FILE=config/config.yaml
run: env
	go run cmd/storage-app/main.go

env:
	@echo "Generating .env from template..."
	@envsubst < $(TEMPLATE_FILE) > $(ENV_FILE)

clean:
	rm -f $(ENV_FILE)
