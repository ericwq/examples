--[[

  ███╗   ██╗███████╗ ██████╗ ██╗   ██╗██╗███╗   ███╗
  ████╗  ██║██╔════╝██╔═══██╗██║   ██║██║████╗ ████║
  ██╔██╗ ██║█████╗  ██║   ██║██║   ██║██║██╔████╔██║
  ██║╚██╗██║██╔══╝  ██║   ██║╚██╗ ██╔╝██║██║╚██╔╝██║
  ██║ ╚████║███████╗╚██████╔╝ ╚████╔╝ ██║██║ ╚═╝ ██║
  ╚═╝  ╚═══╝╚══════╝ ╚═════╝   ╚═══╝  ╚═╝╚═╝     ╚═╝


Neovim init file
Version: 0.42.0 - 2021/12/01
Maintainer: Brainf+ck
Website: https://github.com/brainfucksec/neovim-lua

--]]

-----------------------------------------------------------
-- Neovim provider
-----------------------------------------------------------
vim.g.loaded_python_provider  = 0       -- disable python 2 provider
vim.g.loaded_python3_provider = 0       -- disable python 3 provider
vim.g.loaded_ruby_provider    = 0       -- disable ruby provider
vim.g.loaded_perl_provider    = 0       -- disable perl provider
--vim.g.python3_host_prog       = '/usr/bin/python3'

-----------------------------------------------------------
-- Import Lua modules
-----------------------------------------------------------
require('settings')
require('keymaps')
require('plugins/packer')
require('plugins/nvim-tree')
require('plugins/feline')
require('plugins/indent-blankline')
require('plugins/nvim-treesitter')
require('plugins/nvim-treesitter-context')
require('plugins/nvim-cmp')
require('plugins/nvim-lspconfig')
require('plugins/symbols-outline')
require('plugins/alpha-nvim')
--[[
require('plugins/vista')
--]]
