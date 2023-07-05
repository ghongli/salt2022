macos xcode issues

---

1. `xcrun: error: developer path`

   ```shell
   go install -v golang.org/x/tools/gopls@latest
   runtime/cgo
   # runtime/cgo
   xcrun: error: active developer path ("/Applications/Xcode.app/Contents/Developer") does not exist
   Use `sudo xcode-select --switch path/to/Xcode.app` to specify the Xcode that you wish to use for command line developer tools, or use `xcode-select --install` to install the standalone command line developer tools.
   See `man xcode-select` for more details.
   ```

   尝试了如下的方式：

   ```shell
   mkdir -p /Applications/Xcode.app/Contents/Developer
   go install -v golang.org/x/tools/gopls@latest
   runtime/cgo
   # runtime/cgo
   xcrun: error: invalid active developer path (/Applications/Xcode.app/Contents/Developer), missing xcrun at: /Applications/Xcode.app/Contents/Developer/usr/bin/xcrun
   ```

   更新 Mac OS 后，无法运行 git, gcc 等命令，或卸载了 xcode 应用后，出现 missing xcrun 错误

   ```shell
   # 下载 xcode 命令行工具，但并不可行
   xcode-select --install
   xcode-select: error: command line tools are already installed, use "Software Update" to install updates
   ```

   ```shell
   ln -s /usr/bin/xcrun /Applications/Xcode.app/Contents/Developer/usr/bin/xcrun
   go install -v golang.org/x/tools/gopls@latest
   runtime/cgo
   # runtime/cgo
   xcode-select: error: malformed developer path ("/Applications/Xcode.app/Contents/Developer")
   ```

   ```shell
   type xcrun
   ls -hna /Library/Developer/CommandLineTools/usr/bin/
   ```

   最终的解决方式：

   ```shell
   # 查看真正的 developer path
   xcode-select -p 
   
   sudo xcode-select --switch /Library/Developer/CommandLineTools
   go install -v golang.org/x/tools/gopls@latest
   ```

   或尝试安装 xcode 命令行工具

   ```shell
   1. 打开Apple的开发者下载：https://developer.apple.com/download/more/
   2. 使用 AppleID 登录下，登录完成后，在左边的搜索框中搜索Command Line Tools
   3. 下载适合Mac OS 版本的 xcode 命令行工具，安装好即可
   ```

   

2. 