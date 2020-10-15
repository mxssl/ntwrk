# ntwrk

Simple web app that you can use to check your external IP address.
Hosted on http://sre.monster

Usage:

```sh
curl ntwrk.duckdns.org
```

## Install

If you want to run your own instance of the app.

Create `.env` file with the following content:

```sh
PORT=80
```

Create `docker-compose.yml` file:

```yaml
version: '2'

services:
  ntwrk:
    image: mxssl/ntwrk:0.0.4
    env_file: .env
    restart: always
    network_mode: host
```

Pull and start the app container:

```sh
docker-compose pull
docker-compose up -d
docker-compose logs
```
