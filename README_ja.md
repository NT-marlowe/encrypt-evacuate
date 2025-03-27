# Fuga: ランサムウェアからリアクティブにデータを保護するシステム

## 概要
![System Overview](/assets/system_overview.png)

## 実行環境

### minimum requirements
以下のdependenciesが必要である．必要十分ではない．カーネルバージョンが5.7では動かない (後述)．
- Linux kernel version 5.7 or later (for bpf_link support)
- LLVM 11 or later (clang and llvm-strip)
- libbpf headers
- Linux kernel headers
- Go compiler version supported by ebpf-go's Go module

[参考](https://ebpf-go.dev/guides/getting-started/#ebpf-c-program)

### kashiwa-go
Thu Mar 27 04:40:19 PM JST 2025 時点での依存関係は以下の通り．
> [!NOTE]  
> すべてのバージョンで試したわけではないが，カーネルバージョンが5.xだと，eBPF verifierがプログラムを許可せずエラーが出てしまう．
> 6.xであればverifierの実装が改善され，false positiveが削減されているため実行可能になる．

```bash
$ uname -r
6.3.0-060300-generic


$ clang -v
Ubuntu clang version 14.0.0-1ubuntu1.1
Target: x86_64-pc-linux-gnu
Thread model: posix
InstalledDir: /usr/bin
Found candidate GCC installation: /usr/bin/../lib/gcc/x86_64-linux-gnu/11
Found candidate GCC installation: /usr/bin/../lib/gcc/x86_64-linux-gnu/12
Selected GCC installation: /usr/bin/../lib/gcc/x86_64-linux-gnu/12
Candidate multilib: .;@m64
Selected multilib: .;@m64

$ llvm-strip --version
llvm-strip, compatible with GNU strip
Ubuntu LLVM version 14.0.0

  Optimized build.
  Default target: x86_64-pc-linux-gnu
  Host CPU: icelake-server


$ go version
go version go1.22.5 linux/amd64
```

## 実行方法
上記の依存関係を満たしていることを前提とする．
```bash
$ git clone [URL]
$ cd encrypt-evacuate
$ make run
```
`make run`実行時の処理は`Makefile`にコメントを残している．
