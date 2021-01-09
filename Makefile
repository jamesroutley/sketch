.PHONY: linecount

linecount:
	scc --no-cocomo --count-as "skt:clj" | sed 's/Clojure/Sketch /g'
