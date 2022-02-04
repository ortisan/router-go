# Router App

## Running App

- Start app:

  ```sh
  docker-compose up --build
  ```

- If you need run local for debugging, comment **app** service into **docker-compose.yaml** and start using **run command** or using **debugger (F5 main.go on VSCode)**.

- Use postman collection (**Router.postman_collection.json**) to test apis.

- Grafana is enabled on 3000 and Prometheus on 9000

## TODO

- Jaeger to trace requests
- New app to control healthcheck (all instances of this app is updating cache - race conditions and multiple calls to backends healthcheck)
- Run on Kubernetes
- Performance tests

## Helpful Commands

```sh
# Init go lang project
go mod init github.com/ortisan/<project_name>
# Import module
go get -u <module>
# Run main
go run main.go
# Build project
go build go build -o <sh/exe name>
```
