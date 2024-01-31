package judge

// Option 各个语言的配置
type Option struct {
	fileName  string
	ImageName string
	buildCmd  string
	RunCmd    string
}

var langMap = map[string]Option{
	"go": {
		fileName:  "main.go",
		ImageName: "endmax/go:latest",
		buildCmd:  "go mod init endmax && go build main.go",
		RunCmd:    "{ time -f \"-%K:%e~\" ./main ; } 2>&1",
		//{ time -f "-%K:%e~" ./main ; } 2>&1
		//time -f "-%K:%e~" ./main
	},
	"python": {
		fileName:  "main.py",
		ImageName: "endmax/python:latest",
		buildCmd:  "",
		RunCmd:    "{ time -f \"-%K:%e~\" python3 main.py; } 2>&1",
	},
	"c": {
		fileName:  "main.c",
		ImageName: "endmax/c:latest",
		buildCmd:  "gcc -v main.c -o main",
		RunCmd:    "{ time -f \"-%K:%e~\" ./main;} 2>&1",
	},
}
