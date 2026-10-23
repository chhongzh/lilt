# lilt

[pkg.go.dev](https://pkg.go.dev/github.com/chhongzh/lilt) · `import "github.com/chhongzh/lilt"`

启动、暂停、停止 UI goroutine。整个框架是单例：一块屏幕、一棵树、一个帧循环。

## 运行

| 签名 | 说明 | 详见 |
| --- | --- | --- |
| `func Start(s gfx.Screen, root ui.Node)` | 绑定屏幕，把 `root`（非 nil 时）作为第一屏，启动 UI goroutine；暂停时调用则恢复运行并忽略参数 | [快速开始](../getting-started.md)、[生命周期](../concepts/threading.md#lifecycle) |
| `func Pause()` | 停止更新和绘制，直到再次调用 `Start`；期间的按键在队列里等待，队列满后丢弃 | [生命周期](../concepts/threading.md#lifecycle)、[Recorder](../guides/drivers.md#recorder) |
| `func Stop() error` | 结束 UI goroutine 并丢弃所有页面，之后可以重新 `Start`；返回屏幕 `Flush` 失败导致的错误 | [生命周期](../concepts/threading.md#lifecycle) |
| `var FPS = 60` | 帧率上限，每帧读取一次，运行中可改 | [帧率与电子墨水屏](../guides/hardware.md#fps) |

`Start` 之前在任何 goroutine 上建树都可以；之后只有 UI goroutine 能碰树，其他 goroutine 用 [`ui.Do`](../concepts/threading.md#do)。
