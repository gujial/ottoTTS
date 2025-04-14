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

## 作为 go package 引入
```shell
go get github.com/gujial/ottoTTS@latest
```
> 本方法不会同步资源文件，需要自行下载到工作目录中
## 安装 cli 工具
```shell
go install github.com/gujial/ottoTTS/cli/ottoTTScli@latest
```
> 本方法不会同步资源文件，需要自行下载到工作目录中

## 配置文件说明
`config.toml`文件格式如下
```toml
expression_override = true
Debug = false
```
`expression_override` 为 true 时，会进行短语匹配，否则只匹配单个字符。`Debug` 为 true 时，会在控制台打印调试信息。
## 修改资源文件

### 音频文件说明
放在`assets/sounds`文件夹下，使用汉语拼音（不含声调）+`.wav`的格式命名，以便使用 [go-pinyin](https://github.com/mozillazg/go-pinyin) 库返回的拼音直接找到音频文件。
> 所有 wav 文件的格式应当相同，合并音频时使用第一段音频的格式，不相同的可以使用 ffmpeg 处理

### `dictionary.json`文件说明
`dictionary.json`的文件格式如下
```json
{
  "expressions": [
    {
      "expression": ["<短语1>", "<短语2>", ...],
      "otto": "<音频文件主文件名>"
    },
    ...
  ],
  "letters": [
    {
      "expression": ["<大写字母>", "<小写字母>", "<其他形式>", ...],
      "otto": "<音节1> <音节2> ..."
    },
    ...
  ],
  "letters": [
    {
      "expression": ["<数字1>", "<数字2>", ...],
      "otto": "<音节1> <音节2> ..."
    },
    ...
  ],
}
```
`expressions`中的短语在匹配到后会跳过短语长度再进行接下来的匹配，适用于多音节的发音。其他两个字段只匹配单个字符。在`otto`字段中的音节使用空格隔开。

## 使用到的开源库
- [go-pinyin](https://github.com/mozillazg/go-pinyin)
- [toml](https://github.com/BurntSushi/toml)