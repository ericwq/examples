local null_ls = require "null-ls"
local b = null_ls.builtins

local sources = {
-- If you have performance issues with a diagnostic source, you can configure any it to run on save (not on each change) by overriding method:
-- method = null_ls.methods.DIAGNOSTICS_ON_SAVE,
--
-- should echo 1 if available (and 0 if not)
-- :echo executable("eslint")
--
   -- others
   b.formatting.prettierd.with({ filetypes = { "html", "markdown", "css" ,"yaml", "json", "javascript" }, }),

   -- go
   b.formatting.goimports,
   b.formatting.gofmt,
   b.diagnostics.golangci_lint.with({ diagnostics_format = "[#{c}] #{m} (#{s})", }),

   -- english text
   b.diagnostics.proselint.with({ diagnostics_format = "[#{c}] #{m} (#{s})", }),
   b.completion.spell,

   -- c/c++
   b.formatting.clang_format,
   b.diagnostics.cppcheck.with({ diagnostics_format = "[#{c}] #{m} (#{s})", }),

--[[
--
   -- Lua
   b.formatting.stylua,
   b.diagnostics.luacheck.with { extra_args = { "--global vim" } },

   -- Shell
   b.formatting.shfmt,
   b.diagnostics.shellcheck.with { diagnostics_format = "#{m} [#{c}]" },
--]]
}

local M = {}

M.setup = function()
   null_ls.setup {
      debug = true,
      sources = sources,

      -- format on save
      on_attach = function(client)
         if client.resolved_capabilities.document_formatting then
            vim.cmd "autocmd BufWritePre <buffer> lua vim.lsp.buf.formatting_sync()"
         end
      end,
   }
end

return M
