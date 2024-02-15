package swagger

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/spf13/viper"
)

type SwaggerInfo struct {
	Swagger             string                 `json:"swagger"`
	Info                Info                   `json:"info"`
	Host                string                 `json:"host,omitempty"`
	BasePath            string                 `json:"basePath,omitempty"`
	Schemes             []string               `json:"schemes,omitempty"`
	Consumes            []string               `json:"consumes,omitempty"`
	Produces            []string               `json:"produces,omitempty"`
	Paths               map[string]interface{} `json:"paths,omitempty"`
	Definitions         map[string]interface{} `json:"definitions,omitempty"`
	Parameters          map[string]interface{} `json:"parameters,omitempty"`
	Responses           map[string]interface{} `json:"responses,omitempty"`
	SecurityDefinitions map[string]interface{} `json:"securityDefinitions,omitempty"`
	Security            []map[string][]string  `json:"security,omitempty"`
	Tags                []Tag                  `json:"tags,omitempty"`
	ExternalDocs        ExternalDocs           `json:"externalDocs,omitempty"`
}

type Info struct {
	Title          string  `json:"title"`
	Description    string  `json:"description,omitempty"`
	TermsOfService string  `json:"termsOfService,omitempty"`
	Contact        Contact `json:"contact,omitempty"`
	License        License `json:"license,omitempty"`
	Version        string  `json:"version"`
}

type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

type License struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

type Tag struct {
	Name         string       `json:"name"`
	Description  string       `json:"description,omitempty"`
	ExternalDocs ExternalDocs `json:"externalDocs,omitempty"`
}

type ExternalDocs struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
}

func ReadSwagger(url, prefixPath string) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var swaggerInfo SwaggerInfo
	err = json.Unmarshal(body, &swaggerInfo)
	if err != nil {
		return
	}

	swaggerInfo.Host = "localhost:3000"
	swaggerInfo.BasePath = path.Join(prefixPath, swaggerInfo.BasePath)

	bts, err := json.Marshal(&swaggerInfo)

	if err != nil {
		return
	}

	directory := viper.GetString("app.docs_directory")

	os.MkdirAll(directory, os.ModePerm)

	os.Remove(directory + prefixPath + ".json")
	writeFile(directory+prefixPath+".json", bts)
}

func writeFile(filePath string, b []byte) error {
	fileToWrite, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}

	defer fileToWrite.Close()
	_, err = fileToWrite.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	_, err = fileToWrite.Write(b)
	if err != nil {
		return err
	}
	return fileToWrite.Sync()
}
