# Router App

## Init go project

```sh
go mod init github.com/ortisan/router-go
go get -u github.com/rs/zerolog/log
go get -u github.com/spf13/viper
go get -u github.com/gin-gonic/gin
go get -u github.com/prometheus/client_golang/prometheus/promhttp
go get go.etcd.io/etcd/client/v3
```

## Running App

- Start app:

  ```sh
  docker-compose up --build
  ```

- Use postman collection (**Router.postman_collection.json**) to test apis.

- Grafana is enabled on 3000 and Prometheus on 9000


