# Simple File Server

一个简单的文件服务器，使用 Go 和 Gin 框架构建，支持文件上传、下载和静态文件服务。

## 功能特性

- **文件上传**：支持普通文件上传和分片上传（multipart upload）
- **文件下载**：通过 HTTP GET 请求下载文件
- **静态文件服务**：自动服务数据目录中的文件
- **API 认证**：上传操作需要 admin-api-token 认证
- **跨平台支持**：支持 Linux 和 macOS 的 amd64 和 arm64 架构

## 安装

### 从源码构建

确保你已经安装了 Go 1.22 或更高版本。

```bash
git clone <repository-url>
cd simple-file-server
make
```

构建完成后，二进制文件将在 `build/` 目录中生成。

### 使用 Docker 构建

确保你已经安装了 Docker。

```bash
docker build -t simple-file-server .
docker run -p 60088:60088 --rm \
      -v $(pwd)/data-docker:/data:rw \
      -v $(pwd)/config.json:/config.json simple-file-server
```

### 下载预编译二进制文件

从 release 页面下载适合你平台的二进制文件。

## 配置

服务器通过 `config.json` 文件进行配置：

```json
{
    "debug": false,
    "port": 60088,
    "apiToken": "your-admin-api-token",
    "tempDir": "./temp",
    "dataDir": "./data"
}
```

- `debug`: 是否启用调试模式
- `port`: 服务器监听端口
- `apiToken`: 管理员 API 令牌，用于上传操作
- `tempDir`: 临时文件目录
- `dataDir`: 数据文件存储目录

## 运行

### 使用二进制文件

```bash
./simple-file-server
```

### 使用 Docker

```bash
docker run -p 60088:60088 -v $(pwd)/data:/root/data -v $(pwd)/temp:/root/temp simple-file-server
```

服务器将在配置的端口上启动，并开始监听请求。

## API 文档

### Ping

检查服务器状态。

- **URL**: `/_admin/ping`
- **Method**: GET
- **Response**: `{"code": 0, "msg": "ok", "data": "ok"}`

### 文件上传

上传单个文件。

- **URL**: `/_admin/upload`
- **Method**: POST
- **Headers**:
  - `admin-api-token`: 管理员令牌
- **Form Data**:
  - `file`: 要上传的文件
  - `filePath`: 文件保存路径（必需）
- **Response**: `{"code": 0, "msg": "ok", "data": {"filePath": "path/to/file"}}`

### 分片上传初始化

初始化分片上传。

- **URL**: `/_admin/upload/multipart_init`
- **Method**: POST
- **Headers**:
  - `admin-api-token`: 管理员令牌
  - `Content-Type`: application/json
- **Body**:
  ```json
  {
    "filePath": "example.txt",
    "totalParts": 10,
    "totalSize": 10485760
  }
  ```
- **Response**: `{"code": 0, "msg": "ok", "data": {"uploadId": "123456789"}}`

### 分片上传

上传文件的一个分片。

- **URL**: `/_admin/upload/multipart_upload`
- **Method**: POST
- **Headers**:
  - `admin-api-token`: 管理员令牌
- **Form Data**:
  - `uploadId`: 上传 ID
  - `partNumber`: 分片编号
  - `file`: 分片文件
- **Response**: `{"code": 0, "msg": "ok", "data": "ok"}`

### 分片上传完成

完成分片上传并合并文件。

- **URL**: `/_admin/upload/multipart_end`
- **Method**: POST
- **Headers**:
  - `admin-api-token`: 管理员令牌
  - `Content-Type`: application/json
- **Body**:
  ```json
  {
    "uploadId": "123456789"
  }
  ```
- **Response**: `{"code": 0, "msg": "ok", "data": {"filePath": "example.txt"}}`

### 分片上传中止

中止分片上传并清理临时文件。

- **URL**: `/_admin/upload/abort`
- **Method**: POST
- **Headers**:
  - `admin-api-token`: 管理员令牌
  - `Content-Type`: application/json
- **Body**:
  ```json
  {
    "uploadId": "123456789"
  }
  ```
- **Response**: `{"code": 0, "msg": "ok", "data": "ok"}`

### 检查文件是否存在

检查指定文件是否存在。

- **URL**: `/_admin/has`
- **Method**: POST
- **Headers**:
  - `admin-api-token`: 管理员令牌
  - `Content-Type`: application/json
- **Body**:
  ```json
  {
    "path": "path/to/file.txt"
  }
  ```
- **Response**: `{"code": 0, "msg": "ok", "data": true}` 或 `{"code": 0, "msg": "ok", "data": false}`

### 获取文件大小

获取指定文件的大小。

- **URL**: `/_admin/size`
- **Method**: POST
- **Headers**:
  - `admin-api-token`: 管理员令牌
  - `Content-Type`: application/json
- **Body**:
  ```json
  {
    "path": "path/to/file.txt"
  }
  ```
- **Response**: `{"code": 0, "msg": "ok", "data": {"size": 12345}}` (文件大小字节数)

### 获取文件内容

获取指定文件的内容。

- **URL**: `/_admin/get`
- **Method**: POST
- **Headers**:
  - `admin-api-token`: 管理员令牌
  - `Content-Type`: application/json
- **Body**:
  ```json
  {
    "path": "path/to/file.txt"
  }
  ```
- **Response**: `二进制数据`

### 移动文件

移动文件到新位置。

- **URL**: `/_admin/move`
- **Method**: POST
- **Headers**:
  - `admin-api-token`: 管理员令牌
  - `Content-Type`: application/json
- **Body**:
  ```json
  {
    "from": "old/path/file.txt",
    "to": "new/path/file.txt"
  }
  ```
- **Response**: `{"code": 0, "msg": "ok", "data": "ok"}`

### 删除文件

删除指定文件。

- **URL**: `/_admin/delete`
- **Method**: POST
- **Headers**:
  - `admin-api-token`: 管理员令牌
  - `Content-Type`: application/json
- **Body**:
  ```json
  {
    "path": "path/to/file.txt"
  }
  ```
- **Response**: `{"code": 0, "msg": "ok", "data": "ok"}`

### 文件下载

下载文件。

- **URL**: `/{fileName}`
- **Method**: GET
- **Response**: 文件内容

## 许可证

[Apache 2.0 License](LICENSE)
