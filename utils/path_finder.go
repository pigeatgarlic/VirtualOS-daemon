package utils

import (
	"os"
	"os/exec"
	"strings"
)

func FindProcessPath(dir *string,process string) (string,error){
	cmd := exec.Command("where.exe",process)

	if dir != nil {
		cmd.Dir = *dir
	}

	bytes,err := cmd.Output()
	if err != nil{
		return "",nil
	}
	paths := strings.Split(string(bytes), "\n")
	return paths[0],nil
}

func CurrentDir(exeName string) string {
	exe := os.Args[0]
	splits := strings.Split(exe, "\\")
	path := splits[:len(splits)-1]
	path = append(path, exeName)
	return strings.Join(path, "\\")
}