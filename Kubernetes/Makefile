# === CONFIG ===
REGISTRY=local
AUTH_IMG=auth-service:$(REGISTRY)
POLICY_IMG=policy-service:$(REGISTRY)
THREAT_IMG=threatlog-service:$(REGISTRY)

# === TARGETS ===

## build: Compila las tres imágenes dentro de Minikube
build:
	eval $$(minikube -p minikube docker-env) && \
	echo "🛠️ Building Auth Service..." && \
	docker build -t $(AUTH_IMG) -f ../auth-service/Dockerfile .. && \
	echo "🛠️ Building Policy Service..." && \
	docker build -t $(POLICY_IMG) -f ../policy-service/Dockerfile .. && \
	echo "🛠️ Building ThreatLog Service..." && \
	docker build -t $(THREAT_IMG) -f ../threatlog-service/Dockerfile ..

## apply: Aplica todos los Deployments y Services
apply:
	kubectl apply -f auth-deploy.yaml
	kubectl apply -f policy-deploy.yaml
	kubectl apply -f threatlog-deploy.yaml

## delete: Elimina todos los recursos del cluster
delete:
	kubectl delete -f threatlog-deploy.yaml --ignore-not-found
	kubectl delete -f policy-deploy.yaml --ignore-not-found
	kubectl delete -f auth-deploy.yaml --ignore-not-found

## status: Muestra el estado de pods, servicios y deployments
status:
	kubectl get pods,svc,deploy

## logs: Muestra logs de un servicio específico (usa make logs SVC=auth)
logs:
	@if [ -z "$(SVC)" ]; then \
		echo "❌ Usa make logs SVC=<auth|policy|threatlog>"; \
	else \
		POD=$$(kubectl get pods -l app=$(SVC) -o jsonpath='{.items[0].metadata.name}'); \
		kubectl logs -f $$POD; \
	fi

## url: Obtiene el URL del Auth Service (LoadBalancer)
url:
	minikube service auth-svc --url

## restart: Reinicia los deployments (rolling update)
restart:
	kubectl rollout restart deploy auth-deploy
	kubectl rollout restart deploy policy-deploy
	kubectl rollout restart deploy threatlog-deploy
	kubectl rollout status deploy auth-deploy
	kubectl rollout status deploy policy-deploy
	kubectl rollout status deploy threatlog-deploy

## clean: Limpia imágenes locales
clean:
	eval $$(minikube -p minikube docker-env) && \
	docker rmi $(AUTH_IMG) $(POLICY_IMG) $(THREAT_IMG) || true

.PHONY: build apply delete status logs url restart clean
