# 对象模型



## 对象模型的定义

## 对象的访问

每个 DO 都需要嵌入 BaseModel，并通过如下方式访问

```go
func PluginBOToDO(b *bo.PluginMeta) *do.PluginMeta {
    if b == nil {
        return nil
    }
    return &do.PluginMeta{
       BaseModel: common.BaseModel{ // 通过嵌入类型名初始化
            ID:        b.ID,
            CreatedAt: b.CreatedAt,
            // UpdatedAt 和 DeletedAt 通常由数据库自动维护，不需要从 BO 传递
        },
        Name:      b.Name,
        AuthorID:  b.AuthorID,
        Tags:     b.Tags,
        Type:  b.Type,
        Status:    b.Status,
       
    }
}
```

