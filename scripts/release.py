#!/usr/bin/env python3
"""Release a new version of lilt from dev.

    python scripts/release.py v0.2.0            check, confirm, then release
    python scripts/release.py v0.2.0 --dry-run  check and print the plan only
    python scripts/release.py v0.3.0-rc.1 --yes skip the confirmation prompt

Runs the checks CI runs, merges dev into main with a merge commit, tags main,
pushes both, and fast-forwards dev to main. Pushing the tag starts
.github/workflows/release.yml, where GoReleaser writes the notes from the
commits and publishes the GitHub release. RELEASING.md has the full process.
"""

import argparse
import io
import os
import re
import shutil
import subprocess
import sys
import tarfile
import tempfile
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
REMOTE = "origin"
REPO_URL = "https://github.com/chhongzh/lilt"
VERSION_RE = re.compile(r"^v(\d+)\.(\d+)\.(\d+)(?:-([0-9A-Za-z.-]+))?$")

# Matches the vet step in .github/workflows/ci.yml. driver/sh1106 and
# examples/esp32s3-generic import machine and are covered by the TinyGo build.
VET_HOST = [".", "./anim", "./gfx", "./ui", "./widget", "./driver/silent",
            "./examples/gallery", "./cmd/preview", "./cmd/docgif", "./internal/..."]
VET_WASM = ["./driver/desktop", "./cmd/sim"]


class ReleaseError(Exception):
    pass


def run(*cmd, env=None, check=True):
    """Runs cmd in the repository root and returns its stdout, stripped."""
    full_env = None if env is None else {**os.environ, **env}
    p = subprocess.run(cmd, cwd=ROOT, env=full_env, capture_output=True,
                       text=True, encoding="utf-8", errors="replace")
    if check and p.returncode != 0:
        out = (p.stdout + p.stderr).strip()
        raise ReleaseError(f"`{' '.join(cmd)}` failed:\n{out}")
    return p.stdout.strip()


def step(msg):
    print(f"\n==> {msg}", flush=True)


def parse_version(tag):
    m = VERSION_RE.match(tag)
    if not m:
        return None
    return tuple(int(x) for x in m.groups()[:3]), m.group(4)


def version_key(tag):
    # A release sorts above its own pre-releases, which git's v:refname
    # ordering gets the other way round. Numeric parts of a suffix compare
    # as numbers, so rc.10 comes after rc.2.
    nums, suffix = parse_version(tag)
    parts = [(0, int(p), "") if p.isdigit() else (1, 0, p) for p in (suffix or "").split(".") if p]
    return nums, suffix is None, parts


def latest_tag():
    tags = [t for t in run("git", "tag", "--list", "v*").splitlines() if parse_version(t)]
    return max(tags, key=version_key, default=None)


def check_version(tag, prev):
    new = parse_version(tag)
    if new is None:
        raise ReleaseError(f"{tag} is not vX.Y.Z or vX.Y.Z-suffix")
    if prev is None:
        return
    if version_key(tag) <= version_key(prev):
        raise ReleaseError(f"{tag} does not come after the latest tag {prev}")


def check_repo(tag):
    step("checking the repository")
    if run("git", "branch", "--show-current") != "dev":
        raise ReleaseError("switch to dev first")
    if run("git", "status", "--porcelain"):
        raise ReleaseError("the working tree has uncommitted changes")
    run("git", "fetch", REMOTE, "--tags", "--prune")
    if run("git", "rev-list", "--count", f"dev..{REMOTE}/dev") != "0":
        raise ReleaseError(f"dev is behind {REMOTE}/dev; pull first")
    if run("git", "rev-list", "--count", f"{REMOTE}/dev..dev") != "0":
        raise ReleaseError(f"dev has commits not on {REMOTE}/dev; push dev and let CI pass first")
    # main only ever moves by merging dev, so dev must already contain it.
    merge = subprocess.run(["git", "merge-base", "--is-ancestor", f"{REMOTE}/main", "dev"], cwd=ROOT)
    if merge.returncode != 0:
        raise ReleaseError(f"{REMOTE}/main has commits that dev lacks; merge main into dev first")
    if run("git", "tag", "--list", tag) or run("git", "ls-remote", "--tags", REMOTE, f"refs/tags/{tag}"):
        raise ReleaseError(f"{tag} already exists")


def check_gofmt():
    """Runs gofmt on HEAD as committed. With core.autocrlf the working tree
    holds CRLF files, which gofmt reports even when they are formatted."""
    with tempfile.TemporaryDirectory() as tmp:
        archive = subprocess.run(["git", "-c", "core.autocrlf=false", "archive", "HEAD"],
                                 cwd=ROOT, capture_output=True, check=True).stdout
        with tarfile.open(fileobj=io.BytesIO(archive)) as tar:
            tar.extractall(tmp, filter="data")
        out = subprocess.run(["gofmt", "-l", "."], cwd=tmp, capture_output=True, text=True)
        if out.returncode != 0:
            raise ReleaseError(f"gofmt failed:\n{out.stderr.strip()}")
        if out.stdout.strip():
            raise ReleaseError(f"run gofmt on:\n{out.stdout.strip()}")


def check_code():
    step("gofmt")
    check_gofmt()

    step("go mod tidy")
    run("go", "mod", "tidy")
    # git diff compares content after line-ending conversion; git status
    # would also flag tidy rewriting a CRLF checkout with LF.
    tidy = subprocess.run(["git", "diff", "--quiet", "go.mod", "go.sum"], cwd=ROOT).returncode == 0
    run("git", "checkout", "--", "go.mod", "go.sum")
    if not tidy:
        raise ReleaseError("go.mod or go.sum is not tidy; run `go mod tidy` and commit")

    step("go vet")
    run("go", "vet", *VET_HOST)
    run("go", "vet", *VET_WASM, env={"GOOS": "js", "GOARCH": "wasm"})

    if shutil.which("tinygo"):
        step("tinygo build")
        with tempfile.TemporaryDirectory() as tmp:
            run("tinygo", "build", "-target", "esp32s3-generic",
                "-o", str(Path(tmp) / "gallery.bin"), "./examples/esp32s3-generic")
    else:
        print("\n(tinygo not found; the firmware build is left to CI)")

    if shutil.which("goreleaser"):
        step("goreleaser check")
        run("goreleaser", "check")
    else:
        print("(goreleaser not found; .goreleaser.yaml is checked in CI)")


def plan(tag):
    msg = f"Merge branch 'dev' for {tag}"
    return [
        ["git", "switch", "main"],
        ["git", "merge", "--ff-only", f"{REMOTE}/main"],
        ["git", "merge", "--no-ff", "dev", "-m", msg],
        ["git", "tag", "-a", tag, "-m", tag],
        ["git", "push", REMOTE, "main"],
        ["git", "push", REMOTE, tag],
        ["git", "switch", "dev"],
        ["git", "merge", "--ff-only", "main"],
        ["git", "push", REMOTE, "dev"],
    ]


def main():
    ap = argparse.ArgumentParser(description="Release a new version of lilt from dev.")
    ap.add_argument("version", help="the new version, such as v0.2.0 or v0.3.0-rc.1")
    ap.add_argument("--dry-run", action="store_true", help="run the checks and print the plan only")
    ap.add_argument("--yes", action="store_true", help="do not ask before pushing")
    args = ap.parse_args()

    tag = args.version if args.version.startswith("v") else "v" + args.version
    prev = latest_tag()
    check_version(tag, prev)
    check_repo(tag)
    check_code()

    commits = run("git", "log", "--oneline", "--no-merges", f"{prev}..dev" if prev else "dev")
    step(f"{prev or 'nothing'} -> {tag}, {len(commits.splitlines())} commits")
    print(commits or "(none)")
    if not commits:
        raise ReleaseError("dev has no commits since the latest tag")

    commands = plan(tag)
    step("plan")
    for cmd in commands:
        print("  " + " ".join(f'"{a}"' if " " in a else a for a in cmd))
    if args.dry_run:
        return
    if not args.yes and input(f"\nRelease {tag} and push to {REMOTE}? [y/N] ").strip().lower() != "y":
        raise ReleaseError("stopped, nothing changed")

    for i, cmd in enumerate(commands):
        try:
            run(*cmd)
        except ReleaseError:
            done = [" ".join(c) for c in commands[:i]]
            print("\nDone before the failure:\n  " + ("\n  ".join(done) or "(nothing)"), file=sys.stderr)
            raise

    step(f"released {tag}")
    print(f"GoReleaser is publishing the release: {REPO_URL}/actions/workflows/release.yml")
    print(f"Release page: {REPO_URL}/releases/tag/{tag}")


if __name__ == "__main__":
    try:
        main()
    except ReleaseError as e:
        sys.exit(f"\nrelease: {e}")
    except KeyboardInterrupt:
        sys.exit("\nrelease: interrupted")
