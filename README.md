# github-issue-editor

一个基于 Go 语言的 GitHub Issue 命令行管理工具

## 前置条件

- GITHUB_TOKEN
Settings-Developer Settings-Personal access tokens-Tokens(classic)
需要打开完整repo权限

```shell
# linux / mac
export GITHUB_TOKEN=你的_GitHub_Token
```

```shell
# Windows Powershell
$env:GITHUB_TOKEN="你的_GitHub_Token"
```

如果不使用vscode打开
在此只是举个例子

```shell
export EDITOR=vim
```

## 使用

```shell
go build -o gie.exe main.go
# or
go build -o gie main.go
```

基本语法：`./gie <command> <owner> <repo> [number]`

## 功能

- 创建issue
- 查看issue
- 更新issue
- 关闭issue

### 创建issue

```shell
./gie create myuser myrepo
```

### 查看issue

```shell
./gie read myuser myrepo 1
```

### 更新issue

```shell
./gie update myuser myrepo 1
```

### 关闭issue

```shell
./gie close myuser myrepo 1
```

### Warning

请不要在issue里面丢垃圾，只推荐private自行尝试或不使用
