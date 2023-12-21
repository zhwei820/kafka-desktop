package conf

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/zhwei820/gconv"
)

// ========
type DBType string

const (
	MySQL DBType = "mysql"
)

type KafkaConfig struct {
	Type DBType
	Name string
	Host string
}

var LocalKafkaConfig = KafkaConfig{ //for debug
	Type: MySQL,
	Name: "local",
	Host: "10.171.22.153:30200",
}

type Config struct {
	homeDir string

	Instance map[string]*KafkaConfig
}

func (c *Config) Configs() *Configs {
	res := &Configs{ConfigNames: make(map[string]bool)}
	for key := range c.Instance {
		if strings.Contains(key, "_table_setting") {
			continue
		}
		res.ConfigNames[key] = true
	}
	return res
}

type Configs struct {
	ConfigNames map[string]bool
}

func New() *Config {
	// Get the current user
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	// Access the home directory
	homeDir := currentUser.HomeDir
	return &Config{
		Instance: make(map[string]*KafkaConfig),
		homeDir:  homeDir + "/",
	}
}

func (c *Config) SaveAll() {
	os.Mkdir("kafka_conf", 0755)
	if len(c.Configs().ConfigNames) != 0 {
		saveFileWithHistory(c.homeDir+"kafka_conf/conf_db.json", c.homeDir+"kafka_conf")
		gconv.WritePlainToFile(c.Configs(), c.homeDir+"kafka_conf/conf_db.json")
		for key := range c.Configs().ConfigNames {
			if c.Instance[key] != nil {
				gconv.WritePlainToFile(c.Instance[key], c.homeDir+"kafka_conf/conf_"+key+".json")
			}
		}
	}
}
func (c *Config) LoadAll() {
	keys := &Configs{}
	gconv.LoadPlainJsonFromFile(keys, c.homeDir+"kafka_conf/conf_db.json")
	fmt.Println("keys", gconv.Export(keys))
	for key := range keys.ConfigNames {
		c.Instance[key] = &KafkaConfig{}
		gconv.LoadPlainJsonFromFile(c.Instance[key], c.homeDir+"kafka_conf/conf_"+key+".json")
		c.Instance[key].Name = key
	}
}
