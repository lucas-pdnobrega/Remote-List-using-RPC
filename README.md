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
  
    ```go run .\remote_list_rpc_server.go```

- On a separate terminal instance, run the client and observe the automatic interactions via RPC calls
  
    ```go run .\remote_list_rpc_client.go```
