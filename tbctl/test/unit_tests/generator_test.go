package unit_test

// import (
// 	"testing"
//
// 	"github.com/hanapedia/the-bench/tbctl/pkg/kubernetes/generate"
// 	"github.com/hanapedia/the-bench/tbctl/test/support"
// )

// func TestStatelessManifestGenerator(t *testing.T) {
// 	suc := support.GetServiceUnitConfig("./testdata/stateless.yaml")
// 	errs := generate.GenerateStatelessManifests("./testdata/stateless.yaml", "./testdata/output/stateless", suc)
// 	if errs.Exist() {
// 		errs.Print()
// 		t.Fatal("Failed to generate manifests.")
// 	}
// }
//
// func TestKafkaTopicManifestGenerator(t *testing.T) {
// 	suc := support.GetServiceUnitConfig("./testdata/kafka.yaml")
// 	errs := generate.GenerateBrokerManifests("./testdata/output/kafka", suc)
// 	if errs.Exist() {
// 		errs.Print()
// 		t.Fatal("Failed to generate manifests.")
// 	}
// }
//
// func TestMongoManifestGenerator(t *testing.T) {
// 	suc := support.GetServiceUnitConfig("./testdata/mongo.yaml")
// 	errs := generate.GenerateStatefulManifests("./testdata/output/mongo", suc)
// 	if errs.Exist() {
// 		errs.Print()
// 		t.Fatalf("Failed to generate manifests.")
// 	}
// }
