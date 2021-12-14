-----------------------------------------------------------
-- Treesitter configuration file
-----------------------------------------------------------

-- Plugin: nvim-treesitter
-- https://github.com/nvim-treesitter/nvim-treesitter
-- :TSInstallInfo


require('nvim-treesitter.configs').setup {
  -- one of "all", "maintained" (parsers with maintainers), or a list of languages
  ensure_installed = {"go", "c", "cpp", "yaml", "lua", "json", "dockerfile", "vim", "javascript"},
  sync_install = false, -- install languages synchronously (only applied to `ensure_installed`)
  --ignore_install = { "javascript" }, -- List of parsers to ignore installing
  highlight = {
    enable = true,              -- false will disable the whole extension
    --disable = { "c", "rust" },  -- list of language that will be disabled
    -- Setting this to true will run `:h syntax` and tree-sitter at the same time.
    -- Set this to `true` if you depend on 'syntax' being enabled (like for indentation).
    -- Using this option may slow down your editor, and you may see some duplicate highlights.
    -- Instead of true it can also be a list of languages
    additional_vim_regex_highlighting = false,
  },

  -- enable Indentation
  indent = {enable = false },

  -- enable Incremental selection
  incremental_selection = {
    enable = true,
    keymaps = {
      init_selection = "gnn",
      node_incremental = "grn",
      scope_incremental = "grc",
      node_decremental = "grm",
--[[
      init_selection = '<CR>',
      node_incremental = '<CR>',
      node_decremental = '<BS>',
      scope_incremental = '<TAB>',
--]]
    },
  },
}

-- enable Folding
--vim.wo.foldmethod = 'expr'
--vim.wo.foldexpr = 'nvim_treesitter#foldexpr()'
