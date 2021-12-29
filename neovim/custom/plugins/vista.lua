-----------------------------------------------------------
-- Vista configuration file
-----------------------------------------------------------

-- Plugin: vista.vim
-- https://github.com/liuchengxu/vista.vim


local g = vim.g
local cmd = vim.cmd

g.vista_close_on_jump = 1

-- How each level is indented and what to prepend.
--- This could make the display more compact or more spacious.
--- e.g., more compact: ["▸ ", ""]
--- Note: this option only works for the kind renderer, not the tree renderer
g.vista_icon_indent = '["╰─▸ ", "├─▸ "]'

-- Executive used when opening vista sidebar without specifying it.
--- See all the avaliable executives via `:echo g:vista#executives`.
g.vista_default_executive = 'ctags'

-- Set the executive for some filetypes explicitly. Use the explicit executive
-- instead of the default one for these filetypes when using `:Vista` without
-- specifying the executive.
--[[
g.vista_executive_for = {
   vimwiki = "markdown",
   pandoc = "markdown",
   markdown = "toc",
   terraform = "nvim_lsp",
   rust = "nvim_lsp",
   c = "nvim_lsp",
   cpp = "nvim_lsp",
   go = "nvim_lsp",
}
--]]
-- Ensure you have installed some decent font to show these pretty symbols,
--- then you can enable icon for the kind.
g["vista#renderer#enable_icons"] = 1
--cmd [[let g:vista#renderer#enable_icon = 1]]


-- Change some default icons
--- see: https://github.com/slavfox/Cozette/blob/master/img/charmap.txt
--[[
local t ={
["function"]  = "\u0192"
["variable"]  = "\uf00d"
["prototype"] = "\uf013"
["macro"]     = "\uf00b"
}
g["vista#renderer#icons"] = t
--]]

cmd [[
  let g:vista#renderer#icons = {
  \   "function": "\u0192",
  \   "variable": "\uf00d",
  \   "prototype": "\uf013",
  \   "macro": "\uf00b",
  \ }
]]
