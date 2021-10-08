package resource

import "embed"

//go:embed embeded
var embededResources embed.FS

func readEmbededFileAsString(name string) string {
	data, err := embededResources.ReadFile("embeded/" + name)
	if err != nil {
		return ""
	}

	return string(data)
}
