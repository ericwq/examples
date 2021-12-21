-- This is an example init file , its supposed to be placed in /lua/custom dir
-- lua/custom/init.lua

-- This is where your custom modules and plugins go.
-- Please check NvChad docs if you're totally new to nvchad + dont know lua!!

local hooks = require "core.hooks"

-- MAPPINGS
-- To add new plugins, use the "setup_mappings" hook,

hooks.add("setup_mappings", function(map)

   -- Vista tag-viewer
   map('n', '<C-m>', ':Vista!!<CR>', opt)   -- open/close

   map("n", "<leader>cc", ":Telescope <CR>", opt)
   map("n", "<leader>q", ":q <CR>", opt)
end)

-- NOTE : opt is a variable  there (most likely a table if you want multiple options),
-- you can remove it if you dont have any custom options

-- Install plugins
-- To add new plugins, use the "install_plugin" hook,

-- examples below:

hooks.add("install_plugins", function(use)

  -- tagviewer
  use {
      'liuchengxu/vista.vim',
      event = "BufRead",
      --run before this plugin is loaded.
      --setup =
      --run after this plugin is loaded.
      config = function()
         require("custom.vista")
      end,
   }

--[[

   use {
      "max397574/better-escape.nvim",
      event = "InsertEnter",
   }

   use {
      "user or orgname/reponame",
      --further packer options
   }
]]
end)

-- NOTE: we heavily suggest using Packer's lazy loading (with the 'event' field)
-- see: https://github.com/wbthomason/packer.nvim
-- https://nvchad.github.io/config/walkthrough

-- source a vimscript file
vim.cmd('source ~/.config/nvim/vimrc')

-----------------------------------------------------------
-- Neovim provider
-----------------------------------------------------------
vim.g.loaded_python_provider  = 0       -- disable python 2 provider
vim.g.loaded_python3_provider = 0       -- disable python 3 provider
vim.g.loaded_ruby_provider    = 0       -- disable ruby provider
vim.g.loaded_perl_provider    = 0       -- disable perl provider
--vim.g.python3_host_prog       = '/usr/bin/python3'
