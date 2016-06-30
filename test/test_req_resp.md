```
47.88.78.210:5000/chaincode
```


```
{
  "jsonrpc": "2.0",
  "method": "deploy",
  "params": {
    "type": 1,
    "chaincodeID":{
        "path":"https://github.com/swbsin/test/test/"
    },
    "ctorMsg": {
        "function":"init",
         "args":["wenbin"]
    },
    "secureContext": "jim"
  },
  "id": 21
}
```

```
47.88.78.210:5000/chaincode
```

```
{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID":{
        "name":"b8a8360a39e37535ae11f1cc5e5faa7cd7e1a56e45775e93de8de6c4dc5890670af294adf4966d86ff2bcac9176d5a7a1bb7e289d4748af9c9242230edb6bc8d"
    },
    "ctorMsg": {
        "function":"query"
    },
    "secureContext": "jim"
  },
  "id": 51
}
```

```
47.88.78.210:5000/chaincode
```

```
{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
    "type": 1,
    "chaincodeID":{
        "name":"b8a8360a39e37535ae11f1cc5e5faa7cd7e1a56e45775e93de8de6c4dc5890670af294adf4966d86ff2bcac9176d5a7a1bb7e289d4748af9c9242230edb6bc8d"
    },
    "ctorMsg": {
        "function":"createAccount",
         "args":["yajun"]
    },
    "secureContext": "jim"
  },
  "id": 51
}
```
