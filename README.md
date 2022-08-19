# Demo Saga implementation with Golang and Redis Streams

Implementation idea inspired by following article:

1. [Distributed Transactions in Microservices with Kafka Streams and Spring Boot](https://piotrminkowski.com/2022/01/24/distributed-transactions-in-microservices-with-kafka-streams-and-spring-boot/) - how to implement distributed transaction based on the SAGA pattern with Spring Boot and Kafka Streams

## Description

There are three microservices:

`order-service` - it sends `Order` events to the Redis streams and orchestrates the process of a distributed transaction

`payment-service` - it performs local transaction on the customer account basing on the `Order` price

`stock-service` - it performs local transaction on the store basing on number of products in the `Order`

## Installation And Configuration

### Local development

```sh
$ make help # print help message
# install required tools
$ make tools
$ make up # make down for up and down local environment
$ make run service=order # run order service
```

## License

Saga is primarily distributed under the terms of both the MIT license and the Apache License (Version 2.0).

See [LICENSE-APACHE](LICENSE-APACHE) and/or [LICENSE-MIT](LICENSE-MIT) for details.
