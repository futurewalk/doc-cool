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
  //@Description 支持bytes数据类型方式1
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
   `4、//@Method *:DocCool：controller层对应的请求方法(这里完全支持beego)`  
   `5、//@Url /v1/example/DocCool：请求的url`
   `6、//@Description :描述`   
  
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

    `git clone git@github.com:futurewalk/static.git`

    `注:不建议修改static文件夹名，这样会造成未知错误`
    
+ 下载完成放在项目根路径，下载下来的前端ui结构如下
     
     ```
     |-- static
        |-- dist
           |--css
           |--fonts
           |--js
        |-- extend
           |-- extend.js //扩展js
        |-- proto //存放proto文件包
           |-- example.proto 
        |-- favicon.ico
        |-- index.html
     ```    
+ 把项目中用到的proto文件放在static/proto包下
    
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
         
2、 bytes数据类型支持（由于bytes数据类型比较特殊，所以需要对doc-cool进行扩展，框架提供了扩展方式,这里对三层的包装类型进行扩展，注:只支持到三层） 
     
+ 如以下的proto message结构为
  
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
+ controller对应的注解
  
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

+ 后台对bytes支持进行扩展
    
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
            //这里的意思是，把BytesMessageReq结构体下的Data字段替换成BytesMessageInfo结构体
            plugin.Swap(&protobuf.BytesMessageInfo{})//此处一定要符合某个条件才调用，否则会无限递归
       }
  }
   ```
+ 前端对bytes进行扩展
 
  上边我们已经对bytes数据类型的属性进行了数据替换，前端也需要我们进行数据替换，以免我们在包装的时候找不到对应的bytes数据类型
  从我们下载的前端的目录结构，我们看到有一个extend文件夹，下面有一个extend.js，打开extend.js，代码如下
  
  ```javascript
  var Resolver = {}
  
  //处理bytes数据类型的回调
  Resolver.byteInvoke = function (plugin) {
      if (plugin.StructName == "BytesMessageReq" && plugin.FieldName == "Data" && plugin.Url == '/v1/example/BytesSupport') {
          //这里的意思是，把BytesMessageReq结构体下的Data字段替换成BytesMessageInfo结构体
          plugin.Swap("BytesMessageInfo");
      }
  }
  // 处理处理message包装的数据和当前message不在同一个proto文件的回调(很少用到，这里也建议message包装的数据类型和当前message一个proto文件)
  Resolver.dataInvoke = function (plugin){
  
  }
  //处理请求数据返回之后需要处理返回message对应的bytes字段的回调，使用方式和Resolver.byteInvoke，但是这两个出发的条件不一样。
  Resolver.RspInvoke = function (plugin) {
      if (plugin.StructName == "BytesMessageRsp" && plugin.FieldName == "Data" && plugin.Url == '/v1/example/BytesSupport') {
          plugin.Swap("BytesMessageInfo");
      }
  }
  window.Resolver = Resolver;
  ```
  扩展代码写完，启动项目，相信可以看到生成的对应的Body数据格式和生成的对应的message输入框，如果你在扩展之前就已经启动过项目看了一下效果，
  你会发现，先前对应的bytes字段生成的输入框只有一个，扩展代码写完之后，已经生成bytes替换成对应的结构体的输入框

  如果上面的文档你觉得写得不够明了，没关系，我们有example，如果example和文档都看的还不是很明了，那么我和你其中一个是不适合写代码的。
  如有问题，请联系:wangdequan2829@qq.com。

3、你可能会遇到的几个问题
  
  + 项目启动之后无法找到proto文件
    
    1、cool.protoPackage = protobuf 配置和proto文件的package不一致
   
    2、用到的proto文件没有拷贝的static/proto包下
  
  + 项目启动成功，访问不到，报404
  
    1、在项目启动时没有初始化doc-cool ```go _"example/protobuf/cool"//初始化生成的doc_cool_register.go文件```  
  
  + <font color = "red" size = 6>如果你的问题不是上面其中的几个，那么请提issue吧!</font>
       