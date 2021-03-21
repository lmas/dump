package engine

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVectors(t *testing.T) {
	Convey("Given a newly created vector", t, func() {
		v := NewVector(1, 2)

		Convey("When calling NewVector(1, 2)", func() {
			Convey("The output should be float64", func() {
				So(v.X, ShouldHaveSameTypeAs, 1.0)
				So(v.Y, ShouldHaveSameTypeAs, 2.0)
			})
			Convey("The output should match 1.0, 2.0", func() {
				So(v.X, ShouldEqual, 1.0)
				So(v.Y, ShouldEqual, 2.0)
			})
		})

		Convey("When calling Vector.Str()", func() {
			Convey("The output should be a string", func() {
				So(v.Str(), ShouldHaveSameTypeAs, "<1.0, 2.0>")
			})
			Convey("The output should match <1.0, 2.0>", func() {
				So(v.Str(), ShouldEqual, "<1.0, 2.0>")
			})
		})

		Convey("When calling Vector.Clone()", func() {
			clone := v.Clone()
			Convey("The new vector should have same values as the cloned", func() {
				So(clone.X, ShouldEqual, v.X)
				So(clone.Y, ShouldEqual, v.Y)
			})
			Convey("The new vector should be a new instance", func() {
				So(clone, ShouldNotEqual, v)
			})
		})
	})
}
