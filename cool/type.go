package cool

import (
    "reflect"
    "github.com/astaxie/beego"
    "io"
    "github.com/golang/protobuf/proto"
)

type Annotation struct {
    Id               string
    Url              string
    Body             map[string]interface{}
    Method           string
    RespData         map[string]interface{}
    Controller       string
    Description      string
    ProtoBufControl  string
    RespProtoBufCtrl string
    ProtoBufFileName string
}
type DocController struct {
    beego.Controller
}
type Container struct {
    Key           string
    PkPath        string
    DocCtrl       *DocController
    DocName       string
    CoolDocs      map[string]*Annotation
    RequestHeader map[string]interface{}
}
type ProtoFileContainer struct {
    Package    string
    StructList []string
}
type structField struct {
    reqType   reflect.Type
    reqInst   reflect.Value
    rspType   reflect.Type
    rspInst   reflect.Value
    container map[string]interface{}
    ext       Extension
    url       string
}
type Plugin struct {
    StructName string
    FieldName  string
    Url        string
    swapStruct proto.Message
}
type Ignore struct {
    ProtoFile string
}
type Extension interface {
    Invoke(plugin *Plugin)
}

type Cooler interface {
    //scan project all protobuf struct
    scanProtoBuf(path string)

    //get the controller file path
    getControllers(path string)

    //readFile
    walk(path string, fc func(cnt string))

    //deal go file to contains
    dealGoFile(fp, fln string, ant *Annotation)

    getAnnotation(cnt string, ant *Annotation, sfd *structField) *Annotation

    //scan proto file ,if it assign from proto.message then return it
    getProtoBuf(sp, rp string)

    //analyze all struct and get useful struct
    analyzeStruct(sp, rp string, flt bool)

    //generate register file
    generateProtoBufGoFile() error

    generateStruct(w io.Writer, ptc ProtoFileContainer, pk string) error

    //new struct
    newInstance(object interface{}) (reflect.Type, reflect.Value)

    //get protoStruct field
    getFields(sfd *structField) map[string]interface{}

    //provider a method for API
    CreateDocument()

    //create annotation router
    createRouter()
}
