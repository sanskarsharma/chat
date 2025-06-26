# chat

## Overview

Wedsocket chat server. Completely ephemeral, no logging/storage on server.
Demo [here](https://chat.pohawithpeanuts.com)


## Usage
### Running on local with golang
```bash
go run *.go
```
### Running on docker
```bash
docker build -t chat:v0 .
docker run -d -p 8080:8080 chat:v0
```

### Running with reverse proxy and minimal front-end on docker-compose
```bash
docker-compose up -d 
```

### Deploying using Clouflare {Workers + Containers}

This deploys to your cloudflare account - make sure to edit/remove custom domain in [wrangler.jsonc](wrangler.jsonc) -> `routes` as needed. 

```bash
# install
npm install
npx wrangler deploy

```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to add / update tests as appropriate.
