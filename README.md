## Proof of Concept: Request coalescing middleware for reducing database node hotspots during traffic spikes

# Components

`messages-api`
* Directly performs data mutations on the `scylla` database cluster.
* Forwards query requests to the message-query-service via the `envoy-load-balancer`.

`envoy-load-balancer`
* Routes requests to the appropriate node based on the x-bucket-id in the request header.

`message-query-service`
* For each incoming query, spawns a worker to fetch results from the database and waits for the response.
* Subsequent queries with identical parameters within a specified time frame are subscribed to the existing worker, significantly reducing traffic to the database.
* Once the worker fetches the result from the database, it returns the result to all subscribers.
* The service is horizontally scalable.
* While scaling, some partitions may be reallocated, but this does not impact ongoing workers or results as the service is designed to be stateless.

![request-coalescing-Page-1 drawio](https://github.com/yousuf64/data-service/assets/77720223/c755bcc1-7ee1-45db-8795-918052d2c1c9)

![request-coalescing-Page-1 drawio (1)](https://github.com/yousuf64/data-service/assets/77720223/fcce4b4b-c805-4611-8052-147012913220)
