# Code flow by command

Entry point: **`cmd/cdm/main.go`** (package main). It calls **`cdmcmd.Run()`**. **`internal/cdmcmd/root.go`** (package cdmcmd) defines **`Run()`**, which builds the root Cobra command, registers subcommands, and runs **`fang.Execute(context.Background(), root)`**, which handles parsing and dispatch. Each subcommand’s `RunE` is in **`internal/cdmcmd/*.go`** and is invoked with the parsed args and flags.

Commands that need a repo root first call **`runner.GetConfigRoot(runner.Options{TraverseUpPrompt: true})`** in **`internal/runner/runner.go`**: **`GetConfigRoot`** uses **`config.ExistsInDir(cwd)`** and, if not found, **`config.FindConfigDir`**; if `TraverseUpPrompt` is true and still not found, it calls **`ui.ConfirmTraverseUp()`** (in **`internal/ui/prompt.go`**), then **`config.FindConfigDir(cwd, true, maxTraverseLevels)`** to search up to 5 parent dirs. The returned `root` is the directory that contains `.cdm.toml`.

---

## `cdm apply` (and `cdm apply -d` / `--dry-run`)

1. **`internal/cdmcmd/apply.go`**: **`runApply`** reads the `--dry-run`/`-d` flag, calls **`runner.GetConfigRoot`** for repo root, then **`apply.Run(root, dryRun)`**.
2. **`internal/apply/apply.go`**: **`Run`** loads config via **`config.Load(root)`** and **`config.LoadLocal(root)`**, then **`cfg.ResolvedMappings(local)`** (in **`internal/config/mappings.go`**: **`ResolvedMappings`** merges default + OS-specific mappings and **`substituteVars`** for variable expansion). Missing variables are warned with **`ui.PrintMissingVarsWarning`**.
3. If dry-run, **`ui.PrintDryRunBanner()`** is called.
4. **`state.LoadHashes(root)`** (in **`internal/state/hashes.go`**) loads `.cdm/hashes.toml`.
5. If not dry-run, pre_apply scripts are run: for each path in **`cfg.Scripts.PreApply`**, **`runner.AbsPath(root, script)`** is used, then **`runScript(scriptPath, root)`** (which uses **`exec.Command`** with `cmd`/`sh`).
6. For each resolved mapping, **`Run`** reads the source with **`runner.AbsPath(root, srcRel)`**, hashes with **`state.ContentHash`** (**`internal/state/hash.go`**). It ensures target dir exists (**`os.MkdirAll`**), and if the target file exists and its hash (from **`state.FileHash`**) differs from **`hashes[targetPath]`**, it backs up via **`state.EnsureTempDir`**, **`state.BackupPath`**, writes the backup, and records in `backupManifest`. In dry-run it only appends warnings and skips writes.
7. If not dry-run: **`os.WriteFile(targetPath, srcData, 0644)`**, updates **`hashes[targetPath]`**, and **`state.AppendLog(root, ...)`** (**`internal/state/log.go`**).
8. After the loop: **`state.SaveHashes(root, hashes)`**, **`state.WriteLastBackupManifest(root, backupManifest)`** (**`internal/state/backup.go`**), and post_apply scripts via **`runScript`**.
9. Back in **`runApply`**: warnings and **`result.Errors`** are printed with **`ui.PrintWarn`** / **`ui.PrintError`**; on success, **`ui.PrintSuccess`** for applied count and any backup count.

---

## `cdm diff`

1. **`internal/cdmcmd/diff.go`**: **`runDiff`** gets root via **`runner.GetConfigRoot`**, then **`diff.Run(root)`**.
2. **`internal/diff/diff.go`**: **`Run`** loads **`config.Load`** and **`config.LoadLocal`**, then **`cfg.ResolvedMappings(local)`**; missing vars get **`ui.PrintMissingVarsWarning`**. For each mapping it reads source and target files, then **`udiff.Unified(srcRel, targetPath, targetData, srcData)`** (from **`github.com/aymanbagabas/go-udiff`**) and prints the unified diff.

---

## `cdm status`

1. **`internal/cdmcmd/status.go`**: **`runStatus`** gets root via **`runner.GetConfigRoot`**, then **`status.Run(root)`**.
2. **`internal/status/status.go`**: **`Run`** loads config and **`cfg.ResolvedMappings(local)`**, warns on missing vars, then **`state.LoadHashes(root)`**. For each mapping it gets **`state.FileHash(srcPath)`** and compares to **`hashes[targetPath]`**: if no stored hash it prints "? … (not yet applied)"; if **`srcHash != stored`** it prints "M " + srcRel (modified).

---

## `cdm init` and `cdm init <repo>`

1. **`internal/cdmcmd/init.go`**: **`runInit`** gets cwd with **`os.Getwd()`**, sets `repo := args[0]` if present, and calls **`initcmd.Run(cwd, repo)`**.
2. **`internal/initcmd/init.go`**: **`Run`** — if `cloneRepo` is non-empty, **`cloneAndChdir(&root, cloneRepo)`** runs **`runGitClone(repo, cloneDir)`** ( **`exec.Command("git", "clone", repo, dir)`** ) and sets `root` to the cloned directory.
3. **`Run`** then: **`os.MkdirAll`** for `.cdm/temp`; if `.cdm.toml` is missing, writes a commented template and **`ui.PrintSuccess`**; reads `.gitignore`, and if `.cdm/` or `.local.cdm.toml` are missing, appends them and **`ui.PrintSuccess("updated .gitignore")`**; finally **`ui.PrintSuccess("init done at " + root)`**.

---

## `cdm add <path> [target-path]`

1. **`internal/cdmcmd/add.go`**: **`runAdd`** gets root via **`runner.GetConfigRoot`**, takes `args[0]` as source path and optional `args[1]` as target-path, then **`addcmd.Run(root, sourcePath, targetPath)`**.
2. **`internal/addcmd/add.go`**: **`Run`** resolves **`absSrc`** with **`filepath.Abs(sourcePath)`**, reads the file, and computes **`destRel`** (repo-relative path: basename or, if second arg is a non-absolute, non-`$` path, **`filepath.ToSlash(filepath.Clean(targetDest))`**). It writes the file to **`filepath.Join(root, filepath.Clean(destRel))`** with **`os.MkdirAll`** for the parent dir.
3. **`config.Load(root)`** loads existing config, sets **`cfg.Mappings[destRel]`** to **`targetDest`** or **`"$HOME/." + baseName`**, then **`writeConfig(root, cfg)`** marshals the full config with **`toml.Marshal(cfg)`** and **`os.WriteFile`** to `.cdm.toml.
4. Back in **`runAdd`**, **`ui.PrintSuccess("added " + sourcePath)`**.

---

## `cdm watch`

1. **`internal/cdmcmd/watch.go`**: **`runWatch`** gets root via **`runner.GetConfigRoot`**, builds **`signal.NotifyContext(cmd.Context(), os.Interrupt)`** and defers **`stop()`**, defines **`applyFn := func() (*apply.ApplyResult, error) { return apply.Run(root, false) }`**, and calls **`watch.Run(ctx, root, applyFn)`**.
2. **`internal/watch/watch.go`**: **`Run`** loads **`config.Load`** / **`config.LoadLocal`** and **`cfg.ResolvedMappings(local)`**, collects source paths, and creates **`fsnotify.NewWatcher()`**. It adds each source file’s parent directory with **`watcher.Add(dir)`**. A goroutine loops on **`watcher.Events`** and **`ctx.Done()`**; on event it calls **`scheduleApply`**, which (under a mutex) resets a **`time.AfterFunc(debounceMs, ...)`** that invokes **`applyFn()`**. **`Run`** blocks on **`<-ctx.Done()`** and returns **`ctx.Err()`** (e.g. on Ctrl+C).

---

## `cdm undo-conflict-apply`

1. **`internal/cdmcmd/undo.go`**: **`runUndoConflictApply`** gets root via **`runner.GetConfigRoot`**, then **`undo.Run(root)`**.
2. **`internal/undo/undo.go`**: **`Run`** calls **`state.LoadLastBackupManifest(root)`** (**`internal/state/backup.go`**: reads `.cdm/last_backup.toml`). If empty, **`ui.PrintWarn("no backup manifest…")`** and return. For each **`tempPath → targetPath`** in the manifest it reads the temp file, **`os.MkdirAll`** for the target’s dir, **`os.WriteFile(targetPath, data, 0644)`**, **`ui.PrintSuccess("restored …")`**, and **`state.AppendLog(root, …)`**. Finally **`state.ClearLastBackupManifest(root)`** removes `.cdm/last_backup.toml`.

---

## `cdm help` / `cdm -h` / `cdm --help` and `cdm completion`

Handled by Cobra and Fang: **`fang.Execute`** wires the root command so **`-h`**/ **`--help`** and the built-in **`help`** subcommand show usage. Fang adds a **`completion`** subcommand that generates shell completion scripts (bash/zsh/fish); no custom code in this repo beyond registering the root and subcommands in **`internal/cdmcmd/root.go`**.
