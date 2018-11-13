# gcurl
a command line tool for gRPC, like curl for HTTP  
# 解决什么问题    
- 基于json格式的http请求，我们可以方便的使用curl去模拟客户端发请求调试， 但是gRPC却不能这样操作。
- gcurl支持输入json格式的请求，再转发到gRPC服务器， 就像使用curl调试http一样去调试gRPC接口

> PS  
> 1、新功能持续开发添加中，当前只是一个简单可用版本  
> 2、当前版本需要gRPC server开启 reflect, 后续会支持本地proto模式

# Install  
```shell
	go get -u github.com/lanfang/gcurl

```

# Usage
## 查看方法列表，消息定义等
```
gcurl desc localhost:port serevice(fullyQualifiedName package.service) symbol1 symbol2 
```

## 发起RPC请求
```
gcurl localhost:port method-symbol -d '{"name":"fang.l"}' 
```

## 详细帮助
```
 gcurl --help
a command line tool for gRPC, like curl for HTTP.
you can interact with the rpc server like this:
gcurl host:port method -d '{"username":"gcurl", "password":"gcurl"}' or exec the subcommand

Usage:
  gcurl [flags]
  gcurl [command]

Available Commands:
  desc        desc symbol, show the detail info of the symbol
  help        Help about any command

Flags:
  -d, --data string    request data
  -h, --help           help for gcurl
  -b, --proto string   local proto file
      --version        version for gcurl
```
# TODO
- 支持Stream模式  
- 支持extensions, Any


# Thanks 
- [protoreflect](https://github.com/jhump/protoreflect)