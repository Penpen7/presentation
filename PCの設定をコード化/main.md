---
marp: true
theme: nix-talk
paginate: true
---

<!-- _class: lead -->

# PCの設定をコード化する話

## nix-darwinで宣言的セットアップ

ネオ美魔

---

# 自己紹介

- 名前: Penpen7
- 職種: バックエンドエンジニア
- 好きな言語: Rust, Go
- 仕事: KINTO FACTORYの開発に携わっています
- 最近のマイブーム: PCの設定を全部コードで管理すること（このスライドもコードで管理しています）

---

<!-- _class: lead -->

# 新しいMacが届きました🎉

---

<!-- _class: lead -->

# さて、環境構築するか...😇

---

# 環境構築あるある

- アプリを1個ずつダウンロードしてインストールする
- システム設定をポチポチ開いてキーリピートを最速にする
- `.zshrc`や`.gitconfig`をどこからかコピペする
- 手順書を作ったが、半年後には実態とズレている
- どこを変更したか分からなくなり、秘伝のタレになる
- そして二度と同じ環境は作れない

<div class="talk">
<div class="avatar"></div>
<div class="bubble">これ、毎回やるの本当にだるいんですよね…</div>
</div>

---

# 今回言いたいこと

- PCの設定を「宣言的」にコード化しよう
- nix-darwin を使う
- macOSのシステム設定・アプリ・CLIツール・dotfilesを1つのGitリポジトリで管理する
- 設定は`~/dotfiles`に全部書いてある
- 新しいMacでも1コマンドで同じ環境が手に入る

---

<!-- _class: lead -->

# 「宣言的」ってなに？

---

# 命令的なセットアップ

- 「何を実行するか」を並べる
- 2回流すと壊れることがある
- slackを消したくても、スクリプトを消すだけでは消えない

```sh
brew install slack
brew install discord
defaults write com.apple.dock tilesize 128
```

---

# 宣言的なセットアップ

- 最終的にどうなっていてほしいかを書く
- あとの差分はNixが埋めてくれる
- リストからslackを消すと、勝手にアンインストールされる

```nix
homebrew.casks = [ "slack" "discord" ];
system.defaults.dock.tilesize = 128;
```

<div class="talk">
<div class="avatar"></div>
<div class="bubble">消したいものはリストから消すだけ。これが気持ちいい</div>
</div>

---

# 登場人物

- Nix
  - パッケージマネージャ兼ビルドシステムで、再現性が高い
- nix-darwin
  - macOSのシステム設定・Homebrew・App Storeアプリを宣言的に管理する
- home-manager
  - ユーザー単位のdotfilesやCLIツールを宣言的に管理する
- flake
  - 入力と出力を1ファイルに固定する
- これらを`flake.nix`で束ねている

---

# セットアップの手順

- 新しいMacでやることは3ステップで終わる

```sh
# 1. Nixを入れる
curl -L https://install.determinate.systems/nix \
  | sh -s -- install

# 2. dotfilesをclone
git clone https://github.com/Penpen7/dotfiles \
  ~/dotfiles

# 3. 適用する
nix run nix-darwin -- switch --flake .#work
```

---

# システム設定を管理

- GUIでポチポチしていた設定が、全部コードになる

```nix
system.defaults = {
  dock = {
    orientation = "left";   # Dockは左
    tilesize = 128;         # アイコン大きめ
    autohide = false;
  };
  NSGlobalDomain = {
    KeyRepeat = 1;          # キーリピート最速
    InitialKeyRepeat = 15;  # リピート開始を短く
  };
  trackpad.Clicking = true; # タップでクリック
};
```

---

# アプリを管理

- Homebrew Cask（GUIアプリ）も宣言的に書ける

```nix
homebrew.casks = [
  "slack"
  "discord"
  "alfred"
  "cleanshot"
  "claude"
];
```

---

# アプリを管理

- Mac App Store（mas経由）のアプリまで宣言的に書ける

```nix
homebrew.masApps = {
  "Numbers" = 409203825;
  "Pages"   = 409201541;
  "Skitch"  = 425955336;
};
```

---

# CLIツールを管理

```nix
home.packages = with pkgs; [
  ripgrep   # 高速grep
  bat       # cat の強化版
  eza       # ls の強化版
  fzf       # ファジーファインダ
  ghq lazygit
];
```

---

# dotfilesを管理

- `.zshrc`や`.gitconfig`もリポジトリで管理する
- シンボリックリンクの作成もNixに任せられる

```nix
home.file.".gitconfig".source = ./gitconfig;
```

---

# アプリの設定を管理

- 最近はClaude Codeの設定までコード化している
- 許可コマンド・hooks・statusLineまで宣言的に書ける
- 「あの設定どこだっけ？」がコードを見れば一発でわかる

```nix
programs.claude-code = {
  enable = true;
  settings = {
    model = "opus";
    permissions.allow = [
      "Bash(git status:*)"
      "Bash(git diff:*)"
      "Bash(rg:*)"
    ];
  };
};
```

<div class="talk">
<div class="avatar"></div>
<div class="bubble">ツールの設定までGit管理。これが一番気に入ってます</div>
</div>

---

# 仕事用と個人用を1つに

- 共通モジュールにプロファイル差分を足す構成にする
- 仕事用には`gather`、個人用には`dj`や`game`など、profileごとに足すモジュールを変える
- `switch --flake .#work` か `.#personal` かを選ぶだけでよい

```nix
darwinConfigurations = {
  work     = mkDarwinSystem "work";
  personal = mkDarwinSystem "personal";
};
```

---

# 宣言的にして嬉しいこと

- 再現性: 新しいMac・PC交換でも、cloneして1コマンドで同じ環境が手に入る
- 差分が見える・戻せる: 設定変更はすべてGitのdiffに残り、気に入らなければ`git revert`で戻せる
- 全体像を一望できる: 「自分のPCに何が入っているか」がリポジトリを見ればわかる
- 掃除が効く: リストから消せば、アプリも設定も綺麗に消える（秘伝のタレが育たない）

---

# 正直、しんどいところもある

- Nix言語の学習コストが高い（独特の関数型言語）
- エラーメッセージが解読しにくいことがある
- 初回ビルドや更新に時間がかかる
- すべてのアプリが管理できるわけではない
  - Cisco Packet TracerなどはREADMEに「手動管理」と明記して逃げ道を残している

→ それでも、環境をいつでも作り直せる安心感の方が大きい

---

# まとめ

- PCの設定は「命令」ではなく「宣言」で書ける
- nix-darwin + home-manager で、システム設定・アプリ・CLI・dotfiles・アプリ内設定まで1つのリポジトリにまとまる
- 新しいMacでも1コマンドで同じ環境を再現できる
- 学習コストはあるけど、「秘伝のタレ」から卒業できる
- 興味が出たら、ぜひnix-darwinを触ってみてください！
