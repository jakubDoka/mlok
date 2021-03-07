package pck

import (
	"gobatch/ggl"
	"gobatch/ggl/txt"
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
	templates.Data<PicData, Data>
)*/

// PicData is data related to picture
type PicData struct {
	Name   string
	bounds mat.AABB
	Img    draw.Image
	drawer *txt.Drawer
}

// NPicData creates Slice of PicData from texture paths, it flips the textures
// so they are not upside down
func (d *Data) AddImages(paths ...string) error {
	data := *d
	for _, v := range paths {
		d := PicData{}
		d.Name = v

		img, err := ggl.LoadImage(v)
		if err != nil {
			return err
		}
		ggl.FlipNRGBA(img)

		d.Img = img

		data = append(data, d)
	}

	*d = data
	return nil
}

// AddMarkdown adds markdown font textures into Data
func (v *Data) AddMarkdown(m *txt.Markdown) {
	for _, font := range m.Fonts {
		*v = append(*v, PicData{
			Name:   font.Name,
			Img:    font.Pic,
			drawer: font,
		})
	}
}

// Sheet contains sprite sheet and Regions
type Sheet struct {
	Data

	Root    string
	Pic     *image.NRGBA
	Regions map[string]mat.AABB
}

// NSheet creates new sheet containing textures from given paths
func NSheet(root string, paths ...string) (sh *Sheet, err error) {
	sh = &Sheet{
		Root: root,
	}

	err = sh.AddImages(paths...)
	if err != nil {
		return
	}

	sh.Pack()
	return sh, nil
}

// Pack takes all Data in sheet and translates it into one packed image
// and regions, regarding the current Root the names will be modified, for
// example if root is "root" then "something/root/anything.png" will be saved under
// kay with value "anything" into s.Regions
func (s *Sheet) Pack() {
	if len(s.Data) == 0 {
		return
	}

	if s.Regions == nil {
		s.Regions = map[string]mat.AABB{}
	}
	// cleanup
	for k := range s.Regions {
		delete(s.Regions, k)
	}
	for i := range s.Data {
		d := &s.Data[i]
		d.bounds = mat.FromRect(d.Img.Bounds())
	}

	w, h := Pack(s.Data)

	bounds := image.Rect(0, 0, w, h)

	s.Pic = image.NewNRGBA(bounds)

	for _, d := range s.Data {
		r := d.bounds.ToImage()
		draw.Draw(s.Pic, r, d.Img, d.Img.Bounds().Min, draw.Over)
		name, w, h, ok := DetectSpritesheet(d.Name, s.Root)

		if ok {
			w := float64(w)
			h := float64(h)
			for y, n := d.bounds.Max.Y, 0; y > d.bounds.Min.Y; y -= h {
				for x := d.bounds.Min.X; x < d.bounds.Max.X; x += w {
					n++
					s.Regions[name+strconv.Itoa(n)] = mat.A(x, y-h, x+w, y)
				}
			}
		} else {
			s.Regions[name] = d.bounds
			if d.drawer != nil {
				d.drawer.Region = d.bounds.Min
			}
		}
	}

	s.Regions["All"] = mat.A(0, 0, float64(w), float64(h))

	return
}

// Pack packs rectangles in reasonable way, it tries to achieve
// size efficiency not speed
func Pack(data Data) (width, height int) {
	// to guarantee that rect that is first in the row is highest
	data.Sort(func(a, b PicData) bool {
		return a.bounds.H() > b.bounds.H()
	})

	count := len(data)
	if count == 1 { // useless, bail
		return int(data[0].bounds.W()), int(data[0].bounds.H())
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
			length += data[i].bounds.W()
		}

		// it would result to infinit loop if there is rectangle that is wider then length
		// si increase length and try again
		for i := point; i < count; i++ {
			if data[i].bounds.W() > length {
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
				total += data[current].bounds.W()
				if total > length {
					break
				}
				current++
			}
		}

		// calculating height of final cube for ratio
		var tollness float64
		for _, v := range breakpoints {
			tollness += data[v].bounds.H()
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
			data[j].bounds = data[j].bounds.Moved(offset)
			offset.X += data[j].bounds.W()
		}
		offset.Y += data[best[i]].bounds.H()
		offset.X = 0.0
	}

	return width, height
}

// calcOptimalArea calculates the most optimal rectangle
// that could theoretically be build from given data, side
// will be exactly correct only if all inputted rectangles
// are of same size
func calcOptimalSide(data Data) int {
	// calculate total area
	var area float64
	for _, p := range data {
		area += p.bounds.Area()
	}

	// side of square
	side := math.Sqrt(area)

	// find the break point
	for i, p := range data {
		side -= p.bounds.W()
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
