1. params 为空应如何描述？
2. cfx_getLogs 返回结果应为一个数组
3. 没有 type: array
4. 更新最新的 rust 代码
5. cfx_getLogs filter 参数的处理
6. epochNumber schema 定义
7. 枚举值应该为 下划线分隔
8. 字段名应该驼峰形式
9. RewardInfo schema 定义有问题
10. transaction 没有 space 字段
11. TransactionStatus schema 定义

TODO:
1. VariadicValue Schema生成 （rpc logfilter topics）
2. struct enum 根据是否实现了 Serialize 生成不同格式
3. one of 多一个空schema (BlockHashOrEpochNumber) V
4. 默认枚举camelcase
5. param nullable 问题
6. block.transaction type array
7. 去掉trait中的前后注释

Option<Vec<U256>>

require: true
items: u256

Vec<Option<U256>>

items: {
    type: u256
    nullable: true
}
