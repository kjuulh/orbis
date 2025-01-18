package modelregistry

import "sync"

type Model struct {
	Name     string `json:"name"`
	Schedule string `json:"schedule"`
	Lookup   string `json:"lookup"`
}

type ModelRegistry struct {
	models     []Model
	modelsLock sync.RWMutex
}

func NewModelRegistry() *ModelRegistry {
	return &ModelRegistry{
		models: make([]Model, 0, 128), // We start off with a capacity of 128 models. Makes memory more deterministic
	}
}

func (m *ModelRegistry) GetModels() ([]Model, error) {
	m.modelsLock.RLock()
	defer m.modelsLock.RUnlock()

	return []Model{
		{
			Name:     "69C42481-650D-46E0-9C96-3D61B96565FB",
			Schedule: "* * * * * *",
			Lookup:   "something",
		},
		{
			Name:     "5B0F96E5-BC37-427E-B615-E635156386F0",
			Schedule: "* * * * * *",
			Lookup:   "something",
		},
		{
			Name:     "5A511180-6613-4F8E-9125-2E8FE272E03C",
			Schedule: "* * * * * *",
			Lookup:   "something",
		},
		{
			Name:     "4EFE77E8-082B-4828-8527-635E5B6253D9",
			Schedule: "* * * * * *",
			Lookup:   "something",
		},
		{
			Name:     "E53DCA1E-641B-421A-9DB6-1A6F09F3D96D",
			Schedule: "* * * * * *",
			Lookup:   "something",
		},
		{
			Name:     "63BDC98C-ECBA-44FD-BFAE-056B1C004077",
			Schedule: "* * * * * *",
			Lookup:   "something",
		},
		{
			Name:     "C18A1948-0045-4099-AC58-5B7C587AC0F0",
			Schedule: "* * * * * *",
			Lookup:   "something",
		},
		{
			Name:     "8B87FDB5-A119-43C0-8D15-9B517577A8AE",
			Schedule: "* * * * * *",
			Lookup:   "something",
		},
		{
			Name:     "4E121E78-1CBD-4BC1-8A10-14354B76E553",
			Schedule: "* * * * * *",
			Lookup:   "something",
		},
		{
			Name:     "14AF7CDF-783F-4DFE-8B3D-E4C23C12AEDC",
			Schedule: "* * * * * *",
			Lookup:   "something",
		},
		{
			Name:     "60FE99E9-4EF7-40A5-A19D-9A439BA12B24",
			Schedule: "* * * * * *",
			Lookup:   "something",
		},
		{
			Name:     "622C03C2-CAF7-4708-B26D-D536E3C3F4DD",
			Schedule: "* * * * * *",
			Lookup:   "something",
		},
		{
			Name:     "E7BC4A8D-FFF6-4A8D-A48B-6569340746E4",
			Schedule: "* * * * * *",
			Lookup:   "something",
		},
		{
			Name:     "7CB66BA1-CF1E-4FA6-8C32-F048D88FCE54",
			Schedule: "* * * * * *",
			Lookup:   "something",
		},
		{
			Name:     "82518E22-EFED-4AA8-AC19-FF3D81ECE609",
			Schedule: "* * * * * *",
			Lookup:   "something",
		},
	}, nil

	//return m.models, nil
}
