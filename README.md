# Go simple upload server  
## Start Server
```
./upload-server -secret your_secret_key -port 8080
```
## Curl  
```
curl -X POST -F "file=@file.txt" -F "secret=my_secret_key" http://localhost:8080/upload
```  