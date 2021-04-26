# serviceMonitor
golang开发windows服务管理桌面程序(lxn/walk)

通常小的项目都部署在用户Win server中，数据库、web服务和业务系统各部署一个Windows服务，在Windows的服务管理程序即可管理。

最近一个小项目用户提了个需求，需要一个win服务的管理界面，索性用golang lxn/walk做了一个。现分享出来，供大家参考。
运行环境：windows server 2012

# 运行界面：



代码很简单，定义了三个全局变量，对应三个服务，自行修改为自己项目中用到的服务信息
myService.txt 对应服务显示名称，界面显示使用
myService.serviceName 对应服务名，控制服务状态使用
app.title 是窗口标题

Windows服务的操作用到了github.com/shirou/gopsutil/winservices包，并简单封装了一下。

# 编译

直接运行build.bat即可编译成exe可执行文件，可以单个exe文件部署，比较方便。
为了使生成的exe文件大小小一些，采用了32位的编译，在32位或64位系统中都可以运行。
还可以用upx再进一步压缩生成的exe文件大小，最终文件大小不到2M。

如果缺少rsrc.exe和upx.exe文件，可按以下方式下载安装。

rsrc安装

`go get github.com/akavel/rsrc`
upx下载地址：https://github.com/upx/upx/releases
