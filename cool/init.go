package cool

import (
    "github.com/astaxie/beego"
    "reflect"
    "fmt"
    "net/http"
    "encoding/json"
)

type Cool struct {
    ant                Ant
    config             Config
    structures         Structures
    registerStructures map[string]interface{}
}

var container Cool

func Export() {

    config := Config{}

    err := config.LoadConfig(getPath(separator + "conf"))

    if err != nil {
        panic(err)
    }
    err = config.GenerateRegisterGoFile()

    if err != nil {
        panic(err)
    }

}

type DocController struct {
    beego.Controller
}

func (p *DocController) CreateDocument() {

    var dataRow = make(map[string]interface{})
    dataRow["docName"] = container.config.Get("cool.docName")
    dataRow["docPkPath"] = container.config.Get("cool.protoPackage")
    dataRow["docKey"] = container.config.Get("cool.key")
    dataRow["rows"] = container.structures
    //dataRow["header"] = cooler.container.RequestHeader

    //log.Println(dataRow)
    p.Data["json"] = dataRow
    Access(p.Ctx)
    p.ServeJSON()
}

func CreateDocument(w http.ResponseWriter, r *http.Request) {

    var dataRow = make(map[string]interface{})
    dataRow["docName"] = container.config.Get("cool.docName")
    dataRow["docPkPath"] = container.config.Get("cool.protoPackage")
    dataRow["docKey"] = container.config.Get("cool.key")
    dataRow["rows"] = container.structures
    b, _ := json.Marshal(dataRow)

    fmt.Fprintf(w, string(b[:]))
}

func createRouter() {

    docs := container.ant.Ants
    for _, value := range docs {
        beego.Router(value.Url, container.registerStructures[value.Controller].(beego.ControllerInterface), value.Method)
        fmt.Println("router create success:", value.Url)
    }
    beego.Router("/coolDocument/createDocument", &DocController{}, "*:CreateDocument")
}

func doRegister(structure ...interface{}) {
    for i := 0; i < len(structure); i++ {
        v := reflect.ValueOf(structure[i]).Elem()
        container.registerStructures[v.Type().Name()] = structure[i]
    }
}

func Register(ext Extension, structure ...interface{}) {

    config := Config{}

    container.registerStructures = make(map[string]interface{})

    container.registerStructures["Extension"] = ext

    doRegister(structure...)

    GetIgnoreFile(container.config.Get("cool.ignoreFile"))

    err := config.LoadConfig(getPath(separator + "conf"))

    if err != nil {
        panic(err)
    }
    createRouter()

    //config.ShowBanner()
}
