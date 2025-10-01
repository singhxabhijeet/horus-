# Horus - A Cloud-Native Microservices Application 🚀

![Go](https://img.shields.io/badge/Go-1.22+-blue.svg)
![Docker](https://img.shields.io/badge/Docker-24.0-blue?logo=docker)
![Kubernetes](https://img.shields.io/badge/Kubernetes-v1.31-blue?logo=kubernetes)

"Horus" is a project that demonstrates the evolution of a simple monolithic application into a scalable, resilient, and distributed system using modern cloud-native principles. It features decoupled microservices, asynchronous communication via a message queue, and declarative deployment on Kubernetes.

---

## Architecture

The system is designed with a producer/consumer pattern, decoupled by a message queue.

```
+-----------+     +----------------+     +-----------------+     +-------------------+     +---------+
|           |     |                |     |                 |     |                   |     |         |
|   User    +---->|   Horus API    +---->|    RabbitMQ     +---->|  Notifier Worker  +---->| Discord |
|           |     | (Go Producer)  |     | (Message Queue) |     |   (Go Consumer)   |     |         |
+-----------+     +----------------+     +-----------------+     +-------------------+     +---------+
                  (Writes to PGSQL)
```

1.  The **Horus API** receives a request to monitor a site.
2.  Its background worker performs a health check and publishes the result as a message to a **RabbitMQ** queue.
3.  The **Notifier Worker**, an independent microservice, consumes messages from the queue. If a message indicates a site is down, it sends a formatted alert to a **Discord** channel.

---

## Features

* **Event-Driven Architecture**: Services communicate asynchronously through messages, improving resilience and scalability.
* **Decoupled Microservices**: The `Horus` producer and `Notifier` consumer are completely independent, allowing them to be developed, deployed, and scaled separately.
* **Kubernetes Deployment**: The entire multi-service application is deployed declaratively on a local Kubernetes cluster using `kubectl` and manifest files.
* **CI/CD with GHCR**: A GitHub Actions workflow automatically builds and pushes tagged container images for each microservice to the GitHub Container Registry on every commit to `main`.
* **Real-time Discord Alerts**: The notifier worker sends richly formatted downtime alerts to a Discord channel via webhooks.

---

## Tech Stack

* **Language**: Go
* **Containerization**: Docker, Docker Compose
* **Orchestration**: Kubernetes (Minikube)
* **Message Broker**: RabbitMQ
* **Database**: PostgreSQL
* **CI/CD**: GitHub Actions, GitHub Container Registry (GHCR)
* **Notifications**: Discord Webhooks

---

## Running the Project

There are two ways to run this project: a simple local setup with Docker Compose, or the full deployment on Kubernetes.

### 1. Local Development (with Docker Compose)

This is the quickest way to get the services running on your machine.

**Prerequisites:**
* Git
* Docker and Docker Compose

**Setup:**
1.  Clone the repository:
    ```bash
    git clone [https://github.com/your-username/horus.git](https://github.com/your-username/horus.git)
    cd horus
    ```
2.  Create an environment file named `.env` in the root of the project and paste the following, adding your Discord Webhook URL.
    ```
    DATABASE_URL=postgres://horus:horuspass@localhost:5432/horus_db
    RABBITMQ_URL=amqp://user:password@localhost:5672/
    DISCORD_WEBHOOK_URL="YOUR_WEBHOOK_URL_HERE"
    ```
3.  Build and run the application stack:
    ```bash
    docker-compose up --build
    ```
    The API will be available at `http://localhost:8080`.

### 2. Kubernetes Deployment (with Minikube)

This simulates a production-style deployment on a Kubernetes cluster.

**Prerequisites:**
* Git, Docker, `kubectl`
* [Minikube](https://minikube.sigs.k8s.io/docs/start/)

**Setup:**
1.  **Start Minikube** with a reliable DNS setting:
    ```bash
    minikube start --driver=docker --dns-ip=8.8.8.8
    ```
2.  **Create the Discord Webhook secret.** Replace `YOUR_URL_HERE` with your actual webhook URL.
    ```bash
    kubectl create secret generic discord-webhook --from-literal=DISCORD_WEBHOOK_URL='YOUR_URL_HERE'
    ```
3.  **(If your repo is private)** Create a secret to pull images from GHCR. Replace `YOUR_USERNAME` and `YOUR_PAT`.
    ```bash
    kubectl create secret docker-registry ghcr-secret --docker-server=ghcr.io --docker-username=YOUR_USERNAME --docker-password=YOUR_PAT
    ```
4.  **Deploy the application** by applying the manifest file:
    ```bash
    kubectl apply -f k8s/all-in-one.yml
    ```
5.  **Create the database tables.** Wait for the postgres pod to be `Running` (`kubectl get pods`), then run:
    ```bash
    kubectl exec -it <your-postgres-pod-name> -- psql -U horus -d horus_db
    ```
    Inside the `psql` shell, paste the `CREATE TABLE` SQL from Project 1.

6.  **Access the API.** Get the service URL by running:
    ```bash
    minikube service horus-app-service --url
    ```
    Use the URL provided by this command in Postman.

---

## Key Concepts Demonstrated

* **System Design & Architecture**: Evolving a monolith to a distributed microservices system.
* **Asynchronous Processing**: Using message queues (RabbitMQ) to build resilient, decoupled systems.
* **Cloud-Native Orchestration**: Managing a multi-container application on Kubernetes.
* **Declarative Infrastructure**: Defining the entire application state using Kubernetes manifest files.
* **Advanced CI/CD**: Building and publishing container images to a remote registry (GHCR).
* **Advanced Troubleshooting**: Debugging complex, real-world issues related to DNS, networking, and container builds (`scratch` vs. `alpine`) in a Kubernetes environment.