package assethandlers

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

type Prefab struct {
	ObjectFile string
	StyleFile  string
	Name       string
	Style      string
	Object     string
}

func GeneratePrefabs() []Prefab {

  prefabs := []Prefab{}
	log.Printf("Generating prefabs")

	baseDirectory, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting base directory ")
	}

	prefabBaseDir := path.Join(baseDirectory, "prefabs/")
	prefabDirectories, err := os.ReadDir(prefabBaseDir)

	if err != nil {
		log.Fatalf("Could not read prefab directory: %+v", err)
	}

	for _, v := range prefabDirectories {

		if !v.Type().IsDir() {
			continue
		}

		prefabName := v.Name()
		log.Printf("name: %+v", prefabName)
		objectFilePath := path.Join(prefabBaseDir, prefabName+"/object.html")
		styleFilePath := path.Join(prefabBaseDir, prefabName+"/style.css")

		log.Printf("Loading %+v object at %+v", prefabName, objectFilePath)
		log.Printf("Loading %+v style at %+v", prefabName, styleFilePath)

    objectContent, err := os.ReadFile(objectFilePath)

    if err != nil {
     log.Fatalf("Error loading %+v object",prefabName) 
    }

    styleContent, err := os.ReadFile(styleFilePath)

    if err != nil {
     log.Fatalf("Error loading %+v style",prefabName) 
    }
    prefabs = append(prefabs,Prefab{
      Name: prefabName,
      ObjectFile: objectFilePath,
      StyleFile: styleFilePath,
      Style: string(styleContent),
      Object: string(objectContent),
    })
	}
  return prefabs
}

func TransformPrefabs(prefabs *[]Prefab){
  for i := 0; i < len(*prefabs); i++ {
    prefab := (*prefabs)[i]

    prefab.Style = strings.ReplaceAll(prefab.Style,"\n.","\n .")
    prefab.Style = strings.ReplaceAll(prefab.Style,"\n[","\n [")
    prefab.Style = strings.ReplaceAll(prefab.Style,"\n#","\n #")

    prefab.Style = fmt.Sprintf(".%s-prefab {\n %s \n}",prefab.Name,prefab.Style)

    prefab.Object = fmt.Sprintf("<div class='%s-prefab'>%s</div>",prefab.Name, prefab.Object)

    (*prefabs)[i] = prefab
  }
}

func CreateTestPrefabFiles(prefabs []Prefab){
  for i := 0; i < len(prefabs); i++ {
    htmlFile := fmt.Sprintf(`
    <html>
      <body>
        %s
        <style>
        %s
        </style>
      </body>
    </html>`,prefabs[i].Object,prefabs[i].Style)
    os.WriteFile("/tmp/"+prefabs[i].Name+".html",([]byte)(htmlFile),0664) 
  }

}
