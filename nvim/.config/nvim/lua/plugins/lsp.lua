return {
  {
    'VonHeikemen/lsp-zero.nvim',
    branch = 'v4.x',
    lazy = true,
    config = false,
  },
  {
    'williamboman/mason.nvim',
    lazy = false,
    opts = {},
  },

  -- Autocompletion
  {
    'hrsh7th/nvim-cmp',
    event = 'InsertEnter',
    dependencies = {
      'hrsh7th/cmp-cmdline',
    },
    config = function()
      local cmp = require('cmp')

      cmp.setup({
        sources = {
          {name = 'nvim_lsp'},
        },
        mapping = cmp.mapping.preset.insert({
          ['<D-S>'] = cmp.mapping.complete(),
          ['<D-u>'] = cmp.mapping.scroll_docs(-4),
          ['<D-d>'] = cmp.mapping.scroll_docs(4),
          ['<CR>'] = cmp.mapping.confirm({ select = true}),
        }),
        snippet = {
          expand = function(args)
            vim.snippet.expand(args.body)
          end,
        },
      })
      -- Command-line mode setup for showing command suggestions
      cmp.setup.cmdline(':', {
        mapping = cmp.mapping.preset.cmdline(),
        sources = {
          { name = 'cmdline' }
        }
      })
    end
  },

  -- Syntax highlighting
  {
    'nvim-treesitter/nvim-treesitter',
    build = ':TSUpdate',
    event = { 'BufReadPre', 'BufNewFile' },
    config = function()
      require('nvim-treesitter.configs').setup({
        ensure_installed = { 'go', 'typescript', 'java', 'sql', 'dockerfile', 'dart', 'rust', 'c' },
        highlight = {
          enable = true,
          additional_vim_regex_highlighting = false,
        },
      })
    end,
  },

  -- LSP
  {
    'neovim/nvim-lspconfig',
    cmd = {'LspInfo', 'LspInstall', 'LspStart'},
    event = {'BufReadPre', 'BufNewFile'},
    dependencies = {
      {'hrsh7th/cmp-nvim-lsp'},
      {'williamboman/mason.nvim'},
      {'williamboman/mason-lspconfig.nvim'},
    },
    init = function()
      -- Reserve a space in the gutter
      -- This will avoid an annoying layout shift in the screen
      vim.opt.signcolumn = 'yes'
    end,
    config = function()
      local lsp_defaults = require('lspconfig').util.default_config

      -- Add cmp_nvim_lsp capabilities settings to lspconfig
      -- This should be executed before you configure any language server
      lsp_defaults.capabilities = vim.tbl_deep_extend(
        'force',
        lsp_defaults.capabilities,
        require('cmp_nvim_lsp').default_capabilities()
      )

      -- LspAttach is where you enable features that only work
      -- if there is a language server active in the file
      vim.api.nvim_create_autocmd('LspAttach', {
        desc = 'LSP actions',
        callback = function(event)
          local opts = {buffer = event.buf}

         -- Mappings for LSP features (Command key-based)
          vim.keymap.set('n', 'K', '<cmd>lua vim.lsp.buf.hover()<cr>', opts)               -- Show documentation
          vim.keymap.set('n', 'gd', '<cmd>lua vim.lsp.buf.definition()<cr>', opts)         -- Go to definition
          vim.keymap.set('n', 'gD', '<cmd>lua vim.lsp.buf.declaration()<cr>', opts)        -- Go to declaration
          vim.keymap.set('n', 'gi', '<cmd>lua vim.lsp.buf.implementation()<cr>', opts)     -- Go to implementation
          vim.keymap.set('n', 'go', '<cmd>lua vim.lsp.buf.type_definition()<cr>', opts)    -- Go to type definition
          vim.keymap.set('n', 'gr', '<cmd>lua vim.lsp.buf.references()<cr>', opts)         -- Find references
          vim.keymap.set('n', 'gs', '<cmd>lua vim.lsp.buf.signature_help()<cr>', opts)     -- Show signature help
          vim.keymap.set('n', '<D-r>', '<cmd>lua vim.lsp.buf.rename()<cr>', opts)          -- Rename symbol
          vim.keymap.set({'n', 'x'}, '<D-f>', '<cmd>lua vim.lsp.buf.format({async = true})<cr>', opts)  -- Format document
          vim.keymap.set('n', '<D-a>', '<cmd>lua vim.lsp.buf.code_action()<cr>', opts)     -- Show code actions
        end,
      })

      require('mason-lspconfig').setup({
        ensure_installed = {
          'gopls',
          'ts_ls',
          --'jdtls',
          'pylsp',
          'sqlls',
          'dockerls',
          'docker_compose_language_service',
          'ast_grep',
          'rust_analyzer',
          'clangd'
        },
        handlers = {
          function(server_name)
            -- Configure gopls with semantic tokens
            if server_name == "gopls" then
              require('lspconfig').gopls.setup({
                settings = {
                  gopls = {
                    semanticTokens = true,  -- Enable semantic tokens for gopls
                  },
                },
              })
            else
              require('lspconfig')[server_name].setup({})
            end
          end,
        }
      })
    end
  }
}

