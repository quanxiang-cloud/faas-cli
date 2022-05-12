# Motivation
开发人员能够本地完全的管理Faas,包括代码编写、调试、提交、编译、部署，而不用登陆到protal端。

## Proposal
`faas-cli`使用本地配置文件
```yaml=
url: https://api.quanxiang.dev
user.email: wenttang@yunify.com
user.password: 123456
```
打通用户登陆模块，通过使用命令，直接操作`Faas`的`API`。
主要应该包含以下功能：
- `config`: 管理用户邮箱、密码等配置信息
- `create`: 创建`group`、`project`等资源
- `list`: 获取`group`、`project`等资源列表
- `get`: 获取指定资源详情
- `delete`: 删除指定资源,包括`serving`
- `build`: 打包指定xiang m
- `publish`: 发布指定的项目版本
- `log`: 查看build和serving的日志
- `run`: 本地run,能够访问到低代码平台的沙箱

## Usage

### cli config
> cli config [OPTIONS]
**OPTIONS**
```shell
  -h, --help            帮助
      --user.email      用户邮箱
      --user.password   密码
      --url             代码API URL
 ```

***Example***
```bash=
cli config --url https://api.quanxiang.dev \
    --user.email wenttang@yunify.com  \
    --user.password 123456
```

### cli create
> cli create [RESOURCE] [OPTIONS] 

#### cli create group 
> cli create group [OPTIONS] 
group name,the name can contain only lowercase letters, numbers, and hyphens (-), must start with a lowercase letter, and must end with a lowercase letter or number. The maximum length is 20 characters.
**OPTIONS**
```shell
  -h, --help            帮助
      --app             应用标识，必须是应用的管理员
      --describe        描述
  -b, --binding         如果git group已经存在，则绑定该group
      --binding-only    只绑定，与binding互斥
      --join            加入已经存在的group，必须是应用的管理员
 ```
***Example***
```bash=
cli create group demo  \
   --app demo  \
   -b  \
   --describe "sample group"
```


#### cli create project
> cli create project [OPTIONS] 
```shell
  -h, --help            帮助
      --group           group名字
      --describe        描述
      --alias           别名
  -l, --language        项目语言，默认GO
  -v, --version         语言版本
  -b, --binding         如果git group已经存在，则绑定该group
      --binding-only    只绑定，与binding互斥 
```
***Example***
```bash=
cli create project demo  \
   --group demo  \
   -l GO  \
   -v 1.16  \
   --describe "sample project"
```

#### cli run
> cli run [OPTIONS] 
```shell
  -h, --help            帮助
      --group           group名字
      --src             package路径
```
```bash=
cli run demo  \
   --group demo
   --src   .
```