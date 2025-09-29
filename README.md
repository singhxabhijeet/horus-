# Horus - A Concurrent Web Service Health Checker 🩺

![GitHub Actions CI](https://github.com/singhxabhijeet/horus/actions/workflows/ci.yml/badge.svg)
![Go Version](https://img.shields.io/badge/Go-1.22+-blue.svg)
![Docker](https://img.shields.io/badge/Docker-24.0-blue?logo=docker)

Horus is a backend monitoring service written in Go that periodically checks the health of registered websites and APIs. It's built with professional practices in mind, including full containerization with Docker, a concurrent worker model, and a CI pipeline for automated builds.

---

## Problem Solved

In any distributed system, knowing the health of your services is critical. This project tackles that foundational concept of reliability and monitoring. Instead of a simple CRUD application, Horus solves a real-world problem: providing a lightweight, API-driven way to track service uptime, status codes, and response times.

---

## Features

* **RESTful API**: Endpoints to add and list sites for monitoring.
* **Concurrent Background Worker**: A non-blocking worker, built with Goroutines, checks all registered sites simultaneously every minute.
* **Data Persistence**: Health check results (status, response time, etc.) are stored in a PostgreSQL database.
* **Fully Containerized**: The entire stack (Go application + PostgreSQL database) is defined in Docker Compose for one-command setup and execution.
* **CI Pipeline**: A GitHub Actions workflow automatically builds the application on every push to ensure code integrity.

---

## Tech Stack

* **Language**: **Go**
* **Database**: **PostgreSQL**
* **Containerization**: **Docker** & **Docker Compose**
* **CI/CD**: **GitHub Actions**

---

## Getting Started

You can get the entire application running on your local machine with just two commands.

### Prerequisites

* Git
* Docker and Docker Compose

### Installation & Setup

1.  **Clone the repository:**
    ```bash
    git clone [https://github.com/your-username/horus.git](https://github.com/your-username/horus.git)
    cd horus
    ```

2.  **Create an environment file:**
    Create a file named `.env` in the root of the project and paste the following line into it. This file is ignored by Git.
    ```
    DATABASE_URL="postgres://horus:horuspass@db:5432/horus_db"
    ```

3.  **Run the application:**
    Use Docker Compose to build the images and start the containers.
    ```bash
    docker-compose up --build
    ```
    The API will be available at `http://localhost:8080`.

---

## API Usage

You can interact with the API using any client like `curl` or Postman.

### 1. Add a new site to monitor

* **Endpoint**: `POST /api/sites`
* **Body**:
    ```json
    {
        "url": "[https://www.google.com](https://www.google.com)"
    }
    ```
* **Example `curl` command:**
    ```bash
    curl -X POST -H "Content-Type: application/json" -d '{"url": "[https://www.google.com](https://www.google.com)"}' http://localhost:8080/api/sites
    ```
* **Success Response**: `201 Created`
    ```json
    {
        "id": 1
    }
    ```

### 2. List all monitored sites

* **Endpoint**: `GET /api/sites`
* **Example `curl` command:**
    ```bash
    curl http://localhost:8080/api/sites
    ```
* **Success Response**: `200 OK`
    ```json
    [
        {
            "id": 1,
            "url": "[https://www.google.com](https://www.google.com)",
            "created_at": "2025-09-30T02:18:19Z"
        }
    ]
    ```

---

## Key Concepts Demonstrated

This project was built to demonstrate proficiency in several key areas of modern backend and DevOps engineering:

* **Backend API Development**: Building a clean, RESTful API in Go.
* **Concurrency in Go**: Using Goroutines and Channels to perform concurrent, non-blocking background tasks.
* **Database Management**: Designing a SQL schema and interacting with a relational database (PostgreSQL).
* **Containerization**: Creating reproducible environments with Docker and Docker Compose.
* **Build Optimization**: Implementing a multi-stage `Dockerfile` to produce a minimal (~15MB), secure final image from a `scratch` base.
* **CI/CD Automation**: Establishing an automated build pipeline using GitHub Actions.
