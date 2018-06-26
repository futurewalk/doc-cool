package cool

import (
    "reflect"
    "strings"
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
type structField struct {
    reqType   reflect.Type
    reqInst   reflect.Value
    rspType   reflect.Type
    rspInst   reflect.Value
    container map[string]interface{}
    ext       Extension
    url       string
}

type Ants map[string]*Annotation

type Ant struct {
    Ants
    //extend baser ,it will have all baser interface api
    Base
}

type Annotator interface {
    //extend baser ,it will have all baser interface api
    Baser

    //analyze controller
    //fp:controller file path
    //fn:controller fileName
    AnalyzeController(fp, fn string, ant *Annotation) error

    //analyze annotation in controller
    AnalyzeAnnotation(cnt string, ant *Annotation, sfd *structField) *Annotation
}

func (at *Ant) AnalyzeController(fp, fn string, ant *Annotation) error {
    var (
        ctrl       string
        hasGetCtrl = false
        sfd        = &structField{}
        ext        Extension
    )
    if container.registerStructures["Extension"] != nil {
        ext = container.registerStructures["Extension"].(Extension)
        sfd.ext = ext
    }

    if at.Ants == nil {
        at.Ants = make(map[string]*Annotation)
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
        at.AnalyzeAnnotation(cnt, ant, sfd)
        rs := isNotNull(ant.Id, ant.Url, ant.Method, ant.ProtoBufFileName)
        if rs && (ant.Body != nil || ant.ProtoBufControl != "") {
            at.Ants[ant.Id] = ant
        }
    }
    container.ant = *at
    return at.Walk(fp+separator+fn, readFunc)
}
func (at *Ant) AnalyzeAnnotation(cnt string, ant *Annotation, sfd *structField) *Annotation {
    if Match(cnt, reqDataRxp) {
        req := container.registerStructures[trstr2spas(cnt, reqDataRxp)]
        if req != nil {
            ant.ProtoBufControl = reflect.TypeOf(req).Elem().Name()
            t, value := at.NewInstance(req)
            sfd.reqType = t
            sfd.reqInst = value
        } else {
            ant.ProtoBufControl = trstr2spas(cnt, reqDataRxp)
        }
    }
    if Match(cnt, rspRxp) {
        rsp := container.registerStructures[trstr2spas(cnt, rspRxp)]
        if rsp != nil {
            t, value := at.NewInstance(rsp)
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
    if Match(cnt, descriptionRxp) {
        desc := trstr2spas(cnt, descriptionRxp)
        ant.Description = desc
    }
    if sfd.reqType != nil && sfd.url != "" {
        rc := make(map[string]interface{})
        sfd.container = rc
        at.GetFields(sfd)
        ant.Body = rc
    }

    /*if sfd.rspType != nil {
        rspMap := make(map[string]interface{})
        p.getFields(sfd)
        ant.RespData = rspMap
    }*/

    return ant
}
