# goframe基础开发SDK

---

## 框架目录结构说明

    utils --- 自定义功能组件

    middleware --- 中间件

    funcs --- 自定义全局方法

---

## SDK功能开发进度：

* [X] 1\.  plog日志记录，基于zap实现
* [X] 2\.  Crash处理，系统运行Panic异常告警，目前仅邮件形式告警，可以扩展短信等等方式
* [X] 3\.  接口验签
* [X] 4\.  接口限流，引用第三方包[https://github.com/juju/ratelimit](https://github.com/juju/ratelimit)实现
* [X] 5\.  CORS跨域处理
* [X] 6\.  邮件发送
* [X] 7\.  pupload文件上传组件，支持上传到本地目录，腾讯云COS，可扩展其它
* [X] 8\.  短信能力，目前仅接入腾讯云短信，可扩展其它
* [X] 9\.  Timer定时器，引用第三方包[https://github.com/gogf/gf/v2/os/gtimer](https://github.com/gogf/gf/v2/os/gtimer)实现
* [X] 10\.  Cron定时任务，引用第三方包[https://github.com/gogf/gf/v2/os/gcron](https://github.com/gogf/gf/v2/os/gcron)实现
* [X] 11\.  Redis常用操作能力封装，基于第三方包[github.com/gomodule/redigo/redis]([https://](https://pkg.go.dev/)github.com/gomodule/redigo/redis)实现: string,hash,list,set,zset,expire,scan,geo,bit,transaction,HyperLogLog
* [X] 12\.  Excel文件导入导出,基与第三方包[github.com/xuri/excelize/v2](https://pkg.go.dev/github.com/xuri/excelize/v2)实现
* [X] 13\.  pgraphic生成二维码&图片合成工具
* [X] 14\.  mysql数据库操作方法封装
* [X] 15\.  psnowflake 分布式唯一ID生成工具
* [X] 16\.  prand随机数生成工具
* [X] 17\.  perrors全局错误处理
* [X] 18\.  pelastic组件，当前仅实现日志文档上报，后续可扩展其他能力
* [X] 19\.  TimeZone时区组件
* [X] 20\.  分页组件，包含普通offset偏移量分页，cursor游标分页
* [X] 21\.  pzip压缩，解压缩组件
* [X] 22\.  pfile文件处理组件
* [ ] 23\.  I18N国际化
* [ ] 更多功能持续迭代。。。
