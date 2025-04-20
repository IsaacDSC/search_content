package reader

import (
	"github.com/IsaacDSC/search_content/internal/content/entity"
	"strings"
)

type ListEnterprises map[entity.EnterpriseKey]EnterpriseData

type EnterpriseData map[entity.PathKey]entity.Enterprise

func NewEnterprisesData(pathKey entity.PathKey, data entity.Enterprise) EnterpriseData {
	return EnterpriseData{pathKey: data}
}

func (ed EnterpriseData) Append(pathKey entity.PathKey, data entity.Enterprise) EnterpriseData {
	ed[pathKey] = data
	return ed
}

func (ed EnterpriseData) GetContent(input entity.PathKey) (entity.Video, bool) {
	content, found := ed[input]
	if found {
		return content.Video, true
	}

	inputPaths := input.ToListPaths()

	for _, v := range ed {
		url := v.Url
		paths := entity.NewPathKey(url).ToListPaths()
		totalMatch := 0

		if len(inputPaths) < len(paths) {
			continue
		}

		var length int
		if len(inputPaths) > len(paths) {
			if !strings.Contains(url.Path, "*") {
				continue
			}
			length = len(paths)
		} else {
			length = len(inputPaths)
		}

		for i := 0; i < length; i++ {
			if paths[i] == inputPaths[i] {
				totalMatch++
			}
			if paths[i] == "*" {
				totalMatch++
			}
		}

		if totalMatch == length && (len(paths) != 0 && len(inputPaths) != 0) {
			return v.Video, true
		}
	}

	return entity.Video{}, false
}
