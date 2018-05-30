package cool

import (
	"path/filepath"
	"fmt"
	"os"
	"io/ioutil"
	"bufio"
	"io"
	"github.com/astaxie/beego"
	"reflect"
	"github.com/gogo/protobuf/proto"
	"strings"
	"log"
)

var (
	container      = new(Container)
	protoContainer = map[string][]ProtoFileContainer{}
	imports        = make(map[string]string)
	structs        = make(map[string]interface{})
	protoFile      = make(map[string]string)
	ignoreFile     = map[string]string{}
)

func Export() {
	var doc DocController
	doc.getControllers("./controllers/")
	doc.scanProtoBuf(getProPath(beego.AppConfig.String("cool.protoPath")))
	doc.generateProtoBufGoFile()
	container = new(Container)
}

func (p *Plugin) Remove() {
	p.FieldName = ""
}
func (p *Plugin) Swap(swp proto.Message) {
	p.swapStruct = swp
}

func (p *DocController) scanProtoBuf(path string) {

	dir, _ := filepath.Abs(path)
	fileList, _ := ioutil.ReadDir(dir)

	for _, v := range fileList {

		if v.IsDir() {
			p.scanProtoBuf(dir + backSlash + v.Name())
		} else {
			p.getProtoBuf(dir, v.Name())
		}
	}
}
func (p *DocController) getControllers(path string) {
	fp, _ := filepath.Abs(path)

	dir, _ := filepath.Abs(fp)
	fileList, _ := ioutil.ReadDir(dir)
	for _, v := range fileList {
		if v.IsDir() {
			p.getControllers(dir + backSlash + v.Name())
		} else {
			p.analyzeStruct(fp, v.Name(), false)
			p.dealGoFile(fp, v.Name(), &Annotation{})
		}
	}
}
func (p *DocController) getProtoBuf(sp, rp string) {
	if Match(rp, proFileRxp) {
		if ignoreFile[rp] != "" {
			return
		}
		protoFile[rp] = rp
	}
	if !Match(rp, proGoFileRxp) {
		return
	}
	p.analyzeStruct(sp, rp, true)

}
func (p *DocController) analyzeStruct(sp, rp string, flt bool) {
	var (
		pfc          ProtoFileContainer
		strcList     [] string
		pkn          = ""
		readPath     = sp + separator + rp
		isSelfPackage = false
	)
	if subLst(sp, separator) == beego.AppConfig.String("cool.generatePath") {
		isSelfPackage = true
	}
	p.readFile(readPath, func(cnt string) {

		if Match(cnt, goPackageRxp) && !isSelfPackage{
			pkn = trim(trimSpace(cnt), goFilePackage)
		}
		if Match(cnt, structRxp) {
			var strn = getStructName(cnt)
			if flt {
				for _, ant := range container.CoolDocs {
					if ant.ProtoBufControl == strn {
						strcList = append(strcList, strn)
						imports[sp] = sp
					}
					if ant.RespProtoBufCtrl == strn {
						strcList = append(strcList, strn)
						imports[sp] = sp
					}
				}
			} else {
				strcList = append(strcList, strn)
				imports[sp] = sp
			}
		}
	})
	if len(strcList) == 0 {
		return
	}
	pfc.StructList = strcList

	ct := protoContainer[pkn]
	if ct == nil {
		ct = []ProtoFileContainer{}
		ct = append(ct, pfc)
	} else {
		ct = append(ct, pfc)
	}
	protoContainer[pkn] = ct
}
func (p *DocController) readFile(fp string, fc func(cnt string)) {
	f, err := os.Open(fp)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer f.Close()

	br := bufio.NewReader(f)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		if fc != nil {
			fc(string(a))
		}
	}
}
func (p *DocController) generateProtoBufGoFile() error {
	var (
		err    error
		gnDir  = defaultGeneratePath
		goFile = defaultGeneratePath + generateGoFile
		gnPk   = generatePackage
	)
	if beego.AppConfig.String("cool.generatePath") != "" {
		gnDir = filterPath(beego.AppConfig.String("cool.generatePath"))
		goFile = gnDir + generateGoFile
		gnPk = goFilePackage + " " + getGeneratePk(beego.AppConfig.String("cool.generatePath"))
	}
	if !isExist(gnDir) {
		err = os.Mkdir(gnDir, os.ModePerm)
	}

	if err != nil {
		return err
	}
	fs, err1 := os.Create(goFile)

	defer fs.Close()

	if err1 == nil {
		w := bufio.NewWriter(fs)
		writeLine(w, gnPk)

		write(w, importStr+oneQuoMark+coolPackage+oneQuoMark)

		for _, v := range imports {
			if v == "" {
				continue
			}
			if subLst(v, separator) == gnDir {
				continue
			}
			imp := importStr + oneQuoMark + getImportPath(v) + oneQuoMark
			write(w, imp)

		}
		write(w, "\n")

		write(w, initMethod)
		write(w, coolContainer)

		for key, vp := range protoContainer {
			for _, v := range vp {
				p.generateStruct(w, v, key)
			}
		}
		write(w, resolverInit)
		write(w, coolStartMethod)
		writeLine(w, rightBrace)

		writeLine(w, resolverStruct)
		write(w, invokeMethod)
		write(w, rightBrace)

		w.Flush()
	}
	return nil
}

var temps = make(map[string]string)

func (p *DocController) generateStruct(w io.Writer, ctn ProtoFileContainer, pk string) error {
	sts := ctn.StructList
	for _, value := range sts {
		var sb StringBuilder
		if temps[value] == "" {
			sb.Append("    container").Append("[").Append(oneQuoMark)
			sb.Append(value).Append(oneQuoMark).Append("]").Append("= ").Append("&")
			if isNotNull(pk){
				sb.Append(pk).Append(".")
			}
			sb.Append(value).Append(leftBrace).Append(rightBrace)
			write(w, sb.String())
		}
		temps[value] = value
	}
	return nil
}
func (p *DocController) dealGoFile(fp, fln string, ant *Annotation) {
	var (
		ctrl       string
		hasGetCtrl = false
		sfd        = &structField{}
		ext        Extension
	)
	if structs["Extension"] != nil {
		ext = structs["Extension"].(Extension)
		sfd.ext = ext
	}
	readFunc := func(cnt string) {
		if Match(cnt, structRxp) && !hasGetCtrl {
			ctrl = getStructName(cnt)
			hasGetCtrl = true
		}
		if hasGetCtrl {
			ant.Controller = ctrl
		}

		if Match(cnt, funcRxp) && Match(cnt, ptrFuncRxp+ctrl+methodEndCloseRxp) {
			sfd = &structField{}
			sfd.ext = ext
			ant = &Annotation{}
		}
		p.getAnnotation(cnt, ant, sfd)
		rs := isNotNull(ant.Id, ant.Url, ant.Method, ant.ProtoBufFileName)
		if rs && (ant.Body != nil || ant.ProtoBufControl != "") {
			container.CoolDocs = append(container.CoolDocs, ant)
		}
	}
	p.readFile(fp+separator+fln, readFunc)
}

func (p *DocController) getAnnotation(cnt string, ant *Annotation, sfd *structField) *Annotation {
	if Match(cnt, reqDataRxp) {
		req := structs[trstr2spas(cnt, reqDataRxp)]
		if req != nil {
			ant.ProtoBufControl = reflect.TypeOf(req).Elem().Name()
			t, value := p.newInstance(req)
			sfd.reqType = t
			sfd.reqInst = value
		} else {
			ant.ProtoBufControl = trstr2spas(cnt, reqDataRxp)
		}
	}
	if Match(cnt, rspRxp) {
		rsp := structs[trstr2spas(cnt, rspRxp)]
		if rsp != nil {
			t, value := p.newInstance(rsp)
			ant.RespProtoBufCtrl = reflect.TypeOf(rsp).Elem().Name()
			sfd.rspType = t
			sfd.rspInst = value
		} else {
			ant.RespProtoBufCtrl = trstr2spas(cnt, rspRxp)
		}
	}
	if Match(cnt, reqFileRxp) {
		ant.ProtoBufFileName = trstr2spas(cnt, reqFileRxp)
	}
	if Match(cnt, rspFileRxp) {
		ant.RespProtoBufCtrl = trstr2spas(cnt, rspFileRxp)
	}

	if Match(cnt, methodRxp) {
		ant.Method = trstr2spas(cnt, methodRxp)
	}
	if Match(cnt, coolUrlRxp) {
		url := trstr2spas(cnt, coolUrlRxp)
		ant.Url = url
		ant.Id = strings.Replace(url, "/", "", -1)
		sfd.url = url
	}
	if sfd.reqType != nil && sfd.url != "" {
		rc := make(map[string]interface{})
		sfd.container = rc
		p.getFields(sfd)
		ant.Body = rc
	}

	/*if sfd.rspType != nil {
		rspMap := make(map[string]interface{})
		p.getFields(sfd)
		ant.RespData = rspMap
	}*/

	return ant
}
func (p *DocController) newInstance(itf interface{}) (reflect.Type, reflect.Value) {
	v := reflect.ValueOf(itf).Elem()
	return v.Type(), v
}
func (p *DocController) getFields(sfd *structField) map[string]interface{} {
	var (
		inst = sfd.reqInst
		rt   = sfd.reqType
		c    = sfd.container
		ext  = sfd.ext
		url  = sfd.url
	)

	for i := 0; i < inst.NumField(); i++ {
		field := inst.Field(i)

		if "XXX_unrecognized" == rt.Field(i).Name {
			continue
		}
		if "XXX_NoUnkeyedLiteral" == rt.Field(i).Name {
			continue
		}
		if "XXX_sizecache" == rt.Field(i).Name {
			continue
		}
		c[rt.Field(i).Name] = field.Type().Elem().String()

		plugin := &Plugin{
			StructName: rt.Name(),
			FieldName:  rt.Field(i).Name,
			Url:        url,
		}
		ext.Invoke(plugin)
		if plugin.swapStruct != nil {
			st, sv := p.newInstance(plugin.swapStruct)
			newStf := getStructField(url, ext, st, sv)
			//c[rt.Field(i).Name] = p.recursion(plugin.swapStruct, ext, url)
			c[rt.Field(i).Name] = p.getFields(newStf)
			continue
		}
		if field.Kind() == reflect.Slice {
			slp := sliceType(field.Type().String())
			if structs[slp] == nil {
				continue
			}
			t, v := p.newInstance(structs[slp])
			newStf := getStructField(url, ext, t, v)
			c[rt.Field(i).Name] = p.getFields(newStf)
			continue
		}

		if _, ok := field.Interface().(proto.Message); ok {
			tv := field.Type().Elem()
			newStf := getStructField(url, ext, tv, reflect.New(tv).Elem())
			c[rt.Field(i).Name] = p.getFields(newStf)
		}
	}
	return c
}
func (p *DocController) recursion(cls interface{}, ext Extension, url string) map[string]interface{} {
	dm := make(map[string]interface{})
	t, v := p.newInstance(cls)
	if _, ok := cls.(proto.Message); ok {
		newStf := getStructField(url, ext, t, v)
		p.getFields(newStf)
	}
	return dm
}
func getStructField(url string, ext Extension, t reflect.Type, value reflect.Value) *structField {
	dm := make(map[string]interface{})
	newStf := &structField{
		url:       url,
		ext:       ext,
		reqInst:   value,
		reqType:   t,
		container: dm,
	}
	return newStf
}
func (p *DocController) CreateDocument() {
	docs := container.CoolDocs

	var dataRow = make(map[string]interface{})
	dataRow["docName"] = container.DocName
	dataRow["docPkPath"] = container.PkPath
	dataRow["docKey"] = container.Key
	dataRow["rows"] = docs
	dataRow["protoFile"] = protoFile
	dataRow["header"] = container.RequestHeader

	log.Println(dataRow)
	p.Data["json"] = dataRow
	Access(p.Ctx)
	p.ServeJSON()
}
func (p *DocController) createRouter() {
	container.DocCtrl = p
	container.DocName = beego.AppConfig.String("cool.docName")
	container.PkPath = beego.AppConfig.String("cool.protoPackage")
	container.Key = beego.AppConfig.String("cool.Key")

	docs := container.CoolDocs
	for _, value := range docs {
		log.Println(value.Url)
		beego.Router(value.Url, structs[value.Controller].(beego.ControllerInterface), value.Method)
	}
	beego.Router("/coolDocument/createDocument", p, "*:CreateDocument")
}

func Start(c map[string]interface{}) {
	structs = c
	var doc DocController
	doc.getControllers("./controllers/")
	log.Println("doc-cool starting ...... ")
	GetIgnoreFile(beego.AppConfig.String("cool.ignoreFile"))
	doc.scanProtoBuf(getProPath(beego.AppConfig.String("cool.protoPath")))
	doc.createRouter()
}
