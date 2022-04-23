## 特点

局域网手机与电脑互传文件

- 电脑上传文字/文件后，生成二维码，手机可以扫码下载
- 手机端上传文字/文件后利用 websocket 通知电脑，弹出下载提示

## 技术选择

窗口利用 lorca - 利用本地的 Chrome 等浏览器作为 webview

- 后端：go，web 框架是 gin
- 前端：react

## 基本 api

```go
router.POST("/api/v1/files", controller.FilesController)
router.GET("/api/v1/qrcodes", controller.QrcodesController)
router.GET("/uploads/:path", controller.UploadsController)
router.GET("/api/v1/addresses", controller.AddressesController)
router.POST("/api/v1/texts", controller.TextsController)
```

## 使用

```
git clone git@github.com:physicshi/go_project.git
cd server/frontend/
npm i
npm run build
cd ...
go build .
```
