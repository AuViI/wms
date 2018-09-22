package main

var (
	// binary resources
	// TODO change to map string to []byte
	resources = ResourceReloader{
		"logo.png":        load("logo.png"),
		"logo_invert.png": load("logo_invert.png"),
		"linkgen.html":    load("linkgen.html"),
	}
	// if resource not in here, then it's at: "resource/<name>"
	resourceLocation = ResourceReloader{}
)

type ResourceReloader map[string]string

func (r *ResourceReloader) Get(name string) string {
	if *nc {
		(*r)[name] = load(name)
	}
	return (*r)[name]
}

func (r *ResourceReloader) Has(name string) bool {
	_, ok := (*r)[name]
	return ok
}
