# Overview
- > [GUIDES](https://grpc.io/docs/guides/)
- With gRPC we can define our service once in a `.proto` file and implement clients and servers in any of gRPC’s supported languages, which in turn can be run in environments ranging from servers inside Google to your own tablet - all the complexity of communication between different languages and environments is handled for you by gRPC
    - We also get all the advantages of working with protocol buffers, including efficient serialization, a simple IDL, and easy interface updating
- gRPC (Remote Procedure Call) can use protocol buffers as both its Interface Definition Language (IDL) and as its underlying message interchange format
- In gRPC a client application can directly call **methods** on a server application on a different machine as if it was a local object
- gRPC is based around the idea of defining a service, specifying the methods that can be called remotely with their parameters and return types
    - On the server side, the **server implements this interface** and runs a gRPC server to handle client calls
    - On the client side, the client has a *stub* (referred to as just a client in some languages) that **provides** the same methods as the server
    ![](https://grpc.io/img/landing-2.svg)
- Define gRPC services in ordinary proto files, with RPC method parameters and return types specified as protocol buffer messages
	
    ```proto
    // The greeter service definition.
    service Greeter {
        // Sends a greeting
        rpc SayHello (HelloRequest) returns (HelloReply) {}
    }
    // The request message containing the user's name.
    message HelloRequest {
        string name = 1;
    }
    // The response message containing the greetings
    message HelloReply {
        string message = 1;
    }
    ```

- gRPC also uses `protoc` with a special gRPC plugin to generate code from your proto file. However, with the gRPC plugin, you get generated gRPC client and server code, as well as the regular protocol buffer code for populating, serializing, and retrieving your message types
	
    ```bash
    protoc --go_out=plugins=grpc:. *.proto
    ```

    - https://github.com/golang/protobuf#grpc-support
# Concepts
## Service definition
- RPC lets you define 4 kinds of service method
    1. Unary RPCs where the client sends a single request to the server and gets a single response back, just like a normal function call
    	
        ```go
        rpc SayHello(HelloRequest) returns (HelloResponse){
        }
        ```
    
    2. Server streaming RPCs where the client sends a request to the server and gets a stream to read a sequence of messages back. The client reads from the returned stream until there are no more messages. gRPC guarantees **message ordering** within an individual RPC call
    	
        ```go
        rpc LotsOfReplies(HelloRequest) returns (stream HelloResponse){
        }
        ```
    
    3. Client streaming RPCs where the client writes a sequence of messages and sends them to the server, again using a provided stream. Once the client has finished writing the messages, it waits for the server to read them and return its response. Again gRPC guarantees message ordering within an individual RPC call
    	
        ```go
        rpc LotsOfGreetings(stream HelloRequest) returns (HelloResponse) {
        }
        ```
    
    4. Bidirectional streaming RPCs where both sides send a sequence of messages using a read-write stream. The two streams operate independently, so clients and servers can read and write in whatever order they like: for example, the server could wait to receive all the client messages before writing its responses, or it could alternately read a message then write a message, or some other combination of reads and writes. The order of messages in each stream is preserved
    	
        ```go
        rpc BidiHello(stream HelloRequest) returns (stream HelloResponse){
        }
        ```

## Using the API surface
- Starting from a service definition in a `.proto` file, gRPC provides protocol buffer compiler plugins that generate client- and server-side code
- gRPC users typically call these APIs on the client side and implement the corresponding API on the server side
    - On the server side, the server implements the methods declared by the service and runs a gRPC server to handle client calls. The gRPC infrastructure decodes incoming requests, executes service methods, and encodes service responses
    - On the client side, the client has a local object known as stub (for some languages, the preferred term is client) that implements the same methods as the service. The client can then just call those methods on the local object, wrapping the parameters for the call in the appropriate protocol buffer message type - gRPC looks after sending the request(s) to the server and returning the server’s protocol buffer response(s)
## Synchronous vs. asynchronous
- Synchronous RPC calls that block until a response arrives from the server are the closest approximation to the abstraction of a procedure call that RPC aspires to. On the other hand, networks are inherently asynchronous and in many scenarios it’s useful to be able to start RPCs without blocking the current thread
- The gRPC programming surface in most languages comes in both synchronous and asynchronous flavors
## RPC life cycle
- stub 发送 metadata 和发送 request 的内容是分两步进行的
### Unary RPC
- Once the client calls the method on the stub/client object, the server is notified that the RPC has been invoked with the client’s metadata for this call, the method name, and the specified deadline if applicable (message not arrived yet)
- The server can then either send back its own initial metadata (which must be sent before any response) straight away, or wait for the client’s request message - which happens first is application-specific
- Once the server has the client’s request message, it does whatever work is necessary to create and populate its response. The response is then returned (if successful) to the client together with status details (status code and optional status message) and optional trailing metadata
- If the status is OK, the client then gets the response, which completes the call on the client side
### Server streaming RPC
- A server-streaming RPC is similar to our simple example, except the server sends back a stream of responses after getting the client’s request message. After sending back all its responses, the server’s status details (status code and optional status message) and optional trailing metadata are sent back to complete on the server side. The client completes once it has all the server’s responses
### Client streaming RPC
- A client-streaming RPC is also similar to our simple example, except the client sends a stream of requests to the server instead of a single request. The server sends back a single response, typically but not necessarily after it has received all the client’s requests, along with its status details and optional trailing metadata
### Bidirectional streaming RPC
- In a bidirectional streaming RPC, again the call is initiated by the client calling the method and the server receiving the client metadata, method name, and deadline. Again the server can choose to send back its initial metadata or **wait for the client to start sending requests**
- What happens next depends on the application, as the client and server can read and write in any order - the streams operate completely independently
    - For example, the server could wait until it has received all the client’s messages before writing its responses, or the server and client could “ping-pong”: the server gets a request, then sends back a response, then the client sends another request based on the response, and so on
### Deadlines/Timeouts
- gRPC allows clients to specify how long they are willing to wait for an RPC to complete before the RPC is terminated with the error `DEADLINE_EXCEEDED`. On the server side, the server can query to see if a particular RPC has timed out, or how much time is left to complete the RPC
    - How the deadline or timeout is specified varies from language to language - for example, not all languages have a default deadline, some language APIs work in terms of a deadline (a fixed point in time), and some language APIs work in terms of timeouts (durations of time)
### RPC termination
- In gRPC, both the client and server make **independent and local determinations** of the success of the call, and their conclusions may not match
    - This means that, for example, you could have an RPC that finishes successfully on the server side (“I have sent all my responses!”) but fails on the client side (“The responses arrived after my deadline!“)
    - It’s also possible for a server to decide to complete before a client has sent all its requests
### Cancelling RPCs
- Either the client or the server can cancel an RPC at any time
- A cancellation terminates the RPC immediately so that no further work is done. It is **not an “undo”**: changes made before the cancellation will not be rolled back
### Metadata
- Metadata is information about a particular RPC call (such as authentication details) in the form of **a list of key-value pairs**, where the keys are strings and the values are typically strings (but can be binary data). Metadata is opaque to gRPC itself - it lets the client provide information associated with the call to the server and vice versa
- Access to metadata is language-dependent
### Channels
- A gRPC channel provides a connection to a gRPC server on a specified host and port and is used when creating a client stub (or just “client” in some languages)
- Clients can specify channel arguments to modify gRPC’s default behaviour, such as switching on and off message compression
- A channel has state, including `connected` and `idle`
- How gRPC deals with closing down channels is language-dependent
- Some languages also permit querying channel state
# Tutorial
- https://grpc.io/docs/tutorials/basic/go/
- `route_guide.pb.go` contains 
    - All the protocol buffer code to populate, serialize, and retrieve our request and response message types
    - An interface type (or stub) for clients to **call** with the methods defined in the `RouteGuide` service
    - An interface type for servers to **implement**, also with the methods defined in the `RouteGuide` service
- There are two parts to making our RouteGuide service do its job:
    - Implementing the service interface generated from our service definition: doing the actual “work” of our service
    - Running a gRPC server to listen for requests from clients and dispatch them to the right service implementation
- Return the error “as is” so that it’ll be translated to an RPC status by the gRPC layer
- stub 和 server 的同名函数的签可能不一样（rpc 层在中间会做处理，如 context）
- Note that in gRPC-Go, RPCs operate in a **blocking/synchronous** mode, which means that the RPC call waits for the server to respond, and will either return a response or an error
- We also pass a `context.Context` object which lets us change our RPC’s behaviour if necessary, such as time-out/cancel an RPC in flight