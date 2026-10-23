# 给模型的任务书：更新 lilt 文档站

你要更新 lilt 的文档站（MkDocs + Material，源文件在 `docs/`，发布到 <https://chhongzh.github.io/lilt/>），让它和当前代码一致、完整、讲得清楚。这份文件是完整的任务说明，按顺序做即可。

lilt 是给 TinyGo 用的 GUI 框架，模块路径 `github.com/chhongzh/lilt`。读者是要在单片机上写界面的 Go 程序员：懂 Go，不一定懂 GUI 框架。

---

## 1. 你要交付什么

1. `docs/` 下的 Markdown 页面，覆盖所有公开 API、核心概念和用法，并带例子讲解。
2. 每个组件和每个会动的概念至少一段**动图演示**，由 `cmd/docgif` 用 silent 驱动录制，页面里的示例代码就是录动图的那份代码。
3. `mkdocs.yml` 的 `nav` 与页面一致。
4. `scripts/docs.sh check` 通过（严格模式，链接、锚点、代码片段有任何问题都会失败）。
5. 最后一条消息里的变更报告（格式见第 9 节）。

只改这些路径：`docs/`、`mkdocs.yml`、`cmd/docgif/`。框架代码、README、`driver/sh1106` 都不要动。发现代码有 bug 或者文档注释和行为不符，写进报告，不要顺手修。

---

## 2. 仓库地图

| 路径 | 内容 | 对应页面 |
| --- | --- | --- |
| `lilt.go` | `Start`、`Pause`、`Stop`、`FPS` | `getting-started.md`、`reference/lilt.md` |
| `ui/` | 节点树、布局、层栈、焦点、输入、光标、主题 | `concepts/*`、`reference/ui.md` |
| `widget/` | 标准组件 | `widgets/*`、`reference/widget.md` |
| `anim/` | 弹簧 `Spring`、动画值 `Value`、调度器 | `concepts/animation.md`、`reference/anim.md` |
| `gfx/` | `Screen` 接口、`Canvas`、字体 | `guides/custom-widget.md`、`guides/hardware.md`、`reference/gfx.md` |
| `driver/silent`、`driver/desktop` | 无头屏幕、桌面窗口与 wasm | `guides/drivers.md` |
| `examples/gallery` | 演示所有组件的 gallery | 写例子时的参考 |
| `cmd/docgif` | 文档动图的录制程序和全部 demo | 所有页面的示例代码 |
| `internal/reel` | 按键脚本回放与录制 | 只在写 demo 时用 |
| `DESIGN.md`、`README.md`、`CLAUDE.md` | 设计需求、概览、框架约定 | 背景资料，**以代码为准** |

现有页面结构见 `mkdocs.yml` 的 `nav`。`docs/widgets/button.md` 是写好的范例，新页面照它的结构和语气写。

---

## 3. 工作流程

### 3.1 弄清现在的代码

1. 看自上次更新文档以来代码改了什么：

    ```bash
    git log --oneline -1 -- docs                     # 上次改文档的提交
    git log --oneline <那个提交>..HEAD -- . ':!docs'  # 之后代码的提交
    git diff <那个提交>..HEAD --stat -- . ':!docs'
    ```

    第一次写或者改动很大时，把所有包都过一遍。

2. 列出每个包的公开 API，这是覆盖率的依据：

    ```bash
    go doc -all . ; go doc -all ./ui ; go doc -all ./widget ; go doc -all ./anim ; go doc -all ./gfx
    go doc -all ./driver/silent ; go doc -all ./driver/desktop
    ```

3. 读源码，不要只读文档注释。行为细节（时间常数、默认值、哪个方法会触发回调、谁在哪个 goroutine 上跑）都从代码里确认。例如 `widget/control.go` 的 `loadingDelay`、`widget/button.go` 的 `resultHold`、`ui/theme.go` 的默认主题。

### 3.2 定下要改的页面

对照 API 列表和现有页面，列出：

- 新增的公开类型、函数、方法、选项 → 写到哪一页；
- 改了签名或行为的 → 哪些页面、哪些 demo 要跟着改；
- 删掉的 → 从页面、`nav` 和 demo 里删干净；
- 没讲清楚、没例子、没动图的老内容。

新组件放 `docs/widgets/<名字>.md`，新概念放 `docs/concepts/`，然后加进 `mkdocs.yml` 的 `nav`，`widgets/index.md` 或 `concepts/index.md` 的总览里也要加一行。

### 3.3 写 demo（第 5 节有完整规则）

每页一个 `cmd/docgif/demo_<页面名>.go`。先写 demo、录出来、看过画面，再写页面，这样文字描述的就是读者真正会看到的东西。

### 3.4 录制并检查画面

```bash
scripts/docs.sh shots 'button-*'     # 只录匹配的 demo，静帧写到 .docshots/
```

每个 demo 的 `Shot` 静帧会拼成 `.docshots/<名字>.png`。**逐张打开看**：动作是不是按预期发生，静帧有没有截在动画中途，文字有没有被截断，弹窗有没有挡住要展示的东西。不对就调脚本里的等待时间，重录。GIF 本身在 `docs/assets/demos/`。

### 3.5 写页面（第 4、6 节）

### 3.6 构建

```bash
scripts/docs.sh check    # go vet demo，补录缺失的 GIF，mkdocs build --strict
scripts/docs.sh serve    # 需要时本地预览，http://127.0.0.1:8000/lilt/
```

`check` 必须零警告通过。常见失败：链接写成了站点绝对路径、锚点不存在、`--8<--` 指向的片段名拼错、页面不在 `nav` 里。

MkDocs 通过 `uv` 运行（没有 `uv` 时先 `pip install -r docs/requirements.txt`）。GIF 和 `site/` 都不提交，CI 会重新录制。

### 3.7 收尾

对照第 8 节的清单逐项确认，然后写报告。

---

## 4. 页面怎么写

### 4.1 组件页（`docs/widgets/*.md`）

照 `docs/widgets/button.md`：

1. **一级标题**：类型名。
2. **开头一两句**：它是什么、什么时候用、怎么触发。
3. **按用法分节**，每节：一句话说明 → 动图 → 代码片段 → 解释代码里不显然的地方（为什么用 `ui.Do`、为什么先声明变量等）。
4. **状态与反馈**：加载、失败、回滚等行为，有确定数字的用表格。
5. **`## API` 一节**：构造函数和公开方法的签名代码块（每行注释一句中文），再加一张“选项 → 作用”表。说明默认样式取自主题的哪些字段。
6. 和其他页面相关的地方用相对链接串起来。

### 4.2 概念页（`docs/concepts/*.md`）

讲清楚“是什么、为什么这样设计、怎么用、容易踩的坑”：

- 先给一个最小例子，再展开规则。
- 有流程或层次的东西画 Mermaid 图（帧循环顺序、按键分发路径、层栈）。
- 有视觉效果的概念要有动图：布局变化、弹簧手感对比、焦点移动、页面推入推出。
- 坑放在 `!!! warning` 里，例如 `ui.Do` 不能嵌套、回调里直接改树、Set 不触发 OnChange。

### 4.3 API 速查（`docs/reference/*.md`）

每个包一页。按类别分组（构造、布局选项、树操作、输入……），每组一张表：`签名 | 说明 | 详见`，说明用一句中文，“详见”链接到讲它的概念页或组件页的锚点。页首给出 `pkg.go.dev` 链接：`https://pkg.go.dev/github.com/chhongzh/lilt/<包>`。所有导出的标识符都要出现在速查里。

### 4.4 Markdown 约定

- 动图：`![描述](../assets/demos/<demo 名>.gif){ .demo }`，路径相对当前页面。
- 代码片段：

    ````markdown
    ```go
    --8<-- "cmd/docgif/demo_button.go:button-click"
    ```
    ````

    路径相对仓库根目录。片段原样引入，缩进会自动去掉。
- 只为说明 API 形状、不需要运行的代码（签名列表、板子上的 `main`）可以直接写在页面里，但要能编译：类型、函数名、参数个数和真实代码一致。
- 按键用 `++ok++`、`++back++`、`++menu++`、`++arrow-down++` 这种写法（前三个是 `mkdocs.yml` 里加的自定义键名）；框架里的键名（`ui.KeyOK`）写在代码格式里。
- 链接一律相对路径并带 `.md`：`[线程模型](../concepts/threading.md#control)`。
- 会被其他页面链接的标题加显式锚点：`## 激活 { #activation }`。改标题时保留锚点。
- 提示框只用 `!!! tip`、`!!! note`、`!!! warning`，一页不超过三四个。
- 表格、提示框、Mermaid、标签页（`=== "标签"`）都可以用，按内容需要选，不要为了排版而用。

---

## 5. demo 的规则

demo 在 `cmd/docgif/demo_<页面名>.go`，`package main`，在 `init` 里 `register`。两个页面同名时换一个不冲突的名字（`widgets/input.md` 的 demo 在 `demo_field.go`，因为 `demo_input.go` 属于 `concepts/input.md`），`grep -rho "cmd/docgif/demo_[a-z_]*\.go" docs | sort | uniq -c` 能看出每个文件被哪些页面引用。几个 demo 都要用的录制包装放进 `helpers.go`。

```go
func init() {
	register(demo{
		name:  "button-click",     // 文件名，<页面名>-<用法>，小写连字符
		build: buttonClick,         // 返回根节点，录制前调用一次
		script: func(s *reel.Script) {
			s.Wait(0.5)               // 开场留一点时间
			s.Emit(ui.KeyOK, ui.Press)
			s.Wait(0.3)
			s.Shot("pressed")         // 静帧，名字写这一刻发生了什么
			s.Emit(ui.KeyOK, ui.Release)
			s.Wait(0.5)
			s.Shot("tapped")
		},
	})
}

func buttonClick() ui.Node {
	// --8<-- [start:button-click]
	...读者要看的代码...
	// --8<-- [end:button-click]
	return page("Button", padded(tap))
}
```

- **片段标记**：`// --8<-- [start:名字]` 和 `// --8<-- [end:名字]` 之间的代码会原样出现在页面上，名字一般和 demo 名相同。只框住读者需要的部分；外面的 `page`、`padded` 是录制用的包装（定义在 `helpers.go`），不进片段。片段不能只给一半：读者照抄要能看懂每个变量从哪来。
- **脚本 API**（`internal/reel`）：`Wait(秒)`、`Emit(键, 阶段)`、`Tap(键, 之后等待)`、`Taps(键, 次数, 间隔)`、`Hold(键, 按住秒数)`、`Shot(名字)`。`Tap` 自带 0.08 秒的按下时间。
- **时长**：一个 demo 2 到 6 秒，只演示一件事。开头 `Wait(0.5)` 左右让画面先停一下；动作之间留够让动画落定的时间（一般 0.3 到 0.6 秒）。
- **回调必须在录制结束前返回**。回调在自己的 goroutine 上跑并按真实时间 sleep，脚本末尾要留够时间，否则它会作用到下一个 demo 上。
- **屏幕**默认 128x64，与 SH1106 相同。只有页面本身在讲尺寸或大屏时才设 `w`、`h`。
- **主题、弹簧预设**（包括 `ui.CursorLead` 等光标弹簧）**和 `ui.TextInput`** 每个 demo 结束后会自动恢复，所以演示改主题、调 `anim.Default` 都可以在 `build` 里直接改。组件在构造时读取主题，改主题要放在构造组件之前。
- demo 代码必须遵守框架约定（第 7 节），它就是读者会抄走的代码。用户可见的文字用英文（字体只有 ASCII 字形），简短、像真实产品里的文案，不要 `foo`、`test`。
- 新增或修改 demo 后：`go vet ./cmd/docgif`，`gofmt -l cmd/docgif`（工作区是 CRLF 签出时先去掉 `\r` 再看），再 `scripts/docs.sh shots <名字>` 看静帧。

---

## 6. 写作规范

- **语言**：简体中文。标识符、类型、函数名保持英文原样并放进代码格式。中英文之间不加空格也可以，但同一页保持一致；全角标点。
- **说具体的事**：写“回调超过 150ms 显示加载点”，不写“提供了优雅的加载体验”。不用“强大”“轻松”“无缝”“优雅”“丝滑”这类形容词，不写“让我们来看看”“总之”这类套话，不在结尾总结全文。
- **每段话都要有信息**：删掉复述标题的句子和读者已经知道的铺垫。
- **讲为什么**：设计上的取舍要讲原因，例如为什么整树重排而不做脏标记（`ui/layout.go` 的注释）、为什么回调不在 UI goroutine 上。原因从代码注释、`DESIGN.md` 里找，找不到就不写，不要编。
- **术语统一**：

    | 英文 | 中文 |
    | --- | --- |
    | node / tree | 节点 / 节点树 |
    | widget | 组件 |
    | option | 选项 |
    | layer stack | 层栈 |
    | page（`ui.Push`） | 页面 |
    | popup（`ui.Open`） | 弹窗 |
    | Drawer / Dialog / Toast | 抽屉 / 对话框 / Toast（不译） |
    | cursor | 光标 |
    | focus | 焦点 |
    | activate | 激活 |
    | handler / callback | 回调 |
    | unit | 单位 |
    | spring | 弹簧 |
    | marquee | 跑马灯 |
    | Press / Repeat / Release | 按下 / 连发 / 松开 |

- **不确定就去查代码**。不要描述代码里没有的功能，不要猜参数的默认值。

---

## 7. 框架约定（示例代码必须遵守）

- 可激活的组件在 `KeyOK` 的 `Release` 且未 `Canceled` 时触发；`Press` 只给反馈。
- 运行时不做每帧分配，不用 `strings`，只用 float32。属性一律走 Set/Get。
- `OnClick`、`OnChange` 的回调在自己的 goroutine 上跑。回调里改树必须包在 `ui.Do` 里；`ui.Do` 不能嵌套。
- 代码里调用 Set 改值不触发 `OnChange`。
- 树只属于 UI goroutine。`lilt.Start` 之前在任何地方建树都可以。
- Padding、Margin、Gap、Radius 的单位是“单位”，不是像素；共用的内边距从主题取（`RowPadding`、`ButtonPadding`、`Gap`、`PanelPadding`）。
- 弹窗的边框和内边距由层栈画、由层栈留，弹窗内容自己不要再加。
- 自定义组件不要声明和 `ui.Box`、`ui.Flex` 未导出方法同名的方法（`box`、`set`、`visible`、`split`、`join`、`margins`、`crossAuto`），TinyGo 0.42/0.43 会崩溃。讲自定义组件时要提醒这一点。
- `ui.Push` 的页面铺满整屏，锚点对页面根节点不起作用。

---

## 8. 收尾清单

- [ ] `go doc -all` 列出的每个导出标识符都出现在 `docs/reference/` 里，并且有一处讲解或链接。
- [ ] 每个组件页都有动图、片段和 `## API` 节。
- [ ] 每个片段标记名都被某个页面引用，每个 demo 的 GIF 都被某个页面使用；没有孤儿 demo。
- [ ] 删掉或改名的 API 在页面、片段、demo 里都清理干净了（`grep -rn 旧名字 docs cmd/docgif`）。
- [ ] `.docshots/` 里的静帧都看过，画面和文字描述一致。
- [ ] `gofmt -l cmd/docgif` 无输出，`go vet ./cmd/docgif` 通过。
- [ ] `scripts/docs.sh check` 通过。
- [ ] `mkdocs.yml` 的 `nav`、`widgets/index.md`、`concepts/index.md` 与页面一致。

---

## 9. 最后的报告

用中文，分四段，每段是列表：

1. **改了哪些页面**：每页一句话说新增或修改了什么。
2. **新增或修改的 demo**：名字和展示的内容。
3. **代码和文档对不上的地方**：代码行为、文档注释、README、DESIGN.md 之间的矛盾，或者疑似 bug，附文件和行号。没有就写“无”。
4. **没做完的**：做不到或拿不准的地方，以及原因。

不要提交 git，由维护者审阅后提交。
