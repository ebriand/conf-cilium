.PHONY: build push deploy

build:
	docker image build -t ebriand/conf-cilium-events-frontend:0.2 .

push: build
	docker image push ebriand/conf-cilium-events-frontend:0.2

deploy: push
	kubectl apply -f k8s && kubectl delete -n events pods --all