-----------------------------------------------------------
-- Color schemes configuration file
-----------------------------------------------------------

-- Import colorscheme with:
--- require('colors')


--[[
-- Theme: monokai (classic)
--- See: https://github.com/tanvirtin/monokai.nvim/blob/master/lua/monokai.lua
local _M = {
  bg = '#202328', --default: #272a30
  fg = '#f8f8f0',
  pink = '#f92672',
  green = '#a6e22e',
  cyan = '#78dce8',
  yellow = '#e6db74',
  orange = '#fa8419',
  purple = '#9c64fe',
  red = '#ed2a2a',
}

return _M


--]]
-- Theme: Rosé Pine
--- See: https://github.com/rose-pine/neovim#custom-colours
--- color names are adapted to the format above
local _M = {
  bg = '#111019',
  fg = '#e0def4',
  pink = '#eb6f92',
  green = '#1f1d2e',
  cyan = '#31748f',
  yellow = '#f6c177',
  orange = '#2a2837',
  purple = '#c4a7e7',
  red = '#ebbcba',
}

return _M

--[[
-------------------- Color scheme -------------------------------
vim.g.material_theme_style ='default'   --  default, palenight, ocean, lighter, and darker.
vim.g.material_terminal_italics =1
vim.cmd 'colorscheme material'            -- Put your favorite colorscheme here

--vim.g.sonokai_style = 'andromeda'
--vim.g.sonokai_enable_italic = 1
--vim.g.sonokai_disable_italic_comment = 1
--vim.cmd 'colorscheme sonokai'            -- Put your favorite colorscheme here
--]]
