" Copyright 2011 The Go Authors. All rights reserved.
" Use of this source code is governed by a BSD-style
" license that can be found in the LICENSE file.
"
" indent/go.vim: Vim indent file for Go.
"
" TODO:
" - function invocations split across lines
" - general line splits (line ends in an operator)

if exists("b:did_indent")
  finish
endif
let b:did_indent = 1

" C indentation is too far off useful, mainly due to Go's := operator.
" Let's just define our own.
setlocal nolisp
setlocal autoindent
setlocal indentexpr=s:GoIndent(v:lnum)
setlocal indentkeys+=<:>,0=},0=)

if exists("*GoIndent")
  finish
endif

function! s:GoIndent(lnum) abort
  let l:prevlnum = prevnonblank(a:lnum-1)
  if l:prevlnum == 0
    " top of file
    return 0
  endif

  " grab the previous and current line, stripping comments.
  let l:prevl = substitute(getline(l:prevlnum), '//.*$', '', '')
  let l:thisl = substitute(getline(a:lnum), '//.*$', '', '')
  let l:previ = indent(l:prevlnum)

  let l:ind = l:previ

  for l:synid in synstack(a:lnum, 1)
    if synIDattr(l:synid, 'name') == 'goRawString'
      if l:prevl =~ '\%(\%(:\?=\)\|(\|,\)\s*`[^`]*$'
        " previous line started a multi-line raw string
        return 0
      endif
      " return -1 to keep the current indent.
      return -1
    endif
  endfor

  if l:prevl =~ '[({]\s*$'
    " previous line opened a block
    let l:ind += shiftwidth()
  endif
  if l:prevl =~# '^\s*\(case .*\|default\):$'
    " previous line is part of a switch statement
    let l:ind += shiftwidth()
  endif
  " TODO: handle if the previous line is a label.

  if l:thisl =~ '^\s*[)}]'
    " this line closed a block
    let l:ind -= shiftwidth()
  endif

  " Colons are tricky.
  " We want to outdent if it's part of a switch ("case foo:" or "default:").
  " We ignore trying to deal with jump labels because (a) they're rare, and
  " (b) they're hard to disambiguate from a composite literal key.
  if l:thisl =~# '^\s*\(case .*\|default\):$'
    let l:ind -= shiftwidth()
  endif

  return l:ind
endfunction

" vim: sw=2 ts=2 et
