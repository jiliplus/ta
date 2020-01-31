package ta

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_NewEWMA(t *testing.T) {
	Convey("新建 EWMA", t, func() {
		Convey("参数 N 没有大于 1 的话，会 panic", func() {
			So(func() {
				NewEWMA(1)
			}, ShouldPanicWith, "N in EWMA should bigger than 1")
		})
		Convey("参数 N 为 3 的话", func() {
			e := NewEWMA(3)
			Convey("alpha 应该是 0.5", func() {
				So(e.alpha, ShouldEqual, 0.5)
			})
		})
	})
}

func Test_EWMA_IsInited(t *testing.T) {
	Convey("新建 EWMA", t, func() {
		e := NewEWMA(10)
		Convey("刚刚生成好的 *EWMA 没有完成初始化", func() {
			So(e.IsInited(), ShouldBeFalse)
		})
		e.count = ewmaWarmUpSamples - 1
		Convey("没有达到次数的 *EWMA 也没有完成初始化", func() {
			So(e.IsInited(), ShouldBeFalse)
		})
		e.count = ewmaWarmUpSamples
		Convey("达到次数的 *EWMA 才完成初始化", func() {
			So(e.IsInited(), ShouldBeTrue)
		})
	})
}

func Test_EWMA_move(t *testing.T) {
	Convey("已经预热好的 EWMA", t, func() {
		e := &EWMA{
			count: ewmaWarmUpSamples,
			alpha: 0.5,
			value: 1024,
		}
		Convey("e.move(0) = 512", func() {
			So(e.move(0), ShouldEqual, 512)
		})
		Convey("e.move(1024) = 1024", func() {
			So(e.move(1024), ShouldEqual, 1024)
		})
		Convey("e.move(2048) = 1536", func() {
			So(e.move(2048), ShouldEqual, 1536)
		})
		Convey("e.move(-48) = 488", func() {
			So(e.move(-48), ShouldEqual, 488)
		})
	})
}

func Test_EWMA_Update(t *testing.T) {
	Convey("还差一次就预热好的 EWMA", t, func() {
		e := &EWMA{
			count: ewmaWarmUpSamples - 1,
			alpha: 0.5,
			value: 10240,
		}
		Convey("还没有预热好", func() {
			Convey("count < ewmaWarmUpSamples", func() {
				So(e.count, ShouldBeLessThan, ewmaWarmUpSamples)
			})
			So(e.IsInited(), ShouldBeFalse)
		})
		e.Update(10240)
		Convey("update 10240 后，预热好了，并且数值是和", func() {
			Convey("count 等于了 ewmaWarmUpSamples", func() {
				So(e.count, ShouldEqual, ewmaWarmUpSamples)
			})
			So(e.IsInited(), ShouldBeTrue)
			So(e.Value(), ShouldEqual, 20480)
		})
		e.Update(0)
		Convey("update 0 后，预热好了，并且数值是平均数", func() {
			Convey("count 等于了 1+ewmaWarmUpSamples", func() {
				So(e.count, ShouldEqual, 1+ewmaWarmUpSamples)
			})
			So(e.IsInited(), ShouldBeTrue)
			So(e.Value(), ShouldEqual, 1024)
		})
		e.Update(0)
		Convey("再次 update 0 后，预热好了，并且数值是平均数", func() {
			Convey("count 没有变化", func() {
				So(e.count, ShouldEqual, 1+ewmaWarmUpSamples)
			})
			So(e.IsInited(), ShouldBeTrue)
			So(e.Value(), ShouldEqual, 512)
		})
	})
}
