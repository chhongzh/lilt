# lilt PV

用 [Remotion](https://www.remotion.dev) 做的宣传片，1920×1080，60fps，约 52 秒。

片里 OLED 上的画面都是 gallery 真跑出来的：`capture` 在主机上按脚本按键，60fps 逐帧导出 1-bit 画面和按键时刻，Remotion 再把每个像素放大画出来。片中其余的动效用的是 `src/lilt/spring.ts`，那是 `anim` 包弹簧的移植：同样的 Response / Damping、同样的 1/240 秒子步长、同样的预设，转场照搬 `SlideLeft` 和 `Sheet`，光标照搬 `ui.Cursor` 的前沿 / 后沿。

## 素材

`public/footage` 里的文件由 `capture` 生成，已经提交，改了 gallery 或脚本再重录：

```sh
npm run capture
```

| 文件 | 内容 |
| --- | --- |
| `reel` | 主片：光标、按下与松手、按住确认、按钮状态、Sheet、Chips、横向滚动、软键盘、长列表 |
| `compare` / `compare-stiff` | 同一串按键录两遍，后者 `anim.TimeScale = 40`，弹簧一两帧就落定 |
| `slow` | `anim.TimeScale = 0.25` 录的光标慢放，按键也放慢 4 倍 |
| `logo.json` | 用设备上的 Notosans 12pt 画出的 `lilt` 像素 |

`capture` 按墙上时钟跑（`OnAction` 在后台 sleep），所以录一遍要花和片长一样的时间。

## 预览与渲染

需要 Node 22。字幕字体是系统里的 Noto Sans SC，没装的话会退回微软雅黑。

```sh
npm install
npm run studio
npm run render
```

成片输出到 `out/lilt-pv.mp4`。下载 chrome-headless-shell 太慢的话，可以用本机的 Chrome：

```sh
REMOTION_CHROME="C:/Program Files/Google/Chrome/Application/chrome.exe" npm run render
```

声音只有按键声（`public/sfx`，ffmpeg 合成），对齐录制时真实的按下和松开。背景音乐没有加。
