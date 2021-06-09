# go-realworld

TBD: Go Real World example

## services

- auth
  - should authenticate requests and issue tokens
  - should contain endpoint which ingress will use as proxy to check requests
- users
  - should operate users CRUD
- photos
  - simple photo saving service
  - user can save photos in postgresql db in base64 format
  - user can get all photos
  - user can get photo by id
  - user can delete photo
- analytics
  - should handle requests from other services and log alanytics data into kafka
- client
  - should provide web ui for login and requests calls
  - should provide choice of 3 different communication protocols (rest/grpc/kafka events)
