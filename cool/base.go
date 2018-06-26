package cool

import (
    "os"
    "bufio"
    "io"
    "reflect"
    "github.com/golang/protobuf/proto"
    "io/ioutil"
)

type Baser interface {
    //read file
    Walk(p string, fc func(c string)) error
    //get struct field
    GetFields(sfd *structField) map[string]interface{}
    //new struct
    NewInstance(object interface{}) (reflect.Type, reflect.Value)
    //generate struct, include controller struct and proto struct
    GenerateStructure(w io.Writer, ptc StructureContainer, pk string) error
    //generate cool register file
    GenerateRegisterGoFile() error//
}
type Base struct {}

func (b *Base) Walk(fp string, fc func(cnt string)) error {
    f, err := os.Open(fp)
    if err != nil {
        return err
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
    return nil
}
func (b *Base) WalkAll(fp string, fc func(cnt string)) error {
    f, err := os.Open(fp)
    if err != nil {
        return err
    }
    defer f.Close()

    content,err := ioutil.ReadAll(f)
    fc(string(content[:]))
    return err
}
func (b *Base) GetFields(sfd *structField) map[string]interface{} {
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
        c[rt.Field(i).Name] = field.Type().Elem().String() + "," + getIndexStr(rt.Field(i).Tag.Get("protobuf"), ",", 2)
        plugin := &Plugin{
            StructureName: rt.Name(),
            FieldName:     rt.Field(i).Name,
            Url:           url,
        }
        ext.Invoke(plugin)
        if plugin.swapStructure != nil {
            st, sv := b.NewInstance(plugin.swapStructure)
            newStf := b.GetStructField(url, ext, st, sv)
            c[rt.Field(i).Name] = b.GetFields(newStf)
            continue
        }
        if field.Kind() == reflect.Slice {
            slp := sliceType(field.Type().String())
            if container.structures.StructureContainer[slp] == nil {
                continue
            }
            t, v := b.NewInstance(container.structures.StructureContainer[slp])
            newStf := b.GetStructField(url, ext, t, v)
            c[rt.Field(i).Name] = b.GetFields(newStf)
            continue
        }
        if _, ok := field.Interface().(proto.Message); ok {
            tv := field.Type().Elem()
            newStf := b.GetStructField(url, ext, tv, reflect.New(tv).Elem())
            c[rt.Field(i).Name] = b.GetFields(newStf)
        }
    }
    return c
}
func (b *Base) GetStructField(url string, ext Extension, t reflect.Type, value reflect.Value) *structField {
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
func (b *Base) NewInstance(itf interface{}) (reflect.Type, reflect.Value) {
    v := reflect.ValueOf(itf).Elem()
    return v.Type(), v
}
func (b *Base) GenerateRegisterGoFile() error {
    var (
        err    error
        gnDir  = defaultGeneratePath
        goFile = defaultGeneratePath + generateGoFile
        gnPk   = generatePackage
    )

    if container.config.Get("cool.generatePath") != "" {
        gnDir = getPath(separator + filterPath(container.config.Get("cool.generatePath")))
        goFile = gnDir + separator + generateGoFile
        gnPk = goFilePackage + " " + getGeneratePk(container.config.Get("cool.generatePath"))
    }
    if !isExist(gnDir) {
        err = os.Mkdir(gnDir, os.ModePerm)
    }

    if err != nil {
        return err
    }
    fs, err1 := os.Create(goFile)

    if err1 != nil {
        return err1
    }

    defer fs.Close()

    w := bufio.NewWriter(fs)
    writeLine(w, gnPk)

    write(w, importStr+oneQuoMark+coolPackage+oneQuoMark)

    for _, v := range container.structures.ImportList {
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

    writeLine(w, coolStartMethod)
    writeLine(w, "        &Resolver{},")
    for key, vp := range container.structures.StructureContainer {
        for _, v := range vp {
            b.GenerateStructure(w, v, key)
        }
    }
    writeLine(w, "    )")
    writeLine(w, rightBrace)
    writeLine(w, resolverStruct)
    write(w, invokeMethod)
    write(w, rightBrace)

    w.Flush()

    return nil
}

var temps = make(map[string]string)

func (b *Base) GenerateStructure(w io.Writer, ctn StructureContainer, pk string) error {
    sts := ctn.StructList
    for _, value := range sts {
        var sb StringBuilder
        if temps[value] == "" {
            sb.Append("        &")
            if isNotNull(pk) {
                sb.Append(pk).Append(".")
            }
            sb.Append(value).Append(leftBrace).Append(rightBrace).Append(",")
            writeLine(w, sb.String())
        }
        temps[value] = value
    }
    return nil
}
