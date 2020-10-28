package main

import (
	"github.com/ganeshrvel/go-mtpfs/mtp"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMtpInitialize(t *testing.T) {
	var dev *mtp.Device
	var sid uint32

	Convey("Testing Initialize", t, func() {
		d, err := Initialize(Init{})
		dev = d

		So(err, ShouldBeNil)
		So(d.Timeout, ShouldBeGreaterThan, 1)
	})

	Convey("Testing FetchDeviceInfo", t, func() {
		info, err := FetchDeviceInfo(dev)

		So(err, ShouldBeNil)
		So(info, ShouldNotBeNil)
	})

	Convey("Testing FetchStorages", t, func() {
		storages, err := FetchStorages(dev)

		sid = storages[0].sid

		So(err, ShouldBeNil)
		So(sid, ShouldEqual, 0x10001)
	})

	Dispose(dev)
}
