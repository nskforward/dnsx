run:
	go build -ldflags="-w -s" -o ./build/app_local ./cmd/*
	./build/app_local ca

deploy:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o ./build/app_linux_amd64 ./cmd/*
	@echo "----------------------------------------------------------------------------------------------------------------"
	@echo "------------------------------------------------- DEPLOYMENT ---------------------------------------------------"
	@echo "----------------------------------------------------------------------------------------------------------------"
	$(eval host := ivan@85.209.2.237)
	ssh $(host) sudo systemctl stop dns.service
	scp ./build/app_linux_amd64 $(host):/app/app
	scp -r ./config/config.json $(host):/app/config.json
	ssh $(host) sudo setcap 'cap_net_bind_service=+ep' /app/app
	ssh $(host) sudo chmod +wx /app/app
	ssh $(host) sudo systemctl start dns.service
	ssh $(host) sudo journalctl -u dns.service -n 30 -f