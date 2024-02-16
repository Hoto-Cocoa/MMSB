for /f "delims=" %%A in ('git rev-parse --show-toplevel') do (cd %%A)
copy /Y assets\syso.json syso.json
%USERPROFILE%\go\bin\syso.exe
del syso.json
go build  -ldflags "-s -w" -o "dist/MMSB.exe"
del out.syso
upx --best --lzma "dist/MMSB.exe"
