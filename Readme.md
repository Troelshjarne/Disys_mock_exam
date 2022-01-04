# do the following steps, to start up system
## start replica(s)  (individual terminals)
  1. cd replica/
  2. go run replica.go -port 127.0.0.1:9081
  3. go run replica.go -port 127.0.0.1:9082
  4. go run replica.go -port 127.0.0.1:9083

## start frontend(s)(individual terminals)
  1. cd frontend/
  2. go run frontend.go -port 127.0.0.1:9084
  3. go run frontend.go -port 127.0.0.1:9085

## start client(s)(individual terminals)
  1. cd client/
  2. go run client.go