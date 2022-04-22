
# Multiple API feature
Since swag 1.7.9 we are allowing registration of multiple endpoints into the same server.

Generate documentation for v1 endpoints
```shell
swag i -g main.go -dir api/v1 --instanceName v1
```


Generate documentation for v2 endpoints
```shell
swag i -g main.go -dir api/v2 --instanceName v2
```

Run example
```shell
    go run main.go
```

Now you can access the v1 swagger here [http://localhost:8080/swagger/v1/index.html](http://localhost:8080/swagger/v1/index.html) , 
and v2 swagger here [http://localhost:8080/swagger/v2/index.html](http://localhost:8080/swagger/v2/index.html)

