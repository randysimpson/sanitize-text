# Sanitize-Text

This application was created to sanitize text from a log file, by removing specific dates and other numbers that will affect a NLP application.  Currently the only supported sanitize method is for golang log files using the format that kubernetes components use, i.e. kubelet, kube scheduler, etc.

## Run

To run the container issue:

```
docker run -d -p 8081:8081 randysimpson/sanitize-text:latest
```

## Use

This program has various endpoints to utilize the sanitize text functions:

* `/upload/log` - used to upload a file to the service
* `/load/log` - used to send the data as text inside a body of a request
* `/load/text` - used to send text inside a body of a request

Currently the same sanitize function is used for all methods of sanitizing text.

## Upload

### k8's log file

```sh
curl -i -H "Accept: application/json" -F "file=@test.log" -X POST http://localhost:8081/api/v1/upload/log
```

#### testing

```
ubuntu@master-1:~/code/go/src/sanitize-text$ curl -H "Accept: application/json" -F "file=@test.log" -X POST http://localhost:8081/api/v1/upload/log > output.curl
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100 12617    0  5740  100  6877  1121k  1343k --:--:-- --:--:-- --:--:-- 2464k
ubuntu@master-1:~/code/go/src/sanitize-text$ less output.curl
ubuntu@master-1:~/code/go/src/sanitize-text$ curl -H "Accept: application/json" -F "file=@test2.log" -X POST http://localhost:8081/api/v1/upload/log > output2.curl
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100 40135    0 14257  100 25878   605k  1098k --:--:-- --:--:-- --:--:-- 1781k
ubuntu@master-1:~/code/go/src/sanitize-text$ curl -i -H "Accept: application/json" -H "Content-Type: application/json" -d @output.curl -X POST http://localhost:8080/api/v1/build
HTTP/1.1 100 Continue

HTTP/1.1 201 Created
Content-Type: application/json; charset=UTF-8
Date: Fri, 27 Dec 2019 18:55:57 GMT
Content-Length: 34

{"size":"541","status":"Success"}
ubuntu@master-1:~/code/go/src/sanitize-text$ curl -i -H "Accept: application/json" -H "Content-Type: application/json" -X GET http://localhost:8080/api/v1/predict/%3Cerror%3E
HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Date: Fri, 27 Dec 2019 18:56:07 GMT
Content-Length: 13

{"<date>":1}
ubuntu@master-1:~/code/go/src/sanitize-text$ curl -i -H "Accept: application/json" -H "Content-Type: application/json" -d '{"data":"Somthing is very wrong here"}' -X POST http://localhost:8080/api/v1/entropy
HTTP/1.1 200 OK
Content-Type: application/json; charset=UTF-8
Date: Fri, 27 Dec 2019 18:56:29 GMT
Content-Length: 27

{"entropy":"13831.766209"}

cat output2.curl | jq -r '.data' | sed -e 's/<\/line> /<\/line>\n/g' | while read LINE ; do echo '{"data":"'$LINE'"}';  done | while read LINE ; do echo 'curl -H "Accept: application/json" -H "Content-Type: application/json" -d '"'"''$LINE''"'"' -X POST http://localhost:8080/api/v1/entropy';  done | sh

cat big-output2.json | jq -r '.data' | sed -e 's/<\/line> /<\/line>\n/g' | while read LINE ; do echo '{"data":"'$LINE'"}'; done | while read LINE ; do echo ''"'$LINE'"''; done | xargs -n1 curl -H "Accept: application/json" -H "Content-Type: application/json" -X POST http://localhost:8080/api/v1/entropy -d
```

## Installation

This is a microservice which has been written based on REST-API to allow for deployment from a docker container.

### Manual Build

To manually build the source files you will need to get the external dependencies and then build the binary executable file.

```
go get k8s.io/klog
go get github.com/gorilla/mux
go build
```

# Licence

MIT License

Copyright (©) 2019 - Randall Simpson

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.