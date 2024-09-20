
![Logo](https://ik.imagekit.io/pj3r6oe9k/prasorganic-high-resolution-logo-transparent.svg?updatedAt=1726835541390)

# Prasorganic Order Service

Prasorganic Order Service is one of the components in the Prasorganic Microservice architecture built with Go (Golang). This service supports operations to manage user orders via RESTful API, gRPC, and Message Broker.

## Tech Stack

[![My Skills](https://skillicons.dev/icons?i=go,docker,postgresql,redis,kafka,bash,git&theme=light)](https://skillicons.dev)

## Features

- **Order Management:** Supports operations for creating, retrieving, and deleting orders.

- **RESTful API:** Provides a RESTful API using Fiber with various middleware for managing requests and responses.

- **gRPC:** Utilizes gRPC for inter-service communication, equipped with interceptors for handling requests and responses.

- **Message Broker:** Consumes messages from Kafka for notification processing.

- **Task Management:** Uses Asynq for scheduling and managing order tasks to third-party shippers.

- **Database:** Uses PostgreSQL for data storage with database migration support.

- **Logging:** Logs are recorded using Logrus.

- **Error Handling:** Implements error handling to ensure proper detection and handling of errors, minimizing the impact on both the client and server.

- **System Resilience:** Uses a Circuit Breaker to enhance application resilience and fault tolerance, protecting the system from cascading failures.

- **Configuration and Security:** Employs Viper and HashiCorp Vault for integrated configuration and security management.

- **Testing:** Implements unit testing using Testify.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

This project makes use of third-party packages and tools. The licenses for these
dependencies can be found in the `LICENSES` directory.

## Dependencies and Their Licenses

- `Go:` Licensed under the BSD 3-Clause "New" or "Revised" License. For more information, see the [Go License](https://github.com/golang/go/blob/master/LICENSE).

- `Docker:` Licensed under the Apache License 2.0. For more information, see the [Docker License](https://github.com/docker/docs/blob/main/LICENSE).

- `Docker Compose:` Licensed under the Apache License 2.0. For more information, see the [Docker Compose License](https://github.com/docker/compose/blob/main/LICENSE).

- `PostgreSQL:` Licensed under PostgreSQL License. For more information, see the [PostgreSQL License](https://www.postgresql.org/about/licence/).

- `Redis:` Follows a dual-licensing model with RSALv2 and SSPLv1. For more information, see the [Redis License](https://redis.io/legal/licenses/).

- `RedisInsight:` Licensed under the RedisInsight License. For more information, see the [RedisInsight License](https://github.com/RedisInsight/RedisInsight/blob/main/LICENSE).

- `Kafka:` Licensed under the Apache License 2.0. For more information, see the [Kafka License](https://github.com/apache/kafka/blob/trunk/LICENSE).

- `Zookeeper:` Licensed under the Apache License 2.0. For more information, see the [Zookeper License](https://github.com/apache/zookeeper/blob/master/LICENSE.txt).

- `GNU Make:` Licensed under the GNU General Public License v3.0. For more information, see the [GNU Make License](https://www.gnu.org/licenses/gpl.html).

- `GNU Bash:` Licensed under the GNU General Public License v3.0. For more information, see the [Bash License](https://www.gnu.org/licenses/gpl-3.0.html).

- `Git:` Licensed under the GNU General Public License version 2.0. For more information, see the [Git License](https://opensource.org/license/GPL-2.0).

## Thanks 👍
Thank you for viewing my project.