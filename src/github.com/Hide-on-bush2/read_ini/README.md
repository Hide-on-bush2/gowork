# 读取配置文件包

## 使用方法

首先选择包中实现的监听函数中的一个：

在使用包之前将其import到你的go文件中：

```
import "github.com/Hide-on-bush2/read_ini"
```

```
listenFunc := read_ini.Listen_methods{read_ini.MyListen}
```

然后使用这个监听函数来监听你的配置文件：
```
read_ini.Watch(path, listenFunc)
```

便可以实现监听了

## Demo

在main包中的`main.go`文件中导入这个包并监听处于`../read_ini/test.ini`的配置文件：
```
package main

import "github.com/Hide-on-bush2/read_ini"

func main() {
	hide_on_bush := read_ini.Listen_methods{read_ini.MyListen}
	read_ini.Watch("../read_ini/test.ini", hide_on_bush)

}

```

然后在main目录下go install,然后执行`main`:

![](https://i.loli.net/2020/10/19/QdvjPu1fAwVEYTM.png)