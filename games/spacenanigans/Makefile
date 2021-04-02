SRCJS = js/*.js
BUILTJS = static/game.js
SRCGO = cmd/server/server.go
BUILTGO64 = cmd/server/server64
BUILTGO32 = cmd/server/server32

buildjs: $(SRCJS)
	echo "\"use strict\";" > $(BUILTJS)
	echo "" >> $(BUILTJS)
	echo "// Built: $$(date)" >> $(BUILTJS)
	cat --squeeze-blank $(SRCJS) >> $(BUILTJS)

buildgo64: $(SRCGO)
	env CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o "$(BUILTGO64)" $(SRCGO)

buildgo32: $(SRCGO)
	env CGO_ENABLED=1 GOOS=linux GOARCH=386 go build -o "${BUILTGO32}" ${SRCGO}

buildall: buildjs buildgo64 buildgo32

run: buildall
	./$(BUILTGO64)

clean: $(BUILTJS) $(BUILTGO64) $(BUILTGO32)
	rm "$(BUILTJS)"
	rm "$(BUILTGO64)"
	rm "$(BUILTGO32)"
