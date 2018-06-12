package algorand

import (
	"math/big"
	"mjoy.io/common"
	"fmt"
	"errors"
	"mjoy.io/utils/crypto"
	"mjoy.io/common/types"
	"crypto/ecdsa"
)

var (
	ErrInvalidSig = errors.New("invalid block v, r, s values")
	ErrInvalidChainId = errors.New("invalid chain id for block signer")
)


// SignHeader signs the header using the given signer and private key
func SignCredential(c *CredentialData, s Signer, prv *ecdsa.PrivateKey) (*CredentialSig, error) {
	hash := c.Hash()
	sig, err := crypto.Sign(hash[:], prv)
	if err != nil {
		return nil, err
	}
	credentialSig := &CredentialSig{}
	R, S, V, err := s.SignatureValues(sig)
	credentialSig.R.IntVal.Set(R)
	credentialSig.S.IntVal.Set(S)
	credentialSig.V.IntVal.Set(V)
	credentialSig.Round = c.Round
	credentialSig.Step = c.Step
	return credentialSig, nil
}

// Signer encapsulates transaction signature handling. Note that this interface is not a
// stable API and may change at any time to accommodate new protocol rules.
type Signer interface {
	// Sender returns the sender address of the transaction.
	Sender(cdata *CredentialData, sig *SignatureVal) (types.Address, error)
	// SignatureValues returns the raw R, S, V values corresponding to the
	// given signature.
	SignatureValues(sig []byte) (r, s, v *big.Int, err error)

	// Equal returns true if the given signer is the same as the receiver.
	Equal(Signer) bool
}

type AlgRandSigner struct {
	chainId, chainIdMul *big.Int
}

// NewBlockSigner returns a Signer based on the given chain config
func NewBlockSigner(chainId *big.Int) AlgRandSigner {
	if chainId == nil {
		chainId = new(big.Int)
	}
	return AlgRandSigner{
		chainId:    chainId,
		chainIdMul: new(big.Int).Mul(chainId, common.Big2),
	}
}

func (s AlgRandSigner) Equal(signer Signer) bool {
	bSigner, ok := signer.(AlgRandSigner)
	return ok && bSigner.chainId.Cmp(s.chainId) == 0
}

func (s AlgRandSigner) Sender(cdata *CredentialData, sig *SignatureVal) (types.Address, error) {
	if deriveChainId(&sig.V.IntVal).Cmp(s.chainId) != 0 {
		return types.Address{}, ErrInvalidChainId
	}

	V := &big.Int{}
	if s.chainId.Sign() != 0 {
		V = V.Sub(&sig.V.IntVal, s.chainIdMul)
		V.Sub(V, common.Big35)
	} else{
		V = V.Sub(&sig.V.IntVal, common.Big27)
	}
	address, err :=  recoverPlain(cdata.Hash(), &sig.R.IntVal, &sig.S.IntVal, V, true)
	return address, err
}


// SignatureValues returns a  R S V based given signature. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
func (s AlgRandSigner) SignatureValues(sig []byte) (R, S, V *big.Int, err error) {
	if len(sig) != 65 {
		errStr:=fmt.Sprintf("wrong size for signature: got %d, want 65", len(sig))
		err = errors.New(errStr)
		return nil, nil, nil, err
	}else{
		R = new(big.Int).SetBytes(sig[:32])
		S = new(big.Int).SetBytes(sig[32:64])

		if s.chainId.Sign() != 0 {
			V = big.NewInt(int64(sig[64] + 35))
			V.Add(V, s.chainIdMul)
		} else {
			V = new(big.Int).SetBytes([]byte{sig[64] + 27})
		}
	}

	return R, S, V, nil
}

func recoverPlain(sighash types.Hash, R, S, Vb *big.Int, homestead bool) (types.Address, error) {
	if Vb.BitLen() > 8 {
		return types.Address{}, ErrInvalidSig
	}
	V := byte(Vb.Uint64())
	if !crypto.ValidateSignatureValues(V, R, S, homestead) {
		return types.Address{}, ErrInvalidSig
	}
	// encode the snature in uncompressed format
	r, s := R.Bytes(), S.Bytes()
	sig := make([]byte, 65)
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = V
	// recover the public key from the snature
	pub, err := crypto.Ecrecover(sighash[:], sig)
	if err != nil {
		return types.Address{}, err
	}
	if len(pub) == 0 || pub[0] != 4 {
		return types.Address{}, errors.New("invalid public key")
	}
	var addr types.Address
	copy(addr[:], crypto.Keccak256(pub[1:])[12:])
	return addr, nil
}

// deriveChainId derives the chain id from the given v parameter
func deriveChainId(v *big.Int) *big.Int {
	if v.BitLen() <= 64 {
		v := v.Uint64()
		if v == 27 || v == 28 {
			return new(big.Int)
		}
		return new(big.Int).SetUint64((v - 35) / 2)
	}
	v = new(big.Int).Sub(v, big.NewInt(35))
	return v.Div(v, big.NewInt(2))
}
