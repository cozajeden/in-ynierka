package monitor

import (
	"fmt"
	"time"
)

var storekeeper = make(chan *payload, 1_000_000)

func keep(toKeep *payload) error {
	switch v := toKeep.value.(type) {
	default:
		return fmt.Errorf("unknown type %v found in keep", v)
	case *dataStructurePayload:
		storageDetails := toKeep.value.(*dataStructurePayload)
		dataStructureQuantities[toKeep.target][storageDetails.kind] += storageDetails.quantity
	case int:
		rawQuantities[toKeep.target] += toKeep.value.(int)
	case replacement:
		rawQuantities[toKeep.target] = int(toKeep.value.(replacement))
	case time.Duration:
		times[toKeep.target] = append(times[toKeep.target], toKeep.value.(time.Duration))
	}
	return nil
}
