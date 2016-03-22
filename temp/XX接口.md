请求方式和URL路径：

```
GET /fdpp/p2p
```

请求参数：

```
page              // 第几页，页码从1开始
page_size         // 一页几条数据
select            // 需要获取的字段，以英文逗号分隔（不传时获取所有字段），比如只获取name和rate字段：select=name,rate
sort              // 以哪个字段排序，和 order_by 配合使用
order_by          // 和 sort 配合使用，比如按name升序排列：sort=name&order_by=asc
(数据字段过滤)      // 比如 name=光大传媒&rate=3.
                  // 可以指定取值范围，比如 initial={gte}10{lt}100 表示 initial 大于等于10并且小于100
                  //                    initial={gt}100 表示大于100
                  //                    initial={lte}10 表示小于等于10
```

成功时返回：

```
{
    data: [
        { // 数据内容
            xxx: ???,
            ...
        },
        ...
    ],
    pagination: {
        count:               // 数据总共多少条
        current:             // 当前是第几页
        page_size:           // 一页几条
    }
}
```

P2P产品所有字段：

```
{
    name:           String, // 名称  
    rate:           Number, // 年化收益率  
    max_rate:       Number, // 最高收益
    platform:       String, // 平台
    initial:        Number, // 起投金额
    horizon:        String, // 期限  
    horizon_max:    String, // 期限
    total:          Number, // 总销售额  
    value_date:     Date,   // 起息日
    payment_method: String, // 还款方式
    guarantee:      String, // 担保机构
    url:            String  // URL
}
```
