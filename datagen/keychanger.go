package datagen

import (
	"flag"
	"fmt"

	"gopkg.in/raintank/schema.v1"

	gc "github.com/rakyll/globalconf"
)

type Keychanger struct {
	id        int
	keyPoints []int
	currKey   []int
}

var (
	pointsPerKey int
	syncSwitch   bool
	kcKeyCount   int
	kcKeyPrefix  string
)

func init() {
	modules["keychanger"] = kcNew
	regFlags = append(regFlags, kcRegFlags)

}

func kcNew(id int) Datagen {
	initValue := 0
	keyPoints := make([]int, kcKeyCount)
	currKey := make([]int, kcKeyCount)

	for i := 0; i < kcKeyCount; i++ {
		currKey[i] = 0
		keyPoints[i] = initValue
		if !syncSwitch {
			initValue++
		}
	}

	return &Keychanger{id, keyPoints, currKey}
}

func kcRegFlags() {
	flags := flag.NewFlagSet("key-changer", flag.ExitOnError)
	flags.IntVar(&pointsPerKey, "points-per-key", 10, "number of points per key")
	flags.IntVar(&kcKeyCount, "key-count", 100, "number of keys to generate")
	flags.StringVar(&kcKeyPrefix, "key-prefix", "some.key", "prefix for keys")
	flags.BoolVar(&syncSwitch, "sync-switch", true, "change all keys at once")
	gc.Register("key-changer", flags)
}

func (kc *Keychanger) GetData(ts int64) []*schema.MetricData {
	metrics := make([]*schema.MetricData, kcKeyCount)

	for i := 0; i < kcKeyCount; i++ {
		name := fmt.Sprintf(kcKeyPrefix+"%d.%d.%d", kc.id, i, kc.currKey[i])
		metrics[i] = &schema.MetricData{
			Name:   name,
			Metric: name,
			OrgId:  i,
			Value:  0,
			Unit:   "ms",
			Mtype:  "gauge",
			Tags:   []string{"some_tag", "ok", "k:2"},
			Time:   ts,
		}

		kc.keyPoints[i]++

		if kc.keyPoints[i]%pointsPerKey == 0 {
			kc.keyPoints[i] = 0
			kc.currKey[i]++
		}
	}

	return metrics
}
