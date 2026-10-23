---
hide:
  - toc
---

# 在线体验

下面是 [gallery](https://github.com/chhongzh/lilt/tree/main/examples/gallery) 编译成 WebAssembly 后的样子，和固件、桌面模拟器是同一份代码。先点一下画面让它拿到键盘焦点。

<iframe class="lilt-sim" src="../sim/" title="lilt gallery" loading="lazy"></iframe>

| 键 | 对应 |
| --- | --- |
| ++w++ ++a++ ++s++ ++d++ 或方向键 | 上 左 下 右 |
| ++space++、++enter++ | OK |
| ++j++ | Back |
| ++k++ | Menu |

想全屏玩，直接打开 <a href="../sim/" target="_blank">sim/</a>。本地预览时（`scripts/docs.sh serve`）这里是空的，`scripts/docs.sh build` 才会把 wasm 编进 `site/sim/`。
