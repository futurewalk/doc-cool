package cool

import "strings"

type Configuration interface {
    Set(k, v string)
    Get(key string) string
    LoadConfig(path string) error
    GetIgnoreFile(files string)
}

type Config struct {
    key   string
    value string
    Data  map[string]string
    Base
}

func (cfg *Config) Set(k, v string) {
    data := cfg.Data
    if data == nil {
        data = make(map[string]string)
        cfg.Data = data
    }
    data[k] = v
}
func (cfg *Config) Get(key string) string {
    if cfg.Data == nil {
        return ""
    }
    return cfg.Data[key]
}
func (cfg *Config) LoadConfig(path string) error {
    configPath := path + separator + "app.conf"
    err := cfg.Walk(configPath, func(cnt string) {
        if Match(cnt, coolConfigRxp) {
            eqIdx := strings.Index(cnt, "=")
            cfg.Set(trimSpace(cnt[:eqIdx]), strings.TrimSpace(cnt[eqIdx+1:]))
        }
    })

    if err != nil {
        return err
    }

    container.config = *cfg

    stucture := Structures{
        ImportList:         make(map[string]string),
        StructureContainer: make(map[string][]StructureContainer),
    }
    err = stucture.GetControlStructure(getPath(separator + "controllers"))

    if err != nil {
        return err
    }

    stucture.ScanProtoStructure(getPath(separator + cfg.Get("cool.protoPath")))

    container.structures = stucture

    return err
}
func (cfg *Config) GetIgnoreFile(files string) {
    arr := strings.Split(files, ",")
    for _, value := range arr {
        container.structures.IgnoreFile[value] = value
    }
}
