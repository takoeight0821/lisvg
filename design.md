以下では、\*\*「S 式 → SVG」\*\*のワンパス・パイプラインを持つ“Graphviz／PlantUML 風”ツールを設計する際の典型的な構成と実装ステップを整理します。

---

## 1. 言語設計 — 最低限の S-式 DSL

まずは **ノード／エッジ／属性** を表現できる最小構文から始め、徐々に機能拡張するのが安全です。

```lisp
;; sample.sxd  (sexp-diagram)
(diagram
  (size 800 400)                 ; 描画領域 (幅 高さ)
  (node-style :shape "rect")     ; デフォルト属性
  (edge-style :stroke "#555")

  (nodes                         ; ノード一覧
    (id "A" :label "開始")
    (id "B" :label "処理")
    (id "C" :label "終了" :shape "ellipse"))

  (edges                         ; エッジ一覧
    ("A" "B" :label "next")
    ("B" "C")))
```

*ポイント*

| 要素                        | 役割                   |
| ------------------------- | -------------------- |
| `diagram`                 | ルート。以降の式を評価順に解釈      |
| `(size w h)`              | SVG ビューポート           |
| `(nodes …)` / `(edges …)` | コレクションごとにまとめるとパースが容易 |
| 属性は `:key value` 形式       | 処理系をシンプルに保てる         |

後述のパーサが読み込むと **AST** は「図全体 → ノード集合 → エッジ集合 → 属性マップ」というツリーになります。

---

## 2. パーサ／AST 構築

### 2-1. S 式パーサ

* **Clojure/Lisp** 系を使うなら Reader をそのまま流用。
* **Rust** なら `lexpr`、**Go** なら `github.com/anz-bank/go-sql-parser/sexpr` など既成ライブラリで十分。
* あるいは十数行の再帰下降パーサでも実装可能（トークンは `(` `)` と Atom の 3 種）。

### 2-2. バリデーション

1. ルートが `diagram` か確認
2. ノード ID の重複チェック
3. エッジが参照するノード ID の存在チェック

---

## 3. レイアウト戦略

「**自前実装 vs 既存レイアウトエンジン流用**」で労力が大きく変わります。

| 方針                   | 概要                                        | メリット                                               | デメリット      |
| -------------------- | ----------------------------------------- | -------------------------------------------------- | ---------- |
| **Graphviz に丸投げ**    | AST → DOT 変換 → `dot -Tplain` で座標取得        | 高品質レイアウト／複雑グラフ対応                                   | dot が必須依存  |
| **libgraphviz 埋め込み** | C ライブラリを FFI で呼出                          | バイナリ配布が重い                                          |            |
| **軽量アルゴリズム実装**       | Tree／DAG 向け Sugiyama，Force-Directed などを自前 | 単一バイナリ・軽量                                          | 高度な経路整形は大変 |
| **モダン OSS**          | 例: D2 のレイアウトプラグインを流用する構造                  | コード例多い／活発コミュニティ ([github.com][1], [d2lang.com][2]) | API 安定度に注意 |

MVP 段階では **Graphviz 丸投げ**が推奨です。AST を素直に DOT にシリアライズして座標・ベジェ制御点を受取り、SVG 組立だけ自前で行えば短期間で品質を確保できます。

---

## 4. SVG 生成

1. `<svg width=… height=…>` ルート要素生成
2. ノードごとに

   * 図形 (`<rect>`, `<ellipse>` など)
   * ラベル (`<text>`、`dominant-baseline="middle" text-anchor="middle"`)
3. エッジは `<path>`（Graphviz 由来ならベジェ制御点をそのまま適用）
4. グループ化したい場合は `<g>` と `class` 属性を付与
5. 可能なら **CSS を埋め込む**（色・線幅をスタイルシート化すると diff がきれい）

Go なら `github.com/ajstarks/svgo`, Rust なら `svg` crate で DOM 生成が容易です。svgdx という「SVG 拡張 DSL」も参考になるでしょう ([github.com][3])。

---

## 5. アーキテクチャ例

```
            +-----------+       +-------------+       +-----------+
.sxd ──▶───►│  Parser   │ ───▶──│  Validator  │ ───▶──│  Layout   │
            +-----------+       +-------------+       +-----------+
                                                       │ coords   │
                                                       ▼
            +-----------+       +-------------+       +-----------+
            │  Renderer │ ◀───  │  SVG DOM    │ ◀───  │  Styles   │
            +-----------+       +-------------+       +-----------+
```

* CLI: `sxd compile sample.sxd -o sample.svg`
* Watcher: `sxd watch src/**/*.sxd` → ホットリロード
* 将来的に PlantUML ライクな「シーケンス図方言」を追加するときは **方言ごとに AST→DOT 変換器を差し替え**る設計にすると保守が楽です。

---

## 6. 今後の拡張アイデア

| 機能                    | 実装メモ                                            |
| --------------------- | ----------------------------------------------- |
| **レイアウトプラグイン**        | D2 のように外部バイナリで差替可能にする ([terrastruct.com][4])    |
| **import / include**  | 大規模図をモジュール化                                     |
| **テーマ機構**             | Mermaid に近い CSS-like テーマ宣言で色一括変更                |
| **双方向編集**             | SVG に `data-node-id` を埋めてブラウザ上からドラッグ→座標を逆シリアライズ |
| **エクスポート多様化**         | PNG・PDF 生成には `rsvg-convert` などをラップ              |
| **WebAssembly ランタイム** | フロントエンドで \*.sxd → SVG を即時プレビュー                  |

---

## 7. 参考 OSS・仕様比較（調査結果ダイジェスト）

| ツール / 規格           | 概要                          | ヒント                                             |
| ------------------ | --------------------------- | ----------------------------------------------- |
| **Graphviz DOT**   | de-facto 標準。平易だが属性が多い       | レイアウトだけ借用しても OK                                 |
| **Terrastruct D2** | 2022 年 OSS 化。プラグインでレイアウト交換可 | 文法・レンダラ構成の参考 ([github.com][1], [d2lang.com][2]) |
| **svgdx**          | “SVG 拡張 DSL” とライブラリ         | SVG DOM 操作のアイデア ([github.com][3])               |

---

## 8. 実装ロードマップ（3 フェーズ例）

| フェーズ              | 目標                               | 期間    | 成果物                |
| ----------------- | -------------------------------- | ----- | ------------------ |
| **0. MVP**        | Parser → DOT → Graphviz → SVG 出力 | 1-2 週 | コマンド `sxd compile` |
| **1. ネイティブレイアウト** | 単純ツリー／DAG を自前アルゴリズムで配置           | 2-3 週 | `--layout=simple`  |
| **2. 拡張 DSL**     | サブグラフ、クラスタ、ラベル装飾、テーマ             | 4-6 週 | v0.3.0 リリース        |

---

### 次の具体的アクション

1. **実装言語の選定**（Rust, Go, TypeScript など）
2. GitHub に `sxd` などのリポジトリを作成
3. Parser & DOT 変換器を最小構文で書き切る
4. `sample.sxd` をコミットし SVG が出るところまで動かす
5. Issue に将来拡張をラベル分けして登録—公開ロードマップ化

---

ご不明点や「○○言語でのライブラリ選定」「レイアウトアルゴリズムの詳細」「Web フロントエンドへの組み込み」など、深掘りしたいトピックがあれば気軽にどうぞ!

[1]: https://github.com/terrastruct/d2?utm_source=chatgpt.com "D2 is a modern diagram scripting language that turns text to ... - GitHub"
[2]: https://d2lang.com/?utm_source=chatgpt.com "Home | D2 Documentation"
[3]: https://github.com/codedstructure/svgdx?utm_source=chatgpt.com "svgdx - create SVG diagrams easily - GitHub"
[4]: https://terrastruct.com/?utm_source=chatgpt.com "Terrastruct | Diagramming tools crafted to visualize software ..."