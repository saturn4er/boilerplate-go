package segmentationstorage

import (
	fmt "fmt"

	segmentationservice "github.com/saturn4er/boilerplate-go/example/segmentation/segmentationservice"
	// user code 'imports'
	// end user code 'imports'
)

const (
	someEnumValue1 = "value1"
	someEnumValue2 = "value2"
)

func convertSomeEnumToDB(someEnumValue segmentationservice.SomeEnum) (string, error) {
	result, ok := map[segmentationservice.SomeEnum]string{
		segmentationservice.SomeEnumValue1: someEnumValue1,
		segmentationservice.SomeEnumValue2: someEnumValue2,
	}[someEnumValue]
	if !ok {
		return "", fmt.Errorf("unknown SomeEnum value: %d", someEnumValue)
	}
	return result, nil
}

func convertSomeEnumFromDB(someEnumValue string) (segmentationservice.SomeEnum, error) {
	result, ok := map[string]segmentationservice.SomeEnum{
		someEnumValue1: segmentationservice.SomeEnumValue1,
		someEnumValue2: segmentationservice.SomeEnumValue2,
	}[someEnumValue]
	if !ok {
		return 0, fmt.Errorf("unknown SomeEnum db value: %s", someEnumValue)
	}
	return result, nil
}
