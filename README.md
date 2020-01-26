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
    