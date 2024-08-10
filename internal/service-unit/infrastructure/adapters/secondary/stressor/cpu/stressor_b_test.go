package cpu

import (
	"context"
	"fmt"
	"testing"

	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/service-unit/utils"
)

/*
## Results
Numbers are just for a relative comparison between the stressors

BenchmarkStressCPU-8             1401730               827.5 ns/op
BenchmarkStressCPUB-8               6789            173886 ns/op
BenchmarkStressCPUB_double_p-8      3471            353621 ns/op
(tested locally on MacBook Air M1)

sha256 stressor takes about 210 times more time per call
sha256 stressor scales linearly by payload size
*/

var iter = 10

func BenchmarkStressCPU(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stressCPU(context.Background(), iter)
	}
}

var p = utils.GenerateRandomString(constants.PayloadSizeMap[constants.DefaultPayloadSize])

func BenchmarkStressCPUB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stressCPU_sha256(context.Background(), iter, []byte(p))
	}
}

var pp = fmt.Sprintf("%s%s", p, p)

func BenchmarkStressCPUB_double_p(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stressCPU_sha256(context.Background(), iter, []byte(pp))
	}
}
