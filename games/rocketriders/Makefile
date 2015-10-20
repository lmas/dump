
love=/usr/bin/love
zip=/usr/bin/zip
unzip=/usr/bin/unzip
luac=/usr/bin/luac

builddir=build
distdir=dist

windir=love-0.8.0-win-x86
winzip=love-0.8.0-win-x86.zip
osxapp=love.app
osxzip=love-0.8.0-macosx-ub.zip

game=rocketriders
sources=*.lua engine/*.lua
res=assets/* libs/*

.PHONY : run test love clean win

rungame : test
	$(love) .

runserver : test
	lua server.lua

test :
	$(luac) -p $(sources)

dist : love win osx

love : $(builddir)/$(game).love
	cp $(builddir)/$(game).love $(distdir)/$(game).love

osx : $(builddir)/$(game).app
	cd $(builddir); \
		zip -9 -q -r ../$(distdir)/$(game).osx.zip $(game).app

win : $(builddir)/$(game).exe
	cd $(builddir); \
		cp $(windir)/*.dll .; \
		zip -q ../$(distdir)/$(game).win.zip $(game).exe *.dll; \
		rm *.dll

$(builddir)/$(game).app : $(builddir)/$(game).love
	cd $(builddir); \
	unzip ../$(osxzip); \
	cp -a $(osxapp) $(game).app; \
	cp $(game).love $(game).app/Contents/Resources/; \
	sed -i.bak 's/<string>LÖVE<\/string>/<string>$(game)<\/string>/g' "$(game).app/Contents/Info.plist"; \
	sed -i.bak 's/<string>org\.love2d\.love<\/string>/<string>net\.$(game)<\/string>/g' "$(game).app/Contents/Info.plist"; \
	sed -i.bak '/<key>UTExportedTypeDeclarations<\/key>/,/^\t<\/array>/{d}' "$(game).app/Contents/Info.plist"

$(builddir)/$(game).exe : $(builddir)/$(game).love
	cd $(builddir); \
	unzip ../$(winzip); \
	cat $(windir)/love.exe $(game).love > $(game).exe

$(builddir)/$(game).love : $(sources) $(res)
	mkdir -p $(builddir)
	mkdir -p $(distdir)
	$(zip) -r $(builddir)/$(game).love $(sources) $(res)

clean :
	rm -rf $(distdir)/ $(builddir)/
