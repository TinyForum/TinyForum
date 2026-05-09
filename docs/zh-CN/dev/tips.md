# 提示

## 重新生成 API 文档

```bash
cd backend/docs
npx swagger-markdown -i ./swagger.json
```





```bash
internal/service/
├── upload/                   # 底层上传引擎
│   ├── engine.go             # UploadService 接口 + 实现（依赖 StorageDriver, Registry, AttachmentRepo）
│   └── dto.go                # UploadRequest, UploadResponse
├── attachment/               # 附件业务
│   ├── service.go            # AttachmentService（依赖 upload.Engine）
│   ├── upload.go             # UploadFile 实现（调用 engine.Upload）
│   └── query.go              # GetUserFiles, GetFile, DeleteFile
└── plugin/                   # 插件业务
    ├── service.go            # PluginService（依赖 storage, pluginRepo）
    ├── install.go            # Install 方法
    └── query.go              # ListPlugins, ListUserPlugins
```



请求绑定



```bash
 List(ctx context.Context, queryBO *common.PageQuery[do.PluginMeta]) ([]*do.PluginMeta, int64, error)
```





```bash
	var req request.ListPluginsRequest
	if err := req.Bind(c); err != nil {
		response.HandleError(c, err)
		return
	}
```

