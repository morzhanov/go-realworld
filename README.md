# go-realworld

TBD: Go Real World example

## services

- auth
  - should authenticate requests and issue tokens
  - should contain endpoint which ingress will use as proxy to check requests
- users
  - should operate users CRUD
- pictures
  - simple picture saving service
  - user can save pictures in postgresql db in base64 format
  - user can get all pictures
  - user can get picture by id
  - user can delete picture
- analytics
  - should handle requests from other services and log alanytics data into kafka
- client
  - should provide web ui for login and requests calls
  - should provide choice of 3 different communication protocols (rest/grpc/kafka events)
