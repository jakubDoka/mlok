inp = open("keys.txt")
dest = open("keys.go", "w")
last = ""
values = []
for line in inp:
	if line.strip() == "" or "KEY_KP" in line:
		continue
	name, code = line.split("   ")
	trueName = name.replace("_", " ").title().replace(" ", "").replace("Glfw", "")
	name = trueName.replace("Key", "")
	if name[0].isdigit():
    		name = "_" + name
	comment = ""
	if "/*" in code:
		code, comment = code.split("/*")
		comment = comment[:-3]
	last = trueName
	values.append((trueName, name, comment))
values.append((last, "Last", ""))
inp.close()

dest.write("""package key

import "github.com/go-gl/glfw/v3.3/glfw"

// Key can be any key, keyboard or mouse 
type Key int

// all key constants
const (
""")

for (trueName, name, comment) in values:
	dest.write(f'\t{name} = Key(glfw.{trueName}) //{comment}\n')

dest.write(""")

func (k Key) String() string {
	val, ok := Names[k]
	if !ok {
		return "Invalid"
	}

	return val
}

// Names is helper for Key.String() method
var Names = map[Key]string{
""")

for (trueName, name, comment) in values:
	b = False
	for i in "MouseButtonLast", "MouseButtonLeft", "MouseButtonRight", "MouseButtonMiddle", "Last":
		if i == name:
			b = True
			break
	if b:
		continue
	dest.write(f'\t{name}:"{name.replace("_", "")}",\n')

dest.write("}")

