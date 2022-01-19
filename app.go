package cmd

import (
	"encoding/base64"
	"fmt"
	"os"
)

func NewApp() {
	modulePath := "my-app"
	fmt.Print("Module Path : ")
	fmt.Scan(&modulePath)

	dbDriver := "postgres"
	fmt.Print("Database Driver : ")
	fmt.Scan(&dbDriver)

	useRedis := "true"
	fmt.Print("Use Redis Cache : ")
	fmt.Scan(&useRedis)

	useOauth2 := "true"
	fmt.Print("Use Default OAuth2 Server : ")
	fmt.Scan(&useOauth2)

	fmt.Println("modulePath :", modulePath)
	fmt.Println("dbDriver :", dbDriver)
	fmt.Println("useRedis :", useRedis)
	fmt.Println("useOauth2 :", useOauth2)
	err := writeSwaggerFiles()
	if err != nil {
		fmt.Println("writting swagger files error : ", err.Error())
	}
}

func writeSwaggerFiles() error {
	fmt.Println("writting swagger files...")
	filePath := "docs"

	// reset
	os.RemoveAll(filePath)
	if err := os.MkdirAll(filePath, 0755); err != nil {
		return err
	}

	for fileName, fileContent := range GetSwaggerUiFiles() {
		content, err := base64.StdEncoding.DecodeString(fileContent)
		if err != nil {
			return err
		}
		fmt.Println("writting file :", filePath+"/"+fileName)
		err = os.WriteFile(filePath+"/"+fileName, content, 0755)
		if err != nil {
			return err
		}
	}
	fmt.Println("swagger files has been written")
	return nil
}
