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



Option<Vec<U256>>

require: true
items: u256

Vec<Option<U256>>

items: {
    type: u256
    nullable: true
}
