# 架构
##	启动流程
###初始化
	1.加载配置文件
	2.mkdir 目录文件
	3.加载Meta文件(描述当前有多少索引，多少文档数 多少term 多少G数据)
	4.遍历目录下面的每个目录下面的currentIndex Meta文件 目录名(index_{name}_{m|r}_{1..n}).
	5.恢复log

