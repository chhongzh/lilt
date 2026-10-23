---
hide:
  - navigation
---

# lilt

给 TinyGo 用的 GUI 框架。手感参照 osu!framework 和 iOS：所有会动的东西都由弹簧驱动，按下有反馈，松手才触发。框架与硬件无关，屏幕和按键都是接口，单色、彩色屏都能接。

![gallery 演示](assets/highlight.gif){ .demo }

上面这段是 `go run ./cmd/preview -highlight` 在主机上录的，跑的是和固件同一份 [gallery](https://github.com/chhongzh/lilt/tree/main/examples/gallery) 代码。没有板子也能[在浏览器里玩一下](play.md)。

```go
list := widget.NewList()
ui.Add(list,
	widget.NewListItem("Open a page", widget.OnClick(func() {
		ui.Do(func() { ui.Push(widget.NewLabel("Hello, lilt", widget.TextAlign(0.5))) })
	})),
	widget.NewCheckBox("Wi-Fi", true),
)
lilt.Start(gfx.FromDisplayer(display), list)
```

<div class="grid cards" markdown>

-   :material-rocket-launch: **[快速开始](getting-started.md)**

    装好依赖，在桌面窗口里跑出第一个界面，再烧到板子上。

-   :material-file-tree: **[核心概念](concepts/index.md)**

    节点树、布局、线程模型、输入与焦点、弹簧动画、页面与弹窗。

-   :material-widgets: **[组件](widgets/index.md)**

    每个标准组件的用法、动图和 API。

-   :material-chip: **[接入硬件](guides/hardware.md)**

    实现 `gfx.Screen`、读按键，或者直接用 tinygo 的现成驱动。

</div>

## 特点

- **松手才触发**：按下只给反馈（光标收缩、往方向键那边探一点），松开且没被取消才执行。
- **回调不卡动画**：`OnClick`、`OnChange` 都在自己的 goroutine 上跑，慢的回调自动显示加载点，返回 `error` 自动抖动、回滚。
- **一个光标**：焦点由一个会拉伸、会跨屏跳跃的光标标出，能反色的屏上画成反色。
- **为单片机写**：每帧不分配内存，只用 float32，不依赖 `strings`。
- **主机上可预览**：`driver/silent` 把界面画进内存再导出 GIF，`driver/desktop` 开窗口或编成 wasm。这个站点里的动图都是这样录的。
