build:
	docker-compose build

up:
	@make build
	docker-compose -f docker-compose.yml up

logs:
	docker-compose logs -f

rm:
	docker-compose rm  -sfv

start:
	docker-compose start

stop:
	docker-compose stop

check:
	CGO_ENABLED=0 go vet ./...
	staticcheck -checks=all ./...

test:
   manually cd into module and run:
	go test -v -cover ./...