# goframe基础开发框架

开发框架基于GIN框架，数据库操作依赖 GORM 开源包。

GIN框架文档地址：
[https://gin-gonic.com/zh-cn/docs/introduction/](https://gin-gonic.com/zh-cn/docs/introduction/)

GORM文档地址：
[https://gorm.io/zh\_CN/docs/index.html](https://gorm.io/zh_CN/docs/index.html)

***

## 框架目录结构说明


	utils --- 自定义功能组件

	middleware --- 中间件

	funcs --- 自定义全局方法

	pconstants --- 常量定义

	structs --- 公用结构体


***

## 框架功能开发进度：

* [x] 1\.  日志记录，基于zap实现
* [x] 2\.  Crash处理，系统运行Panic异常告警，目前仅邮件形式告警，可以扩展短信等等方式
* [x] 3\.  接口验签
* [x] 4\.  接口限流，引用第三方包[https://github.com/juju/ratelimit](https://github.com/juju/ratelimit)实现
* [x] 5\.  CORS跨域处理
* [x] 6\.  邮件发送
* [x] 7\.  文件上传，支持上传到本地目录，腾讯云COS，可扩展其它
* [x] 8\.  短信能力，目前仅接入腾讯云短信，可扩展其它
* [x] 9\.  Timer定时器，引用第三方包[https://github.com/gogf/gf/v2/os/gtimer](https://github.com/gogf/gf/v2/os/gtimer)实现
* [x] 10\.  Cron定时任务，引用第三方包[https://github.com/gogf/gf/v2/os/gcron](https://github.com/gogf/gf/v2/os/gcron)实现
* [x] 11\.  Redis常用操作能力封装，基于第三方包[github.com/gomodule/redigo/redis]([https://](https://pkg.go.dev/)github.com/gomodule/redigo/redis)实现: string,hash,list,set,zset,expire,scan,geo
* [x] 12\.  Excel文件导入导出,基与第三方包[github.com/xuri/excelize/v2](https://pkg.go.dev/github.com/xuri/excelize/v2)实现
* [x] 13\.  生成二维码&图片合成工具
* [x] 14\.  mysql数据库操作方法封装
* [x] 15\.  snowflake 分布式唯一ID生成工具
* [x] 16\.  随机数生成工具
* [ ] 17\.  微信小程序用户授权登录机制
* [ ] 18\.  I18N国际化
* [ ] 更多功能持续迭代。。。