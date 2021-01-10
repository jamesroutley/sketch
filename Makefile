.PHONY: linecount

linecount:
	scc --no-cocomo --count-as "skt:clj" | sed 's/Clojure/Sketch /g'

generate:
	go run scripts/bind-module-data/main.go

format:
	find sketch -type f -name "*.skt" -exec go run main.go format -w {} \;
