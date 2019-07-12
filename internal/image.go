package syncd

import (
	"fmt"
	"strings"
)

type Image struct {
	Registry   string
	Repository string
	Tag        string
}

func (i Image) URI() string {
	return i.Registry + "/" + i.Repository
}

func ParseImage(repo string, tag string) (*Image, error) {
	arr := strings.Split(repo, "/")
	if len(arr) != 2 {
		return nil, fmt.Errorf("malformed repo: %s", repo)
	}

	return &Image{
		Registry:   arr[0],
		Repository: arr[1],
		Tag:        tag,
	}, nil
}
