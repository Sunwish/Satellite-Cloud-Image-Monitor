# Satellite Cloud Image Monitor

SCIM用于从国家卫星气象中心采集卫星云图（FY-4B中国区域云图），每15分钟采集一次，每次采集数据为1小时30分钟前的最新云图。云图将存储至目录`./archived`内。当通过Docker镜像使用时，可挂载`/app/archived`至主机目录，以便后续查用。在所采集图像的文件命名中，时间为UTC时间。

## 使用方式一：构建运行

notifyBaseUrl与notifyKey参数为可选，如需配置异常提醒，请在启动参数配置此两参数，具体参见Pushdeer相关文档。

```zsh
go build -o main
./main -notifyBaseUrl http://notify.example.com -notifyKey your_notify_key
```

## 使用方式二：Docker

### 拉取镜像

```zsh
docker pull sunwish/satellite_cloud_image_monitor_amd64:latest
```

> arm64版本可拉取`sunwish/satellite_cloud_image_monitor:v1.0.0`
>
> 理论上不再更新arm架构的镜像，如有更新需求烦请自行构建。

### 运行容器

notifyBaseUrl与notifyKey参数为可选，如需配置异常提醒，请在启动参数配置此两参数，具体参见Pushdeer相关文档。

```zsh
docker run -v /host/path:/app/archived satellite_cloud_image_monitor_amd64:latest ./main -notifyBaseUrl http://notify.example.com -notifyKey your_notify_key
```
