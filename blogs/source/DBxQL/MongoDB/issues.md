MongoDB Issues

---

1. mdb SCRAM 认证失败

   使用SCRAM，MongoDB验证所提供的用户凭证 [`name`](https://docs.mongodb.com/v4.2/reference/system-users-collection/#admin.system.users.user) , [`password`](https://docs.mongodb.com/v4.2/reference/system-users-collection/#admin.system.users.credentials) 和 [`authentication database`](https://docs.mongodb.com/v4.2/reference/system-users-collection/#admin.system.users.db) 。身份验证数据库是创建用户的数据库，它与用户名一起用于标识用户。

   mdb 连接串中必须包含 username, passowrd, authentication database(e.g. admin)。

2. 

[0]: https://docs.mongoing.com/ "MongoDB-CN-Manual"