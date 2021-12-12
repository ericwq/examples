-------------------- Neovim -------------------------------
vim.g.loaded_python_provider  = 0       -- disable python 2 provider
vim.g.loaded_python3_provider = 0       -- disable python 3 provider
vim.g.loaded_ruby_provider    = 0       -- disable ruby provider
vim.g.loaded_perl_provider    = 0       -- disable perl provider
--vim.g.python3_host_prog       = '/usr/bin/python3'

-------------------- Plugins -------------------------------
require('basic')
require('keybindings')
require('plugins')


-------------------- Color scheme -------------------------------
vim.g.material_theme_style ='default'   --  default, palenight, ocean, lighter, and darker.
vim.g.material_terminal_italics =1
vim.cmd 'colorscheme material'            -- Put your favorite colorscheme here

--vim.g.sonokai_style = 'andromeda'
--vim.g.sonokai_enable_italic = 1
--vim.g.sonokai_disable_italic_comment = 1
--vim.cmd 'colorscheme sonokai'            -- Put your favorite colorscheme here

-------------------- Plugin setup -------------------------------
require('plugins/nvim-tree')
