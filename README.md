#devtools

开发工具套件。


### Build

编译 Windows 版本时，可以用以下命令，生成有图形化控制台的版本：

```
gox -osarch="windows/amd64" -ldflags="-H windowsgui"
```
