package main

import ini "gopkg.in/ini.v1"

// LoadConfig loads is a wrapper around ini.MapTo function, which enables mysql type bools.
func LoadConfig(path string) (*Configuration, error) {
	cfg, err := ini.LoadSources(ini.LoadOptions{AllowBooleanKeys: true}, path)
	cfg.NameMapper = ini.TitleUnderscore
	Config := new(Configuration)
	err = cfg.MapTo(Config)
	return Config, err
}

// Configuration is the struct in which all the configuration is loaded.
type Configuration struct {
	LogFileLocation  string
	DataFileLocation string
	MaxWorkerThreads int
	Source           string
	Destination      string
	Saveguard
	Scantype
}

type Scantype struct {
	Blake2bBased   bool
	TimestampBased bool
}

type Saveguard struct {
	SaveguardMaxPercentage int
	SaveguardEnable        bool
}
