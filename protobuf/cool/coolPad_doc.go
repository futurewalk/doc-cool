package protobufCool

import "github.com/futurewalk/doc-cool/cool"
import "doc-cool/controllers"
import "doc-cool/protobuf"

func Init() {
    container := make(map[string]interface{})
    container["CoolController"]= &controllers.CoolController{}
    container["Request"]= &protobuf.Request{}
    container["Response"]= &protobuf.Response{}
    container["Extension"] = &Resolver{} 
    cool.Start(container)
}

type Resolver struct{}

func (p *Resolver) Invoke(plugin *cool.Plugin) {
}
