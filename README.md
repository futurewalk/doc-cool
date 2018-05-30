# doc-cool
1、基于beego支持protobuf的自动化文档管理工具

- 使用注解
  
      //@ReqData SettingReq
      //@RespData SettingResp
      //@ReqProtoFile device_setting.proto
      //@Method *:GetPlanStep
      //@Url /business/deviceset/getPlanStep
      func (this *DeviceSetting) GetPlanStep()  {
      
      }
  >> 说明:  
  >> 1、//@ReqData SettingReq：请求的protobuf结构体名  
  >> 2、//@RespData SettingResp：返回的protobuf结构体名       
  >> 3、//@ReqProtoFile device_setting.proto：请求protobuf结构体所在的protobuf文件  
  >> 4、//@Method *:GetPlanStep：controller层对应的请求方法(这里完成支持beego)  
  >> 5、//@Url /business/deviceset/getPlanStep：请求的url   
- conf文件配置
  >cool.protoPath = protobuf  
  >cool.docName = coolPad web API  
  >cool.protoPackage = protobuf.device  
  >cool.generatePath = protobuf/cool  
  >cool.ignoreFile = push.proto
- 生成注册文档
- 初始化doc-cool