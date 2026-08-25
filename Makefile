.PHONY: build test run docker-build kind-load deploy

build:
	go build ./...

test:
	go test ./...

run:
	@echo "cmd/server and cmd/worker not implemented yet — Week 1 build"

docker-build:
	docker build -t interview-ready-2027-server -f cmd/server/Dockerfile .
	docker build -t interview-ready-2027-worker -f cmd/worker/Dockerfile .

kind-load:
	kind load docker-image interview-ready-2027-server --name interview-ready
	kind load docker-image interview-ready-2027-worker --name interview-ready

deploy:
	kubectl apply -f deploy/k8s/
