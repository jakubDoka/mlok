import json
src = open("colors.json")
dist = open("colors.go", "x")

dist.write("""package rgba

import "gobatch/mat"

// color "constants"
var (
""")
data = json.load(src)
for i in data:
    name = i.replace("_", " ").title().replace(" ", "")
    dat = data[i]["rgb"]
    for i in range(len(dat)):
        dat[i] = round(dat[i]/255, 4)
    dist.write(f'{name} = mat.RGB({dat[0]}, {dat[1]}, {dat[2]})\n')
dist.write(")\n\n")
dist.write("// Colors contains colors with simple names mapped to the string literals\n")
dist.write("var Colors = map[string]mat.RGBA{\n")
for i in data:
    if "_" in i:
        continue
    dist.write(f'"{i.lower()}": {i.title()},\n')
dist.write("}")
dist.close()
