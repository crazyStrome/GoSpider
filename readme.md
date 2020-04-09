#  爬虫

##  简介

使用Go语言实先一个简单的爬虫，使用Redis和Mysql进行持久化处理和查询。目前只实现了网页的一次遍历，递归遍历可以进一步实现。

##  结构

![image-20200409205925938](https://i.loli.net/2020/04/09/RIPJNMc6X1a4T2D.png)

图中带箭头细线是函数调用关系，粗线箭头表示数据流向关系，使用channel实现的不同goroutine之间信息的传递。

## database操作

分别使用Redis和MySQL进行持久化操作。

因为Redis是基于内存的，所以查询比较迅速。

而MySQL进行本地化保存比较方便，可以存储大量数据。

###  mysql.go

该文件中主要有一个save函数实现将数据保存在数据库中。

```go
func Save(url string, title string) (id int64, err error)
```

将链接和网页的title保存在数据库中，目前title使用空字符串表示。

###  redis.go

该文件中实现了数据库的插入操作和查询操作。

在Redis中使用Hash保存数据。

* 测试链接是否存在

	```go
	func HKExists(key string, field string) (bool, error)
	```

* 将指定链接的计数加一

	```go
	func HIncr(key string, field string) (int64, error)
	```

* 当链接不存在数据库时，添加到数据库

	```go
	func HSetNX(key string, field string) (bool, error)
	```

##  链接提取

extract包中实现了网页中链接的提取。

* Content结构体

	使用一个Content结构体表示需要处理的信息载体

	```go
	// Content contents the full path url and the base path
	type Content struct {
		URL string
		Dir string
	}
	```

	有两个字段，其中URL表示需要处理的链接，Dir字段表示该连接的基地址。

	比如URL为```http://www.baidu.com/a```，那么Dir字段应该为```http://www.baidu.com```。

* Extract函数

	```go
	func Extract(c Content) (res []Content, err error)
	```

	对链接对应的网站进行解析，返回网页中包含的链接，如果链接是相对路径，则和基地址组成绝对路径。

* 获取基地址

	```go
	// GetDirURL return the base url of input url
	// for example input http://www.baidu.com/haha.js
	// return http://www.baidu.com
	func GetDirURL(full string) string {
		dir := path.Dir(full)
		// dir="http:" || dir="https:"
		if len(dir) == 5 || len(dir) == 6 {
			return full
		}
		return dir
	}
	```

##  主逻辑处理

main.go是主逻辑处理的文件。

* worker函数

	```go
	func worker(content extract.Content, data chan<- []extract.Content)
	```

	worker 负责网页中链接的提取，把结果传入到data中

* controller函数

	```go
	// controller 负责向worker中分配任务，一个controller中最多可以运行MAXGOROUTINES个worker
	// controller 把上一层的receive传给worker，使worker直接把结果传给上一层
	// controller等待所有子worker 都退出之后再退出
	func controller(contents []extract.Content, receive chan []extract.Content) 
	```

* saver函数

	```go
	// saver 负责将上层中的数据存入MySQL和Redis数据库
	func saver(data <-chan extract.Content)
	```

* identifier函数

	```go
	// identifier 负责将传入的数据分类鉴别，如果在Redis已经存在的，就舍弃掉
	// 如果不存在，就储存在Redis中和MySQL中，传入saver中
	func identifier(receive <-chan []extract.Content)
	```

* main函数

	main函数把controller解析的数据传给identifier进行处理。

