package data

//DBType represents a datasource type
type DBType string

const (
	POSTGRES DBType = "postgres"
	MONGO    DBType = "mongo"
	EDGEDB   DBType = "edgedb"
)

//Properties represents a datasource properties
type Properties struct {
	Type     DBType `yaml:"type"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
}

//WithProperties represents an object having properties
type WithProperties struct {
	Properties *Properties
}

//IsValid checks if properties is valid
func (p Properties) IsValid() bool {
	return p.Host != "" && p.Database != ""
}

//IsValid checks if type is valid
func (w WithProperties) IsValid() bool {
	return w.Properties.IsValid()
}
