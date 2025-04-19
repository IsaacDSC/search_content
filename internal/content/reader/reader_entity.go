package reader

import (
	"github.com/IsaacDSC/search_content/internal/content/entity"
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
		paths := entity.NewPathKey(v.Url).ToListPaths()
		totalMatch := 0

		var length int
		equalResult := len(inputPaths) - len(paths)
		if equalResult < 0 {
			length = equalResult * -1
		} else if equalResult == 0 {
			length = len(inputPaths)
		} else {
			length = equalResult
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
