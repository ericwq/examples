-----------------------------------------------------------
-- File manager configuration file
-----------------------------------------------------------

-- Plugin: nvim-tree
-- https://github.com/kyazdani42/nvim-tree.lua

-- Keybindings are defined in `keymapping.lua`:
--- https://github.com/kyazdani42/nvim-tree.lua#keybindings

-- Note: options under the g: command should be set BEFORE running the
--- setup function:
--- https://github.com/kyazdani42/nvim-tree.lua#setup
--- See: `help NvimTree`
local g = vim.g

g.nvim_tree_quit_on_open = 1              --0 by default, closes the tree when you open a file
g.nvim_tree_indent_markers = 1            --0 by default, this option shows indent markers when folders are open
g.nvim_tree_git_hl = 1                    --0 by default, will enable file highlight for git attributes (can be used without the icons).
g.nvim_tree_highlight_opened_files = 1    --0 by default, will enable folder and file icon highlight for opened files/directories.
g.nvim_tree_disable_window_picker = 1     --0 by default, will disable the window picker.
g.nvim_tree_respect_buf_cwd = 1           --0 by default, will change cwd of nvim-tree to that of new buffer's when opening nvim-tree.
--g.nvim_tree_width_allow_resize  = 1

g.nvim_tree_show_icons = {
  git = 1,
  folders = 1,
  files = 1,
  folder_arrows = 1,
}
--If 0, do not show the icons for one of 'git' 'folder' and 'files'
--1 by default, notice that if 'files' is 1, it will only display
--if nvim-web-devicons is installed and on your runtimepath.
--if folder is 1, you can also tell folder_arrows 1 to show small arrows next to the folder icons.
--but this will not work when you set indent_markers (because of UI conflict)

g.nvim_tree_icons = {
	default = "‣ "
}

-- following options are the default
-- each of these are documented in `:help nvim-tree.OPTION_NAME`
require('nvim-tree').setup {
  open_on_setup = false,
  update_cwd = true,
  view = {
    width = 32,
    auto_resize = true
  },
  filters = {
    dotfiles = true,
    custom = { '.git', 'node_modules', '.cache', '.bin' },
  },
  git = {
    enable = true,
    ignore = true,
  },
}
