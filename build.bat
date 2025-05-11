cd /d %~dp0
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build -o bootstrap jobsheet-go-aws2
powershell compress-archive -Force -Path bootstrap,template,font -DestinationPath lambda-handler.zip
aws lambda update-function-code --function-name jobsheet-function --zip-file fileb://lambda-handler.zip