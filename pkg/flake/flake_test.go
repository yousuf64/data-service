package flake

import (
	"fmt"
	"testing"
)

func TestMakeBucket(t *testing.T) {
	fmt.Println(MakeBucket())
	fmt.Println(BucketFromSonyflake(20613037266305025))
	fmt.Println(BucketsFromSonyflakeToNow(20613037266305025))
}
