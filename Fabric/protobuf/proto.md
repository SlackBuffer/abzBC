- [Language Guide (proto3)](https://developers.google.com/protocol-buffers/docs/proto3)
- With protocol buffers, you write a `.proto` description of the data structure you wish to store. From that, the **protocol buffer compiler creates a class** that implements automatic encoding and parsing of the protocol buffer data with an efficient binary format
    - The generated class provides getters and setters for the fields that make up a protocol buffer and takes care of the details of reading and writing the protocol buffer as a unit
    - Importantly, the protocol buffer format supports the idea of extending the format over time in such a way that the code can still read data encoded with the old format
- https://github.com/protocolbuffers/protobuf/tree/master/examples
## [Download protocol buffers](https://developers.google.com/protocol-buffers/docs/downloads)
- [Download the compiler](https://github.com/protocolbuffers/protobuf/releases)
- Install the Go protocol buffers `go get -u github.com/golang/protobuf/protoc-gen-go
`
- >https://github.com/protocolbuffers/protobuf/issues/5131
## Define protocol format
- Start with a `.proto` file
- Definitions: Add a message for each data structure that need to be serialize, then specify a name and a type for each field in the message
- The `.proto` file starts with a package declaration, which helps to prevent naming conflicts between different projects
    - In Go, the package name is used as the Go package, unless you have specified a `go_package`. Even if you do provide a `go_package`, you should still define a normal package as well to avoid name collisions in the Protocol Buffers name space as well as in non-Go languages
- Many standard simple data types are available as field types, including `bool`, `int32`, `float`, `double`, and `string`
    - You can also add further structure to your messages by using other message types as field types
- You can even define message types nested inside other messages
- You can also define enum types if you want one of your fields to have one of a predefined list of values
- The " = 1", " = 2" markers on each element identify the unique "tag" that field uses in the binary encoding
    - Tag numbers 1-15 require one less byte to encode than higher numbers, so as an optimization you can decide to use those tags for the commonly used or repeated elements, leaving tags 16 and higher for less-commonly used optional elements
    - Each element in a repeated field requires re-encoding the tag number, so repeated fields are particularly good candidates for this optimization
- If a field value isn't set, a default value is used: zero for numeric types, the empty string for strings, false for bools
- For embedded messages, the default value is always the "default instance" or "prototype" of the message, which has none of its fields set
    - Calling the accessor to get the value of a field which has not been explicitly set always returns that field's default value
- If a field is `repeated`, the field may be repeated any number of times (including zero). The order of the repeated values will be preserved in the protocol buffer
    - Think of repeated fields as dynamically sized arrays
## Compile protocol buffers
- Now run the compiler, specifying the source directory (where your application's source code lives – the current directory is used if you don't provide a value), the destination directory (where you want the generated code to go; often the same as `$SRC_DIR`), and the path to your `.proto`
	
    ```bash
    # protoc -I=$SRC_DIR --go_out=$DST_DIR $SRC_DIR/addressbook.proto
    protoc --go_out=. address_book.proto
    ```

## The protocol buffer application
- Generate `addressbook.pb.go` gives you the following useful types
    - An `AddressBook` structure with a `People` field
    - A `Person` structure with fields for `Name`, `Id`, `Email` and `Phones`
    - A `Person_PhoneNumber` structure, with fields for `Number` and `Type`
    - The type `Person_PhoneType` and a value defined for each value in the `Person.PhoneType` enum
    - > https://developers.google.com/protocol-buffers/docs/reference/go-generated
- Create an instance of `Person`
	
    ```go
    p := pb.Person {
        Id:    1234,
        Name:  "John Doe",
        Email: "jdoe@example.com",
        Phones: []*pb.Person_PhoneNumber {
                {Number: "555-4321", Type: pb.Person_HOME},
        },
    }
    ```

    - https://github.com/protocolbuffers/protobuf/blob/master/examples/list_people_test.go
## Writing a message
- The whole purpose of using protocol buffers is to serialize your data so that it can be parsed elsewhere
- In Go, you use the `proto` library's Marshal function to serialize your protocol buffer data
    - A pointer to a protocol buffer message's struct implements the `proto.Message` interface
    - Calling `proto.Marshal` returns the protocol buffer, encoded in its wire format
## Reading a message
- To parse an encoded message, you use the proto library's `Unmarshal` function
    - Calling this parses the data in `buf` as a protocol buffer and places the result in `pb`
## Extending a protocol buffer
- Sooner or later after you release the code that uses your protocol buffer, you will undoubtedly want to "improve" the protocol buffer's definition
- If you want your new buffers to be backwards-compatible, and your old buffers to be forward-compatible – and you almost certainly do want this – then there are some rules you need to follow
- In the new version of the protocol buffer:
    - you must not change the tag numbers of any existing fields
    - you may delete fields
    - you may add new fields but you must use fresh tag numbers (i.e. tag numbers that were never used in this protocol buffer, not even by deleted fields)
    - > [Some exceptions](https://developers.google.com/protocol-buffers/docs/proto3#updating)
- If you follow these rules, old code will happily read new messages and simply ignore any new fields. To the old code, singular fields that were deleted will simply have their default value, and deleted repeated fields will be empty. New code will also transparently read old messages
- New fields will not be present in old messages, so you will need to do something reasonable with the default value