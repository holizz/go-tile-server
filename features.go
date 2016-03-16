package tiles

type Feature struct {
	MinZoom int64
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
	"borders2": {
		MinZoom: 0,
		Tags: []Tag{
			{"admin_level", "2"},
		},
	},
	"borders3": {
		MinZoom: 6,
		Tags: []Tag{
			{"admin_level", "3"},
		},
	},
	"borders4": {
		MinZoom: 8,
		Tags: []Tag{
			{"admin_level", "4"},
		},
	},
	"borders5": {
		MinZoom: 10,
		Tags: []Tag{
			{"admin_level", "5"},
		},
	},
	"borders6": {
		MinZoom: 12,
		Tags: []Tag{
			{"admin_level", "6"},
		},
	},
	"borders7": {
		MinZoom: 12,
		Tags: []Tag{
			{"admin_level", "7"},
		},
	},
	"borders8": {
		MinZoom: 13,
		Tags: []Tag{
			{"admin_level", "8"},
		},
	},
	"borders9": {
		MinZoom: 14,
		Tags: []Tag{
			{"admin_level", "9"},
		},
	},
	"all-other-roads": {
		MinZoom: 12,
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
		MinZoom: 9,
		Tags: []Tag{
			{"highway", "primary"},
			{"highway", "secondary"},
			{"highway", "tertiary"},
		},
	},
	"major-major-roads": {
		MinZoom: 6,
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
