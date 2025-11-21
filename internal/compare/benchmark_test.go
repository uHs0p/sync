package compare

import (
    "fmt"
    "testing"
    "terraformsync/internal/collect"
)

func makeSet(n int, prefix string) collect.ResourceSet {
    rs := collect.ResourceSet{}
    for i := 0; i < n; i++ {
        k := fmt.Sprintf("aws_s3_bucket:%s-%d", prefix, i)
        rs[k] = collect.Resource{Type: "aws_s3_bucket", Name: fmt.Sprintf("%s-%d", prefix, i), Attrs: map[string]string{"bucket": fmt.Sprintf("%s-%d", prefix, i)}}
    }
    return rs
}

func BenchmarkThreeWay100(b *testing.B) {
    code := makeSet(100, "c")
    state := makeSet(100, "c")
    deployed := makeSet(150, "c")
    for i := 0; i < b.N; i++ {
        _ = ThreeWay(code, state, deployed)
    }
}

func BenchmarkThreeWay1000(b *testing.B) {
    code := makeSet(1000, "c")
    state := makeSet(1000, "c")
    deployed := makeSet(1200, "c")
    for i := 0; i < b.N; i++ {
        _ = ThreeWay(code, state, deployed)
    }
}