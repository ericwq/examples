-- This is the chadrc file , its supposed to be placed in /lua/custom dir
-- lua/custom/chadrc.lua

local M = {}

-- make sure you maintain the structure of `core/default_config.lua` here,

M.plugins = {
   status = {
      dashboard = true,
   },
   options = {
      lspconfig = {
         setup_lspconf = "custom.plugins.lspconfig",
      },
   },
   default_plugin_config_replace = {
      nvim_treesitter = "custom.plugins.treesitter",
      dashboard = "custom.plugins.dashboard",
   },
}

-- changing theme and UI
M.ui = {
   theme = "gruvchad",
   italic_comments = true,
   transparency =true,
}

return M
