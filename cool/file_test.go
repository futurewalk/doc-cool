package cool

import (
    "testing"
    "fmt"
    "os"
    "log"
)

func TestExportLogo(t *testing.T)  {
    workPath, err := os.Getwd()
    log.Println("获取相对路径:",workPath)

    cfg := Config{}
    configPath := "banner.txt"
    err = cfg.WalkAll(configPath, func(cnt string) {
        fmt.Println(cnt)
    })
    if err != nil {
        panic(err)
    }
}
