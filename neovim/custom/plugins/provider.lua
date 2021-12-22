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
