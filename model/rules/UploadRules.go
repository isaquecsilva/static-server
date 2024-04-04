package rules

import (
	"encoding/xml"
	"os"
)

const defaultRulesFile string = `<Rules>
	<FileTypes>
		<FileType>*.png</FileType>
		<FileType>*.jpg</FileType>
		<FileType>*.mp4</FileType>
	</FileTypes>
	<MaxFileSize>100MB</MaxFileSize>
</Rules>`

type UploadRules struct {
	FileTypes struct {
		FileTypeList []string `xml:"FileType"`
	} `xml:"FileTypes"`
	MaxFileSize       string `xml:"MaxFileSize"`
	MaxUploadsPerHour uint8  `xml:"MaxUploadsPerHour"`
}

func LoadUploadRulesFromFile(rulesXmlFile string) (*UploadRules, error) {
	file, err := os.Open(rulesXmlFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var rules UploadRules
	if err := xml.NewDecoder(file).Decode(&rules); err != nil {
		return nil, err
	}

	return &rules, nil
}

func CreateDefaultRulesFileTemplate() error {
	return os.WriteFile("rules.xml", []byte(defaultRulesFile), 0755)
}
