package yaml_test

import (
	"errors"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/crypto-com/chainindex/internal/filereader"
	"github.com/crypto-com/chainindex/internal/filereader/yaml"
)

type AnyInterface struct {
	Version string `yaml:"version"`
}

var _ = Describe("YamlConfigReader", func() {
	Describe("FromFile", func() {
		It("should implement Reader interface", func() {
			yamlReader, _ := yaml.FromFile("./assets/any_valid_config.yml")

			var _ filereader.Reader = yamlReader
		})

		It("should return Error when file does not exists", func() {
			_, err := yaml.FromFile("unexisted-file.yaml")

			Expect(err).NotTo(BeNil())
			Expect(errors.Is(err, filereader.ErrFileNotFound)).To(BeTrue())
		})
	})

	Describe("FromIOReader", func() {
		It("should implement Reader interface", func() {
			const anyValidYAMLConfig = `
version: 1.0.0
`
			ioReader := strings.NewReader(anyValidYAMLConfig)

			yamlReader := yaml.FromIOReader(ioReader)

			var _ filereader.Reader = yamlReader
		})
	})

	Describe("Read", func() {
		It("should return Error when the yaml syntax is invalid", func() {
			anyInvalidYAML := "invalid"
			ioReader := strings.NewReader(anyInvalidYAML)

			yamlReader := yaml.FromIOReader(ioReader)

			var value AnyInterface
			err := yamlReader.Read(&value)

			Expect(errors.Is(err, filereader.ErrReadFile)).To(BeTrue())
		})

		It("should return Error when the yaml contain unknown fields", func() {
			anyYAMLWithUnknownField := `
unknown: unknown
`
			ioReader := strings.NewReader(anyYAMLWithUnknownField)

			yamlReader := yaml.FromIOReader(ioReader)

			var config AnyInterface
			err := yamlReader.Read(&config)

			Expect(errors.Is(err, filereader.ErrReadFile)).To(BeTrue())
		})

		Context("FromFile", func() {
			It("should return value based on file", func() {
				yamlReader, fromFileErr := yaml.FromFile("./assets/any_valid_config.yml")
				Expect(fromFileErr).To(BeNil())

				var value AnyInterface
				readErr := yamlReader.Read(&value)
				Expect(readErr).To(BeNil())

				expectedValue := AnyInterface{
					Version: "1.0.0",
				}
				Expect(value).To(Equal(expectedValue))
			})
		})

		Context("FromIOReader", func() {
			It("should return value parsed from the io.Reader", func() {
				const anyValidYAMLConfig = `
version: 1.0.0
`
				ioReader := strings.NewReader(anyValidYAMLConfig)

				yamlReader := yaml.FromIOReader(ioReader)

				var value AnyInterface
				readErr := yamlReader.Read(&value)
				Expect(readErr).To(BeNil())

				expectedValue := AnyInterface{
					Version: "1.0.0",
				}
				Expect(value).To(Equal(expectedValue))
			})
		})
	})
})
