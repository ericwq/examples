-- This is an example chadrc file , its supposed to be placed in /lua/custom dir
-- lua/custom/chadrc.lua

local M = {}

-- make sure you maintain the structure of `core/default_config.lua` here,
-- example of changing theme:

M.plugins = {
   default_plugin_config_replace = {
      nvim_treesitter = "custom.treesitter",
   },
}

M.ui = {
   theme = "gruvchad",
   italic_comments = true,
   transparency =true,
}

return M
