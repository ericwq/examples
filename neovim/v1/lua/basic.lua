-- utf8
vim.g.encoding = "UTF-8"
vim.o.fileencoding = 'utf-8'
-- jk移动时光标下上方保留8行
vim.o.scrolloff = 8 -- Lines of context
-- Round indent
vim.o.sidescrolloff = 8 -- Columns of context
-- 使用相对行号
vim.wo.number = true -- Show line numbers
--vim.wo.relativenumber = true -- Relative line numbers
-- 高亮所在行
vim.wo.cursorline = true
-- 显示左侧图标指示列
vim.wo.signcolumn = "yes"
-- 右侧参考线，超过表示代码太长了，考虑换行
-- vim.wo.colorcolumn = "80"
-- 缩进2个空格等于一个Tab
vim.o.tabstop = 4 -- Number of spaces tabs count for
vim.bo.tabstop = 4
vim.o.softtabstop = 4
vim.o.shiftround = true
-- >> << 时移动长度
vim.o.shiftwidth = 2 -- Size of an indent
vim.bo.shiftwidth = 2
-- 新行对齐当前行，空格替代tab
vim.o.expandtab = true
vim.bo.expandtab = true
vim.o.autoindent = true
vim.bo.autoindent = true
vim.o.smartindent = true -- Insert indents automatically
-- 搜索大小写不敏感，除非包含大写
vim.o.ignorecase = true -- Ignore case
vim.o.smartcase = true -- Do not ignore case with capitals
vim.o.hlsearch = false -- 搜索不要高亮
vim.o.incsearch = true -- 边输入边搜索
vim.o.showmode = false -- 使用增强状态栏后不再需要 vim 的模式提示
vim.o.cmdheight = 2 -- 命令行高为2，提供足够的显示空间
vim.o.autoread = true -- 当文件被外部程序修改时，自动加载
vim.bo.autoread = true
vim.o.wrap = false -- Disable line wrap
vim.wo.wrap = false -- 禁止折行
vim.o.whichwrap = 'b,s,<,>,[,],h,l' -- 行结尾可以跳到下一行
-- 允许隐藏被修改过的buffer
vim.o.hidden = true --Enable background buffers
vim.o.mouse = "a" -- 鼠标支持
vim.o.backup = false -- 禁止创建备份文件
vim.o.writebackup = false
vim.o.swapfile = false
vim.o.updatetime = 300 -- smaller updatetime 
vim.o.timeoutlen = 100 -- 等待mappings
-- split window 从下边和右边出现
vim.o.splitbelow = true -- Put new windows below current
vim.o.splitright = true -- Put new windows right of current
vim.g.completeopt = "menu,menuone,noselect,noinsert" -- 自动补全不自动选中
vim.o.background = "dark" -- 样式
vim.o.termguicolors = true -- True color support
vim.opt.termguicolors = true
vim.o.list = true -- 不可见字符的显示，这里只把空格显示为一个点
--vim.o.listchars = "space:·" -- Show some invisible characters
vim.o.wildmenu = true -- 补全增强
vim.o.shortmess = vim.o.shortmess .. 'c' -- Dont' pass messages to |ins-completin menu|
vim.o.pumheight = 10
vim.o.showtabline = 2 -- always show tabline

--opt.completeopt = {'menuone', 'noinsert', 'noselect'}  -- Completion options (for deoplete)
--opt.expandtab = true                -- Use spaces instead of tabs
--opt.hidden = true                   
--opt.ignorecase = true               
vim.opt.joinspaces = false              -- No double spaces with join
--opt.list = true                     
--opt.number = true                   
-- opt.relativenumber = true           
--vim.opt.scrolloff = 4                   
--opt.shiftround = true               
----opt.shiftwidth = 4                  
--opt.sidescrolloff = 8               
--opt.smartcase = true                
--opt.smartindent = true              
--opt.splitbelow = true               
--opt.splitright = true               
--opt.tabstop = 4                     
--opt.termguicolors = true            
vim.opt.wildmode = {'list', 'longest'}  
--opt.wrap = false                    
vim.opt.laststatus=2

