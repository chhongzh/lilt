# 快速开始

这一页从零写出一个能跑的界面：先在桌面窗口里调好，再接到真实的屏幕和按键上。

## 准备

| 工具 | 版本 | 用途 |
| --- | --- | --- |
| Go | 1.27 及以上 | lilt 用到了泛型方法 |
| TinyGo | 0.42 及以上 | 编译到单片机 |

## 安装

在你的模块里运行：

```bash
go get github.com/chhongzh/lilt@latest
```

`@latest` 取最新的正式版本。想固定在某个版本，就把它换成版本号：

```bash
go get github.com/chhongzh/lilt@v0.1.0
```

所有已发布的版本，以及每个版本改了什么，都在 [GitHub Releases](https://github.com/chhongzh/lilt/releases) 里。

!!! warning "lilt 还在 0.x 阶段"

    lilt 迭代很快。1.0 之前，次版本号升级（例如 v0.1 到 v0.2）可能带来不兼容的 API 改动，补丁版本只修 bug。`go.mod` 会记下你用的确切版本，只有你自己运行 `go get` 才会升级；升级前先读 Releases 里的说明，并定期看看有没有新版本：

    ```bash
    go list -m -u github.com/chhongzh/lilt
    ```

    这条命令会显示当前用的版本，以及有没有更新的版本可以升级。

## 第一个界面

界面是一棵节点树。下面这棵树的根是一个 `List`，里面三行：一行点开新页面，两行是开关。

![第一个界面](assets/demos/start-hello.gif){ .demo }

```go
--8<-- "cmd/docgif/demo_start.go:start-hello"
```

几件值得注意的事：

- 组件都由 `NewXxx(参数, 选项...)` 创建，选项是 `ui.Option`，可以层层嵌套成一个表达式，见[节点树](concepts/tree.md)。
- `ui.Add` 把子节点挂到父节点下。`List` 会把它们放进内部那一列。
- `OnClick` 的回调**不在** UI goroutine 上执行，所以改树（这里是 `ui.Push` 一个新页面）要包在 `ui.Do` 里，见[线程模型](concepts/threading.md)。
- 页面总是铺满整个屏幕，所以文字居中用的是 `TextAlign(0.5)`。按 Back 关闭页面。

## 在桌面上跑

`driver/desktop` 开一个窗口当屏幕，键盘当按键，不需要 cgo：

```go
package main

import (
	"log"

	"github.com/chhongzh/lilt/driver/desktop"
)

func main() {
	// Run 占住主 goroutine，直到窗口关闭。
	if err := desktop.Run(desktop.Config{Title: "hello", Scale: 6}, hello()); err != nil {
		log.Fatal(err)
	}
}
```

++w++ ++a++ ++s++ ++d++ 或方向键移动，++space++ 是 OK，++j++ 是 Back，++k++ 是 Menu。仓库里自带的 gallery 可以这样跑：

```bash
go run ./cmd/sim
```

## 接到板子上

在板子上，入口是 `lilt.Start`。它绑定屏幕、把根节点作为第一屏，然后在自己的 goroutine 里跑帧循环；调用方接着读按键就行。

```go
package main

import (
	"machine"
	"time"

	"github.com/chhongzh/lilt"
	"github.com/chhongzh/lilt/gfx"
	"github.com/chhongzh/lilt/ui"
	"tinygo.org/x/drivers/ssd1306"
)

func main() {
	machine.I2C0.Configure(machine.I2CConfig{})
	dev := ssd1306.NewI2C(machine.I2C0)
	dev.Configure(ssd1306.Config{Width: 128, Height: 64})

	lilt.Start(gfx.FromDisplayer(dev), hello())

	pins := []struct {
		key ui.Key
		pin machine.Pin
	}{
		{ui.KeyUp, machine.GPIO1}, {ui.KeyDown, machine.GPIO2},
		{ui.KeyOK, machine.GPIO3}, {ui.KeyBack, machine.GPIO4},
	}
	for _, p := range pins {
		p.pin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}
	keys := ui.NewKeyPoller() // 消抖，按住时加速连发
	const dt = 5 * time.Millisecond
	for {
		for _, p := range pins {
			// 按键接地，按下时引脚为低电平。
			keys.Poll(p.key, !p.pin.Get(), float32(dt.Seconds()))
		}
		time.Sleep(dt)
	}
}
```

引脚号要换成你自己的接线。完整的、实际跑过的固件见 [examples/esp32s3-generic](https://github.com/chhongzh/lilt/tree/main/examples/esp32s3-generic)，接其他屏幕和按键的办法见[接入硬件](guides/hardware.md)。

```bash
tinygo build -target esp32s3-generic -o gallery.bin ./examples/esp32s3-generic
```

## 不开窗口也能看

`driver/silent` 是一块只画进内存的屏幕，可以导出 PNG 和 GIF，适合测试和出图。这个站点里的所有动图都是用它录的，见[驱动与模拟器](guides/drivers.md#silent)。

## 下一步

- [核心概念](concepts/index.md)：理解节点树、布局和线程模型，写复杂界面前值得读一遍。
- [组件](widgets/index.md)：每个标准组件的用法和动图。
