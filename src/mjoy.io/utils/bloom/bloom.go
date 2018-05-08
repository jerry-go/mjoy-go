package bloom

import(
	"math/big"
	"mjoy.io/common/types"
	"mjoy.io/utils/crypto"
)

type BloomByte interface {
	Bytes() []byte
}

func calculateBloom(b []byte) *big.Int {
	b = crypto.Keccak256(b[:])

	r := new(big.Int)

	for i := 0; i < 6; i += 2 {
		t := big.NewInt(1)
		b := (uint(b[i+1]) + (uint(b[i]) << 8)) & 2047
		r.Or(r, t.Lsh(t, b))
	}

	return r
}

func CreateBloom(topics []BloomByte) types.Bloom {
	bin := new(big.Int)
	for _, topic := range topics {
		bin.Or(bin, calculateBloom(topic.Bytes()[:]))
	}

	return types.BytesToBloom(bin.Bytes())
}

func BloomLookup(bin types.Bloom, topic BloomByte) bool {
	bloom := bin.Big()
	cmp := calculateBloom(topic.Bytes()[:])

	return bloom.And(bloom, cmp).Cmp(cmp) == 0
}
