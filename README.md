# doc-cool
1、基于beego支持protobuf的自动化文档管理工具

- 下载安装
  
  `go get -t github.com/futurewalk/doc-cool/cool`

- 使用注解
  
      //@ReqData ExampleMessageReq
      //@RespData ExampleReqMessageRsp
      //@ReqProtoFile example.proto
      //@Method *:DocCool
      //@Url /v1/example/DocCool
      func (p *ExampleController) DocCool() {
          reqData := &protobuf.ExampleMessageReq{}
          rspData := &protobuf.ExampleReqMessageRsp{}
        
          err := proto.Unmarshal(p.Ctx.Input.RequestBody,reqData)
        
          if err == nil{
              reqData.Data = rspData.Data
          }
          byte,err := proto.Marshal(rspData)
          p.Ctx.ResponseWriter.Write(byte)
      }
  > 说明:  
  > `1、//@ReqData SettingReq：请求的protobuf结构体名 `   
  > `2、//@RespData SettingResp：返回的protobuf结构体名`          
  > `3、//@ReqProtoFile device_setting.proto：请求protobuf结构体所在的protobuf文件`  
  > `4、//@Method *:GetPlanStep：controller层对应的请求方法(这里完成支持beego)`  
  > `5、//@Url /business/deviceset/getPlanStep：请求的url`   
- conf文件配置
  >`cool.protoPath = protobuf ` //protobuf的根目录   
  >`cool.docName = coolPad web API`  //API的名称   
  >`cool.protoPackage = protobuf.device`  //protobuf生成go文件的package  
  >`cool.generatePath = protobuf/cool`  //生成注册文件的路径  
  >`cool.ignoreFile = push.proto`//需要忽略扫描的protobuf文件  
- 生成结构体注册文件

      package main
  
      import (
  	    "testing"
  	    "doc-cool/cool"
      )
  
      func TestExport(t *testing.T)  {
       //调用此方法生成结构体注册文件，生成的路径是conf配置的ool.generatePath
       //默认生成路径是protobuf/cool
        cool.Export()
      }
- 初始化doc-cool

      package main
      //项目启动时初始化
      import (
          "github.com/astaxie/beego"
          _"example/protobuf/cool"//初始化doc-cool,只有初始化才能生成路由
      )
      
      func main() {
          //beego配置访问前端页面  
          beego.BConfig.WebConfig.DirectoryIndex = true
          beego.BConfig.WebConfig.StaticDir["/static"] = "static"
          beego.Run()
      }