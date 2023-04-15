lsblk 查看块设备

---

### 安装

#### Centos

```shell
yum install util-linux
```

#### 查看

```shell
whereis lsblk
rpm -qf /usr/bin/lsblk
lsblk --version
```

### 用途

以树形列出所有块设备，输出信息各字段的含义：

```wiki
NAME   :设备的名称
MAJ:MIN:主要设备号和次要设备号
RM:是否可移动设备，值为1时表示可移动
SIZE:设备的容量大小
RO:   是否只读，值为1时表示只读
TYPE: disk: 磁盘
part:   分区
rom:   光盘   
MOUNTPOINT:设备的挂载点,通常是目录
```

### 常用命令

```shell
# 列出所有设备，包括空设备
lsblk -a

# 列出设备的权限和所属的组
# -m参数：显示owner:所属用户,group:设备所属的组,mode:访问模式
lsblk -m

# 只列出指定的设备
lsblk /dev/vda

# 列出scsi设备
lsblk -S

# 以字节显示大小
# -b:size的显示单位用字节
lsblk -b
```

