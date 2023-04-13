### kube-api
---

A WebApi Wrapper that helps you in executing following operations of kubernetes:
  - Run a Image
  - Get logs of any running pod
  - List All the available pods

Steps To Run The Webserver:
---
```golang
go run main.go
```
By Default, the application will listen on port *8092*,  
To modify that behavior, you can pass an inline param for port like below:

```go
go run main.go 8093
```
Now, the applciation will listen on port 8093