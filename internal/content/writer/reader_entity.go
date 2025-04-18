package writer

type ListEnterprises map[EnterpriseKey]EnterpriseData

type EnterpriseData map[PathKey]Enterprise

func NewEnterprisesData(pathKey PathKey, data Enterprise) EnterpriseData {
	return EnterpriseData{pathKey: data}
}

func (ed EnterpriseData) Append(pathKey PathKey, data Enterprise) EnterpriseData {
	ed[pathKey] = data
	return ed
}
