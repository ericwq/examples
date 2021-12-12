return require('packer').startup(function()
  -- Packer can manage itself
  use 'wbthomason/packer.nvim'
    -- gruvbox theme
    use {
        "ellisonleao/gruvbox.nvim",
        requires = {"rktjmp/lush.nvim"}
    }

    -- monokai theme
    use {
      "crusoexia/vim-monokai"
    }

    -- material theme
    use {
      "kaicataldo/material.vim"
    }

    -- nvim-tree
    use {
    'kyazdani42/nvim-tree.lua',
    requires = { 'kyazdani42/nvim-web-devicons',} -- optional, for file icon 
    }

end)
