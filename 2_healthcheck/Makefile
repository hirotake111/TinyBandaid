
up:
	docker-compose up

down:
	docker-compose down

build:
	docker build -t tinybandiad:latest .

serve:
	docker run --rm --name go-toy-lb -p 3000:3000 go-toy-lb

test:
	go test ./...