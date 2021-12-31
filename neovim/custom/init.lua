-- This is the init file , its supposed to be placed in /lua/custom dir
-- lua/custom/init.lua
-- This is where your custom modules and plugins go.
-- Please check NvChad docs if you're totally new to nvchad + dont know lua!!
local hooks = require "core.hooks"

-- MAPPINGS
-- To add new plugins, use the "setup_mappings" hook,
hooks.add("setup_mappings", function(map)

    -- Vista tag-viewer
    map('n', '<C-m>', ':Vista!!<CR>', opt) -- open/close
    -- Searches for the string under your cursor in your current working directory
    map("n", "<leader>fs", ":Telescope grep_string<CR>", opt)
    map("n", "<leader>xx", ":q <CR>", opt)
end)
-- NOTE : opt is a variable  there (most likely a table if you want multiple options),
-- you can remove it if you dont have any custom options

-- Install plugins
-- To add new plugins, use the "install_plugin" hook,
hooks.add("install_plugins", function(use)

    -- A fast and lua alternative to filetype.vim
    -- https://github.com/nathom/filetype.nvim
    -- use ':echo &filetype' to detect the corrrect file type
    -- use `:set filetype=langname` to set file type.
    use {'nathom/filetype.nvim', event = "VimEnter"}

    -- tagviewer
    use {
        'liuchengxu/vista.vim',
        event = "BufRead",
        -- run before this plugin is loaded.
        -- setup =
        -- run after this plugin is loaded.
        config = function()
            require("custom.plugins.vista")
        end
    }

    -- null-ls
    use {
        "jose-elias-alvarez/null-ls.nvim",
        after = "nvim-lspconfig",
        requires = {"nvim-lua/plenary.nvim"},
        config = function()
            require("custom.plugins.null-ls").setup()
        end
    }

    -- treesitter context
    use {
        'romgrk/nvim-treesitter-context',
        after = "nvim-treesitter",
        config = function()
            require("custom.plugins.treesitter-context").setup()
        end
    }

    -- which-key
    use {
        "folke/which-key.nvim",
        config = function()
            require("which-key").setup {
                -- your configuration comes here
                -- or leave it empty to use the default settings
                -- refer to the configuration section below
            }
        end
    }

end)

-- try to call the customized provider

pcall(require, "custom.plugins.provider")
