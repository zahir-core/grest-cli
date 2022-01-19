package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const default_swagger_json_url string = `url: "https://petstore.swagger.io/v2/swagger.json",`
const additional_swagger_ui_config string = `url: window.location.href.split("#")[0] + "openapi.json",
        defaultExpanded: false,
        syntaxHighlight: {
          activated: true,
          theme: "tomorrow-night"
        },`

var filename string = "./files.go"
var files_go string = `package cmd

func GetSwaggerUiFiles() map[string]string {
	return map[string]string{[files]
	}
}
`

func main() {
	resetSwaggerGoFile()
	downloadLastVersionSwaggerUI()
	updateSwaggerGoFile()
}

func resetSwaggerGoFile() {
	_, err := os.Stat(filename)
	if err == nil {
		os.Remove(filename)
	}
}

func downloadLastVersionSwaggerUI() {
	// Downloading last version of swagger-ui
	fmt.Println("Downloading last version of swagger-ui : swagger-ui-master.tar.gz")
	resp, err := http.Get("https://codeload.github.com/swagger-api/swagger-ui/tar.gz/master")
	if err != nil {
		log.Fatalf("downloadLastVersionSwaggerUI: http.Get(url) failed: %s", err.Error())
	}
	defer resp.Body.Close()

	// Extracting tar.gz files
	fileReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		log.Fatalf("downloadLastVersionSwaggerUI: gzip.NewReader(resp.Body) failed: %s", err.Error())
	}
	defer fileReader.Close()

	// Extracting tarred files for dist folder only
	tarBallReader := tar.NewReader(fileReader)
	for {
		header, err := tarBallReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("downloadLastVersionSwaggerUI: tarBallReader.Next() failed: %s", err.Error())
		}

		// get the individual filename and extract to the current directory
		filename := header.Name
		if strings.HasPrefix(filename, "swagger-ui-master/dist") {
			filename = strings.Replace(filename, "swagger-ui-master/dist", "swagger-ui", 1)
			switch header.Typeflag {
			case tar.TypeDir:
				// handle directory
				fmt.Println("Creating directory :", filename)
				err = os.MkdirAll(filename, os.FileMode(header.Mode)) // or use 0755 if you prefer
				if err != nil {
					log.Fatalf("downloadLastVersionSwaggerUI: os.MkdirAll(filename, mode) failed: %s", err.Error())
				}
			case tar.TypeReg:
				// handle normal file
				fmt.Println("Untarring :", filename)
				if filename == "swagger-ui/index.html" {
					index_byte, err := io.ReadAll(tarBallReader)
					if err != nil {
						log.Fatalf("downloadLastVersionSwaggerUI: io.ReadAll(tarBallReader) failed: %s", err.Error())
					}
					index_string := strings.Replace(string(index_byte), default_swagger_json_url, additional_swagger_ui_config, -1)
					err = os.WriteFile(filename, []byte(index_string), 0755)
					if err != nil {
						log.Fatalf("downloadLastVersionSwaggerUI: os.WriteFile(filename, []byte(index_string), mode) failed: %s", err.Error())
					}
				} else {
					writer, err := os.Create(filename)
					if err != nil {
						log.Fatalf("downloadLastVersionSwaggerUI: os.Create(filename) failed: %s", err.Error())
					}
					io.Copy(writer, tarBallReader)
					err = os.Chmod(filename, os.FileMode(header.Mode))
					if err != nil {
						log.Fatalf("downloadLastVersionSwaggerUI: os.Chmod(filename, mode) failed: %s", err.Error())
					}
					writer.Close()
				}
			default:
				fmt.Printf("Unable to untar type : %c in file %s", header.Typeflag, filename)
			}
		}
	}
}

func updateSwaggerGoFile() {
	fmt.Println("Updating files.go")
	files, err := os.ReadDir("./swagger-ui")
	if err != nil {
		fmt.Println(err)
	}
	text := ""
	for _, f := range files {
		content, err := os.ReadFile("./swagger-ui/" + f.Name())
		if err != nil {
			fmt.Println(err)
		}
		text = text + "\n\t\t\"" + f.Name() + "\": \"" + base64.StdEncoding.EncodeToString(content) + "\","
	}
	fileText := strings.Replace(files_go, "[files]", text, -1)
	err = os.WriteFile(filename, []byte(fileText), 0755)
	if err != nil {
		fmt.Println(err)
	}
}
