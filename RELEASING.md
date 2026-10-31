# 发布流程

这份文档写给维护者，说明怎么从 `dev` 发出一个新版本。日常只需要一条命令：

```bash
python scripts/release.py v0.2.0
```

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

预发布 tag 带 `-` 后缀，GitHub Release 会标成 pre-release，`go get ...@latest` 也不会选中它。

## 提交信息

发布说明由 [GoReleaser](https://goreleaser.com) 从上一个 tag 以来的提交生成，所以提交信息要遵循 [Conventional Commits](https://www.conventionalcommits.org/zh-hans/)。规则写在 [.goreleaser.yaml](.goreleaser.yaml)：

| 前缀 | 进发布说明的哪一组 |
| --- | --- |
| `feat:`、`feat(ui):` | 新功能 |
| `fix:` | 修复 |
| `perf:` | 性能 |
| `refactor:` | 重构 |
| 类型后加 `!`，如 `feat(ui)!:` | 不兼容的改动，排在最前 |
| `docs:`、`chore:`、`ci:`、`test:`、`style:`、合并提交 | 不出现 |
| 其他 | 其他 |

提交信息就是用户看到的更新说明，写清楚改了什么，而不是改了哪个文件。不兼容的改动一定要带 `!`。

## 发布

在 `dev` 上，工作区干净、已经推送、CI 是绿的，然后：

```bash
python scripts/release.py v0.2.0 --dry-run
```

`--dry-run` 只做检查，打印这次会发布的提交和要执行的命令，什么都不改。确认无误后去掉 `--dry-run` 再跑一次，脚本会在推送前再问一次。

脚本依次：

1. 检查版本号格式，并且比最新的 tag 大；
2. 检查当前在 `dev`、工作区干净、和 `origin/dev` 一致、包含 `origin/main` 的全部提交、tag 还不存在；
3. 跑 CI 会跑的检查：gofmt、`go mod tidy`、host 和 wasm 的 `go vet`，装了 TinyGo 就编一次固件，装了 GoReleaser 就校验配置；
4. 用 `--no-ff` 把 `dev` 合并进 `main`，在合并提交上打注释 tag；
5. 先推 `main`，再推 tag；
6. 切回 `dev`，快进到 `main` 并推送，下次合并不会多出无关的差异。

中途失败时，脚本会列出失败前已经完成的命令，按需手动收尾。

推送之后的事由 GitHub Actions 完成：

- 推 `main` 重新部署文档站；
- 推 tag 触发 [release workflow](.github/workflows/release.yml)：跑完整 CI，GoReleaser 按上面的分组生成发布说明并创建 GitHub Release，最后请求一次 `proxy.golang.org`，让 pkg.go.dev 尽快收录。

最后看一眼 [Releases](https://github.com/chhongzh/lilt/releases) 页面。发布说明可以在 GitHub 上直接编辑补充，比如写一段升级指南。

## 手动发布

没有 Python 时，脚本做的事等于下面这些命令：

```bash
git switch main
git merge --ff-only origin/main
git merge --no-ff dev -m "Merge branch 'dev' for v0.2.0"
git tag -a v0.2.0 -m "v0.2.0"
git push origin main
git push origin v0.2.0
git switch dev
git merge --ff-only main
git push origin dev
```

## 出了问题

- **已经推送的 tag 不要移动或删除。** module proxy 会永久缓存每个版本的内容，挪 tag 只会让不同的人拿到不同的代码。
- 发出去的版本有严重问题时，修好后发一个新的补丁版本。必要时在 `go.mod` 里用 `retract` 标出坏版本，并随那次发布一起推出去：

  ```go
  retract v0.2.0 // Select 打开后无法关闭，请用 v0.2.1
  ```

- release workflow 失败而 tag 已经推送：在 Actions 页面修好原因后重新运行这个 workflow，不需要重新打 tag。
