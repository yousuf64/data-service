package flake

import (
	"fmt"
	"github.com/sony/sonyflake"
	"time"
)

var epoch int64 = 1672531200000
var epochTime = time.UnixMilli(epoch)
var bucketSize int64 = 1000 * 60 * 60 * 24 * 10
var sf = sonyflake.NewSonyflake(sonyflake.Settings{
	StartTime: epochTime,
	MachineID: func() (uint16, error) {
		return 1, nil
	},
	CheckMachineID: nil,
})

func NextId() (uint64, error) {
	return sf.NextID()
}

func MakeBucket() uint64 {
	//t, _ := time.Parse(time.DateOnly, "2023-01-11")
	//ts := t.UnixMilli() - epoch
	fmt.Println(NextId())
	ts := time.Now().UnixMilli() - epoch
	return uint64(ts / bucketSize)
	//id, err := sf.NextID()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(sonyflake.MachineID(id))
	//fmt.Println(sonyflake.ElapsedTime(id))
	//fmt.Println(sonyflake.SequenceNumber(id))
}

func BucketFromSonyflake(i uint64) uint64 {
	d := sonyflake.ElapsedTime(i)
	return uint64(d.Milliseconds() / bucketSize)
}

func BucketsFromSonyflakeToNow(i uint64) (buckets []uint64) {
	start := BucketFromSonyflake(i)
	now := MakeBucket()

	for x := start; x <= now; x++ {
		buckets = append(buckets, x)
	}
	return
}

func MakeBucketsFrom(startId uint64) {
	sonyflake.Decompose(startId)
}
