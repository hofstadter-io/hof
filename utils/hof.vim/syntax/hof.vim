if exists("b:current_syntax")
  finish
endif

syntax case match

hi HofGold ctermfg=215 guifg=#cba44f
hi HofGoldUnder cterm=underline ctermfg=215 gui=underline guifg=#cba44f

syn match     hofDefO /^[a-zA-Z0-9_-]\+\>\ze\s:/
hi def link hofDefO HofGold
syn match     hofDefC /^[a-zA-Z0-9_-]\+\>\ze\s::/
hi def link hofDefC HofGoldUnder

syn match     hofGen /^[a-zA-Z0-9_-]\+\>\ze\s<-\s/
hi def link hofGen Title
syn match hofGenDef /<-/
hi def link hofGenDef Title

syn match     hofType /[a-zA-Z0-9_-]\+\>\ze\./
hi def link hofType Tag

syn match hofInt /\<[0-9]\+\>/
syn match hofDec /\<[0-9]\+\.[0-9]\+\>/
hi def link hofInt Number
hi def link hofDec Number

syn match hofBool /\<true\|false\>/
hi def link hofBool Number

syn match     hofPackage           /^package\>/
syn match     hofImport            /^import\>/

hi def link     hofPackage           Statement
hi def link     hofImport            Statement

syn match hofEllipse "\.\.\."
hi def link     hofEllipse           Question

syn match hofOper /:/
syn match hofOper /::/
syn match hofOper "&"
syn match hofOper "|"

hi def link hofOper Function

syn region      hofString            start=+"+ skip=+\\\\\|\\"+ end=+"+
syn region      hofRawString         start=+`+ end=+`+

hi def link     hofString            String
hi def link     hofRawString         String

" Comments; their contents
syn keyword     hofTodo              contained TODO FIXME XXX BUG Deprecated
syn cluster     hofCommentGroup      contains=@hofTodo

syn region      hofComment           start="//" end="$" contains=@hofCommentGroup
syn region      hofComment           start="/\*" end="\*/" contains=@hofCommentGroup

hi def link     hofComment           Comment
hi def link     hofTodo              Todo

let b:current_syntax = "hof"

" vim: sw=2 ts=2 et

