package cool

import (
    "github.com/golang/protobuf/proto"
)

type Plugin struct {
    StructureName string
    FieldName     string
    Url           string
    swapStructure proto.Message
}
type Extension interface {
    Invoke(plugin *Plugin)
}

func (p *Plugin) Remove() {
    p.FieldName = ""
}
func (p *Plugin) Swap(swp proto.Message) {
    p.swapStructure = swp
}