package main

var realResources = []string{
	"logo.png",
	"logo_invert.png",
	"admin/index.html",
	"admin/linkgen.html",
	"admin/usercount.js",
	"admin/navbar.js",
}

var (
	// binary resources
	// TODO change to map string to []byte
	resources = (func() ResourceReloader {
		var r ResourceReloader
		r = make(ResourceReloader, len(realResources))
		for _, v := range realResources {
			r[v] = load(v)
		}
		return r
	})()

	// if resource not in here, then it's at: "resource/<name>"
	resourceLocation = ResourceReloader{
		"linkgen.html": "admin/linkgen.html",
		"admin":        "admin/index.html",
		"admin/":       "admin/index.html",
	}
)

type ResourceReloader map[string]string

func (r *ResourceReloader) Get(name string) string {
	if *nc {
		for _, v := range realResources {
			if v == name {
				(*r)[name] = load(name)
				break
			}
		}
	}
	return (*r)[name]
}

func (r *ResourceReloader) Has(name string) bool {
	_, ok := (*r)[name]
	return ok
}
