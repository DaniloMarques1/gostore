# Go Key Value Store

It's a slow key-value store. We have a server that receives messages specifying a
key and its value. We store the value as we receive the way we received it.

Syntax the server is expecting is

```
op=delete;key=name;
```

Our key-value store supports three operations (op)

```
store
read
delete
```
