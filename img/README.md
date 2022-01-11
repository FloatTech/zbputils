# imgfactory

> API 详见代码

### 处理中图像
```go
type ImgFactory struct {
	Im *image.NRGBA
	W  int
	H  int
}
```

### 加载图片
```go
func Load(path string) image.Image
```

### 载入图片第一帧作底图
```go
func LoadFirstFrame(path string, w, h int) *ImgFactory
```

### 加载图片每一帧图片
```go
func LoadAllFrames(path string, w, h int) []*image.NRGBA
```
