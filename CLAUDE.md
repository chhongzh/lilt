# Untitled UI Framework

一个`osu!framework` like的给tinygo用的gui框架，实现非常丝滑的视觉效果。

## 代码质量要求

1. 不编写无意义的Testing文件
2. 使用`anti-slop`控制代码质量
3. 及时提交git, 但是不要以claude身份署名

## 框架约定

- 所有可激活的组件都在 `KeyOK` 的 `Release` 且 `!Canceled` 时触发；`Press` 只认领按键、让光标收缩（iOS 同款）。
- 运行时不做每帧分配，不用 `strings`，只用 float32。属性一律走 Set/Get，变化后调用 `ui.Relayout()` 或 `ui.Invalidate()`。
- 主机上 `go build ./...` 会因 `driver/sh1106` 依赖 `machine` 失败，按包编译：`./anim ./gfx ./ui ./widget ./examples/gallery ./cmd/preview`。
- 改动 UI 后用 `go run ./cmd/preview -o /tmp/pv` 看 `shots.png` 验证，再跑 `tinygo build -target esp32s3-generic ./examples/esp32s3-sh1106`。
- TinyGo 0.42 的坑：嵌入 `ui.Box`（或 `ui.Flex`）的类型如果声明了和它们未导出方法同名的方法（`box`、`set`、`visible`、`split`、`join`、`margins`、`crossAuto`），转成接口时编译器会崩溃（`index out of range [0] with length 0`）。组件里避开这些名字。
- `driver/sh1106` 由用户之后重写，不要改动。
