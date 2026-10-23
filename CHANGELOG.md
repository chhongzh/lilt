# 更新日志

这里记录 lilt 每个版本的变化，格式参照 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.1.0/)，版本号遵循[语义化版本](https://semver.org/lang/zh-CN/)。1.0 之前，次版本号升级（例如 0.1 到 0.2）可能带来不兼容的改动，补丁版本只修 bug。发布流程见仓库里的 `RELEASING.md`。

## [Unreleased]

## [0.1.0] - 2026-10-23

第一个公开版本。

### 新增

- `lilt`：`Start`、`Pause`、`Stop` 管理 UI goroutine，`FPS` 限制帧率。
- `anim`：Response / Damping 写法的弹簧，预设 `Default`、`Smooth`、`Snappy`、`Bouncy`、`Interactive`；`Spring.Ease` 指数缓动，`Value.Bump` 给静止的值一个冲量。
- `gfx`：`Screen` 与 `Inverter` 接口，带裁剪和平移的 `Canvas`，字体；`FromDisplayer` 直接包装 `tinygo.org/x/drivers` 的驱动，驱动有 `FillRectangle` 时用它填充。
- `ui`：节点树，Flex、Grid、Stack 布局，页面和弹窗共用的层栈，焦点与光标，按下反馈、松手触发的输入模型，`KeyPoller` 消抖和连发，`ui.Do` 从其他 goroutine 改树，跟随主题字体缩放的间距单位。
- `widget`：Label（放不下时滚动）、List、Scroll（可嵌套，行间错开跟随）、Button（`HoldToConfirm`）、CheckBox、Radio、Select、Chip、Choice、Toggles、Input 与软键盘、Drawer、Dialog、Toast。`OnClick` 和 `OnChange` 在各自的 goroutine 上执行，超过 150ms 显示加载点，返回 error 时抖动并恢复旧值。
- `driver/desktop`：用 Ebitengine 开窗口的模拟器，也能编成 wasm 在浏览器里跑。
- `driver/silent`：只画进内存的单色屏幕，导出 PNG 和 GIF。
- `driver/sh1106`：示例固件用的 SH1106 驱动。之后会重写，API 不保证稳定。
- `examples/gallery` 组件演示，`examples/esp32s3-generic` ESP32-S3 + SH1106 固件。
- `cmd/preview` 按脚本录 GIF 和截图拼图，`cmd/sim` 在桌面或浏览器里跑 gallery。

### 已知限制

- 主机上 `go build ./...` 会因为 `driver/sh1106` 依赖 `machine` 而失败，请按包单独编译。
- 除了 SH1106，README 里列出的屏幕驱动都没有实测过。

[Unreleased]: https://github.com/chhongzh/lilt/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/chhongzh/lilt/releases/tag/v0.1.0
