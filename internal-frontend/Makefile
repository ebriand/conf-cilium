.PHONY: build push

build:
	docker image build -t ebriand/conf-cilium-internal-frontend:0.1 .

push: build
	docker image push ebriand/conf-cilium-internal-frontend:0.1

deploy: push
	kubectl apply -f k8s