-----------------------------------------------------------
-- Neovim LSP configuration file
-----------------------------------------------------------

-- Plugin: nvim-lspconfig
-- for language server setup see: https://github.com/neovim/nvim-lspconfig
--

local M = {}

M.setup_lsp = function(attach, capabilities)
   local lspconfig = require "lspconfig"

   --[[
   -- lspservers with default config
   -- https://github.com/neovim/nvim-lspconfig/blob/master/doc/server_configurations.md

   local servers = { "gopls", "clangd", "sumneko_lua" }

   -- Set settings for language servers below

   local ts_settings = function(client)
     client.resolved_capabilities.document_formatting = true
     --ts_settings(client)
   end

   for _, lsp in ipairs(servers) do
      lspconfig[lsp].setup {
         on_attach = attach,
         capabilities = capabilities,
	 ts_settings = ts_settings,
         flags = {
            debounce_text_changes = 150,
         },
      }
   end
--]]

lspconfig.efm.setup {
    init_options = {documentFormatting = true},
    filetypes = {"lua"},
    settings = {
        rootMarkers = {".git/"},
        languages = {
            lua = {
                {
                    formatCommand = "lua-format -i --no-keep-simple-function-one-line --no-break-after-operator --column-limit=150 --break-after-table-lb",
                    formatStdin = true
                }
            }
        }
    }
}

   --https://github.com/neovim/nvim-lspconfig/blob/master/doc/server_configurations.md#sumneko_lua
   local runtime_path = vim.split(package.path, ';')
 table.insert(runtime_path, "lua/?.lua")
table.insert(runtime_path, "lua/?/init.lua")

lspconfig.sumneko_lua.setup {
        on_attach = function(client, bufnr)
         client.resolved_capabilities.document_formatting = false
         vim.api.nvim_buf_set_keymap(bufnr, "n", "<space>fm", "<cmd>lua vim.lsp.buf.formatting()<CR>", {})
       end,
  settings = {
    Lua = {
      runtime = {
        -- Tell the language server which version of Lua you're using (most likely LuaJIT in the case of Neovim)
        version = 'LuaJIT',
        -- Setup your lua path
        path = runtime_path,
      },
      diagnostics = {
        -- Get the language server to recognize the `vim` global
        globals = {'vim'},
      },
      workspace = {
        -- Make the server aware of Neovim runtime files
        library = vim.api.nvim_get_runtime_file("", true),
      },
      -- Do not send telemetry data containing a randomized but unique identifier
      telemetry = {
        enable = false,
      },
    },
  },
}
--[[
--https://github-wiki-see.page/m/jose-elias-alvarez/null-ls.nvim/wiki/Avoiding-LSP-formatting-conflicts
   -- typescript

 lspconfig.tsserver.setup {
      on_attach = function(client, bufnr)
         client.resolved_capabilities.document_formatting = false
         vim.api.nvim_buf_set_keymap(bufnr, "n", "<space>fm", "<cmd>lua vim.lsp.buf.formatting()<CR>", {})
      end,
   }

-- the above tsserver config will remvoe the tsserver's inbuilt formatting
-- since I use null-ls with denofmt for formatting ts/js stuff.
-- ]]

end

return M
