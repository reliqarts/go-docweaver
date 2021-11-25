package docweaver

import (
	yml "gopkg.in/yaml.v3"
	"os"
)

type source struct {
	Key string
	Url string
}
type sourceConfig struct {
	Sources []source
}

func readSources() (sc *sourceConfig, err error) {
	yaml, err := os.ReadFile(GetSourcesFilePath())
	if err != nil {
		return nil, err
	}

	if err := yml.Unmarshal(yaml, &sc); err != nil {
		return nil, err
	}

	return sc, nil
}
