# 极简版抖音
本项目由 秒速五字节 小队进行开发

## 运行说明
- 需要自行安装 ffmpeg 并配置环境变量
- 需要新建 `config/hide.go` 文件并填入以下内容
```go

package config​
​
const (​
    LoginHead    = "user:passwd@tcp(host:port)/schema_name" // 填入数据库连接​
    COSURL       = "https://xxxxx.cos.ap-guangzhou.myqcloud.com"​// 填入腾讯云 COS 地址
    COSSecretID  = ​
    COSSecretKey = 
    Path         = "./"// ffmpeg 项目目录​
)

```



## 抖音项目服务端简单介绍

具体功能内容参考飞书说明文档
[‌‍⁡‍‍⁣‍​【秒速五字节队】结业项目答辩汇报文档 - 飞书云文档](https://sygubn5jgz.feishu.cn/docx/QJYNdnc3soxIPJxDiqLc3iDXn2e)

工程无其他依赖，直接编译运行即可

```shell
go build && go run .
```


### 测试

test 目录下为不同场景的功能测试case，可用于验证功能实现正确性

其中 common.go 中的 _serverAddr_ 为服务部署的地址，默认为本机地址，可以根据实际情况修改

测试数据写在 demo_data.go 中，用于列表接口的 mock 测试