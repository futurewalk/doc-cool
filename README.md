# doc-cool
1、基于beego支持protobuf的自动化文档管理工具

- 下载安装
  
  `go get -t github.com/futurewalk/doc-cool/cool`

- 使用注解(在controller层使用注解)

  ```go
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
  ```
  > 说明:  
   `1、//@ReqData ExampleMessageReq：请求的protobuf结构体名 `   
   `2、//@RespData ExampleReqMessageRsp：返回的protobuf结构体名`          
   `3、//@ReqProtoFile example.proto：请求protobuf结构体所在的protobuf文件`  
   `4、//@Method *:DocCool：controller层对应的请求方法(这里完成支持beego)`  
   `5、//@Url /v1/example/DocCool：请求的url`   
  
- conf文件配置
  >`cool.protoPath = protobuf ` //protobuf的根目录   
  >`cool.docName = coolPad web API`  //API的名称   
  >`cool.protoPackage = protobuf.device`  //protobuf生成go文件的package  
  >`cool.generatePath = protobuf/cool`  //生成注册文件的路径  
  >`cool.ignoreFile = push.proto`//需要忽略扫描的protobuf文件  
  
- 生成初始化入口和结构体注册文件
  ```go
  package main

  import (
    "testing"
    "doc-cool/cool"
  )

  func TestExport(t *testing.T)  {
   //调用此方法生成结构体注册文件，生成的路径是conf配置的ool.generatePath
   //默认生成路径是protobuf/cool,生成的文件名是doc_cool_register.go
    cool.Export()
  }
  ```
- 下载doc-cool前端ui放置在项目当中。
    `go clone git@github.com:futurewalk/static.git`
    
    + 下载完成放在项目根路径
    
    
- 初始化doc-cool并启动项目(项目启动成功通过http://host:port/static可进行访问)
     
  ```go
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
  ```
         
- bytes数据类型支持（由于bytes数据类型比较特殊，所以需要对doc-cool进行扩展，框架提供了扩展方式,这里对三层的包装类型进行扩展，注:只支持到三层） 
     
  1、 如一下的proto message结构为
  
  ```proto
  //请求的message
  message BytesMessageReq {
      required string docCoolName = 1;
      required bytes data = 2;//bytes包装的message为BytesMessageInfo
  }
  //返回的message
  message BytesMessageRsp {
      required string bytesName = 1;
      required bytes data = 2; //bytes包装的message为BytesMessageInfo
  }
  //包装的bytes message
  message BytesMessageInfo {
      required string respName = 1;
      required int64 respTime = 2;
      required NextMessage data = 3;
  }
  //bytes又包装的一层message
  message NextMessage {
      required string nextName = 1;
      required int64 nextTime = 2;
  }
  ```
  2、controller对应的注解
  
  ```go
      
  //@ReqData BytesMessageReq
  //@RespData BytesMessageRsp
  //@ReqProtoFile example.proto
  //@Method *:BytesSupport
  //@Url /v1/example/BytesSupport
  func (p *ExampleController) BytesSupport() {
    var b []byte
    reqData := &protobuf.BytesMessageReq{}
    rspData := &protobuf.BytesMessageRsp{
        Data: b,
    }
  
    err := proto.Unmarshal(p.Ctx.Input.RequestBody, reqData)
    if err == nil {
        data := reqData.Data
        rspData.Data = data
    }
    rspData.BytesName = reqData.DocCoolName
    byte, err := proto.Marshal(rspData)
    cool.Access(p.Ctx)
    p.Ctx.ResponseWriter.Write(byte)
  }
  ```

  3、对bytes支持进行扩展
    
  打开我们生成的doc_cool_register.go文件，滚动到最下边，可以看到如下代码
    
  ```go
  type Resolver struct{}

  //Plugin提供了三个属性，
  //1、StructName 结构体名称
  //2、FieldName  结构体的某个属性名
  //3、Url 当前属于哪个url
  //很多时候根据StructName和FieldName就可以判断了，但是有时候某几个方法会公用几个结构体，可想而知，这样判断显然不行，因此，这里还提供了一个url进行更加严谨的判断。
  func (p *Resolver) Invoke(plugin *cool.Plugin) {
       //此处进行扩展，扩展的方式也非常简单。从上边的BytesMessageReq可得知，
       //当属性是data的时候，我们希望替换成BytesMessageInfo，那么我们可以这样写
       //注:这里是根据message生成的结构体，并不是message,因此注意属性大小写问题，message的属性名是小写的，生成的结构体也是大写的，所以这里建议message用大写，以免混淆       
       
       if plugin.StructName == "BytesMessageReq" && plugin.FieldName == "Data" && plugin.Url == "/v1/example/BytesSupport"{
            plugin.Swap(&protobuf.BytesMessageInfo{})//此处一定要符合某个条件才调用，否则会无限递归
       }
  }
   ```


     