# chd

- ディレクトリ内のファイル名が先頭から指定した長さまで一致する場合、サイズの小さい方を削除するプログラム
- コマンドライン上でディレクトリ指定可能
    - `--dir <directory>`
        - デフォルト値は`C:\Users`
        - 絶対パスでも相対パスでもかまわない
    - `--length <length>`
        - 比較する先頭からの長さ
        - デフォルト値は15
- 設定ファイルでディレクトリ指定可能
    - `--config <config file>`
    - 設定ファイルの書式は以下の通り
        ```
        {
            "dir": "C:\\Users\\username\\working_directory",
            "length": 15
        }
        ```
- 環境変数でディレクトリ指定可能
    - `PROCESS_DIR`
        - デフォルト値は`C:\Users`
        - 絶対パスでも相対パスでもかまわないが現実的には絶対パス指定
    - `--length <length>`
        - 比較する先頭からの長さ
        - デフォルト値は15
   - `PROCESS_LENGTH`
        - 比較する先頭からの長さ
        - デフォルト値は15


※GitHub Copilotの練習用として作成した
