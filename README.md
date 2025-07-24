a demo microservice application with 3 services (simple order, item and user services), built using Golang for rest apis, postgresql for DB, docker for building image and pushing to docker hub and basic Kubernetes local setup on top of LIMA virtual machines using kubeadm to test basic scaling and other features.
Note: LIMA vm supports 2 type of vm, vz - for macOS and qemu for others OS. Above setup work for vz type VM (only on macOS silicon). Alternatively, You can get idea how and with minor config changes, you can setup on other OS too.

Automation scripts includes lima vm provisioning, k8 cluster setup and starting the app and proxy server without any manual setup.



How to run:
1. clone the repo
2. cd to ms-demo/lima-vm-config
3. run `bash k8infraprovision.sh` - this will create and start VM and setup K8 cluster, control plane on master and join worker and copy the automation script which will start k8 app and setup proxy automatically
4. run `limactl shell master` - go inside master
5. cd to start/
6. run `chmod +x start_k8_app_and_proxy_server.sh`
7. run `bash start_k8_app_and_proxy_server.sh` - this will clone app repo from github and start k8 app and
 ask if you want proxy server so that you can communicate from host... () and run proxy in background







# Microservices Demo Application

This is a simple microservices-based application built with Go, consisting of three services (User, Item, and Order) and an API Gateway.

## Prerequisites

- Go 1.16 or later
- Docker and Docker Compose
- PostgreSQL (will be run via Docker)

## Services

1. **User Service** (Port 8081)
   - Login endpoint
   - Get user details endpoint
   - Pre-configured with sample users

2. **Item Service** (Port 8082)
   - Add item endpoint
   - Get item details endpoint

3. **Order Service** (Port 8083)
   - Create order endpoint
   - Communicates with User and Item services

4. **API Gateway** (Port 8080)
   - Routes requests to appropriate services
   - Provides a unified API endpoint

## Getting Started

1. Start the databases:
   ```bash
   docker-compose up -d
   ```

2. Start each service (in separate terminals):
   ```bash
   # Terminal 1 - API Gateway
   cd api-gateway
   go run main.go

   # Terminal 2 - User Service
   cd user-service
   go run main.go

   # Terminal 3 - Item Service
   cd item-service
   go run main.go

   # Terminal 4 - Order Service
   cd order-service
   go run main.go
   ```

## API Endpoints

### User Service
- Login: `POST /api/login`
  ```json
  {
    "username": "user1",
    "password": "pass1"
  }
  ```

### Item Service
- Add Item: `POST /api/items`
  ```json
  {
    "name": "Example Item",
    "quantity": 10
  }
  ```

### Order Service
- Create Order: `POST /api/orders`
  ```json
  {
    "user_id": 1,
    "item_id": 1,
    "quantity": 2
  }
  ```

## Sample Users
The application comes with two pre-configured users:
1. Username: `user1`, Password: `pass1`
2. Username: `user2`, Password: `pass2`

## Testing the Application

1. First, login with a user:
   ```bash
   curl -X POST http://localhost:8080/api/login -H "Content-Type: application/json" -d '{"username": "user1", "password": "pass1"}'
   ```

2. Add an item:
   ```bash
   curl -X POST http://localhost:8080/api/items -H "Content-Type: application/json" -d '{"name": "Test Item", "quantity": 100}'
   ```

3. Create an order:
   ```bash
   curl -X POST http://localhost:8080/api/orders -H "Content-Type: application/json" -d '{"user_id": 1, "item_id": 1, "quantity": 2}'
   ``` 