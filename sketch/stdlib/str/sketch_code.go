package str

// Code generated by sketch/scripts/bind-module-data DO NOT EDIT
//
// Sources:
// string.skt

const SketchCode = `
; string.skt
(defn
  join
  "Returns a new string made by concatenating the items in 'elements',
    placing 'separator' between each one."
  (elements separator)
  (cond
    ((empty? elements) "") ; Special behaviour if called with an empty list
    ((empty? (rest elements)) (first elements)) ; Recursion base case
    ("else" (+ (first elements) separator (join (rest elements) separator)))))

(export-as string (join))
`
