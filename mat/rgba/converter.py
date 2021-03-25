import json

src = open("colors.json", "r")
dist = open("colors.go", "w")

dist.write("""package rgba

import "github.com/jakubDoka/mlok/mat"

// color "constants"
var (
""")
data = json.load(src)
for i in data:
    dat = data[i]
    dist.write(f'{i} = mat.RGB({dat[0]}, {dat[1]}, {dat[2]})\n')
dist.write(")\n\n")
dist.write("// Colors contains colors with simple names mapped to the string literals\n")
dist.write("var Colors = map[string]mat.RGBA{\n")
for i in data:
    if sum(1 for c in i if c.isupper()) > 1:
        continue
    dist.write(f'"{i.lower()}": {i},\n')
dist.write("}")
dist.close()
