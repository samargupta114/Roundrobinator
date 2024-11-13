# Roundrobinator
### ( Application API Mirror Via Round Robin Algorithm ) 

Roundrobinator is a load balancing service that distributes incoming HTTP requests to multiple application instances using the round-robin algorithm. It acts as a proxy, forwarding requests to  application API  `/mirror`, 
ensuring load is balanced across multiple servers. The application includes health-check mechanisms and supports graceful shutdowns in GO.

# Key Components:
1. Round Robin Load Balancer: Balances incoming requests between multiple application API servers.
2. Health Check: Provides an endpoint to check if the service is healthy.
3. Graceful Shutdown: Ensures that the service shuts down gracefully when required.

# Features
1. Round-Robin Load Balancing: Distributes HTTP requests evenly across multiple application instances.
2. Health Checks: Exposes a `/health` endpoint to monitor service health.
3. Dynamic Routing: Automatically routes requests to the correct backend `/mirror` endpoint.
4. Graceful Shutdown: Handles graceful shutdown when the service is stopped.
5. Customizable Configurations: Backend routes and server configurations are adjustable via configuration files.

## Requirements
1. Go 1.18+ (for building and running the application)
2. An HTTP client to test APIs (such as curl or Postman)

## Setup Project

1. **Clone the repository:**
   ```sh
   git clone https://github.com/samargupta114/Roundrobinator.git
   cd Roundrobinator
   ```
2. **Resolve dependencies:**
   ```sh
   go mod tidy
   go mod vendor
   ```

3. Set config path
   **ROUND_ROBIN_CONF_PATH**="/Users/samargupta/Roundrobinator/config.json"

4. **Run the application:**
   ```sh
   go build main.go
   go run main.go
   ```
   The Application API will be running on host `localhost` at ports `8081` , `8082` , `8083`.
   The Server API for round robin will be running on host `localhost` at ports `8080` with `/route` EP.

## Current Implementation

### Application API
The `Application API` runs on three servers as of now with the ability to response the same
payload as given.

### RoundRobin
Implemented as `RoundRobin` to route to diff servers of `Application API`.

### Healthcheck
Configured with a configurable ticker for periodic health checks, triggering goroutines at the specified intervals. ( configurable through app config)

### Alerts
Currently, alerts are added as comments and not implemented using any library.

## Enhancement

1. **Metrics**:
    - Push metrics from the service for monitoring.
2. **Alerts**:
    - Push alerts from the service for alerting and monitoring.
