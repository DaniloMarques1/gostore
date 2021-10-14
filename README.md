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

## TODO

* [ ] Create an example of usage
* [ ] Create a "lib" that interact with the key-value server
