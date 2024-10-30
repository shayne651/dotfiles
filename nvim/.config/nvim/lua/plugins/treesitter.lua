return {{
    'nvim-treesitter/nvim-treesitter',
    build = ':TSUpdate',
    event = { 'BufReadPre', 'BufNewFile' },
    config = function ()
      local configs = require('nvim-treesitter.configs')

      configs.setup({
          ensure_installed = { 'go', 'lua','dart', 'typescript', 'javascript', 'html', 'python', 'java', 'sql', 'dockerfile' },
          sync_install = false,
          highlight = {
            enable = true,
            additional_vim_regex_highlighting = false,
          },
          indent = {
            enable = true
          },
        })
    end
 }}
