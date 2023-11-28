## About

This repo is forked from [temporal-select-signal-tutorial-code](https://github.com/unijad/temporal-select-signal-tutorial-code) and refactored to make it work e2e

## Setup

To get started, clone this repository or download the source code.

Next, install the dependencies using go mod:

`go mod tidy`

Start the Temporal Server either through Docker or manually on your machine.

## Docker Setup

For local testing and development, you can use Docker to run Temporal. The Temporal team provides an official Docker image that allows running the entire Temporal stack locally.Ensure Docker is installed on your machine by downloading it from the official website.Once Docker is installed, utilize the official docker-compose setup provided by Temporal to run the complete Temporal stack locally. Run the following command to acquire the docker-compose setup files:

`curl -o docker-compose.yml https://raw.githubusercontent.com/temporalio/docker-compose/master/docker-compose-cas.yml`

This command downloads the docker-compose setup for Temporal Community Edition, which incorporates a Cassandra data store.Now, start the services by running the following command:

`docker-compose up`

## Running the Example

### bring up the application
```shell
 go run cmd/shopper/main.go -conf ./config/dev.yaml
```

### add products to a cart and submit the order
```shell
curl -X POST http://localhost:8086/shopping/carts\?products\=1,2,3
curl -X POST http://localhost:8086/shopping/orders?products=1,2,3
```

### signal the order to continue the process 
```shell
curl -X POST http://localhost:8086/shopping/orders/signal?orderId={orderID}&signalName=confirmInvoice&status=confirmed
curl -X POST http://localhost:8086/shopping/orders/signal?orderId={orderID}&signalName=confirmShipping&status=confirmed
```