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
badd +98 backend/main.ts
badd +5 backend/deno.json
badd +6 redis.Dockerfile
badd +1054 redis.conf
badd +31 docker-compose.yml
badd +25 .gitignore
badd +50 worker/main.go
badd +6 worker/Dockerfile
badd +9 .env.public
badd +10 README.md
badd +19 worker/http_server.go
badd +44 Justfile
badd +4 docker-compose-dev.yml
badd +40 worker/redis.go
badd +1 .env
badd +70 worker/api.go
badd +24 worker/utils/env.go
badd +12 worker/utils/work_index.go
badd +2 .ignore
badd +28 worker/cache.go
badd +86 worker/processing.go
badd +29 worker/integration_test.go
badd +48 backend/integration_test.ts
badd +53 backend/ws_client_test.ts
badd +10 frontend/vite.config.ts
badd +6 frontend/src/main.ts
badd +44 frontend/src/stores/websocket.ts
badd +1 frontend/src/stores/counter.ts
badd +9 frontend/package.json
badd +121 frontend/src/App.vue
badd +6 frontend/src/components/NavBar.vue
badd +7 frontend/src/styles/main.sass
badd +37 frontend/src/styles/base.sass
badd +7 frontend/Dockerfile
badd +3 frontend/.dockerignore
badd +1 backend/Dockerfile
badd +9 .env.docker
badd +1 frontend/src/types/api.ts
badd +25 frontend/src/components/ClassComponent.vue
badd +67 frontend/src/components/SpellsComponent.vue
badd +26 frontend/src/components/SpellInfo.vue
argglobal
%argdel
edit frontend/src/components/SpellInfo.vue
argglobal
balt frontend/src/components/SpellsComponent.vue
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
let s:l = 30 - ((29 * winheight(0) + 27) / 55)
if s:l < 1 | let s:l = 1 | endif
keepjumps exe s:l
normal! zt
keepjumps 30
normal! 04|
lcd ~/work/pingflow
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
nohlsearch
doautoall SessionLoadPost
unlet SessionLoad
" vim: set ft=vim :
