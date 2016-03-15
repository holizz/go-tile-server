package tiles

import (
	//"fmt"
	"github.com/golang/geo/s2"
)

type S2Index map[s2.CellID]([]FeatureRef)

func (si S2Index) AddWay(w Way, fname string, data *OsmData) {
	c := s2.EmptyCap()

	for _, node := range w.GetNodes(data.Nodes) {
		c = c.AddPoint(s2.PointFromLatLng(s2.LatLngFromDegrees(node.Lat(), node.Lon())))
	}

	if c.IsEmpty() {
		return
	}

	rc := &s2.RegionCoverer{MaxLevel: 30, MaxCells: 10}
	cu := rc.FastCovering(c)
	for _, cid := range cu {
		si[cid] = append(si[cid], FeatureRef{w.Id, ItemTypeWay, fname})
		for l := cid.Level(); l > 0; l-- {
			cid = cid.Parent(l - 1)
			if _, ok := si[cid]; !ok {
				si[cid] = make([]FeatureRef, 0)
			}
		}
	}
}

func (si S2Index) GetFeatures(nwPt, sePt Pointer, zoom int64, data *OsmData) []FeatureRef {

	r := s2.RectFromLatLng(s2.LatLngFromDegrees(nwPt.Lat(), nwPt.Lon()))
	r = r.AddPoint(s2.LatLngFromDegrees(sePt.Lat(), sePt.Lon()))

	rc := &s2.RegionCoverer{MaxLevel: 30, MaxCells: 10}

	cu := rc.Covering(r)

	visitCid := make(map[s2.CellID]bool)
	visitF := make(map[int64]bool)
	ret := make([]FeatureRef, 0)

	for _, cid := range cu {
		if v, ok := si[cid]; ok {
			ret = si.VisitDown(cid, v, visitF, ret)

			for l := cid.Level(); l > 0; l-- {
				cid = cid.Parent(l - 1)
				ret = si.VisitUp(cid, visitCid, visitF, ret)
			}
		}
	}
	//fmt.Println( len(ret))
	return ret
}

func (si S2Index) VisitUp(cid s2.CellID, visitCid map[s2.CellID]bool, visitF map[int64]bool, ret []FeatureRef) []FeatureRef {
	if visitCid[cid] {
		return ret
	}
	visitCid[cid] = true

	for _, f := range si[cid] {
		if !visitF[f.Id] {
			ret = append(ret, f)
			visitF[f.Id] = true
		}
	}
	return ret
}

func (si S2Index) VisitDown(cid s2.CellID, fr []FeatureRef, visitF map[int64]bool, ret []FeatureRef) []FeatureRef {

	for _, f := range fr {
		if !visitF[f.Id] {
			ret = append(ret, f)
			visitF[f.Id] = true
		}
	}

	if !cid.IsLeaf() {
		chs := cid.Children()
		for i := 0; i < 4; i++ {
			if v, ok := si[chs[i]]; ok {
				ret = si.VisitDown(chs[i], v, visitF, ret)

			}
		}

	}

	return ret

}
