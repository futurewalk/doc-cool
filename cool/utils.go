package cool

import (
    "regexp"
    "strings"
    "io"
    "os"
    "fmt"
    "bytes"
    "github.com/astaxie/beego/context"
    "runtime"
)

const (
    ostype              = runtime.GOOS
    funcRxp             = `(?i:^fun).*(){`
    reqDataRxp          = `//@ReqData`
    rspRxp              = `//@RespData`
    reqFileRxp          = `//@ReqProtoFile`
    rspFileRxp          = `//@RespProtoFile`
    methodRxp           = `//@Method`
    coolUrlRxp          = `//@Url`
    descriptionRxp      = `//@Description`
    proGoFileRxp        = `(?i:^[a-z]).*pb.go`
    proFileRxp          = `(?i:^[a-z]).*.proto`
    pathDotRxp          = `(?i:^./).*`
    ptrFuncRxp          = `(?i:^func).*`
    methodEndCloseRxp   = `()`
    structRxp           = `(?i:^type).*struct {`
    goPackageRxp        = `(?i:^package).*`
    goType              = `type`
    goStruct            = `struct`
    goFilePackage       = `package`
    leftBrace           = `{`
    rightBrace          = `}`
    defaultGeneratePath = `./protobuf/cool`
    generateGoFile      = `/doc_cool_register.go`
    generatePackage     = `package protobuf_cool`
    coolPackage         = `github.com/futurewalk/doc-cool/cool`
    coolContainer       = `    container := make(map[string]interface{})`
    coolStartMethod     = `    cool.Start(container)`
    newLine             = "\n"
    initMethod          = `func init() {`
    backSlash           = `/`
    importStr           = `import `
    oneQuoMark          = `"`
    resolverInit        = `    container["Extension"] = &Resolver{} `
    resolverStruct      = `type Resolver struct{}`
    invokeMethod        = `func (p *Resolver) Invoke(plugin *cool.Plugin) {`
)

var (
    src       = `\src\`
    separator = `\`
)

func fs() {
    if ostype == "linux" {
        src = `/src/`
        separator = `/`
    }
}
func init() {
    fs()
}

//regexp match this string
func Match(content string, regexps string) bool {
    reg := regexp.MustCompile(regexps)
    return reg.MatchString(content)
}
func trimSpace(value string) string {
    return trim(value, " ")
}
func trim(v string, rt ...string) string {
    for _, vrt := range rt {
        v = strings.Replace(v, vrt, "", -1)
    }
    return v
}
func trstr2spas(cont, rts string) string {
    newStr := strings.Replace(trimSpace(cont), rts, "", -1)
    return newStr
}

//check file is Exist
func isExist(filename string) bool {
    var exist = true
    if _, err := os.Stat(filename); os.IsNotExist(err) {
        exist = false
    }
    return exist
}

//write content to file
func write(w io.Writer, lineStr string) (int, error) {
    return fmt.Fprintln(w, lineStr)
}
func writeLine(w io.Writer, lineStr string) (int, error) {
    return fmt.Fprintln(w, lineStr+newLine)
}

type StringBuilder struct {
    bytes.Buffer
}

func (sb *StringBuilder) Append(str string) *StringBuilder {
    sb.WriteString(str)
    return sb
}
func getImportPath(path string) string {

    arrayStr := strings.Split(path, src)
    importPath := strings.Replace(arrayStr[1], separator, backSlash, -1)
    return importPath
}
func getStructName(content string) string {
    return trim(trimSpace(content), goType, goStruct, leftBrace)
}
func isNotNull(values ...string) bool {
    for _, v := range values {
        if v == "" {
            return false
        }
    }
    return true
}
func getProPath(path string) string {
    if !Match(path, pathDotRxp) {
        return "./" + path
    }
    return path
}
func filterPath(path string) string {
    if strings.Index(path, ".") > -1 {
        path = strings.Replace(path, ".", "/", -1)
    }
    return path
}
func getGeneratePk(value string) string {
    if strings.Index(value, ".") > -1 {
        value = strings.Replace(value, ".", "_", -1)
    }
    return strings.Replace(value, "/", "_", -1)
}
func GetIgnoreFile(files string) {
    arr := strings.Split(files, ",")
    for _, value := range arr {
        ignoreFile[value] = value
    }
}
func sliceType(sn string) string {
    if !strings.Contains(sn, ".") {
        return ""
    }
    arr := strings.Split(sn, ".")
    if len(arr) > 0 {
        return arr[1]
    }
    return ""
}
func subLst(s, p string) string {
    lst := strings.LastIndex(s, p)
    s = string([]byte(s)[lst+1:])
    return s
}
func getIndexStr(s, p string, index int) string {
    arr := strings.Split(s, p)
    return arr[index]
}

func Access(ctx *context.Context) {
    ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")                        //允许访问源
    ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS") //允许post访问
    ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "*")                       //header的类型
    ctx.ResponseWriter.Header().Set("Access-Control-Max-Age", "1728000")
    ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
    ctx.ResponseWriter.Header().Set("Content-type", "application/x-protobuf")
}
