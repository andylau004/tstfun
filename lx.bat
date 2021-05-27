
::compile_gotest.bat
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64

go build
rm -rf .\1.zip

"C:\Program Files (x86)\VRV\CEMS\VRVNAC\7z.exe" a 1.zip tstfun



