
## 使用

`proxy`: 支持 `socks5` `http`

下载：有的网站会限流，会导致失败；失败的图片不会保存；重新下载即可（多次下载的时候遇到磁盘已有文件就跳过）

支持两个网站
- https://www.tuwenhanman.com
- https://bingmh.com

界面截图
![example](example.png)



## 编译

`go build -ldflags="-H windowsgui -X main.Version=v1.0" .`

## syso的生成
- go-winres
```
go-winres  init
go-winres  make
```
- rsrc
```
rsrc -manifest main.manifest -ico .\winres\icon16.ico
```