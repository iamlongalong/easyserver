## easyserver

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

- [ ] 一个简单的 dashboard，用于创建 token、查看 token、删除 token
- [ ] dashboard 中可以 list、delete 文件及目录、上传文件及目录

其他 todo:

- [ ] 使用的 demo 放一个 gif 图

> 以上是最基本的功能


## 一些使用场景

### 自建简单的文件服务器

若你有一台云服务器，然后你希望能够比较方便地把一个目录共享出来，可以快速地 `查看文件`、`下载到本地` 这类操作，就可以使用这个工具。

但目前还没有弄 dashboard，所以还不是很好用。

如果，你还想上传文件到服务器 (通过 rest api 接口)，那么也可以用这个工具。

- [ ] 梳理更多的应用场景
