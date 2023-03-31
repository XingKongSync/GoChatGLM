# GoChatGLM

## 简介

ChatGLM-6B 是一个开源的、支持中英双语的对话语言模型，详情请参见[原代码仓库](https://github.com/THUDM/ChatGLM-6B)。

这是一个用 Go 语言编写的服务端，实现了会话管理，允许多用户访问同一个 ChatGLM 实例。

## 使用方法

1. 阅读 ChatGLM 官方仓库的说明，配置好 Python 环境，并下载 ChatGLM 模型
2. 我对官方的 api.py 接口进行了修改，请下载这个仓库的源码 https://github.com/XingKongSync/ChatGLM-6B
3. 参照上述仓库中的 bat 文件，传入正确的参数来启动 api.py
4. 编译并启动 GoChatGLM，默认端口号为3001

可以通过传入`port`参数来配置端口号，通过传入`www`来配置静态文件服务的目录，示例如下。

```shell
go run main.go -port=3000 -www="D:\My Projects\react\my-chatglm\build"
```



## 接口说明

### 创建会话

**接口地址：**

`http://127.0.0.1:3001/api/session`

**请求参数：**

`无`

**返回示例：**

```json
{
    "code": 0,
    "message": "成功",
    "session": "4315db8e-ccff-42bd-abbc-446ca5796fbe"
}
```



### 与 ChatGLM 对话

**接口地址：**

`http://127.0.0.1:3001/api/chat`

**请求参数：**

```json
{
    "prompt": "你好",
    "session": "4315db8e-ccff-42bd-abbc-446ca5796fbe"
}
```

**返回示例：**

```json
{
    "code": 0,
    "message": "成功",
    "response": "你好👋！我是人工智能助手 ChatGLM-6B，很高兴见到你，欢迎问我任何问题。"
}
```



### 与 ChatGLM 对话（打字机效果）

**接口地址：**

`http://127.0.0.1:3001/api/streamchat`

**请求参数：**

```json
{
    "prompt": "你好",
    "session": "4315db8e-ccff-42bd-abbc-446ca5796fbe"
}
```

