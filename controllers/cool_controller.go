package controllers

import (
	"github.com/astaxie/beego"
	"doc-cool/protobuf"
	"github.com/gogo/protobuf/proto"
	"log"
)

type CoolController struct {
	beego.Controller
}

//@ReqData Request
//@RespData Response
//@ReqProtoFile base.proto
//@Method *:DeviceInfo
//@Url /v1/cool/testcool
func (c *CoolController) Test()  {
	reqData := &protobuf.Request{}

	body := c.Ctx.Input.RequestBody
	err := proto.Unmarshal(body,reqData)

	log.Println(reqData,err)
}
