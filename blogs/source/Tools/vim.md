Vim Skips

---

1. 插入模式，复制出来的格式不正确

   因为 vim 自动缩进了，导致格式不正确。

   在 vim 视图，输入 `:set paste`，可以使 vim 进入 paste 模式，这时整段复制粘贴，就可以了。

   快捷键方法：切换 paste 开关的选项 - pastetoggle，可以通过它来绑定一个快捷键，即可实现单键控制`激活/退出 paste 模式`。

   ```shell
   :set pastetoggle=<F5>
   ```

1. 快捷键

   ```wiki
   - gg 文件开始
   - j  向下
   - k 向上
   - h 向左
   - i  向右
   - o 行首
   - ^ 行首第一个字符
   - $ 行尾
   - G 文件末
   - yy 复制一行
   - x   剪切
   - d   删除
   - p   粘贴
   - !<history pid> 直接使用某个历史命令，history pid 历史命令的标识
   ```
   
   