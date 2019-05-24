graphviz:
	sudo apt install graphviz

cpu:
	go tool pprof --pdf cpu.prof > cpu.pdf

mem:
	go tool pprof --pdf -alloc_space mem.prof > mem.pdf

bench:
	go test -cpuprofile cpu.prof -memprofile mem.prof -bench .

ui-cpu:
	go tool pprof -http=:8080 cpu.prof

ui-mem:
	go tool pprof -http=:8080 mem.prof
