let SessionLoad = 1
let s:so_save = &g:so | let s:siso_save = &g:siso | setg so=0 siso=0 | setl so=-1 siso=-1
let v:this_session=expand("<sfile>:p")
silent only
silent tabonly
cd ~/work/pingflow
if expand('%') == '' && !&modified && line('$') <= 1 && getline(1) == ''
  let s:wipebuf = bufnr('%')
endif
let s:shortmess_save = &shortmess
if &shortmess =~ 'A'
  set shortmess=aoOA
else
  set shortmess=aoO
endif
badd +6 frontend/src/App.vue
badd +10 frontend/package.json
badd +9 frontend/src/components/HelloWorld.vue
badd +6 frontend/vite.config.ts
badd +10 frontend/.prettierrc
badd +11 frontend/src/main.ts
badd +13 frontend/src/stores.ts
badd +29 backend/main.ts
badd +3 backend/deno.json
badd +6 redis.Dockerfile
badd +1054 redis.conf
badd +15 docker-compose.yml
badd +10 .gitignore
badd +2 worker/main.go
badd +6 worker/Dockerfile
badd +1 frontend/env.d.ts
badd +10 .env.public
badd +10 README.md
badd +37 worker/http_server.go
badd +17 Justfile
badd +4 docker-compose-dev.yml
badd +40 worker/redis.go
badd +1 .env
badd +11 worker/api.go
badd +24 worker/utils/env.go
badd +12 worker/utils/work_index.go
badd +2 .ignore
badd +28 worker/cache.go
badd +17 worker/processing.go
badd +29 worker/integration_test.go
badd +37 backend/integration_test.ts
argglobal
%argdel
edit backend/main.ts
argglobal
balt backend/integration_test.ts
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let &fdl = &fdl
let s:l = 29 - ((28 * winheight(0) + 27) / 55)
if s:l < 1 | let s:l = 1 | endif
keepjumps exe s:l
normal! zt
keepjumps 29
normal! 0
tabnext 1
if exists('s:wipebuf') && len(win_findbuf(s:wipebuf)) == 0 && getbufvar(s:wipebuf, '&buftype') isnot# 'terminal'
  silent exe 'bwipe ' . s:wipebuf
endif
unlet! s:wipebuf
set winheight=1 winwidth=20
let &shortmess = s:shortmess_save
let s:sx = expand("<sfile>:p:r")."x.vim"
if filereadable(s:sx)
  exe "source " . fnameescape(s:sx)
endif
let &g:so = s:so_save | let &g:siso = s:siso_save
doautoall SessionLoadPost
unlet SessionLoad
" vim: set ft=vim :
