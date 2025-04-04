# otto 文字转语音引擎
---
## 构建
直接通过`make`构建
```shell
make
```
同
```shell
make full-build
```

## 使用方法
通过命令行传递参数
```shell
./ottoTTScli <str>
```
或通过标准输入传递字符串
```shell
./ottoTTScli
输入需要转换的字符串
我是电棍，你好你好
```
转换后的 wav 文件存放在当前工作目录中的`otto.wav`文件中