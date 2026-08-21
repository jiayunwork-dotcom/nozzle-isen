# nozzle-isen

nozzle-isen 是理想气体等熵喷管核算的 Go 命令行内核：读取一个 JSON 算例（滞止温度 T0、滞止压力 p0、热容比 γ、气体常数 R、面积比 A_e/A*、喉道面积 A*，可选设计膨胀支与背压 pb），用同一 γ 与同一滞止状态完成喉道壅塞判定、面积–马赫反解与壅塞/未壅塞质量流量核算。

## 构建 / 运行 / 测试

```text
go build ./...     # 编译
go run .
go test ./...      # 测试
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
