## description

cmd_server executes requested command

Usage :  curl -i -X POST -d 'path=./tenant1&cmd=''["date"]''' 192.168.40.158:5000/generate
Parameters: 
	path - directory to create. it this case ./tenant1, 
	date - command to execute in json array format

## to build 

go build cmd_server.go


## to run

Windows
	cmd_server.exe 
Linux
	./cmd_server


## settings
you can overwrite default port 4000
in file `settings.txt` write value, e.g. `5000`