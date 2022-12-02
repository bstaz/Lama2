// Package `import` helps developers switch
// to Lama2 from various popular tools. The
// conversion may not be perfect, but it should
// help teams get started with minimal manual
// effort
package importer

import (
	"fmt"
	"os"

	"github.com/HexmosTech/gabs/v2"
	"github.com/rs/zerolog/log"
)

// create 3 structures:
// Folder
// Request
// Environ

// folderMap
// requestMap
// environMap

type Folder struct {
	Name  string
	Ident string
}

type Request struct {
	TheURL       string
	Name         string
	Method       string
	Auth         string
	RequestType  string
	RawModeData  string
	MultiData    map[string][]string
	ParentFolder string
	HeaderData   map[string]string
	Ident        string
}

type Environ struct {
	Name   string
	Values map[string]string
	Ident  string
}

var (
	folderMap  map[string]Folder
	requestMap map[string]Request
	environMap map[string]Environ
)

func init() {
	folderMap = make(map[string]Folder)
	requestMap = make(map[string]Request)
	environMap = make(map[string]Environ)
}

func generateFolderMap(foldersList *gabs.Container) {
	fmt.Println("The folder", foldersList)
}

// PostmanConvert takes in a Postman data file
// and generates a roughly equivalent Lama2 repository.
// Collections and subcollections become folders.
// Requests become files, environments get stored in
// `l2.env` while file attachments get copied relative
// to the API file
func PostmanConvert(postmanFile string, targetFolder string) {
	fmt.Println(postmanFile, targetFolder)
	contents, e := os.ReadFile(postmanFile)
	if e != nil {
		log.Fatal().Msg(e.Error())
	}
	pJSON, e2 := gabs.ParseJSON(contents)
	if e2 != nil {
		log.Fatal().Msg(e.Error())
	}
	coll := pJSON.S("collections")
	for _, child := range coll.Children() {
		foldersList := child.S("folders")
		for _, folder := range foldersList.Children() {
			fmt.Println(folder)
			generateFolderMap(folder)
			fmt.Println("===")
		}
	}
}