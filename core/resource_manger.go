package core

import (
	"fmt"
	"log"
)

type LoaderParam interface{}
type ResourceType string

type ResourceData interface{}
type Resource struct {
	Type   ResourceType
	Uri    string
	Data   ResourceData
	Empty  bool
	Unload func()
}

type ResourceLoader interface {
	CanLoad(resourceType ResourceType, uri string, param LoaderParam) bool
	Load(uri string, param LoaderParam) (Resource, error)
}

type EmptyResourceLoader struct{}

func (r EmptyResourceLoader) CanLoad(resourceType ResourceType, uri string, param LoaderParam) bool {
	return false
}

func (r EmptyResourceLoader) Load(uri string, param LoaderParam) (Resource, error) {
	return Resource{
		Empty: true,
	}, nil
}

func (r EmptyResourceLoader) Unload(resource Resource) error {
	return nil
}

type ResourceManager struct {
	resourceLoaders []ResourceLoader
	resources       map[string]Resource
}

func (r *ResourceManager) AddLoader(loader ResourceLoader) *ResourceManager {
	r.resourceLoaders = append(r.resourceLoaders, loader)
	return r
}

func (r *ResourceManager) GetLoader(resourceType ResourceType, uri string, param LoaderParam) (ResourceLoader, error) {
	for _, loader := range r.resourceLoaders {
		if loader.CanLoad(resourceType, uri, param) {
			return loader, nil
		}
	}

	return EmptyResourceLoader{}, fmt.Errorf("no resource loader registred that could handle: \"%s/%v\"", resourceType, param)
}

func (r *ResourceManager) PreloadReource(resourceType ResourceType, uri string, param LoaderParam) *ResourceManager {
	loader, err := r.GetLoader(resourceType, uri, param)
	if err != nil {
		log.Fatalf("FATAL! %s", err)
	}

	rsc, err := loader.Load(uri, param)

	if err != nil {
		log.Printf("FAILED! preloading of resource failed: %s(%v). %q\n", resourceType, param, err)
	}

	// Add the resource inven if failed
	r.resources[rsc.Uri] = rsc
	if !rsc.Empty {
		log.Printf("LOADED! resource loaded: %s(%v) -> %q\n", resourceType, param, rsc.Uri)
	} else {
		log.Printf("LOADED! empty resource loaded: %s(%v) -> %q\n", resourceType, param, rsc.Uri)
	}

	return r
}

func (r *ResourceManager) GetResource(resourceUri string) Resource {
	rsc, ok := r.resources[resourceUri]

	if !ok {
		log.Fatal("resource not found. At this point an empty resource should exist")
	}

	return rsc
}

// func (r *ResourceManager) AddDefaultLoaders() *ResourceManager {
// 	return r.
// 		AddLoader(RT_TEXTURE, NewFileTextureLoader()).
// 		AddLoader(RT_MESH_2D, NewProceduralMesh2dLoader()).
// 		AddLoader(RT_SHADER, NewShaderLoader())
// }

func (r *ResourceManager) Unload() {
	for _, r := range r.resources {
		if r.Unload != nil {
			r.Unload()
		}
	}

	r.resources = make(map[string]Resource)
}

func NewResourceManager() *ResourceManager {
	return &ResourceManager{
		resourceLoaders: make([]ResourceLoader, 0),
		resources:       make(map[string]Resource),
	}
}
