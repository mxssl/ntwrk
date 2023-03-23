# ntwrk

![Docker Image Version (latest semver)](https://img.shields.io/docker/v/mxssl/ntwrk)

Simple web app that you can use to check your external IP address.
Hosted on [http://sre.monster](http://sre.monster)

Usage:

```sh
curl sre.monster
```

## Install

If you want to run your own instance of the app.

Create `.env` file with the following content:

```sh
# app starts on this port
PORT=80
# mode can be "native" or "cloudflare"
# if you use cloudflare mode then the app answers with a value of HTTP header "CF-Connecting-IP"
# https://support.cloudflare.com/hc/en-us/articles/200170986-How-does-Cloudflare-handle-HTTP-Request-headers-
MODE=native
```

Create `docker-compose.yml` file:

```yaml
version: '2.4'

services:
  ntwrk:
    image: mxssl/ntwrk:0.1.5
    env_file: .env
    restart: always
    # for native mode you need to use host network mode
    network_mode: host
```

Pull and start the app container:

```sh
docker compose pull
docker compose up -d
docker compose logs
```
