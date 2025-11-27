

# DistributedSecurity – Distributed Computing

**Autor:** Efren Flores Porras  

---

## 🧩 Descripción General

**DistributedSecurity** es un sistema distribuido basado en **microservicios comunicados por gRPC**, desplegado en un **cluster de Kubernetes (Minikube)**.  

El proyecto implementa tres microservicios:

| Microservicio  | Tipo de Service | Puerto | Descripción |
|----------------|-----------------|---------|--------------|
| **Auth Service** | LoadBalancer | `8080` | Servicio expuesto al exterior; maneja autenticación y conexión con Policy y ThreatLog. |
| **Policy Service** | ClusterIP | `8082` | Valida políticas de acceso. |
| **ThreatLog Service** | ClusterIP | `8081` | Registra intentos o eventos de seguridad. |

Todos los servicios se registran dinámicamente en **Consul**, reportan su estado y se comunican por **gRPC**.

---

## ⚙️ Arquitectura

```

[Cliente externo]
│
▼
[Auth Service :8080]  <--->  [Policy Service :8082]
│
▼
[ThreatLog Service :8081]

```

---

## 📦 Requisitos previos

```markdown
- **Go 1.23+**
- **Docker**
- **Minikube + kubectl**
- **protoc** y plugins de gRPC:
  ```bash
  sudo apt install -y protobuf-compiler
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

---

## 🧱 Estructura del Proyecto

```
DistributedSecurity/
├── auth-service/           # LoadBalancer (externo)
├── policy-service/         # ClusterIP (interno)
├── threatlog-service/      # ClusterIP (interno)
├── pkg/proto/              # Archivos .proto y stubs gRPC
├── pkg/discovery/consul/   # Registro en Consul
└── Kubernetes/             # Manifests y Makefile
```

---

## 🚀 Despliegue en Kubernetes

### 1️⃣ Iniciar Minikube

```bash
minikube start
```

### 2️⃣ Compilar imágenes dentro del Docker de Minikube

```bash
cd Kubernetes
make build
```

### 3️⃣ Aplicar los Deployments y Services

```bash
make apply
```

### 4️⃣ Verificar estado

```bash
make status
```

Ejemplo de salida esperada:

```
NAME                        READY   UP-TO-DATE   AVAILABLE   AGE
auth-deploy                 2/2     2            2           1m
policy-deploy               2/2     2            2           1m
threatlog-deploy            2/2     2            2           1m
```

---

## 🌐 Acceso al Auth Service

Obtener el URL expuesto:

```bash
make url
```

Probar con **grpcurl**:

```bash
grpcurl -plaintext <URL>:<PUERTO> auth.AuthService/Login
```

> ⚠️ Nota: el comando `grpcurl` debe ejecutarse desde tu host (instálalo con `go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest`).

---

## 🧩 Comunicación entre Servicios

* `AuthService` llama a `PolicyService` por gRPC (`CheckPolicy()`).
* Si el usuario no cumple las políticas, `AuthService` envía un evento gRPC a `ThreatLogService` (`LogEvent()`).

Logs esperados:

```
[Auth] Login attempt user=blocked
[Policy] Checking policy for user=blocked
[Auth] Access denied user=blocked reason=user is blocked
[ThreatLog] Received event: user=blocked type=unauthorized_login
```

---

## 🔁 Rolling Update

Puedes actualizar la versión del servicio (por ejemplo, cambiar el mensaje de login) y ejecutar:

```bash
make build
make restart
```

---

## 🔍 Diagnóstico rápido

Ver logs de un servicio:

```bash
make logs SVC=auth
```

Eliminar todo del cluster:

```bash
make delete
```

---

## 🧾 Capacidades del Cluster

* **Cluster:** Minikube (1 nodo)
* **Microservicios:** 3 (Auth, Policy, ThreatLog)
* **Recursos asignados:**

  * CPU: 100m (request), 300m (limit)
  * Memoria: 64Mi (request), 128Mi (limit)
* **Total estimado:** ~1 CPU y ~400 MiB RAM
* **Carga esperada:** 10 usuarios concurrentes / 100 RPS promedio.

---

## 📅 Presentación

Durante la demo se mostrará:

1. **Estado del cluster:**

   ```bash
   kubectl get pods,svc
   ```

2. **URL del Auth Service (LoadBalancer).**

3. **Invocación gRPC:** demostración del flujo `Auth → Policy → ThreatLog`.

4. **Escalamiento de réplicas:**

   ```bash
   kubectl scale deploy auth-deploy --replicas=3
   ```

5. **Logs de eventos en tiempo real.**

---

## ✅ Comandos rápidos del Makefile

| Comando                    | Descripción                                             |
| -------------------------- | ------------------------------------------------------- |
| `make build`               | Construye todas las imágenes Docker dentro de Minikube. |
| `make apply`               | Aplica todos los Deployments y Services.                |
| `make status`              | Muestra el estado actual del cluster.                   |
| `make logs SVC=<servicio>` | Muestra los logs de un pod específico.                  |
| `make url`                 | Muestra el URL del Auth Service (LoadBalancer).         |
| `make restart`             | Reinicia todos los Deployments.                         |
| `make delete`              | Elimina todos los recursos del cluster.                 |

---

## 🏁 Conclusión

Este proyecto demuestra un sistema distribuido simple con:

* Comunicación **gRPC** entre microservicios.
* Registro dinámico en **Consul**.
* Despliegue orquestado en **Kubernetes**.
* Escalabilidad horizontal mediante réplicas.

```


---

¿Quieres que le agregue también un **diagrama visual del flujo** (Auth → Policy → ThreatLog) en formato `.png` o `.drawio` para incluirlo en el README y mostrarlo en tu presentación?
```
