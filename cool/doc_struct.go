package cool

import (
    "io/ioutil"
)

type Architecture interface {
    CreateRouter()
    AnalyzeStructure(sp, rp string, flt bool) error
    GetControlStructure(path string) error
    ScanProtoStructure(path string)
    GetProtoStructure(sp, rp string) error
}
type Structures struct {
    Ant
    Package            string
    ImportList         map[string]string
    StructureContainer map[string][]StructureContainer
    IgnoreFile         map[string]string
    ProtoFile          map[string]string
}
type StructureContainer struct {
    Package    string
    StructList []string
}

func (s *Structures) CreateRouter() {

}
func (s *Structures) ScanProtoStructure(path string) {

    fileList, _ := ioutil.ReadDir(path)

    if s.IgnoreFile == nil {
        s.IgnoreFile = make(map[string]string)
    }
    if s.ProtoFile == nil {
        s.ProtoFile = make(map[string]string)
    }

    for _, v := range fileList {
        if v.IsDir() {
            s.ScanProtoStructure(path + separator + v.Name())
        } else {
            err := s.GetProtoStructure(path, v.Name())

            if err != nil {
                panic(err)
            }
        }
    }
}
func (s *Structures) GetProtoStructure(sp, rp string) error {
    if Match(rp, proFileRxp) {
        if s.IgnoreFile[rp] != "" {
            return nil
        }
        s.ProtoFile[rp] = rp
    }
    if !Match(rp, proGoFileRxp) {
        return nil
    }
    return s.AnalyzeStructure(sp, rp, true)

}
func (s *Structures) AnalyzeStructure(sp, rp string, flt bool) error {
    var (
        stc           StructureContainer
        structureList [] string
        pkn           = ""
        readPath      = sp + separator + rp
        isSelfPackage = false
    )

    if subLst(sp, separator) == container.config.Get("cool.generatePath") {
        isSelfPackage = true
    }
    err := s.Walk(readPath, func(cnt string) {

        if Match(cnt, goPackageRxp) && !isSelfPackage {
            pkn = trim(trimSpace(cnt), goFilePackage)
        }
        if Match(cnt, structRxp) {
            structName := getStructName(cnt)
            if flt {
                for _, ant := range container.ant.Ants {
                    if ant.ProtoBufControl == structName {
                        structureList = append(structureList, structName)
                        s.ImportList[sp] = sp
                    }
                    if ant.RespProtoBufCtrl == structName {
                        structureList = append(structureList, structName)
                        s.ImportList[sp] = sp
                    }
                }
            } else {
                structureList = append(structureList, structName)
                s.ImportList[sp] = sp
            }
        }
    })
    if err != nil {
        return err
    }
    if len(structureList) >= 0 {
        stc.StructList = structureList

        ct := s.StructureContainer[pkn]
        if ct == nil {
            ct = []StructureContainer{}
            ct = append(ct, stc)
        } else {
            ct = append(ct, stc)
        }
        s.StructureContainer[pkn] = ct
    }
    return nil
}
func (s *Structures) GetControlStructure(path string) error {
    fileList, err := ioutil.ReadDir(path)

    if err != nil {
        return err
    }

    for _, v := range fileList {
        if v.IsDir() {
            s.GetControlStructure(path + separator + v.Name())
        } else {
            err := s.AnalyzeStructure(path, v.Name(), false)

            if err != nil {
                panic(err)
            }

            err = s.AnalyzeController(path, v.Name(), &Annotation{})

            if err != nil {
                panic(err)
            }

        }
    }
    return nil
}
