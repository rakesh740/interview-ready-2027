.PHONY: build test run docker-build kind-load deploy

SERVICES := api dispatcher worker

build:
	go build ./...

test:
	go test ./...

run:
	@echo "cmd/api, cmd/dispatcher, cmd/worker not implemented yet — Week 1 build"

docker-build:
	@for svc in $(SERVICES); do \
		docker build -t interview-ready-2027-$$svc -f cmd/$$svc/Dockerfile . ; \
	done

kind-load:
	@for svc in $(SERVICES); do \
		kind load docker-image interview-ready-2027-$$svc --name interview-ready ; \
	done

deploy:
	kubectl apply -f deploy/k8s/
