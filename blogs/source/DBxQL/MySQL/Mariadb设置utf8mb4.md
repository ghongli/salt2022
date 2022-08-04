MariaDB 设置 utf8mb4

---

> 连接串添加：useUnicode=true&characterEncoding=UTF-8&character_set_server=utf8mb4&character_set_database=utf8mb4

#### 验证方法

```mysql
SHOW VARIABLES WHERE Variable_name LIKE 'character\_set\_%' OR Variable_name LIKE 'collation%';
```

#### 设置方法

```
vim /etc/my.conf

# mariaDB的设置方法：
[mysqld]
character_set_server=utf8mb4 
collation-server=utf8mb4_unicode_ci 
init_connect='SET NAMES utf8mb4' 
skip-character-set-client-handshake=true 

# mysql的设置方法：
[client]
default-character-set = utf8mb4
[mysql]
default-character-set = utf8mb4
[mysqld]
character-set-client-handshake = FALSE
character-set-server = utf8mb4
collation-server = utf8mb4_unicode_ci
init_connect='SET NAMES utf8mb4'
```

`/etc/init.d/mariadb reload`

