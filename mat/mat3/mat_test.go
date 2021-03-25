package mat3

import "testing"

func TestMatrixMul(t *testing.T) {
	testCases := []struct {
		desc      string
		a, b, out Mat
	}{
		{
			desc: "identity",
			a:    IM,
			b:    IM,
			out:  IM,
		},
		{
			desc: "no change",
			a:    Mat{V(0, 1, 2), V(4, 0, 3), V(6, 3, 0), V(0, 7, 8)},
			b:    IM,
			out:  Mat{V(0, 1, 2), V(4, 0, 3), V(6, 3, 0), V(0, 7, 8)},
		},
		{
			desc: "move",
			a:    IM.Moved(V(10, 10, 10)),
			b:    IM.Moved(V(5, 10, 15)),
			out:  IM.Moved(V(15, 20, 25)),
		},
		{
			desc: "scaled",
			a:    IM.ScaleComp(V(10, 10, 10)),
			b:    IM.ScaleComp(V(5, 10, 15)),
			out:  IM.ScaleComp(V(50, 100, 150)),
		},
		{
			desc: "scale move",
			a:    IM.ScaleComp(V(10, 10, 10)).Moved(V(1, 1, 1)),
			b:    IM.ScaleComp(V(5, 10, 15)).Moved(V(15, 10, 5)),
			out:  IM.ScaleComp(V(50, 100, 150)).Moved(V(20, 20, 20)),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			res := tC.a.Mul(tC.b)
			if !res.Approx(tC.out, 8) {
				t.Error(res)
			}
		})
	}
}

func TestProjectUnproject(t *testing.T) {
	testCases := []struct {
		desc     string
		mat      Mat
		inp, out Vec
	}{
		{
			desc: "move",
			mat:  IM.Moved(V(10, 4, 2)),
			inp:  V(2, -2, 3),
			out:  V(12, 2, 5),
		},
		{
			desc: "scale",
			mat:  IM.ScaleComp(V(2, 1, 3)).Moved(V(10, 4, 2)),
			inp:  V(2, -2, 3),
			out:  V(14, 2, 11),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			res := tC.mat.Project(tC.inp)
			if res != tC.out {
				t.Error(res, "prj")
			}
			res = tC.mat.Unproject(res)
			if res != tC.inp {
				t.Error(res, "un")
			}
		})
	}
}
