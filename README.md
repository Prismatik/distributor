# Distributor

## What?

This is a Docker container that distributes traffic to other Docker containers. It contains an nginx reverse proxy that you can configure entirely from environment variables. That way, you can expose a single port from your Docker cluster to the internet, then use this to route traffic internally via the links system.

It's primarily intended for use with Amazon's ECS service and an ELB. An ELB can't map traffic to more than one kind of container in an ECS service. With Distributor, the ELB can map traffic only to the Distributor container, then you can route to other containers from here.

```
      Users

         +--+ "Take me to ohai.example.com"                   Distributor
         |  |
      +--+--+--+            +-------------+  Port   +-------------------------------+
+-----+        +-----+      |             |  80     |                               |
+-----+        +-----+      |  ELB        +------>  |  LISTEN_PORT=80               |
      |        |     +----> |             |         |  DISTRIBUTOR_OHAI=ohai,3000   |
      |        |            +-------------+         |  DISTRIBUTOR_THERE=there,3000 |
      +-+----+-+                                    |  DISTRIBUTOR_LOL=lol,4000     |
      | |    | |                                    |  DOMAIN=.example.com          |
      | |    | |                                    +----+--------------------------+
      | |    | |                                         |                       
      +-+    +-+                                         |                       
                                                         |                       
                                                      +--v-----+ +--------+ +--------+
                                                      |        | |        | |        |
                                                      | ohai   | | there  | |  lol   |
                                                      |        | |        | |        |
                                                      |        | |        | |        |
                                                      |        | |        | |        |
                                                      |        | |        | |        |
                                                      +--------+ +--------+ +--------+

                                                                Containers
```

## How?

Distributor is configured via environment variables. The ones it needs are:

* LISTEN_PORT - The port to listen for incoming connections on
* HEALTH_PORT - The port to respond to `/health` on
* DOMAIN - The domain suffix to append to incoming names
* DISTRIBUTOR_{{NAME,PORT}} - Names and ports for containers

So a valid command would be:

```
docker run \
-e LISTEN_PORT=80 \
-e HEALTH_PORT=3000 \
-e DOMAIN=example.com \
-e DISTRIBUTOR_OHAI=ohai,4000 \
Prismatik/distributor
```

Or in docker-compose speak:

```
version: '2'
services:
  test:
    image: Prismatik/distributor
    environment:
      HEALTH_PORT: 3001
      LISTEN_PORT: 3000
      DOMAIN: .example.com
      DISTRIBUTOR_OHAI: ohai,80
      DISTRIBUTOR_LOL: lol,80
    ports:
      - "3000-3001:3000-3001"
    links:
      - ohai
      - lol
  ohai:
    image: nginx
    volumes:
      - /tmp/ohai:/usr/share/nginx/html:ro
  lol:
    image: nginx
    volumes:
      - /tmp/lol:/usr/share/nginx/html:ro
```

**protip**: For great success, you probably want to configure your `DOMAIN` as a wildcard DNS record on a subdomain, with the whole thing pointing to your ELB. ie:

`*.testcluster.example.com --> CNAME some.elb.domain`

That way you can drop new services into your cluster without needing to fool with your DNS.
