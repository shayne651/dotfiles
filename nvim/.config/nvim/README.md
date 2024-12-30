# Neovim Key Mappings

---

### General Editor Settings

| Setting                          | Description                        | Example Use                                    |
|----------------------------------|------------------------------------|------------------------------------------------|
| `vim.cmd("set expandtab")`       | Convert tabs to spaces            | Automatically converts all tabs to spaces      |
| `vim.cmd("set tabstop=2")`       | Set tab width to 2 spaces         | Sets tab width for consistency                 |
| `vim.cmd("set shiftwidth=2")`    | Indent width to 2 spaces          | Sets indentation level when using `>` or `<`   |
| `vim.opt.number = true`          | Show line numbers                 | Displays line numbers on the side              |
| `vim.opt.relativenumber = true`  | Show relative line numbers        | Enables easy line jumping                      |

---

### File Explorer and Finder

| Key Mapping       | Description                           | Example Use                                      |
|-------------------|---------------------------------------|--------------------------------------------------|
| `<leader>p`       | Find files with Telescope            | Press `<leader>p` to open file finder            |
| `<leader>g`       | Live grep with Telescope             | Press `<leader>g` to search within files         |
| `<leader>fb`      | Open buffer list in Telescope        | Press `<leader>fb` to view all buffers           |
| `<leader>fh`      | Show help tags in Telescope          | Press `<leader>fh` to browse help topics         |
| `<leader>l`       | Open Neo-tree file explorer          | Press `<leader>l` to toggle file explorer        |

---

### Harpoon (Quick Navigation)

| Key Mapping       | Description                           | Example Use                                      |
|-------------------|---------------------------------------|--------------------------------------------------|
| `<leader>a`       | Add file to Harpoon                  | Press `<leader>a` to bookmark the current file   |
| `<leader>s`       | Open Harpoon quick menu              | Press `<leader>s` to see and access bookmarks    |
| `<leader>1-4`     | Open Harpoon file at index 1-4       | Press `<leader>1` to go to first bookmark, etc.  |
| `<C-n>`           | Go to the next file in Harpoon       | Press `<C-n>` to go to the next bookmarked file  |
| `<C-p>`           | Go to the previous file in Harpoon   | Press `<C-p>` to go to the previous bookmarked file |

---

### LSP (Language Server Protocol) Navigation and Actions

| Key Mapping       | Description                           | Example Use                                      |
|-------------------|---------------------------------------|--------------------------------------------------|
| `K`               | Show LSP hover information           | Press `K` to view documentation for symbol under cursor |
| `gd`              | Go to definition                     | Press `gd` to jump to the definition             |
| `gD`              | Go to declaration                    | Press `gD` to jump to the declaration            |
| `gi`              | Go to implementation                 | Press `gi` to jump to the implementation         |
| `go`              | Go to type definition                | Press `go` to see type definition                |
| `gr`              | Find all references                  | Press `gr` to list references                    |
| `gs`              | Show signature help                  | Press `gs` to see function/method signature      |

---

### LSP (Language Server Protocol) Edits and Formatting

| Key Mapping       | Description                           | Example Use                                      |
|-------------------|---------------------------------------|--------------------------------------------------|
| `<D-r>`           | Rename symbol                        | Press `<D-r>` to rename the symbol under cursor  |
| `<D-f>`           | Format document                      | Press `<D-f>` to format the current file         |
| `<D-a>`           | Show code actions                    | Press `<D-a>` to see available code actions      |

---

### DAP (Debug Adapter Protocol)

| Key Mapping       | Description                           | Example Use                                      |
|-------------------|---------------------------------------|--------------------------------------------------|
| `<Leader>dc`      | Start debugging                      | Press `<Leader>dc` to start debug session        |
| `<F10>`           | Step over                            | Press `<F10>` to step over the current line      |
| `<F11>`           | Step into                            | Press `<F11>` to step into function              |
| `<F12>`           | Step out                             | Press `<F12>` to step out of the current function|
| `<Leader>dt`      | Toggle breakpoint                    | Press `<Leader>dt` to set or remove a breakpoint |
| `<Leader>B`       | Set a conditional breakpoint         | Press `<Leader>B` to set breakpoint with condition|
| `<Leader>lp`      | Log point                            | Press `<Leader>lp` to set a log point            |
| `<Leader>dr`      | Open DAP REPL                        | Press `<Leader>dr` to open REPL for debugging    |
| `<Leader>dl`      | Run last debug configuration         | Press `<Leader>dl` to rerun the last debug config|
| `<Leader>dh`      | Show hover (DAP UI)                  | Press `<Leader>dh` to see variable value on hover|
| `<Leader>dp`      | Show preview (DAP UI)                | Press `<Leader>dp` to preview expressions        |
| `<Leader>df`      | Show DAP frames                      | Press `<Leader>df` to open frames view           |
| `<Leader>ds`      | Show DAP scopes                      | Press `<Leader>ds` to open scopes view           |

---

### Git Integration

| Key Mapping       | Description                           | Example Use                                      |
|-------------------|---------------------------------------|--------------------------------------------------|
| `<Leader>gp`      | Preview Git hunk                     | Press `<Leader>gp` to preview changes in hunk    |
| `<Leader>gb`      | Toggle Git blame                     | Press `<Leader>gb` to toggle line blame info     |
| `<Leader>gpl`     | Pull from Git                        | Press `<Leader>gpl` to execute `git pull`        |

---

