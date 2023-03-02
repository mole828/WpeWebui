package wpe

import (
	"encoding/json"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Project struct {
	Id string `json:"id"`
	//Contentrating string
	//Description   string
	File string `json:"file"`
	//General       interface{}
	Preview string `json:"preview"`
	//Tags    []string
	Title string `json:"title"`
	Type  string `json:"type"`
	//Version int
}

func LoadJson(file *os.File) (Project, error) {
	info, err := file.Stat()
	if info.IsDir() || err != nil {
		return Project{}, err
	}
	var data Project
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		return Project{}, err
	}
	return data, nil
}

func IndexProjectDir(root string) Project {
	var data Project
	if err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		//log.Printf("FullPath: %s, %d, %t", FullPath, info.Size(), info.Name() == "project.json")

		if info.Name() == "project.json" {

			p, err := os.Open(path)
			if err != nil {
				return err
			}
			loadJson, err := LoadJson(p)
			if err != nil {
				return err
			}
			data = loadJson
			id := filepath.Base(filepath.Dir(path))
			//log.Printf("id: %s", id)
			data.Id = id
			if err = p.Close(); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		log.Panicln(err)
	}
	return data
}

func IterWorkshopContent(root string, callback func(project Project)) error {
	c := strings.Count(root, "\\")
	return filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if strings.Count(path, "\\") != c+1 {
			return nil
		}
		project := IndexProjectDir(path)
		if project.Type == "video" {
			//log.Printf("%s: %s", FullPath, project.Title)
			callback(project)
		}
		return nil
	})
}

type ProjectIndex struct {
	root       string
	projects   []Project
	projectMap map[string]Project
}

func New(root string) *ProjectIndex {
	index := ProjectIndex{}
	index.root = root
	index.projects = []Project{}
	index.projectMap = make(map[string]Project)
	go func() {
		err := IterWorkshopContent(root, func(project Project) {
			index.projects = append(index.projects, project)
			index.projectMap[project.Id] = project
			// log.Printf("%d %s", len(index.projects), index.FullPath(project.Id))
		})
		if err != nil {
			panic(err)
		}
	}()
	return &index
}

func (it ProjectIndex) List() []Project {

	return it.projects
}

func (it ProjectIndex) Map() map[string]Project {
	return it.projectMap
}

func (it ProjectIndex) Find(id string) Project {
	return it.projectMap[id]
}

func (it ProjectIndex) FullPath(id string) string {
	project := it.Find(id)
	fullPath := it.root + "/" + project.Id + "/" + project.File
	return fullPath
}
