# 发布流程

这份文档写给维护者，说明怎么从 `dev` 发出一个新版本。

## 分支

- `dev`：日常开发都在这里，提交随时可能不稳定。
- `main`：只放已经发布的代码，每次发布从 `dev` 合并一次。文档站从 `main` 部署。
- 版本号是 `main` 上的注释 tag，格式 `vX.Y.Z`。Go 的 module proxy 和 pkg.go.dev 都按 tag 识别版本。

## 版本号

遵循[语义化版本](https://semver.org/lang/zh-CN/)。1.0 之前的约定：

| 改动 | 升哪一位 | 例子 |
| --- | --- | --- |
| 不兼容的 API 改动，或者新增功能 | 次版本 | `v0.1.3` → `v0.2.0` |
| 只修 bug，API 不变 | 补丁版本 | `v0.2.0` → `v0.2.1` |
| 想先让人试用 | 预发布后缀 | `v0.3.0-rc.1` |

预发布 tag 带 `-` 后缀，release workflow 会把它标成 pre-release，`go get ...@latest` 也不会选中它。

## 步骤

下面以 `v0.2.0` 为例，换成实际的版本号。

### 1. 在 dev 上准备

把 [CHANGELOG.md](CHANGELOG.md) 里 `[Unreleased]` 下的条目移到新的版本标题下，写上日期，再更新文件底部的两个对比链接：

```md
## [Unreleased]

## [0.2.0] - 2026-11-20

### 变更

- ...

[Unreleased]: https://github.com/chhongzh/lilt/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/chhongzh/lilt/compare/v0.1.0...v0.2.0
```

标题必须是 `## [0.2.0]` 这个格式。release workflow 按它截取发布说明，找不到就失败。

然后在本地跑一遍 CI 会跑的检查：

```bash
gofmt -l .
go mod tidy && git diff --exit-code go.mod go.sum
go vet . ./anim ./gfx ./ui ./widget ./driver/silent ./examples/gallery ./cmd/preview ./cmd/docgif ./internal/...
GOOS=js GOARCH=wasm go vet ./driver/desktop ./cmd/sim
tinygo build -target esp32s3-generic -o /tmp/gallery.bin ./examples/esp32s3-generic
scripts/docs.sh check
```

都通过后提交：

```bash
git commit -am "chore: release v0.2.0"
git push origin dev
```

等 `dev` 上的 CI 变绿再继续。

### 2. 合并到 main

用 `--no-ff` 保留一个合并提交，`main` 的历史里每个版本都是一个节点：

```bash
git switch main
git pull --ff-only origin main
git merge --no-ff dev -m "Merge branch 'dev' for v0.2.0"
```

### 3. 打 tag 并推送

```bash
git tag -a v0.2.0 -m "v0.2.0"
git push origin main
git push origin v0.2.0
```

推送 `main` 会重新部署文档站。推送 tag 会触发 [release workflow](.github/workflows/release.yml)，它依次：

1. 跑完整的 CI；
2. 从 CHANGELOG 截取 `[0.2.0]` 这一节作为发布说明，创建 GitHub Release；
3. 请求一次 `proxy.golang.org`，让 pkg.go.dev 尽快收录新版本。

### 4. 同步回 dev

```bash
git switch dev
git merge --ff-only main
git push origin dev
```

这样 `dev` 也包含那个合并提交，下次合并不会多出无关的差异。

### 5. 检查

- GitHub 的 Releases 页面出现了新版本，说明和 CHANGELOG 一致；
- `go list -m github.com/chhongzh/lilt@v0.2.0` 能查到；
- 文档站已经更新。

## 出了问题

- **已经推送的 tag 不要移动或删除。** module proxy 会永久缓存每个版本的内容，挪 tag 只会让不同的人拿到不同的代码。
- 发出去的版本有严重问题时，修好后发一个新的补丁版本。必要时在 `go.mod` 里用 `retract` 标出坏版本，并在那次发布里一起推出去：

  ```go
  retract v0.2.0 // Select 打开后无法关闭，请用 v0.2.1
  ```

- release workflow 失败而 tag 已经推送：在 Actions 页面修好原因后重新运行这个 workflow，不需要重新打 tag。
