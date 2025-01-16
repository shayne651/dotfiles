-- converts tabs to spaces
vim.cmd("set expandtab")

-- sets tab to 2 spaces
vim.cmd("set tabstop=2")
vim.cmd("set softtabstop=2")
vim.cmd("set shiftwidth=2")

-- set leader
vim.g.mapleader = " "

-- show hidden files/folders
vim.g.netrw_list_hide = ""

-- set line numbers and relative line numbers
vim.opt.number = true
vim.opt.relativenumber = true

-- install lazy vim
require("config.lazy")

-- setup telescope (fuzzy finding)
local builtin = require('telescope.builtin')
vim.keymap.set('n', '<leader>p', function()
                                      require('telescope.builtin').find_files({ hidden = true })
                                 end,
  { desc = 'Telescope find files'})
vim.keymap.set('n', '<leader>g', builtin.live_grep, { desc = 'Telescope live grep'})
vim.keymap.set('n', '<leader>fb', builtin.buffers, { desc = 'Telescope buffers'})
vim.keymap.set('n', '<leader>fh', builtin.help_tags, { desc = 'Telescope help tags'})

-- setup neo-tree (file explorer)
vim.keymap.set ('n', '<leader>l', ':Neotree filesystem reveal left<CR>')

-- setup harppon (file hotkey navigation)
local harpoon = require("harpoon")
harpoon:setup()
vim.keymap.set('n', '<leader>a', function() harpoon:list():add() end)
vim.keymap.set('n', '<leader>s', function() harpoon.ui:toggle_quick_menu(harpoon:list()) end)
vim.keymap.set('n', '<leader>1', function() harpoon:list():select(1) end)
vim.keymap.set('n', '<leader>2', function() harpoon:list():select(2) end)
vim.keymap.set('n', '<leader>3', function() harpoon:list():select(3) end)
vim.keymap.set('n', '<leader>4', function() harpoon:list():select(4) end)
vim.keymap.set('n', '<C-n>', function() harpoon:list():next() end)
vim.keymap.set('n', '<C-p>', function() harpoon:list():previous() end)

-- setup lsp highlighting
vim.cmd [[
  hi link @lsp.type.variable.go Normal
  hi link @lsp.type.parameter.go Identifier
  hi link @lsp.type.function.go Function
  hi link @lsp.type.keyword.go Keyword
]]

-- run gofmt on save
vim.api.nvim_create_autocmd("BufWritePre", {
  pattern = "*.go",
  callback = function()
    vim.lsp.buf.format({ async = false })  -- Format using LSP before saving
  end,
})

