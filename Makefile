out:
	mkdir -p out

Gauges.app: out
	rm -rf ./out/Gauges.app
	cp -r ./internal/tray/assets/Gauges.app ./out/Gauges.app
	CGO_ENABLED=1 go build -o ./out/Gauges.app/Contents/gauges_tray ./cmd
