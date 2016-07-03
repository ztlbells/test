##数字货币简单商业应用案例
#### 功能描述
该智能合约实现一个简单的商业应用案例，即数字货币的发行与转账。在这之中一共分为三种角色：中央银行，商业银行，企业。其中中央银行可以发行一定数量的货币，企业之间可以进行相互的转账。主要实现如下的功能：
- 初始化中央银行及其发行的货币数量
- 新增商业银行，同时央行并向其发行一定数量的货币
- 新增企业
- 商业银行向企业转给一定数量的数字货币
- 企业之间进行相互的转账
- 查询企业、银行、交易信息

#### function及各自实现的功能：
- `init`  初始化中央银行，并发行一定数量的货币
- `invoke`   调用合约内部的函数
- `query`   查询相关的信息
- `createBank`   新增商业银行，同时央行向其发行一定数量的货币
- `createCompany`   新增企业
- `issueCoin` 央行再次发行一定数量的货币
- `issueCoinToCp`  商业银行向企业转一定数量的数字货币
- `transfer`   企业之间进行相互转账
- `getCompanys`   获取所有的公司信息
- `getBanks`    获取所有的商业银行信息
- `getTransactions` 获取所有的交易记录
- `getCompanyById`   获取某家公司信息
- `getBankById`   获取某家银行信息
- `getTransactionBy` 获取某笔交易记录


#### 数据结构设计
- centerBank 中央银行
  - Name 名称
  - TotalNumber 发行货币总数额
  - RestNumber 账户余额
- bank  商业银行
  - Name 名称
  - TotalNumber 收到货币总数额
  - RestNumber 账户余额
  - ID 银行id
- company 企业
  - Name 名称
  - Number  账户余额
  - ID 企业id
- transaction 交易内容
  - FromType 发送方角色
  - FromID 发送方ID
  - ToType  接收方角色
  - ToID 接收方ID
  - Time  交易时间
  - Number 交易数额
  - ID 交易ID
 
#### 接口设计
`createBank`:


request参数:
```

```

response参数:
```

```
#### 其它
对于查询请求，为了兼顾读写速度，将一些信息备份存放在非区块链数据库上也是一个较好的选择。
