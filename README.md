# chat

## Overview

Wedsocket chat server, without database.


## Usage
### Running locally
```bash
git clone https://github.com/sanskarsharma/chat.git
cd chat
go run *.go
```
### Running on docker
```bash
git clone https://github.com/sanskarsharma/chat.git
cd chat
docker build -t chat:v-local .
docker run -d -p 8080:8080 chat:v-local
```

### Running with reverse proxy and minimal front-end on docker-compose
```
docker-compose up -d 
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to add / update tests as appropriate.
