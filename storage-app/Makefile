ENV_FILE=.env

run: env
	go run main.go

env:
	@echo "Generating .env from template..."
	@envsubst < $(TEMPLATE_FILE) > $(ENV_FILE)

clean:
	rm -f $(ENV_FILE)
