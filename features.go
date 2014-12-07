package tiles

type Feature struct {
	MinZoom int
	Tags    []Tag
}

type Tag struct {
	Key, Val string
}

var mapFeatures = map[string]Feature{
	"coastline": {
		MinZoom: 0,
		Tags: []Tag{
			{"natural", "coastline"},
		},
	},
	"all-other-roads": {
		MinZoom: 14,
		Tags: []Tag{
			{"highway", "unclassified"},
			{"highway", "residential"},
			{"highway", "service"},
			{"highway", "motorway_link"},
			{"highway", "trunk_link"},
			{"highway", "primary_link"},
			{"highway", "secondary_link"},
			{"highway", "living_street"},
			{"highway", "pedestrian"},
			{"highway", "track"},
			{"highway", "bus_guideway"},
			{"highway", "raceway"},
			{"highway", "road"},
		},
	},
	"major-ish-roads": {
		MinZoom: 12,
		Tags: []Tag{
			{"highway", "primary"},
			{"highway", "secondary"},
			{"highway", "tertiary"},
		},
	},
	"major-major-roads": {
		MinZoom: 10,
		Tags: []Tag{
			{"highway", "motorway"},
			{"highway", "trunk"},
		},
	},
	"buildings": {
		MinZoom: 14,
		Tags: []Tag{
			{"building", "*"},
		},
	},
}
