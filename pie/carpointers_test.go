package pie

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var carPointerA = &car{"a", "green"}
var carPointerB = &car{"b", "blue"}
var carPointerC = &car{"c", "gray"}
var carPointerEmpty = &car{}

var carPointersContainsTests = []struct {
	ss       carPointers
	contains *car
	expected bool
}{
	{nil, carPointerA, false},
	{nil, carPointerEmpty, false},
	{nil, nil, false},
	{carPointers{carPointerA, carPointerB, carPointerC}, carPointerA, true},
	{carPointers{carPointerA, carPointerB, carPointerC}, carPointerB, true},
	{carPointers{carPointerA, carPointerB, carPointerC}, carPointerC, true},
	{carPointers{carPointerA, carPointerB, carPointerC}, &car{"a", "green"}, false},
	{carPointers{carPointerA, carPointerB, carPointerC}, &car{"A", ""}, false},
	{carPointers{carPointerA, carPointerB, carPointerC}, &car{}, false},
	{carPointers{carPointerA, carPointerB, carPointerC}, &car{"d", ""}, false},
	{carPointers{carPointerA, carPointerEmpty, carPointerC}, carPointerEmpty, true},
	{carPointers{carPointerA, nil, carPointerC}, nil, true},
	{carPointers{carPointerA, carPointerEmpty, carPointerC}, nil, false},
}

func TestCarPointers_Contains(t *testing.T) {
	for _, test := range carPointersContainsTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCarPointers(t, &test.ss)()
			assert.Equal(t, test.expected, test.ss.Contains(test.contains))
		})
	}
}

var carPointersSelectTests = []struct {
	ss                carPointers
	condition         func(*car) bool
	expectedSelect    carPointers
	expectedUnselect  carPointers
	expectedTransform carPointers
}{
	{
		nil,
		func(s *car) bool {
			return s.Name == ""
		},
		nil,
		nil,
		nil,
	},
	{
		carPointers{carPointerA, carPointerB, carPointerC},
		func(s *car) bool {
			return s.Name != "b"
		},
		carPointers{carPointerA, carPointerC},
		carPointers{carPointerB},
		carPointers{&car{"A", "green"}, &car{"B", "blue"}, &car{"C", "gray"}},
	},
}

func TestCarPointers_Select(t *testing.T) {
	for _, test := range carPointersSelectTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCarPointers(t, &test.ss)()
			assert.Equal(t, test.expectedSelect, test.ss.Select(test.condition))
		})
	}
}

func TestCarPointers_Unselect(t *testing.T) {
	for _, test := range carPointersSelectTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCarPointers(t, &test.ss)()
			assert.Equal(t, test.expectedUnselect, test.ss.Unselect(test.condition))
		})
	}
}

func TestCarPointers_Transform(t *testing.T) {
	for _, test := range carPointersSelectTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCarPointers(t, &test.ss)()
			assert.Equal(t, test.expectedTransform, test.ss.Transform(func(c *car) *car {
				return &car{
					Name:  strings.ToUpper(c.Name),
					Color: c.Color,
				}
			}))
		})
	}
}

var carPointersFirstAndLastTests = []struct {
	ss             carPointers
	first, firstOr *car
	last, lastOr   *car
}{
	{
		nil,
		&car{},
		&car{"default1", "unknown"},
		&car{},
		&car{"default2", "unknown"},
	},
	{
		carPointers{&car{"foo", "red"}},
		&car{"foo", "red"},
		&car{"foo", "red"},
		&car{"foo", "red"},
		&car{"foo", "red"},
	},
	{
		carPointers{carPointerA, carPointerB},
		carPointerA,
		carPointerA,
		carPointerB,
		carPointerB,
	},
	{
		carPointers{carPointerA, carPointerB, carPointerC},
		carPointerA,
		carPointerA,
		carPointerC,
		carPointerC,
	},
}

func TestCarPointers_FirstOr(t *testing.T) {
	for _, test := range carPointersFirstAndLastTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCarPointers(t, &test.ss)()
			assert.Equal(t, test.firstOr, test.ss.FirstOr(&car{"default1", "unknown"}))
		})
	}
}

func TestCarPointers_LastOr(t *testing.T) {
	for _, test := range carPointersFirstAndLastTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCarPointers(t, &test.ss)()
			assert.Equal(t, test.lastOr, test.ss.LastOr(&car{"default2", "unknown"}))
		})
	}
}

func TestCarPointers_First(t *testing.T) {
	for _, test := range carPointersFirstAndLastTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCarPointers(t, &test.ss)()
			assert.Equal(t, test.first, test.ss.First())
		})
	}
}

func TestCarPointers_Last(t *testing.T) {
	for _, test := range carPointersFirstAndLastTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCarPointers(t, &test.ss)()
			assert.Equal(t, test.last, test.ss.Last())
		})
	}
}

var carPointersStatsTests = []struct {
	ss       carPointers
	min, max *car
	len      int
}{
	{
		nil,
		&car{},
		&car{},
		0,
	},
	{
		carPointers{},
		&car{},
		&car{},
		0,
	},
	{
		carPointers{&car{"foo", "red"}},
		&car{"foo", "red"},
		&car{"foo", "red"},
		1,
	},
	{
		carPointers{&car{"bar", "yellow"}, &car{"Baz", "black"}, &car{"qux", "cyan"}, &car{"foo", "red"}},
		&car{"Baz", "black"},
		&car{"qux", "cyan"},
		4,
	},
}

func TestCarPointers_Len(t *testing.T) {
	for _, test := range carPointersStatsTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCarPointers(t, &test.ss)()
			assert.Equal(t, test.len, test.ss.Len())
		})
	}
}

var carPointersJSONTests = []struct {
	ss         carPointers
	jsonString string
}{
	{
		nil,
		`[]`, // Instead of null.
	},
	{
		carPointers{},
		`[]`,
	},
	{
		carPointers{&car{"foo", "red"}},
		`[{"Name":"foo","Color":"red"}]`,
	},
	{
		carPointers{&car{"bar", "yellow"}, &car{"Baz", "black"}, &car{"qux", "cyan"}, &car{"foo", "red"}},
		`[{"Name":"bar","Color":"yellow"},{"Name":"Baz","Color":"black"},{"Name":"qux","Color":"cyan"},{"Name":"foo","Color":"red"}]`,
	},
}

func TestCarPointers_JSONString(t *testing.T) {
	for _, test := range carPointersJSONTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCarPointers(t, &test.ss)()
			assert.Equal(t, test.jsonString, test.ss.JSONString())
		})
	}
}

var carPointersSortTests = []struct {
	ss        carPointers
	sorted    carPointers
	reversed  carPointers
	areSorted bool
}{
	{
		nil,
		nil,
		nil,
		true,
	},
	{
		carPointers{},
		carPointers{},
		carPointers{},
		true,
	},
	{
		carPointers{&car{"foo", "red"}},
		carPointers{&car{"foo", "red"}},
		carPointers{&car{"foo", "red"}},
		true,
	},
	{
		carPointers{&car{"bar", "yellow"}, &car{"Baz", "black"}, &car{"foo", "red"}},
		carPointers{&car{"Baz", "black"}, &car{"bar", "yellow"}, &car{"foo", "red"}},
		carPointers{&car{"foo", "red"}, &car{"Baz", "black"}, &car{"bar", "yellow"}},
		false,
	},
	{
		carPointers{&car{"bar", "yellow"}, &car{"Baz", "black"}, &car{"qux", "cyan"}, &car{"foo", "red"}},
		carPointers{&car{"Baz", "black"}, &car{"bar", "yellow"}, &car{"foo", "red"}, &car{"qux", "cyan"}},
		carPointers{&car{"foo", "red"}, &car{"qux", "cyan"}, &car{"Baz", "black"}, &car{"bar", "yellow"}},
		false,
	},
	{
		carPointers{&car{"Baz", "black"}, &car{"bar", "yellow"}},
		carPointers{&car{"Baz", "black"}, &car{"bar", "yellow"}},
		carPointers{&car{"bar", "yellow"}, &car{"Baz", "black"}},
		true,
	},
}

func TestCarPointers_Reverse(t *testing.T) {
	for _, test := range carPointersSortTests {
		t.Run("", func(t *testing.T) {
			defer assertImmutableCarPointers(t, &test.ss)()
			assert.Equal(t, test.reversed, test.ss.Reverse())
		})
	}
}