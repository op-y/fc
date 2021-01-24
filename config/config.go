package config

import (
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

var BackendType []string = []string{"oss", "fastdfs", "local"}

var Cfg *Config

type Config struct {
	AppLog       string             `yaml:"appLog"`
	GinLog       string             `yaml:"ginLog"`
	Port         string             `yaml:"port"`
	Backend      string             `yaml:"backend"`
	Verification VerificationConfig `yaml:"verification"`
	PosLog       PosLogConfig       `yaml:"posLog"`
	DbBackup     DbBackupConfig     `yaml:"dbBackup"`
	BillFile     BillFileConfig     `yaml:"billFile"`
	Oss          OssConfig          `yaml:"oss"`
	FastDfs      FastDfsConfig      `yaml:"fastDfs"`
	Local        LocalConfig        `yaml:"local"`
}

type VerificationConfig struct {
	Enable bool   `yaml:"enable"`
	Method string `yaml:"method"`
	Salt   string `yaml:"salt"`
}

type PosLogConfig struct {
	Prefix string `yaml:"prefix"`
}

type DbBackupConfig struct {
	Prefix string `yaml:"prefix"`
}

type BillFileConfig struct {
	Prefix string `yaml:"prefix"`
	Target string `yaml:"target"`
}

type OssConfig struct {
	Endpoint        string `yaml:"endpoint"`
	Bucket          string `yaml:"bucket"`
	AccessKeyId     string `yaml:"accessKeyId"`
	AccessKeySecret string `yaml:"accessKeySecret"`
}

type FastDfsConfig struct {
	FdfsConf string `yaml:"fdfsConf"`
}

type LocalConfig struct {
	Root string `yaml:"root"`
}

func (cfg *Config) Print() {
	log.Println("fc configuration show as blelow:")
	log.Printf("port: %s\n", cfg.Port)
	log.Printf("backend: %s\n", cfg.Backend)
	log.Printf("verification enable: %t\n", cfg.Verification.Enable)
	log.Printf("verification method: %s\n", cfg.Verification.Method)
	log.Printf("verification salt: %s\n", cfg.Verification.Salt)
	log.Printf("pos log prefix: %s\n", cfg.PosLog.Prefix)
	log.Printf("db backup prefix: %s\n", cfg.DbBackup.Prefix)
	log.Printf("bill file prefix: %s\n", cfg.BillFile.Prefix)
	log.Printf("bill file target: %s\n", cfg.BillFile.Target)
	log.Printf("oss endpoint: %s\n", cfg.Oss.Endpoint)
	log.Printf("oss bucket: %s\n", cfg.Oss.Bucket)
	log.Printf("oss access key ID: ***\n")
	log.Printf("oss access key secret: ***\n")
	log.Printf("fast dfs config file: %s\n", cfg.FastDfs.FdfsConf)
	log.Printf("local storage root: %s\n", cfg.Local.Root)
}

func (cfg *Config) Check() bool {
	if cfg.AppLog == "" {
		log.Println("fc need a app log file")
		return false
	}

	if cfg.GinLog == "" {
		log.Println("fc need a gin log file")
		return false
	}

	if cfg.Port == "" {
		log.Println("fc need a valid port")
		return false
	}

	if cfg.Backend == "" {
		log.Println("fc need a backend")
		return false
	} else {
		hit := false
		for _, t := range BackendType {
			if cfg.Backend == t {
				hit = true
				break
			}
		}
		if !hit {
			log.Println("fc need a valid backend: oss/fastdfs/local")
			return false
		}
	}

	if cfg.Verification.Enable {
		if cfg.Verification.Method == "" {
			log.Println("verification need method")
			return false
		}

		if cfg.Verification.Salt == "" {
			log.Println("verification need salt")
			return false
		}
	}

	if cfg.PosLog.Prefix == "" {
		log.Println("fc need a pos log prefix")
		return false
	}

	if cfg.DbBackup.Prefix == "" {
		log.Println("fc need a pos db backup prefix")
		return false
	}

	if cfg.BillFile.Prefix == "" {
		log.Println("fc need a bill file prefix")
		return false
	}

	if cfg.BillFile.Target == "" {
		log.Println("fc need a bill file target")
		return false
	}

	if cfg.Backend == "oss" {
		if cfg.Oss.Endpoint == "" {
			log.Println("OSS backend need a endpoint")
			return false
		}
		if cfg.Oss.Bucket == "" {
			log.Println("OSS backend need a bucket")
			return false
		}
		if cfg.Oss.AccessKeyId == "" {
			log.Println("OSS backend need a accessKeyId")
			return false
		}
		if cfg.Oss.AccessKeySecret == "" {
			log.Println("OSS backend need a accessKeySecret")
			return false
		}
	}

	if cfg.Backend == "fastdfs" {
		if _, err := os.Stat(cfg.FastDfs.FdfsConf); err != nil {
			log.Println("Fast DFS backend need a fastdfs.conf")
			return false
		}
	}

	if cfg.Backend == "local" {
		if err := os.MkdirAll(cfg.Local.Root, 0755); err != nil {
			log.Printf("fail to create local root directory: %s", err.Error())
			return false
		}

		if err := os.MkdirAll(cfg.Local.Root+cfg.PosLog.Prefix, 0755); err != nil {
			log.Printf("fail to create backup directory: %s", err.Error())
			return false
		}

		if err := os.MkdirAll(cfg.Local.Root+cfg.DbBackup.Prefix, 0755); err != nil {
			log.Printf("fail to create backup directory: %s", err.Error())
			return false
		}

		if err := os.MkdirAll(cfg.Local.Root+cfg.BillFile.Prefix, 0755); err != nil {
			log.Printf("fail to create backup directory: %s", err.Error())
			return false
		}
	}

	return true
}

func LoadConfig(conf string) *Config {
	cfg := new(Config)
	buf, err := ioutil.ReadFile(conf)
	if err != nil {
		log.Printf("fail to read config file: %s", err.Error())
		panic("config error")
	}

	if err := yaml.Unmarshal(buf, cfg); err != nil {
		log.Printf("fail to unmarshal yaml config: %s", err.Error())
		panic("config error")
	}

	if cfg == nil {
		log.Println("nil config")
		panic("config error")
	}

	if !cfg.Check() {
		log.Println("some mistakes in config file")
		panic("config error")
	}

	return cfg
}
