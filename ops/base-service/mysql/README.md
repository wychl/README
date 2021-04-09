# mysql

## 数据库和实例概念

- 数据库：物理操作文件系统或其他形式文件类型的集合；
- 实例：MySQL 数据库由后台线程以及一个共享内存区组成；


## 权限管理

- 查询帐户名  

  ​

  ```
  USE mysql;
  SELECT user FROM user;
  ```

  ​

- 创建账户 

  - `GRANT SELECT ,UPDATE ON *.* TO 'testUser'@'localhost'' identified BY 'testpwd'`
  - `CREATE USER myuser IDENTIFIED BY 'mypassword';`

- 修改账户名  `RENAME USER myuser TO newuser;`

- 删除账户  `DROP USER myuser;`

- 查看权限 `SHOW GRANTS FOR myuser;`

- 授予权限 `GRANT SELECT, INSERT ON mydatabase.* TO myuser;`

  - 整个服务器，使用 GRANT ALL 和 REVOKE ALL；
  - 整个数据库，使用 ON database.*；
  - 特定的表，使用 ON database.table；
  - 特定的列；
  - 特定的存储过程。

- 删除权限   `REVOKE SELECT, INSERT ON mydatabase.* FROM myuser;`

- 更改密码 `SET password FOR myuser = Password('new_password');`

## 理解事务的4种隔离级别

数据库事务的隔离级别有4种，由低到高分别为Read uncommitted 、Read committed 、Repeatable read 、Serializable 。

而且，在事务的并发操作中可能会出现脏读，不可重复读，幻读。下面通过事例一一阐述它们的概念与联系。

- Read uncommitted

读未提交，顾名思义，就是一个事务可以读取另一个未提交事务的数据。

事例：老板要给程序员发工资，程序员的工资是3.6万/月。但是发工资时老板不小心按错了数字，按成3.9万/月，该钱已经打到程序员的户口，但是事务还没有提交，就在这时，

程序员去查看自己这个月的工资，发现比往常多了3千元，以为涨工资了非常高兴。但是老板及时发现了不对，马上回滚差点就提交了的事务，将数字改成3.6万再提交。

分析：实际程序员这个月的工资还是3.6万，但是程序员看到的是3.9万。他看到的是老板还没提交事务时的数据。这就是脏读。

那怎么解决脏读呢？Read committed！读提交，能解决脏读问题。

- Read committed

读提交，顾名思义，就是一个事务要等另一个事务提交后才能读取数据。

事例：程序员拿着信用卡去享受生活（卡里当然是只有3.6万），当他埋单时（程序员事务开启），收费系统事先检测到他的卡里有3.6万，就在这个时候！！程序员的妻子要把钱

全部转出充当家用，并提交。当收费系统准备扣款时，再检测卡里的金额，发现已经没钱了（第二次检测金额当然要等待妻子转出金额事务提交完）。程序员就会很郁闷，明明卡里是有钱的…

分析：这就是读提交，若有事务对数据进行更新（UPDATE）操作时，读操作事务要等待这个更新操作事务提交后才能读取数据，可以解决脏读问题。但在这个事例中，出现了一个

事务范围内两个相同的查询却返回了不同数据，这就是不可重复读。

那怎么解决可能的不可重复读问题？Repeatable read ！

- Repeatable read

重复读，就是在开始读取数据（事务开启）时，不再允许修改操作

事例：程序员拿着信用卡去享受生活（卡里当然是只有3.6万），当他埋单时（事务开启，不允许其他事务的UPDATE修改操作），收费系统事先检测到他的卡里有3.6万。这个时候他的妻子不能转出金额了。接下来收费系统就可以扣款了。

分析：重复读可以解决不可重复读问题。写到这里，应该明白的一点就是，不可重复读对应的是修改，即UPDATE操作。但是可能还会有幻读问题。

因为幻读问题对应的是插入INSERT操作，而不是UPDATE操作。

- 什么时候会出现幻读？

事例：程序员某一天去消费，花了2千元，然后他的妻子去查看他今天的消费记录（全表扫描FTS，妻子事务开启），看到确实是花了2千元，就在这个时候，程序员花了1万买了一部电脑，即新增INSERT了一条消费记录，并提交。当妻子打印程序员的消费记录清单时（妻子事务提交），发现花了1.2万元，似乎出现了幻觉，这就是幻读。

- 那怎么解决幻读问题？Serializable！

Serializable 序列化

Serializable 是最高的事务隔离级别，在该级别下，事务串行化顺序执行，可以避免脏读、不可重复读与幻读。但是这种事务隔离级别效率低下，比较耗数据库性能，一般不使用。


## 常见数据库的事务隔离默认级别

| 数据库  | 隔级别 |
| ------------- | ------------- |
| Mysql  | 可重复读（Repeatable Read）|
| Oracle | 读提交（Read Committed） |
|SQLServer| 读提交（Read Committed）|
|DB2| 读提交（Read Committed）|
|PostgreSQL| 读提交（Read Committed）|

## 各隔离级别对各种异常的控制能力

| 隔离级别  | 更新丢失 |脏读 |不可重复读 |幻读 |
| --- | --- | --- | --- | ---|
| RU（读未提交） |避免 | | | |
| RC（读提交） |避免 | 避免| | |
| RR（可重复读） | 避免|避免 |避免 | |
| S（串行化） | 避免| 避免|避免 |避免 |

## 事务四个特性

1. 原子性（Atomicity）

  　　原子性是指事务包含的所有操作要么全部成功，要么全部失败回滚，事务的操作如果成功就必须要完全应用到数据库，如果操作失败则不能对数据库有任何影响。

2. 一致性（Consistency）

  　　一致性是指事务必须使数据库从一个一致性状态变换到另一个一致性状态，也就是说一个事务执行之前和执行之后都必须处于一致性状态。

  　　拿转账来说，假设用户A和用户B两者的钱加起来一共是5000，那么不管A和B之间如何转账，转几次账，事务结束后两个用户的钱相加起来应该还得是5000，这就是事务的一致性。

3. 隔离性（Isolation）

  　　隔离性是当多个用户并发访问数据库时，比如操作同一张表时，数据库为每一个用户开启的事务，不能被其他事务的操作所干扰，多个并发事务之间要相互隔离。

  　　即要达到这么一种效果：对于任意两个并发的事务T1和T2，在事务T1看来，T2要么在T1开始之前就已经结束，要么在T1结束之后才开始，这样每个事务都感觉不到有其他事务在并发地执行。

　　关于事务的隔离性数据库提供了多种隔离级别，稍后会介绍到。

4. 持久性（Durability）

　　持久性是指一个事务一旦被提交了，那么对数据库中的数据的改变就是永久性的，即便是在数据库系统遇到故障的情况下也不会丢失提交事务的操作。