package yaml

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/crypto-com/chainindex/internal/filereader"
)

type YAMLReader struct {
	yamlFile io.Reader
}

func FromFile(filePath string) (*YAMLReader, error) {
	var err error

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file %s: %w", filePath, filereader.ErrFileNotFound)
	}

	return FromIOReader(file), nil
}

func FromIOReader(reader io.Reader) *YAMLReader {
	return &YAMLReader{
		yamlFile: reader,
	}
}

func (reader *YAMLReader) Read(value interface{}) error {
	var err error

	decoder := yaml.NewDecoder(reader.yamlFile)
	decoder.KnownFields(true)

	if err = decoder.Decode(value); err != nil {
		return fmt.Errorf("error decoding YAML: %v: %w", err, filereader.ErrReadFile)
	}

	return nil
}
