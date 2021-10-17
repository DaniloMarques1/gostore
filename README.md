# Go Key Value Store

It's a slow key-value store. We have a server that receives messages specifying a
key and its value. We store the value as we receive the way we received it.

Syntax the server is expecting is

```
op=store;key=name;value=Danilo;
op=read;key=name;
op=delete;key=name;
op=list;
```

Our key-value store supports four operations (op)

```
store
read
delete
list
```

## How to use with telnet

Build and run the server by executing

```
go build . && ./gostore
```

After that, start a tcp connection

```
telnet localhost 5000
```

Then all you need is to start send messages with the operations

```
op=store;key=user;value={"user": "Danilo", "age": 22};
code=0;message=Value stored successfully

op=read;key=user;
code=0;message={"user": "Danilo", "age": 22}

op=list;
code=0;message=[{"user": "Danilo", "age": 22}]

op=keys;
code=0;message=[user]

op=delete;key=user;
code=0;message=Value removed successfully
```
