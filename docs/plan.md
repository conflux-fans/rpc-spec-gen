# 接口测试
1. conflux-rust -> 描述文档生成
2. 人工 check 描述文档
3. sdk 反序列化再序列化后 json <-比较-> conflux-rust json
4. 提供case测试到所有情况

## openrpc 文件夹结构

- doc template
- schemas
    - basetype
      - bool
      - string
    - cfx_types
      - U256
      - H64

# 逻辑测试