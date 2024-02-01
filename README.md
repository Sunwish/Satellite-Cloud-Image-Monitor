**引自国家卫星气象中心公告**

> 我中心定于2024年2月1日至3月5日，对风云四号B星实施位置切换和业务调试。3月5日8时，风云四号B星在东经105度恢复业务服务。切换期间，风云四号B星地面应用系统将暂停地面接收与实时数据服务。3月5日对风云四号A星实施位置切换。当日8时起，风云四号A星地面应用系统将停止地面接收与实时数据服务。

# Satellite Cloud Image Monitor

SCIM用于从[国家卫星气象中心](http://www.nsmc.org.cn/nsmc/cn/home/index.html)采集卫星云图（默认采集FY-4B中国区域云图），每15分钟采集一次，新采集的云图存储至镜像目录`./archived`内。当通过Docker镜像使用时，可挂载`/app/archived`至主机目录，以便后续查用。在所采集图像的文件命名中，时间为UTC时间。

*推荐采用Docker部署，若需同时采集多类云图数据，部署多个容器分别配置xmlURL参数即可。*

## 一、部署方式

### 1.1 构建运行

```zsh
go build -o main
```

```zsh
./main [可选参数]
```

### 1.2 Docker部署

#### STEP 1. 拉取镜像

```zsh
docker pull sunwish/satellite_cloud_image_monitor_amd64:latest
```

> arm64版本可拉取`sunwish/satellite_cloud_image_monitor:v1.0.0`
>
> 理论上不再更新arm架构的镜像，如有更新需求烦请自行构建。

#### STEP 2. 运行容器

```zsh
docker run -v /host/path:/app/archived satellite_cloud_image_monitor_amd64:latest ./main [可选参数]
```

## 二、可选参数

| 参数 | 说明 | 默认值 |
| --- | --- | --- |
| -xmlURL | 监控的云图XML地址，默认为FY-4B中国区域云图XML地址 | `http://img.nsmc.org.cn/CLOUDIMAGE/FY4B/AGRI/GCLR/SEC/xml/FY4B-china-72h.xml` |
| -checkCount | 每轮检查xml中的图像个数 | `5` |
| -dateFormat | 存储目录按时间分级格式 | `2006/20060102/` |
| -filterReg | 图像名过滤正则，匹配的图像名将被跳过检查和下载，且不占用checkCount个数 | `.*thumb.*` |
| -notifyBaseUrl | Pushdeer通知接口baseURL（例：`http://notify.example.com/message/`） | 空 |
| -notifyKey | Pushdeer通知接口key | 空 |
| -notifyPrefix | Pushdeer通知自定义前缀 | `[CloudMonitor]` |

## 三、XML清单

| 卫星 | 云图数据 | XML地址 |
| --- | --- | --- |
| FY-4B | 中国区域云图 | `http://img.nsmc.org.cn/CLOUDIMAGE/FY4B/AGRI/GCLR/SEC/xml/FY4B-china-72h.xml` |
| FY-4B | 全圆盘云图 | `http://img.nsmc.org.cn/PORTAL/NSMC/XML/FY4B/FY4B_AGRI_IMG_DISK_GCLR_NOM.xml` |
| FY-4A | 中国区域云图 | `http://img.nsmc.org.cn/PORTAL/NSMC/XML/FY4A/FY4A_AGRI_IMG_REGI_MTCC_GLL.xml` |
| FY-4A | 全圆盘云图 | `http://img.nsmc.org.cn/PORTAL/NSMC/XML/FY4A/FY4A_AGRI_IMG_DISK_MTCC_NOM.xml` |
| FY-4A | 闪电云图 | `http://img.nsmc.org.cn/PORTAL/NSMC/XML/FY4A/LMI.xml` |
| FY-4A | 南海区云图 | `http://img.nsmc.org.cn/PORTAL/NSMC/XML/FY4A/FY4A_AGRI_IMG_REGI_SCS_GLL_C002.xml` |
| FY-4A | 西北太平洋云图 | `http://img.nsmc.org.cn/PORTAL/NSMC/XML/FY4A/FY4A_AGRI_IMG_REGI_PAC_GLL_C002.xml` |
| FY-4A | 兰勃托投影云图 | `http://img.nsmc.org.cn/PORTAL/NSMC/XML/FY4A/FY4A_AGRI_IMG_REGI_PCC_LBT_C012.xml` |
| FY-4A | 麦卡托投影云图 | `http://img.nsmc.org.cn/PORTAL/NSMC/XML/FY4A/FY4A_AGRI_IMG_REGI_PCC_MCT_C012.xml` |
| FY-4A | 中国气象地理区划云图 | `http://img.nsmc.org.cn/CLOUDIMAGE/FY4A/MTCC/SEC/xml/china-72h.xml` |
| FY-3D | MERSI全球影像 | `http://img.nsmc.org.cn/PORTAL/NSMC/XML/FY3D/GLOBAL.xml` |
| FY-3D | MERSI北极地区影像 | `http://img.nsmc.org.cn/PORTAL/NSMC/XML/FY3D/POLAR_NORTH.xml` |
| FY-3D | MERSI南极地区影像 | `http://img.nsmc.org.cn/PORTAL/NSMC/XML/FY3D/POLAR_SOUTH.xml` |
| FY-3D | MERSI“一带一路”区域影像 | `http://img.nsmc.org.cn/PORTAL/NSMC/XML/FY3D/BR.xml` |
| FY-3D | MERSI中国区影像 | `http://img.nsmc.org.cn/PORTAL/NSMC/XML/FY3D/CHINA.xml` |
| FY-3D | MERSI中亚影像 | `http://img.nsmc.org.cn/PORTAL/NSMC/XML/FY3D/CA.xml` |
| FY-3D | MERSI非洲影像 | `http://img.nsmc.org.cn/PORTAL/NSMC/XML/FY3D/AF.xml` |
| FY-3D | MERSI大洋洲影像 | `http://img.nsmc.org.cn/PORTAL/NSMC/XML/FY3D/OC.xml` |
| FY-3D | MERSI北美洲影像 | `http://img.nsmc.org.cn/PORTAL/NSMC/XML/FY3D/NA.xml` |
| FY-3D | MERSI南美洲影像 | `http://img.nsmc.org.cn/PORTAL/NSMC/XML/FY3D/SA.xml` |
| FY-2H | 圆盘图 | `http://img.nsmc.org.cn/PORTAL/NSMC/XML/FY2H/ETV_NOM.xml` |
| FY-2H | 区域云图 | `http://img.nsmc.org.cn/PORTAL/NSMC/XML/FY2H/ETV_SEC.xml` |
| FY-2H | “丝绸之路经济带”云图 | `http://img.nsmc.org.cn/PORTAL/NSMC/XML/FY2H/P1_IR1.xml` |
| FY-2H | “21世纪海上丝绸之路”云图 | `http://img.nsmc.org.cn/PORTAL/NSMC/XML/FY2H/P2_IR1.xml` |
| FY-2H | 东非与西亚云图 | `http://img.nsmc.org.cn/PORTAL/NSMC/XML/FY2H/P3_IR1.xml` |
| FY-2H | 东欧与中亚云图 | `http://img.nsmc.org.cn/PORTAL/NSMC/XML/FY2H/P4_IR1.xml` |
| FY-2G | 双星动画 | `http://img.nsmc.org.cn/PORTAL/NSMC/XML/FY2/FY2_WXCL.xml` |
| FY-2G | 中国区域云图 | `http://img.nsmc.org.cn/PORTAL/NSMC/XML/FY2G/FY2G_LAN_CLC_GRA.xml` |
| FY-2G | 海区云图 | `http://img.nsmc.org.cn/PORTAL/NSMC/XML/FY2G/FY2G_SEA_CLC_GRA.xml` |
| FY-2G | 圆盘图 | `http://img.nsmc.org.cn/PORTAL/NSMC/XML/FY2G/FY2G_GLB_CLC_GRA.xml` |
