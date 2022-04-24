Microsoft SQLServer

---

#### Quickstarts

##### Run and Connect - Docker

###### Pull and run the container image

1. Pull the SQL Server Linux container image from Microsoft Container Registry.

   ```bash
   # docker pull mcr.microsoft.com/mssql/server:2019-latest
   docker pull mcr.microsoft.com/mssql/server:2017-latest
   ```

   

2. To run the container image with Docker.

   ```bash
   docker run -e "ACCEPT_EULA=Y" -e "SA_PASSWORD=<YourStrong@Passw0rd>" \
      -p 1433:1433 --name sql1 --hostname sql1 \
      -d mcr.microsoft.com/mssql/server:2017-latest
   ```

   

3. 



Notes:

- GA = General Availability - when the version is officially available and supported.
- CU = Cumulative Update - a periodic release that includes bug fixes, security fixes, and occasionally a small feature.
- CTP = Community Technology Preview - a preview release that precedes the GA of an upcoming new version.
- GDR = General Distribution Release - a release that contains ONLY security fixes.



[0]: https://docs.microsoft.com/zh-cn/sql/t-sql/language-reference?view=sql-server-ver15	"Transact-SQL 引用 - SQL Server"
[1]: https://hub.docker.com/_/microsoft-mssql-server
[2]: https://github.com/Microsoft/mssql-docker
[3]: https://docs.microsoft.com/en-us/sql/linux/quickstart-install-connect-docker "快速入门 - Linux Container"
[4]: https://docs.microsoft.com/en-us/sql/linux/sql-server-linux-configure-docker "最佳实践 - Docker"

