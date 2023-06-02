## easyserver

## 这是个啥？

如果你想上传一个文件到服务器，你有什么方式？

- ftp
- scp ?
- ???

ftp 是常用的。你可以选择传统的 ftp server，例如 vsftpd.

但 ftp server 上传文件，认证方式是通过用户名和密码，你需要为每个用户设置一个密码，然后把密码告诉用户，用户才能上传文件。

如果你想让用户上传文件，但不想让用户知道密码，那么你可以使用 token 的方式。

于是，这个工具就诞生了。

这是一个可以支持 user-pass、token、annymous 认证的 http server，拥有如下能力：

- 上传文件
- 下载文件
- 查看文件目录

对，都是很基础的功能，但是，这个工具可以让你更方便地使用这些功能。

比如，在 浏览器中，直接访问 `http://xxx.xxx.xx/_dash`， 你就能看到类似下面的页面：

![dashboard](./doc/dash_info1.jpeg)

像这样上传:

![upload](./doc/dash_info2.png)

像这样做认证:

![auth](./doc/dash_info3.png)

可以像这样查看文件:

![view](./doc/dash_info4.png)

快玩起来吧 🍻🍻🍻

## 先跑起来吧

```bash
easyserver serve .
```

```bash
# 设置 host 和 port
easyserver serve . --addr 0.0.0.0:8080
```

```bash
# 设置 https 证书 和 domain
easyserver serve . --https domain.longalong.cn:./certs/xx.pem:./certs/xx.key
```

```bash
# 设置 user 和 password
easyserver serve . --user admin:passadmin --user readuser:passxx:r:/data/img
```

```bash
# 使用 config 文件
easyserver --config ./config.yaml
```

```bash
# 查看文件 
curl http://127.0.0.1:8080/data/img/xx.jpg

# 需要 basic auth
curl http://admin:passadmin@127.0.0.1:8080/data/img

# 使用 token
curl http://127.0.0.1:8080/data?token=xxxxx
```

```bash
# 上传文件
curl -F "file=@./xx.jpg" http://127.0.0.1:8080/data/longimg.jpg

# basic auth
curl -u admin:passadmin -F "file=@./xx.jpg" http://127.0.0.1:8080/another/longxxx.jpg

# token
curl -F "file=@./xx.db" 'http://127.0.0.1:8080/xx.db?token=xxxxx'
```

```bash
# 生成 token
curl http://admin:easyadmin123@127.0.0.1:8080/_token -H "Content-Type: application/json" -d  '
{
    "path_roles": [
        {"path": "/pulic", "mode": "r"},
        {"path": "/path/to/write", "mode": "w"}
    ],
    "duration": "12h",
    "count_limit": 20
}
'
## {"token":"ec75Aef6FC9e"}

# 删除 token
curl -X DELETE 'http://admin:easyadmin123@127.0.0.1:8080/_token?token=xxxxxx'

# 删除 该 user 所有 token
curl -X DELETE 'http://admin:easyadmin123@127.0.0.1:8080/_token?all=true'
```

## 如何安装

```bash
# 如果你用 golang 环境
go install github.com/iamlongalong/easyserver

# 如果你通过 github 下载
# 1. 到页面 https://github.com/iamlongalong/easyserver/releases
# 2. 下载一个正确的版本
# 3. 移动到 PATH 下，eg: sudo mv easyserver-darwin-amd64 /usr/local/bin/easyserver
# 4. 开始使用  easyserver serve .

# 如果你想通过脚本一键安装，可以使用 (当然，前提是你能访问下面的地址……)
sudo bash -c "$(curl -fsSL https://raw.githubusercontent.com/iamlongalong/easyserver/master/update-easyserver.sh)"

# 所以，给一个国内的地址
sudo bash -c "$(curl -fsSL https://static.longalong.cn/scripts/get-easyserver.sh)"

# 如果你想通过 docker 安装
docker run -p 8080:8080 -v `pwd`:/data --rm --name easyserver -itd iamlongalong/easyserver easyserver serve /data

# then just enjoy your life ~

```


## 已经实现的功能

先做一个最基本的版本，仅实现下面的基本功能：

- [x] 提供一个简单的 http server，可以指定端口、ip、根目录
- [x] 可以指定一个上传目录，可以上传文件到指定目录
- [x] 可以指定 user 和 password，可以通过 basic auth 访问
- [x] 可以为 user 指定目录的权限，可以指定为只读、只写、读写
- [x] user 可以生成 token，可以指定 token 的有效期、可用次数、可用路径、可用操作类型(r/w)
- [x] 可以通过 token 访问，可以指定 token 的有效期、write可用次数、可用路径、可用操作类型(r/w)
- [x] 可以指定 https 的证书文件
- [x] 可以通过 config 文件配置，也可以通过命令行参数配置

后面大概率会加上的功能：

- [x] 一个简单的 dashboard，用于创建 token、查看 token、删除 token (update: 2023-06-02)
- [ ] dashboard 中可以 list、delete 文件及目录、上传文件及目录
- [ ] 支持生成自签证书

其他 todo:

- [ ] 使用的 demo 放一个 gif 图

> 以上是最基本的功能


## 一些使用场景

### 自建简单的文件服务器

若你有一台云服务器，然后你希望能够比较方便地把一个目录共享出来，可以快速地 `查看文件`、`下载到本地` 这类操作，就可以使用这个工具。

但目前还没有弄 dashboard，所以还不是很好用。

如果，你还想上传文件到服务器 (通过 rest api 接口)，那么也可以用这个工具。

- [ ] 梳理更多的应用场景


## 关于界面的基本梳理

使用这个工具，最直接的场景，是在浏览器端直接打开

存在 3 种鉴权方式: basic auth, token

目前采用的 在 response header 中添加 `WWW-Authenticate` 的方式，让浏览器弹出 basic auth 的窗口。这是最基础的做法。

如果使用的 dashboard 页面，这种方式则不完全适用于 token 的模式，因此，要兼容使用 dashboard 时不弹出弹窗，而是由用户自己决定如何认证。

使用流程如下：
- 打开 dashboard 页面
- 初始打开页面时，鉴权(`/_api/user`)，若 401 ，则可设置鉴权类型 (① basic auth ② token ③ annymous)，若正常，则直接进入主页面
- 设置鉴权类型后，页面未重新加载前不再主动弹鉴权窗口

dashboard 需要提供如下能力：
- 更好看的 list dir
  - 支持排序
- 更好看的 查看文本、图片
- 更好看的 download
  - 复制下载链接
  - 批量下载?
- 更好用的 token 管理 (增删改查)
- 更好看的 upload 界面 (这个最重要)
  - 好用的选择方式
  - 批量上传
  - 断点续传

### 还是第一期做最需要的功能吧

- [x] 认证
- [x] upload
- [x] list

### 技术选型

为了简单，直接用 vue + element plus 的方式，先做一版吧，跑起来看看先 🏃🏻‍♀️🏃🏻‍♀️🏃🏻‍♀️

嗯，加了一个简单的，通过访问 `/_dash` 进入 dashboard 页面

