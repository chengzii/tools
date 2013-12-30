tools
=====
#linuxinfo.go
是一个登陆显示系统环境信息的工具(目前仅支持linux环境)
a) 编译过程： go build linuxinfo.go  会在当前目录下生成二进制程序   linuxinfo
b) 部署过程： cp linuxinfo /etc/profile.d/linuxinfo
			  chmod +x /etc/profile.d/linuxinfo
			  vi /etc/profile 在文件最后一行增加 /etc/profile.d/linuxinfo
c) 执行过程： 每次登陆linux环境时，都会显示系统信息，内容包括：（os类型，cpu逻辑核数，系统语言环境，"/"目录挂载磁盘的使用情况，内存使用情况）

存放一些工作中制作的小工具
