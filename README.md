# ntwrk

This is a simple app that you can use to check your external IP address. Hosted on [http://sre.monster](http://sre.monster)

Usage:

```sh
curl sre.monster
```

## Install

If you want to run your own instance of the app, create a `.env` file with the following content:

```sh
# The app starts on this port
PORT=80
# The mode can be set to either "native" or "cloudflare"
# If you use "cloudflare" mode, then the app responds with the value of the HTTP header "CF-Connecting-IP"
# https://developers.cloudflare.com/fundamentals/reference/http-request-headers/#cf-connecting-ip
MODE=native
```

Create `docker-compose.yml` file:

```yaml
version: '3'

services:
  ntwrk:
    image: mxssl/ntwrk:0.1.7
    env_file: .env
    restart: always
    # For "native" mode, you need to use the host network mode
    network_mode: host
```

Pull and start the container:

```sh
docker compose pull
docker compose up -d
docker compose logs
```
