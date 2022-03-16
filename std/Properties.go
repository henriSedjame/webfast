package std

import (
	"fmt"
	"github.com/hsedjame/webfast/data"
	"github.com/hsedjame/webfast/utils"
	"github.com/hsedjame/webfast/web"

	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

const PropsFile = "application"

type Properties struct {
	ActiveProfiles Profiles              `yaml:"profiles"`
	Datasource     *data.Properties      `yaml:"datasource"`
	Server         *web.ServerProperties `yaml:"server"`
	Cors           *web.CorsProperties   `yaml:"cors"`
}

type Profiles struct {
	Active string `yaml:"active"`
}

func AppDefaultProperties() *Properties {
	return &Properties{
		Server: &web.ServerProperties{
			Host: "localhost",
			Port: 8080,
		},
		Cors: &web.CorsProperties{
			AllowedHeaders: "*",
			AllowedMethods: "*",
			AllowedOrigins: "*",
		},
	}
}

// Load : Load properties from a yaml file
func (props *Properties) Load(classpath string) error {

	if err := props.doLoad(fmt.Sprintf("%s/%s.%s",
		classpath, PropsFile, utils.PropertiesFileExtension)); err != nil {
		return err
	}

	profiles := props.ActiveProfiles.Active

	for _, profile := range strings.Split(profiles, ",") {
		profile = strings.TrimSpace(profile)

		if profile != "" {

			fileLocation := fmt.Sprintf("%s/%s-%s.%s",
				classpath, PropsFile, profile, utils.PropertiesFileExtension)

			if err := props.doLoad(fileLocation); err != nil {
				return err
			}
		}
	}

	return nil
}

func (props *Properties) doLoad(fileLocation string) error {

	if file, err := os.OpenFile(fileLocation, os.O_RDONLY, 0777); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	} else {
		defer file.Close()

		if bytes, err := ioutil.ReadAll(file); err != nil {
			return err
		} else {
			if err := yaml.Unmarshal(bytes, props); err != nil {
				return err
			}
		}
	}
	return nil
}
