" gomod.vim: Vim filetype plugin for Go module file.
"
" based by from https://github.com/fatih/vim-go/blob/79ea9ef26807eda0b55809d0521993bcecfa09e5/ftplugin/gomod.vim

if exists("b:did_ftplugin")
  finish
endif
let b:did_ftplugin = 1

setlocal tabstop=2
setlocal softtabstop=2
setlocal shiftwidth=2

setlocal comments=s1:/*,mb:*,ex:*/,://
setlocal commentstring=//\ %s

" vim: sw=2 ts=2 et
