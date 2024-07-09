## IMGUR 包
全部使用 `golang` 的包体，未使用`lib`库。 性能上有待验证。

## 基础用法

### 链式调用 (推荐)

链式调用看起来非常间接哦~！！

```go
package main

import "gateway/lib/imgur"

func main() {
	imgA := imgur.Load("9961593.png").
		Resize(160, 160)
	//装载背景图
	imgur.Load("workingID.png").
		Insert(imgA, 310, 640).
		Insert(imgA, 310, 1240).
	    Save("out.png")
}
```

### 非链式调用

```go
package main

import "gateway/lib/imgur"

func main() {
	imgA := imgur.Load("9961593.png")
	imgA.Resize(160, 160)
	//装载背景图
	img := imgur.Load("workingID.png")
	img.Insert(imgA, 310, 640)
	img.Insert(imgA, 310, 1240)
	img.Save("out.png")
}

```
这个例子和上个例子是等效的。

## 读取图像
读取图像非常简单，只要调用 `imgo.Load()` 方法即可。
```go
imgur.Load("9961593.png")
```

该方法不仅可以读取本地文件，还支持以下输入格式。

* 文件系统中图像的绝对路径或相对路径
* 图片的 URL
* Base64 编码的图像数据
* *os.File 类型的文件实例
* 实现了 image.Image 接口的类型实例
* *imgo.Image 类型的实例

实际上在内部实现了如下方法：

```go
Load(source interface{}) *Image
loadFromString(source string) (i *Image)
LoadFromFile(file *os.File) (i *Image)
LoadFromPath(path string) (i *Image)
imgur.LoadFromUrl(url string) (i *Image)
LoadFromImgo(i *Image) *Image
```

## 错误处理

为了支持简洁的链式调用，大部分方法没有返回错误信息的参数。

当发生异常情况时，会收集这些错误信息，并打印在命令行。在发生错误的方法之后在发起新的调用 方法是不会生效的。

但仍然有办法去判断调用的方法是否发生错误，如下例所示。

```go
package main

import (
    "gateway/lib/imgur"
    "fmt"
)

func main() {
	err := imgur.Load("9961593.png").Resize(160, 160).Save("out.png").Error
	if err != nil {
		fmt.Println("error:", err.Error())
	} else {
		fmt.Println("success")
	}
}
```

嘿嘿，里面有个 `Error` 属性，用来记录方法调用过程中的错误信息，当该属性不为 `nil` ，则表示发生了错误。

## 主要功能方法 

### Blur
使用均值模糊的方式模糊图像。

#### 参数
| 参数名 | 类型 | 说明
|:--- | :--- |:---
| ksize | `int` | 卷积核尺寸

#### 返回值
`*imgo.Image` 类型的实例。

#### 例子
```go
package main

import (
	"gateway/lib/imgur"
)

func main() {
	imgur.Load("9961593.png").
        Blur(5).
        Save("out.png")
}
```

### Bounds
返回当前图像的 Bounds 。

#### 参数
无

#### 返回值
`image.Rectangle` 类型

#### 例子
```go
package main

import (
	"gateway/lib/imgur"
	"fmt"
)

func main() {
	img := imgur.Load("workingID.png")
	fmt.Println(img.Bounds())
}
```
输出：

```
(0,0)-(750,1448)
```

### Canvas
创建画布。

#### 参数
| 参数名       | 类型 | 说明 |
|:----------|:--- | :--- |
| width     | `int` | 画布宽度
| height    | `int` | 画布高度
| fillColor | `color.Color` | 画布颜色。非必需参数，默认为透明。

#### 返回值
`*imgo.Image` 类型的实例。

#### 例子
```go
package main

import (
	"gateway/lib/imgur"
	"golang.org/x/image/colornames"
)

func main() {
	imgur.Canvas(750, 1500, colornames.Deepskyblue).
		Save("out.png")
}
```

### Circle
在图像中绘制圆。

#### 参数

| 参数名    | 类型 | 说明 |
|:-------|:--- | :--- |
| x      | `int` | 圆心坐标 x 轴
| y      | `int` | 圆心坐标 y 轴
| radius | `int` | 圆半径
| c      | `color.Color` | 圆的颜色

#### 返回值

`*imgo.Image` 类型的实例。

#### 例子
```go
package main

import (
	"gateway/lib/imgur"
	"golang.org/x/image/colornames"
)

func main() {
	imgur.Canvas(750, 1500, colornames.Deepskyblue).
		Circle(150, 150, 100, colornames.Black).
		Save("out.png")
}
```

### Color2Hex
颜色转十六进制。

#### 参数

| 参数名 | 类型 | 说明 |
|:----|:--- | :--- |
| c   | `color.Color` | 颜色

#### 返回值
string

#### 例子
```go
package main

import (
	"fmt"
	"gateway/lib/imgur"
	"golang.org/x/image/colornames"
)

func main() {
	fmt.Println(imgur.Color2Hex(colornames.White))
}
```

#### 输出：

```
#FFFFFF
```

### Crop
裁剪图像。

#### 参数

| 参数名    | 类型 | 说明 |
|:-------|:--- | :--- |
| x      | `int` | 裁剪区域左上角坐标的 x 轴
| y      | `int` | 裁剪区域左上角坐标的 y 轴
| width  | `int` | 裁剪区域的宽度
| height | `int` | 裁剪区域的高度

#### 返回值
`*imgo.Image` 类型的实例。

#### 例子
```go
package main

import (
	"gateway/lib/imgur"
)

func main() {
	imgur.Load("workingID.png").
		Crop(50, 50, 50, 50).
		Save("out.png")
}
```

### Ellipse
在图像中绘制椭圆。

#### 参数
| 参数名    | 类型 | 说明 |
|:-------|:--- | :--- |
| x      | `int` | 椭圆中点坐标 x 轴
| y      | `int` | 椭圆中点坐标 y 轴
| width  | `int` | 椭圆宽度
| height | `int` | 椭圆高度
| c      | `color.Color` | 椭圆颜色

#### 返回值
`*imgo.Image` 类型的实例。

#### 例子
```go
package main

import (
	"gateway/lib/imgur"
	"golang.org/x/image/colornames"
)

func main() {
	imgur.Canvas(300, 300, colornames.Deepskyblue).
		Ellipse(150, 150, 150, 100, colornames.Black).
		Save("out.png")
}
```

### Extension
返回当前图像的文件扩展。

#### 参数
无

#### 返回值
`string`

#### 例子
```go
package main

import (
	"fmt"
	"gateway/lib/imgur"
)

func main() {
	fmt.Println(imgur.Load("workingID.png").Extension())
}
```

#### 输出：

```
png
```

### Filesize
返回当前图像的文件大小，单位为字节。只有从文件路径加载的图像该属性才有效，否则返回 `0`。

#### 参数
无

#### 返回值
int

#### 例子
```go
package main

import (
	"fmt"
	"gateway/lib/imgur"
)

func main() {
	fmt.Println(imgur.Load("workingID.png").Filesize())
}
```

#### 输出：

```
255493
```

### Flip
将图像镜像翻转。

#### 参数
| 参数名         | 类型              | 说明 |
|:------------|:----------------| :--- |
| flipType    | `imgur.FlipType` | 水平或垂直镜像翻转。水平镜像翻转使用 `imgur.Horizontal` ， 垂直镜像翻转使用 `imgur.Vertical` 。


#### 返回值
`*imgo.Image` 类型的实例。

#### 例子

```go
// 水平镜像翻转
package main

import (
	"gateway/lib/imgur"
)

func main() {
	imgur.Load("workingID.png").
        Flip(imgur.Horizontal).
        Save("out.png")
}
```

```go
// 垂直镜像翻转
package main

import (
	"gateway/lib/imgur"
)

func main() {
	imgur.Load("workingID.png").
		Flip(imgur.Vertical).
		Save("out.png")
}
```

### GaussianBlur
高斯模糊。

#### 参数
| 参数名   | 类型        | 说明      |
|:------|:----------|:--------|
| ksize | `int`     | 高斯卷积核尺寸 
| sigma | `float64` | 高斯卷积核标准差

#### 返回值
`*imgo.Image` 类型的实例。

#### 例子

```go
package main

import (
	"gateway/lib/imgur"
)

func main() {
	imgur.Load("workingID.png").
		GaussianBlur(5, 10).
		Save("out.png")
}
```

### Grayscale
将图像转为 8 位颜色的灰度图。

#### 参数
无

#### 返回值
`*imgo.Image` 类型的实例。

#### 例子

```go
package main

import (
	"gateway/lib/imgur"
)

func main() {
	imgur.Load("workingID.png").
		Grayscale().
		Save("out.png")
}
```
转为灰度图后仍然可以继续编辑
```go
package main

import (
	"gateway/lib/imgur"
)

func main() {
	imgur.Load("workingID.png").
		Grayscale().
		Resize(100, 0).
		Save("out.png")
}
```
将转为灰度图的图像插入到其他图像上，并不会影响其他图像的色彩空间。下例输出的仍为 RGBA 图像。
```go
package main

import (
	"gateway/lib/imgur"
    "golang.org/x/image/colornames"
)

func main() {
    gopher := imgur.Load("9961593.png").Grayscale()
	imgur.Canvas(300, 300, colornames.Blueviolet).
        Insert(gopher, 50, 50).
        Save("out.png")
}
```

### Height
获取当前图像的高度。

#### 参数
无

#### 返回值
`int`

#### 例子

```go
package main

import (
	"fmt"
	"gateway/lib/imgur"
)

func main() {
	fmt.Println(imgur.Load("workingID.png").Height())
}
```

#### 输出

```
1448
```

### HttpHandler
作为 HTTP 响应。

#### 例子

```go
package main

import (
	"gateway/lib/imgur"
    "net/http"
)

func main() {
    http.HandleFunc("/imgur", imgur.Load("gopher.png").HttpHandler)
	http.ListenAndServe(":9927", nil)
}
```

### Insert
在已有图像上插入一个图像。

#### 参数
| 参数名    | 类型                | 说明             |
|:-------|:------------------|:---------------|
| source | `interface{}`     | 待插入的图像         |
| x      | `int`             | 待插入图像左上角坐标 x 轴 | 
| y      | `int`             | 待插入图像左上角坐标 y 轴 |

`source` 参数支持的类型    

| 支持类型                          | 数据类型                |
|:------------------------------|:--------------------|
| 文件系统中图像的绝对路径或相对路径             | `string` 或 `[]byte` | 
| 图片的 URL                       | `string` 或 `[]byte` |
| Base64 编码的图像数据                | `string` 或 `[]byte` |
| `*os.File` 类型的文件实例            | `*os.File`          | 
| 实现了 `image.Image` 接口的类型实例     | `image.Image` 以及实现了 `image.Image` 接口的类型 |
| `*imgo.Image` 类型的实例           |  `*imgo.Image`  |

#### 返回值
`*imgo.Image` 类型的实例。

#### 例子
##### 图像文件路径
```go
package main

import (
	"gateway/lib/imgur"
	colors "golang.org/x/image/colornames"
)

func main() {
	imgur.Canvas(500, 500, colors.Deepskyblue).
        Insert("9961593.png", 100, 100).
        Save("out.png")
}
```
##### 图像 URL
```go
package main

import (
	"gateway/lib/imgur"
	colors "golang.org/x/image/colornames"
)

func main() {
    url := "https://www.baidu.com/img/flexible/logo/pc/result.png"
	imgur.Canvas(500, 500, colors.Deepskyblue).
        Insert(url, 100, 100).
        Save("out.png")
}
```
##### Base64 编码的图像数据
```go
package main

import (
	"gateway/lib/imgur"
	colors "golang.org/x/image/colornames"
)

func main() {
    base64Img := imgur.Load("9961593.png").ToBase64()
	imgur.Canvas(500, 500, colors.Deepskyblue).
        Insert(base64Img, 50, 50).
        Save("out.png")
}
```
##### `*os.File` 类型的文件实例
```go
package main

import (
	"gateway/lib/imgur"
	colors "golang.org/x/image/colornames"
    "os"
)

func main() {
    file, err := os.Open("9961593.png")
    if err != nil {
        panic(err)
    }
    defer file.Close()

	imgur.Canvas(500, 500, colors.Deepskyblue).
        Insert(file, 100, 100).
        Save("out.png")
}
```
##### 实现了 `image.Image` 接口的类型实例
```go
package main

import (
	"gateway/lib/imgur"
	colors "golang.org/x/image/colornames"
    "image/png"
    "os"
)

func main() {
    file, err := os.Open("9961593.png")
    if err != nil {
        panic(err)
    }
    defer file.Close()
    
    img, err := png.Decode(file)
    if err != nil {
        panic(err)
    }

	imgur.Canvas(500, 500, colors.Deepskyblue).
        Insert(img, 100, 100).
        Save("out.png")
}
```
##### `*imgo.Image` 类型的实例
```go
package main

import (
	"gateway/lib/imgur"
    colors "golang.org/x/image/colornames"
)

func main() {
    img := imgur.Load("9961593.png")
	imgur.Canvas(500, 500, colors.Deepskyblue).
        Insert(img, 100, 100).
        Save("out.png")
}
```

### Load
加载图像。

#### 参数
| 参数名    | 类型                | 说明             |
|:-------|:------------------|:---------------|
| source | `interface{}`     | 待插入的图像         |

`source` 参数支持的类型

| 支持类型                          | 数据类型                |
|:------------------------------|:--------------------|
| 文件系统中图像的绝对路径或相对路径             | `string` 或 `[]byte` | 
| 图片的 URL                       | `string` 或 `[]byte` |
| Base64 编码的图像数据                | `string` 或 `[]byte` |
| `*os.File` 类型的文件实例            | `*os.File`          | 
| 实现了 `image.Image` 接口的类型实例     | `image.Image` 以及实现了 `image.Image` 接口的类型 |
| `*imgo.Image` 类型的实例           |  `*imgo.Image`  |

#### 返回值
`*imgo.Image` 类型的实例。

#### 例子
##### 图像文件路径
```go
package main

import (
	"gateway/lib/imgur"
)

func main() {
	imgur.Load("workingID.png").
        Save("out.png")
}
```
##### 图像 URL
```go
package main

import (
	"gateway/lib/imgur"
)

func main() {
    url := "https://www.baidu.com/img/flexible/logo/pc/result.png"
	imgur.Load(url).
        Save("out.png")
}
```
##### Base64 编码的图像数据
```go
package main

import (
	"gateway/lib/imgur"
	colors "golang.org/x/image/colornames"
)

func main() {
    base64Img := imgur.Load("workingID.png").ToBase64()
	imgur.Load(base64Img).
        Save("out.png")
}
```
##### `*os.File` 类型的文件实例
```go
package main

import (
	"gateway/lib/imgur"
	colors "golang.org/x/image/colornames"
    "os"
)

func main() {
    file, err := os.Open("workingID.png")
    if err != nil {
        panic(err)
    }
    defer file.Close()

	imgur.Load(file).
        Save("out.png")
}
```
##### 实现了 `image.Image` 接口的类型实例
```go
package main

import (
	"gateway/lib/imgur"
    "image/png"
    "os"
)

func main() {
    file, err := os.Open("workingID.png")
    if err != nil {
        panic(err)
    }
    defer file.Close()
    
    img, err := png.Decode(file)
    if err != nil {
        panic(err)
    }

	imgur.Load(img).
        Save("out.png")
}
```
##### `*imgo.Image` 类型的实例
```go
package main

import (
	"gateway/lib/imgur"
)

func main() {
    img := imgur.Load("workingID.png")
	imgur.Load(img).
        Save("out.png")
}
```

### LoadFromBase64
通过图像的 Base64 编码字符串读取图像。

#### 参数
| 参数名    | 类型        | 说明               |
|:-------|:----------|:-----------------|
| base64Str | `string`  | Base64 编码的图像字符串  |

#### 返回值
`*imgo.Image` 类型的实例。

#### 例子

```go
package main

import (
	"gateway/lib/imgur"
)

func main() {
    base64Img := imgur.Load("workingID.png").ToBase64()
	imgur.Load(base64Img).
		Save("out.png")
}
```

### LoadFromFile
从文件实例加载图像。推荐使用 `Load` 方法加载图像。

#### 参数
| 参数名    | 类型          | 说明                  |
|:-------|:------------|:--------------------|
| file | `*os.File`  | `*os.File` 类型的文件实例  |

#### 返回值
`*imgo.Image` 类型的实例。

#### 例子

```go
package main

import (
	"gateway/lib/imgur"
    "os"
)

func main() {
    file, err := os.Open("workingID.png")
    if err != nil {
        panic(err)
    }
    defer file.Close()

	imgur.LoadFromFile(file).
        Save("out.png")
}
```

### LoadFromImage
从 `image.Image` 实例加载图像。推荐使用 `Load` 方法加载图像。

#### 参数
| 参数名    | 类型          | 说明                        |
|:-------|:------------|:--------------------------|
| img | `image.Image`  | 实现了 `image.Image` 接口的类型实例 |

#### 返回值
`*imgo.Image` 类型的实例。

#### 例子

```go
package main

import (
	"gateway/lib/imgur"
    "image/png"
    "os"
)

func main() {
    file, err := os.Open("workingID.png")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    img, err := png.Decode(file)
    if err != nil {
        panic(err)
    }

	imgur.LoadFromImage(img).
        Save("out.png")
}
```

### LoadFromImgo
从 `*imgo.Image` 实例加载图像。推荐使用 `Load` 方法加载图像。

#### 参数
| 参数名  | 类型          | 说明                        |
|:-----|:------------|:--------------------------|
| i    | `*imgo.Image`  | `*imgo.Image` 实例 |

#### 返回值
`*imgo.Image` 类型的实例。

#### 例子

```go
package main

import (
	"gateway/lib/imgur"
)

func main() {
    img := imgur.Load("workingID.png")
	imgur.LoadFromImgo(img).
        Save("out.png")
}
```

### LoadFromPath
从文件路径加载图像。推荐使用 `Load` 方法加载图像。

#### 参数
| 参数名  | 类型       | 说明     |
|:-----|:---------|:-------|
| path | `string` | 图像文件路径 |

#### 返回值
`*imgo.Image` 类型的实例。

#### 例子

```go
package main

import (
	"gateway/lib/imgur"
)

func main() {
	imgur.LoadFromPath("workingID.png").
        Save("out.png")
}
```

### LoadFromUrl
从 URL 加载图像。推荐使用 `Load` 方法加载图像。

#### 参数
| 参数名 | 类型       | 说明     |
|:----|:---------|:-------|
| url | `string` | 图像 URL |

#### 返回值
`*imgo.Image` 类型的实例。

#### 例子

```go
package main

import (
	"gateway/lib/imgur"
)

func main() {
    url := "https://www.baidu.com/img/flexible/logo/pc/result.png"
	imgur.LoadFromUrl(url).
        Save("out.png")
}
```

### Line
在图像中绘制直线。

#### 参数
| 参数名 | 类型       | 说明     |
|:----|:---------|:-------|
| x1  | `int` | 直线端点坐标 x 轴|
| y1  | `int`  | 直线端点坐标 y 轴 |
| x2  | `int`  | 直线另一个端点坐标 x 轴 |
| y2  | `int`  | 直线另一个端点坐标 y 轴 |
| c   | `color.Color` | 直线颜色 |

#### 返回值
`*imgo.Image` 类型的实例。

#### 例子

```go
package main

import (
	"gateway/lib/imgur"
	"golang.org/x/image/colornames"
)

func main() {
	imgur.Canvas(300, 300, colornames.White).
		Line(10, 10, 200, 250, colornames.Black).
		Save("out.png")
}
```


### MainColor
返回图像的主要颜色。

#### 参数
无

#### 返回值
`color.RGBA` 类型的颜色

#### 例子

```go
package main

import (
    "fmt"
	"gateway/lib/imgur"
)

func main() {
    r, g, b, a := imgur.Load("workingID.png").
        MainColor().
        RGBA()
    fmt.Println(r, g, b, a)
}
```

### Mimetype
返回当前图像的 Mimetype。

#### 参数

#### 返回值

#### 例子

```go
package main

import (
    "fmt"
	"gateway/lib/imgur"
)

func main() {
    img := imgur.Load("gopher.png")
    fmt.Println(img.Mimetype())
}
```

#### 输出
```
image/png
```

### Mosaic
给图像打马赛克。

#### 参数
| 参数名 | 类型       | 说明            |
|:----|:---------|:--------------|
| size | `int` | 马赛克像素大小       |
| x1  | `int` | 马赛克区域左上角坐标 x 轴    |
| y1  | `int`  | 马赛克区域左上角坐标 y 轴    |
| x2  | `int`  | 马赛克区域右下角坐标 x 轴 |
| y2  | `int`  | 马赛克区域右下角坐标 y 轴 |

#### 返回值
`*imgo.Image` 类型的实例。

#### 例子

```go
package main

import "gateway/lib/imgur"

func main() {
	imgur.Load("workingID.png").
        Mosaic(5, 60, 50, 120, 100).
        Save("out.png")
}
```

### PickColor
返回图像指定坐标的颜色。

#### 参数
| 参数名  | 类型       | 说明      |
|:-----|:---------|:--------|
| x    | `int` | 坐标 x 轴  |
| y    | `int`  | 坐标 y 轴  |

#### 返回值
`*imgo.Image` 类型的实例。

#### 例子
```go
package main

import (
    "fmt"
	"gateway/lib/imgur"
)

func main() {
    r, g, b, a := imgur.Load("workingID.png").
        PickColor(50, 50).
        RGBA()
    fmt.Println(r, g, b, a)
}
```

### Pixel
在图像中绘制像素点。

#### 参数
| 参数名 | 类型       | 说明        |
|:----|:---------|:----------|
| x   | `int` | 像素点坐标 x 轴 |
| y   | `int`  | 像素点坐标 y 轴 |
| c   | `color.Color`  | 像素点颜色 |

#### 返回值
`*imgo.Image` 类型的实例。

#### 例子

```go
package main

import (
	"gateway/lib/imgur"
	"golang.org/x/image/colornames"
)

func main() {
    imgur.Canvas(300, 300, colornames.White).
		Pixel(100, 100, colornames.Black).
		Save("out.png")
}
```

### Pixelate
将图像像素化。

#### 参数
| 参数名 | 类型       | 说明        |
|:----|:---------|:----------|
| size | `int` | 像素大小       |

#### 返回值
`*imgo.Image` 类型的实例。

#### 例子

```go
package main

import (
	"gateway/lib/imgur"
)

func main() {
	imgur.Load("workingID.png").
        Pixelate(5).
        Save("out.png")
}
```

### RadiusBorder
给图像绘制圆角。

#### 参数
| 参数名 | 类型       | 说明        |
|:----|:---------|:----------|
| radius | `int` | 圆角半径       |

#### 返回值
`*imgo.Image` 类型的实例。

#### 例子

```go
package main

import (
	"gateway/lib/imgur"
)

func main() {
	imgur.Canvas(300, 300, color.White).
		BorderRadius(30).
		Save("out.png")
}
```

### Rectangle
在图像中绘制矩形。

#### 参数
| 参数名    | 类型            | 说明          |
|:-------|:--------------|:------------|
| x      | `int`         | 矩形左上角坐标 x 轴 |
| y      | `int`         | 矩形左上角坐标 y 轴 |
| width  | `int`         | 矩形宽度        |
| height | `int`         | 矩形高度        |
| c      | `color.Color` | 直线颜色        |

#### 返回值
`*imgo.Image` 类型的实例。

#### 例子

```go
package main

import (
	"gateway/lib/imgur"
	"golang.org/x/image/colornames"
)

func main() {
	imgur.Canvas(300, 300, colornames.White).
        Rectangle(100, 100, 100, 150, colornames.Black).
        Save("out.png")
}
```

### Resize
重置图像的尺寸。

#### 参数
| 参数名    | 类型    | 说明           |
|:-------|:------|:-------------|
| width  | `int` | 图像目标宽度       |
| height | `int` | 图像目标高度       |

`width` 和 `height` 参数不允许同时为 0，但允许其中一个为 0，意思为保持原图像的比例，按照非 0 参数通过原图像的比例计算出另一个参数。

#### 返回值
`*imgo.Image` 类型的实例。

#### 例子
```go
package main

import (
	"gateway/lib/imgur"
)

func main() {
	imgur.Load("9961593.png").
        Resize(100, 0).
        Save("out.png")
}
```

### Rotate
顺时针旋转图像。图像旋转之后图像的尺寸可能会发生变化。

#### 参数
| 参数名    | 类型    | 说明     |
|:-------|:------|:-------|
| angle  | `int` | 旋转的角度 |

#### 返回值
`*imgo.Image` 类型的实例。

#### 例子

```go
package main

import (
	"gateway/lib/imgur"
)

func main() {
	imgur.Load("9961593.png").
        Rotate(90).
        Save("out.png")
}
```

### Save
保存图像为文件。

#### 参数
| 参数名    | 类型    | 说明                                                                          |
|:-------|:------|:----------------------------------------------------------------------------|
| path  | `string` | 图像保存路径                                                                      |
| quality | `int` | 图像的质量。可选参数。<br/>只有当 path 中图片文件扩展为 jpg 或 jpeg 时，该参数才生效。<br/>有效范围为 (0, 100)，默认 100 |

支持的输出文件格式：
- png
- jpeg
- bmp
- tiff

#### 返回值
`*imgo.Image` 类型的实例。

#### 例子
##### PNG
```go
package main

import (
	"gateway/lib/imgur"
)

func main() {
	imgur.Load("workingID.png").
        Save("out.png") // 质量为 50
}
```
##### JPEG
```go
package main

import (
	"gateway/lib/imgur"
)

func main() {
	imgur.Load("workingID.png").
        Save("out.jpg", 50) // 质量为 50
}
```

### String
返回当前图像的字符串。

#### 参数
无

#### 返回值
`string`

#### 例子

```go
package main

import (
    "fmt"
	"gateway/lib/imgur"
)

func main() {
    img := imgur.Load("9961593.png")
    fmt.Println(img)
}
```

### Text
在图像中插入文字。

#### 参数
| 参数名   | 类型 | 说明     |
|:------|:--- |:-------|
| label | `string` | 文字内容   |
| x     | `int` | 插入文字的左上角坐标的 x 轴 |
| y     | `int` | 插入文字的左上角坐标的 y 轴 |
| fontPath | `string` | 字体文件路径 |
| fontColor | `color.Color` | 文字颜色   |
| fontSize | `float64` | 文字大小 |
| dpi | `float64` | 文字 DPI |

#### 返回值
`*imgo.Image` 类型的实例。

#### 例子
```go
package main

import (
	"gateway/lib/imgur"
    "golang.org/x/image/colornames"
)

func main() {
    fontPath := "/System/Library/Fonts/Supplemental/Arial.ttf"
	imgur.Canvas(500, 500, colornames.Deepskyblue).
        Text("Hello World", 50, 50, fontPath, colornames.Chocolate, 50, 100).
        Save("out.png")
}

```

### TextWordWrap
在图像中插入文字。可以换行，可以指定位置。 
> 文字换行暂试用 “ ” ，请看下面的演示代码。
> > ``str := "比如这是一段 StringWrapped WordWrap AlignLeft的文字"`` 这样换行
#### 参数
| 参数名    | 类型 | 说明 |
|:-------|:--- | :--- |
| label | `string` | 文字内容   |
| x     | `float64` | 插入文字的左上角坐标的 x 轴 |
| y     | `float64` | 插入文字的左上角坐标的 y 轴 |
| fontPath | `string` | 字体文件路径 |
| fontColor | `color.Color` | 文字颜色   |
| fontSize | `float64` | 文字大小 |
| align | `string` | 文字整体位置 |

`align`参数，包含3个位置：
- `AlignLeft` -- 左上顶点
- `AlignCenter` -- 上中位置
- `AlignRight` -- 右上顶点

#### 返回值
`*imgo.Image` 类型的实例。

#### 例子
```go
package main

import (
	"gateway/lib/imgur"
	"image/color"
	"time"
)

var fontPath = "/System/Library/Fonts/Supplemental/Arial Unicode.ttf"
var layout = "20060102150405"

func main() {
	alignLeftStr()
	//alignRightStr()
	//alignCenterStr()
}

func alignLeftStr() {
	str := "比如这是一段 StringWrapped WordWrap AlignLeft的文字"
	now := time.Now()
	outFileName := "./out-pic/out-" + now.Format(layout) + ".png"
	img := imgur.Load("./src_pic/workingID.png")
	img.TextWordWrap(str, 20, 240, fontPath, color.White, 80, "AlignLeft").Save(outFileName)
}

func alignRightStr() {
	str := "比如这是一段 StringWrapped WordWrap AlignRight的文字"
	now := time.Now()
	outFileName := "./out-pic/out-" + now.Format(layout) + ".png"
	img := imgur.Load("./src_pic/workingID.png")
	img.TextWordWrap(str, 60, 240, fontPath, color.White, 80, "AlignRight").Save(outFileName)
}

func alignCenterStr() {
	str := "比如这是一段 StringWrapped WordWrap AlignCenter的文字"
	now := time.Now()
	outFileName := "./out-pic/out-" + now.Format(layout) + ".png"
	img := imgur.Load("./src_pic/workingID.png")
	w := img.Width()
	img.TextWordWrap(str, float64(w/2), 400, fontPath, color.White, 80, "AlignCenter").Save(outFileName)
}

```

### Thumbnail
生成略缩图。

#### 参数
| 参数名    | 类型 | 说明 |
|:-------|:--- | :--- |
| width  | `int` | 略缩图宽度 |
| height | `int` | 略缩图高度 |

#### 返回值
`*imgo.Image` 类型的实例。

#### 例子

```go
package main

import (
    "gateway/lib/imgur"
)

func main() {
    imgur.Load("9961593.png").
		Thumbnail(80, 80).
        Save("out.png")
}
```

### ToBase64
返回图像的 PNG 格式的 Base64 编码字符串。

#### 参数
无

#### 返回值
`string`

#### 例子

```go
package main

import "gateway/lib/imgur"

func main() {
	Base64Img := imgur.Load("workingID.png").ToBase64()
}
```

#### 输出

`base64Img` 字符串为 `data:image/png;base64,......`


### ToImage
返回当前图像的 `image.Image` 类型的实例。

#### 参数
无

#### 返回值
`image.Image` 类型的实例。

#### 例子

```go
package main

import (
	"fmt"
	"gateway/lib/imgur"
)

func main() {
	img := imgur.Load("workingID.png").ToImage()
	fmt.Println(img.Bounds())
}
```

#### 输出

```
(0,0)-(750,1448)
```

### Width
返回当前图像的宽度。

#### 参数
无

#### 返回值
`int`

#### 例子

```go
package main

import (
	"fmt"
	"gateway/lib/imgur"
)

func main() {
	fmt.Println(imgur.Load("workingID.png").Width())
}
```

#### 输出

```
750
```


