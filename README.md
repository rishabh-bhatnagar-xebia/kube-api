### kube-api
---

A WebApi Wrapper that helps you in executing following operations of kubernetes:
  - Create a pod with a given image
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


----
# Api Definitions

## Endpoints

---

### POST `/create`  
Creates a pod with the given image name

#### Parameters  
- `image` (required): Name of the image to be used for creating a pod
- `pod_name` (required): Name of the pod
- `namespace` (optional): Namespace where the pod will be created. Defaults to `default`

#### Response
| status_code |            response            |             description              |
| :---------: | :----------------------------: | :----------------------------------: |
|     500     | `{"Error": "<error-message>"}` | If there was an error creating a pod |
|     204     |               -                |        If the pod was created        |

---

### GET `/list?namespace={}`  
Lists all the pods in the cluster

#### Parameters  
- `namespace` (optional): url path parameter; pods from this namespace will be listed. Defaults to "default"

#### Response
| status_code |          response          |                                                         description                                                         |
| :---------: | :------------------------: | :-------------------------------------------------------------------------------------------------------------------------: |
|     200     | [<pod definition>] \| null | if the namespace is invalid or there are no pods, it returns a null. Else, a list of pod status and meta-params is returned |

---

### GET `/logs`  
Streams the logs of the given pod

#### Parameters  
- `pod_name` (required): Name of the pod
- `namespace` (optional): Namespace where the pod will be created. Defaults to `default`

#### Response
| status_code |            response            |                    description                     |
| :---------: | :----------------------------: | :------------------------------------------------: |
|     500     | `{"Error": "<error-message>"}` | If there was an error fetching the logs of the pod |
|     200     |               -                |         log stream from pod in plain text          |