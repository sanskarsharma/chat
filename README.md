# chat

## Overview

Wedsocket chat server. Completely ephemeral, no logging/storage on server.


## Usage
### Running on local with golang
```bash
go run *.go
```
### Running on docker
```bash
docker build -t chat:v-local .
docker run -d -p 8080:8080 chat:v-local
```

### Running with reverse proxy and minimal front-end on docker-compose
```bash
# note : check docker-compose.yaml and modify as required before running this
docker-compose up -d 
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to add / update tests as appropriate.
