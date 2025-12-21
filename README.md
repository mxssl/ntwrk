# ntwrk

This is a simple app that you can use to check your external IP address. Hosted on [http://sre.monster](http://sre.monster)

Usage:

IPv4:

```sh
curl -4 sre.monster
```

IPv6:

```sh
curl -6 sre.monster
```

## Install

If you want to run your own instance of the app, create a `.env` file with the following content:

```sh
# The app starts on this port
PORT=80

# The mode can be set to "native", "cloudflare", or "proxy"
# - "native" - extracts IP from the RemoteAddr field
# - "cloudflare" - responds with the value of the "CF-Connecting-IP" HTTP header
#   https://developers.cloudflare.com/fundamentals/reference/http-request-headers/#cf-connecting-ip
# - "proxy" - for using with reverse proxies (nginx, HAProxy, etc.) that set the
#   "X-Forwarded-For" header. Returns the first IP from the header (the original client IP).
MODE=native
```

Create `docker-compose.yml` file:

```yaml
services:
  ntwrk:
    image: mxssl/ntwrk:0.1.13
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
