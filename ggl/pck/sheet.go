package pck

import (
	"gobatch/ggl"
	"gobatch/mat"
	"image"
	"image/draw"
	"math"
	"strconv"

	"github.com/jakubDoka/gogen/dirs"
	"github.com/jakubDoka/gogen/str"
)

/*imp(
	github.com/jakubDoka/gogen/templates
)*/

/*gen(
	templates.Vec<PicData, Vec>
)*/

// PicData is data related to picture
type PicData struct {
	Name   string
	Bounds mat.AABB
	Img    draw.Image
}

// NPicData creates Slice of PicData from texture paths, it flips the textures
// so they are not upside down
func NPicData(paths ...string) (Vec, error) {
	data := make(Vec, len(paths))
	for i, v := range paths {
		d := &data[i]
		d.Name = v

		img, err := ggl.LoadImage(v)
		if err != nil {
			return nil, err
		}
		ggl.FlipNRGBA(img.(*image.NRGBA))

		d.Img = img.(draw.Image)
		d.Bounds = mat.FromRect(img.Bounds())
	}
	return data, nil
}

// Sheet contains sprite sheet and Regions
type Sheet struct {
	Pic     *image.NRGBA
	Regions map[string]mat.AABB
}

// NSheet calls NPicData on paths and then NSheetFromData with data and root
func NSheet(root string, paths ...string) (sh *Sheet, err error) {
	data, err := NPicData(paths...)
	if err != nil {
		return nil, err
	}
	return NSheetFromData(data, root), nil
}

// NSheetFromData creates sprite shit from given pic data, names of regions
// will be truncated by root so "hello/root/something.png" will be put
// under "something" in regions
func NSheetFromData(data Vec, root string) (sh *Sheet) {
	sh = &Sheet{Regions: map[string]mat.AABB{}}

	if len(data) == 0 {
		return
	}

	w, h := Pack(data)

	bounds := image.Rect(0, 0, w, h)

	sh.Pic = image.NewNRGBA(bounds)

	for _, d := range data {
		r := d.Bounds.ToImage()
		draw.Draw(sh.Pic, r, d.Img, d.Img.Bounds().Min, draw.Over)
		name, w, h, ok := DetectSpritesheet(d.Name, root)

		if ok {
			w := float64(w)
			h := float64(h)
			for y, n := d.Bounds.Max.Y, 0; y > d.Bounds.Min.Y; y -= h {
				for x := d.Bounds.Min.X; x < d.Bounds.Max.X; x += w {
					n++
					sh.Regions[name+strconv.Itoa(n)] = mat.A(x, y-h, x+w, y)
				}
			}
		} else {
			sh.Regions[name] = d.Bounds
		}
	}

	sh.Regions["All"] = mat.A(0, 0, float64(w), float64(h))

	return
}

// Pack packs rectangles in reasonable way, it tries to achieve
// size efficiency not speed
func Pack(data Vec) (width, height int) {
	// to guarantee that rect that is first in the row is highest
	data.Sort(func(a, b PicData) bool {
		return a.Bounds.H() > b.Bounds.H()
	})

	count := len(data)
	if count == 1 { // useless, bail
		return int(data[0].Bounds.W()), int(data[0].Bounds.H())
	}

	var (
		point       = calcOptimalSide(data)
		lowestRatio = math.MaxFloat64
		best        []int
	)
o:
	for point < count {
		var length float64
		for i := 0; i < point; i++ {
			length += data[i].Bounds.W()
		}

		// it would result to infinit loop if there is rectangle that is wider then length
		// si increase length and try again
		for i := point; i < count; i++ {
			if data[i].Bounds.W() > length {
				point++
				continue o
			}
		}

		var (
			current     int
			breakpoints []int
		)

		// finding breakpoints
		for current < count {
			var total float64
			breakpoints = append(breakpoints, current)
			for current < count {
				total += data[current].Bounds.W()
				if total > length {
					break
				}
				current++
			}
		}

		// calculating height of final cube for ratio
		var tollness float64
		for _, v := range breakpoints {
			tollness += data[v].Bounds.H()
		}

		// idk, it just kinda works like this
		ratio := length*tollness/300 + math.Abs(length-tollness)

		// deciding if we should stop, as long as ratio is decreasing continue, when it goes up, stop
		if ratio < lowestRatio {
			lowestRatio = ratio
			best = breakpoints
			width = int(length)
			height = int(tollness)
		} else if ratio > lowestRatio {
			break
		}

		point++
	}

	// modifying pic data according to best breakpoints
	best = append(best, count)
	offset := mat.ZV
	for i := 0; i < len(best)-1; i++ {
		// best can look like [0, 3, 6, 9] if there are 2 breakpoints on 3 and 6
		for j := best[i]; j < best[i+1]; j++ {
			data[j].Bounds = data[j].Bounds.Moved(offset)
			offset.X += data[j].Bounds.W()
		}
		offset.Y += data[best[i]].Bounds.H()
		offset.X = 0.0
	}

	return width, height
}

// calcOptimalArea calculates the most optimal rectangle
// that could theoretically be build from given data, side
// will be exactly correct only if all inputted rectangles
// are of same size
func calcOptimalSide(data Vec) int {
	// calculate total area
	var area float64
	for _, p := range data {
		area += p.Bounds.Area()
	}

	// side of square
	side := math.Sqrt(area)

	// find the break point
	for i, p := range data {
		side -= p.Bounds.W()
		if side <= 0 {
			return i
		}
	}

	return 0
}

// DetectSpritesheet parses name of texture into spritesheet if name is in format:
//
// 	name_width_height.ext
//
// where width hight are parameters of one sheet cell
func DetectSpritesheet(path, root string) (name string, w, h int, ok bool) {
	path, _ = str.SplitToTwo(path, '.') // remove extencion

	path = dirs.NormPath(path) // fix slashes

	// remove unimportant part of path
	if root != "" {
		parts := str.RevSplit(path, "/"+root+"/", 2)
		path = parts[len(parts)-1]
	}

	// split to name, width, height
	parts := str.RevSplit(path, "_", 3)

	name = path
	// parse numbers
	if len(parts) == 3 {
		var err error
		w, err = strconv.Atoi(parts[1])
		if err != nil {
			return
		}
		h, err = strconv.Atoi(parts[2])
		if err != nil {
			return
		}

		name = parts[0]
		ok = true
		return
	}

	return
}
