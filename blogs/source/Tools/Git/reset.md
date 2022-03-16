---
title: git reset 回退提交
date: 2022-03-16 19:31:10
---

[TOC]

### `git reflog`

记录每一次的操作。

### `git log`

显示从最近到最远的提交日志，如果输出的信息太多，看得眼花缭乱，可以加上 `--pretty=oneline` ，会把每条记录以一行的形式输出：

```assembly
git log --pretty=oneline
git log --oneline
```

### `git HEAD`

Git 中，`HEAD` 表示当前版本，也就是最新的一次提交，上一个版本就是 `HEAD^`，上上个版本就是 `HEAD^^`，当然往上 100 个版本写 100 个 `^`，显然是不现实的，可以写成 `HEAD~100`。

### `git reset 将当前 <HEAD> 重置为指定状态`，有以下几种模式：

#### `--soft`

不删除工作区改动代码，撤销 `commit`，但不撤销 `git add .`，即 `git status`是绿色的状态。

```shell
git reset --soft HEAD^
```

#### `--mixed`

重置索引，但不重置工作树，更改后的文件标记为未提交(`git add`)的状态，即不删除工作区改动代码，撤销 `commit`，撤销 `git add .`，是默认参数。

```shell
git reset --mixed HEAD^
# default op
git reset HEAD^
```

#### `--hard`

删除工作区改动代码，撤销 `commit`，撤销 `git add .`，注意完成这个操作后，就恢复到了上一次的 `commit` 状态，从指定的 `<commit>`往后，工作树中的任何变化都会被丢弃。

```shell
git reset --hard HEAD^
```

这里如果误操作，在 `git commit`、git `add` 的情况下，是可以回退恢复的！！

#### `--merge`

重置索引并更新工作树中在 `<commit>`、`<HEAD>` 之间不同的文件，但保留那些在索引和工作树之间不同的文件，即那些未被添加的修改。如果一个在`<commit>`和索引之间不同的文件有未分阶段的变化，重置将被终止。

也就是说，`--merge` 做的是类似于`git read-tree -u -m <commit>` 的事情，但会转发未合并的索引条目。

```shell
git reset --merge HEAD^
```

#### `--keep`

重置索引并更新工作树中在 `<commit>`、`<HEAD>` 之间不同的文件，如果一个在`<commit>`和索引之间不同的文件有本地修改，重置将被终止。

```shell
git reset --keep HEAD^
```

#### `--[no-]recurse-submodules`

当工作树被更新时，使用 `--recurse-submodules` 也将根据超级项目中记录的提交，递归地重置所有活动的子模块的工作树，同时也将子模块的 HEAD 设置为在该提交中被分离。

### `git reset`、`git revert`区别

`git revert` 后会多出一条 `commit`，提醒别人，这里有回撤操作。

`git reset` 直接将这前 `commit` 删除，非 `git reset --hard` 操作是不会删除修改代码，如果 `remote repo`已经有之前代码，需要强推 `git push -f`。

### Issues

#### `commit` 注释写错了，只想改下注释

```shell
git commit --amend
```

此时会进入默认 vim 编辑器，修改注释，保存就好。

#### 撤销误提交的 `commit`

```shell
git reset --mixed HEAD^
# or
git reset --mixed HEAD~1
# or
git reset --soft HEAD^
```

#### `git reset --hard HEAD^` 误操作后恢复

> 注意此操作会将 `<commit>`记录一起消除！Git 可以指定回到未来的某个版本，但需要知道 `commit_id`，如果不知道可以使用 `git reflog` 查询。

- 没有 commit, add，找不回来了，Git 本地缓存没有

- 没有 commit，但是有 add 操作

  ```shell
  git fsck --lost-found
  # 在项目 .git/lost-found/other 目录下 add 过的文件
  
  # 找回本地仓库里边最近 add 的 60 个文件
  find .git/objects -type f | xargs ls -lt | sed 60q
  ```

- 执行过 commit，使用 `git reflog` 确认回退的 commit

  ```shell
  git reflog
  git reset --hard  HEAD@{7}
  ```

  
