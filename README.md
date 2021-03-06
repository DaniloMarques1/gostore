# Go Key Value Store

It's a slow key-value store. We have a server that receives messages specifying a
key and its value. We store the value as we receive the way we received it.

Syntax the server is expecting is

```
op=store;key=name;value=Danilo;
op=read;key=name;
op=delete;key=name;
op=list;
op=keys;
op=replace;key=name;value=Fitz;
```

Our key-value store supports six operations (op)

```
store
read
delete
list
keys
replace
```

## How to use with telnet

Build and run the server by executing

```
go build . 
./gostore # default to port 5000 or
./gostore --port 8080
```

After that, start a tcp connection

```
telnet localhost 5000
```

Then all you need to do is to start send messages with the operations

```
op=store;key=user;value={"user": "Danilo", "age": 22};
RESPONSE -> code=0;message=Value stored successfully

op=read;key=user;
RESPONSE -> code=0;message={"user": "Danilo", "age": 22}

op=list;
RESPONSE -> code=0;message=[{"user": "Danilo", "age": 22}]

op=keys;
RESPONSE -> code=0;message=[user]

op=delete;key=user;
RESPONSE -> code=0;message=Value removed successfully

op=replaced;key=user;value={"user": "Fitz", "age": 19};
RESPONSE -> code=0;message=Value replaced successfully
```

### Examples

In the examples folder, you can find three "libs" that implement a communication
with the server. They are implemented using go, typescript and python. Along
with the libs, there is also a main file showing how to use the lib.
