Antilope Backend
-

A small go backend including user creation, issuing JWTs and limiting access to certain content.

###  Routes:
|Method    |  Route           | Description |
|-          |-                  |-          |
|GET    |  /api/ping/           | Get a ping               |
|POST   |  /api/users/          | Create a user            | 
|GET    |  /api/users/          | Get all users (debug)    |
|DELETE |  /api/users/          | Delete all users (debug) |
|POST   |  /api/users/login/    | Get a JWT                |
|POST   |  /api/users/logout/   | Invalidate a JWT         |
|GET    |  /api/secrets/        | Get a secret             |
|POST   |  /api/secrets/        | Create a secret          |
    
    
### Redis

    docker pull redis

#####    Volume

    docker volume create redis-data
    
    
##### Running

    docker run -d \
      -h redis \
      -e REDIS_PASSWORD=redis \
      -v redis-data:/data \
      -p 6379:6379 \
      --name redis \
      --restart always \
      redis:5.0.5-alpine3.9 /bin/sh -c 'redis-server --appendonly yes --requirepass ${REDIS_PASSWORD}'
      
##### Remove

    docker rm -f redis
    docker volume rm redis-data
    
    
##### Stopping:

`$ docker stop redis` or `$ docker rm redis`