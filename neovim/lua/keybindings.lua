vim.g.mapleader = ","
vim.g.maplocalleader = ","

local map = vim.api.nvim_set_keymap
local opt = {
    noremap = true,
    silent = true
}

-- nvimTree
map('n', '<C-n>', ':NvimTreeToggle<CR>', opt)
