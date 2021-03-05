package gate

import (
	"testing"
)

func TestCleanup(t *testing.T) {
	f := &FazerBase{
		Fazes: []FazeRunner{
			func(tIdx, count int) {},
		},
	}
	g := Gate{}
	g.Add(f)
	g.Add(&FazerBase{})
	g.Add(f)
	g.Add(&FazerBase{})

	g2 := Gate{}
	g2.Add(f)
	g2.Add(f)

	g.Run()
	g.CleanWait()

	for i := range g.threads {
		if g.threads[i].idx != g2.threads[i].idx {
			t.Errorf("\n%#v\n%#v", g.threads[i], g2.threads[i])
		}

	}

	g.threads = g.threads[len(g.threads):4]
	e := []*wrapper{{make(chan bool), make(chan bool), 3, nil}, {make(chan bool), make(chan bool), 1, nil}}
	for i := range g.threads {
		if g.threads[i].idx != e[i].idx || g.threads[i].Fazer != e[i].Fazer {
			t.Errorf("\n%#v\n%#v", g.threads[i], e[i])
		}

	}
}

func TestRun(t *testing.T) {
	f1 := FazerBase{
		Fazes: []FazeRunner{
			func(tIdx, count int) {},
			func(tIdx, count int) {},
			func(tIdx, count int) {},
		},
	}
	f2 := f1

	g := Gate{}

	g.Add(&f1)
	g.Add(&f2)

	for i := 0; i < 9; i++ {
		g.Run()
		g.Wait()
	}

	if f1.Cursor != 0 || f2.Cursor != 0 {
		t.Error(f1, f2)
	}

	f3 := f2
	g.Add(&f3)

	for i := 0; i < 9; i++ {
		g.Run()
		g.Wait()
	}

	if f3.Cursor != 0 {
		t.Error(f1, f2, f3)
	}

	g.Clear()

	g.Add(&f1)
	g.Run()
	g.Wait()
	g.Add(&f2)
	g.Run()
	g.Wait()
	g.Add(&f3)

	for i := 0; i < 9; i++ {
		g.Run()
		g.Wait()
	}

	if f1.Cursor != 2 || f2.Cursor != 1 || f3.Cursor != 0 {
		t.Error(f1, f2, f3)
	}
}

func BenchmarkGate(b *testing.B) {
	g := Gate{}
	for i := 0; i < 4; i++ {
		g.Add(&FazerBase{
			Fazes: []FazeRunner{
				func(tIdx, count int) {},
			},
		})
	}

	for i := 0; i < b.N; i++ {
		g.Run()
		g.Wait()
	}
}
