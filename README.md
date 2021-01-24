# fc
fc 是通过HTTP协议接收上传文件的小应用。

主要特性如下：
- 使用HTTP协议访问
- 提供MD5和SHA256摘要验证
- 上传文件可以保存 本地、FastDFS、阿里云OSS
- 上传成功后返回文件下载地址

使用方法:
1. git clone https://github.com/op-y/fc.git
2. cd fc
3. go build -o bin/fc
4. 根据实际情况修改config.json和fastdfs.conf配置
5. ./control start

注意事项:
- 在外网环境需要考虑安全问题，尽量使用HTTPS协议，密钥需要妥善保存并且定时更新。
- 文件上传请求性能相对较低(磁盘IO频繁)，必要时做水平扩展加反向代理。
- FastDFS 和 local 后端特性集成中...
