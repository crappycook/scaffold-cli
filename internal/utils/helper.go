package utils

import (
	"fmt"
	"os"
)

func GetProjectName(dir string) string {
	modFile, err := os.Open(dir + "/go.mod")
	if err != nil {
		fmt.Println("go.mod does not exist", err)
		return ""
	}
	defer modFile.Close()

	var moduleName string
	_, err = fmt.Fscanf(modFile, "module %s", &moduleName)
	if err != nil {
		fmt.Println("read go mod error: ", err)
		return ""
	}
	return moduleName
}
