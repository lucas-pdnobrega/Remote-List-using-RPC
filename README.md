### Remote Integer List RPC Service
This is an experiment in how to use Remote Procedure Calls to access a common 2-dimensional list array service. 
It allows clients to create, modify, and query lists concurrently over TCP.

### Tech Used

<p>
  <a href="https://skillicons.dev">
    <img src="https://skillicons.dev/icons?i=git,go,vscode" />
  </a>
</p>


### Features
- Create and remove integer lists remotely

- Append and remove elements from any list

- Retrieve elements and list sizes

- Thread-safe access using mutex locks

- Concurrent client handling via Go’s built-in RPC and goroutines

### Persistent Storage Strategy: 
The integer lists are saved to and loaded from a JSON file

### How it works
The Server exposes RemoteList methods over RPC on TCP port 5000, afterwards Clients connect and consume exposed methods like CreateList, Append, Get, Remove, and Size
State is thusly saved/loaded from a local JSON file (log.json), which supports multiple clients simultaneously without data races

### Usage
- Run the server (from repository root)
  
    ```go run .\rpc_server.go```

- On a separate terminal instance, run the client and observe the automatic interactions via RPC calls
  
    ```go run .\rpc_client.go```
    ```go run .\rpc_client_concurrent.go```


### Docker Usage
- Build the base Docker image (from repository root)
    ```
    docker build -t go-rpc-base -f Dockerfile.base .
    ```

## Monolithic Execution
- If running monolithically, build the monolith Docker
    ```
    docker build -t rpc-monolith -f Dockerfile.monolith .
    ```

- After building, run the server image in detached mode
    ```
    docker run -d --name my-rpc-app -p 5000:5000 rpc-monolith
    ```

- After the server is running, execute any of the client commands
    ```
    docker exec my-rpc-app ./client
    ```
    ```
    docker exec my-rpc-app ./client_concurrent
    ```


## Microsservice-oriented Execution
- If running in a microsservice-oriented way, build both the server and client images...
    ```
    docker build -t rpc-server -f Dockerfile.server .
    docker build -t rpc-client -f Dockerfile.client .
    ```

- Afterwards, create the docker network...
    ```
    docker network create --subnet=192.168.27.0/24 prog-dist
    ```

- Then, run the server container
    ```
    docker run -d --name meu-servidor --net prog-dist --ip 192.168.27.2 -p 5000:5000 rpc-server
    ```

- With the server running, run the client as well, using the container name
    ```
    docker run --rm --net prog-dist rpc-client meu-servidor:5000
    ```