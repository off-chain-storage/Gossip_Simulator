package ecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"log"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
	gcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/pkg/errors"
)

// libp2p의 crypto.PrivKey를 ecdsa.PrivateKey 객체로 변환
func ConvertFromInterfacePrivKey(privkey crypto.PrivKey) (*ecdsa.PrivateKey, error) {
	// crypto.PrivKey가 Secp256k1 알고리즘을 사용하는지 Type assertion을 시도
	secpKey, ok := privkey.(*crypto.Secp256k1PrivateKey)
	if !ok { // 실패할 경우(Secp256k1 타입이 아닐 경우) 오류 메세지 반환
		return nil, errors.New("could not cast to Secp256k1PrivateKey")
	}

	// Raw 키 데이터 추출
	rawKey, err := secpKey.Raw()
	if err != nil {
		return nil, err
	}

	// 새로운 Ecdsa Private Key Instance 생성
	privKey := new(ecdsa.PrivateKey)
	// rawKey를 big.Int로 변환(ECDSA Private Key의 실제 숫자 값)
	k := new(big.Int).SetBytes(rawKey)
	// 개인 키의 'D' 필드에 숫자 값 할당
	privKey.D = k
	// 사용할 곡선을 Secp256k1으로 설정
	privKey.Curve = gcrypto.S256()
	// 개인 키로부터 공개 키의 좌표를 계산하여 할당
	privKey.X, privKey.Y = gcrypto.S256().ScalarBaseMult(rawKey)
	// 개인 키 반환
	return privKey, nil
}

// 일반적으로 블록체인에서 사용되는 ecdsa.PrivateKey 객체를 libp2p의 crypto.PrivKey 인터페이스로 변환
func ConvertToInterfacePrivkey(privkey *ecdsa.PrivateKey) (crypto.PrivKey, error) {
	// ecdsa.PrivateKey 객체에서 'D'필드 바이트 배열로 추출(개인키의 값)
	privBytes := privkey.D.Bytes()

	// ECDSA 개인키는 길이가 일정하지 않을 수 있으므로 libp2p의 Secp256k1 개인키 형식에 맞추기 위해
	// 배열의 길이가 32 바이트 미만일 경우, 앞쪽을 0으로 채워서 길이를 맞춤
	if len(privBytes) < 32 {
		privBytes = append(make([]byte, 32-len(privBytes)), privBytes...)
	}

	// 바이트 배열을 Secp256K1 개인키 형식으로 변환하여 반환
	return crypto.UnmarshalSecp256k1PrivateKey(privBytes)
}

func ConvertToInterfacePubkey(pubkey *ecdsa.PublicKey) (crypto.PubKey, error) {
	xVal, yVal := new(btcec.FieldVal), new(btcec.FieldVal)
	if xVal.SetByteSlice(pubkey.X.Bytes()) {
		return nil, errors.Errorf("X value overflows")
	}
	if yVal.SetByteSlice(pubkey.Y.Bytes()) {
		return nil, errors.Errorf("Y value overflows")
	}
	newKey := crypto.PubKey((*crypto.Secp256k1PublicKey)(btcec.NewPublicKey(xVal, yVal)))
	// Zero out temporary values.
	xVal.Zero()
	yVal.Zero()
	return newKey, nil
}

// Convert from *ecdsa.PrivateKey to []byte
func ConvertToByteEcdsaPrivKey(privKey *ecdsa.PrivateKey) []byte {
	privKeyBytes := privKey.D.Bytes()

	return privKeyBytes
}

// Convert from []byte to *ecdsa.PrivateKey
func ConvertToEcdsaPrivKeyByte(privKeyBytes []byte) (privKey *ecdsa.PrivateKey) {
	privKey = new(ecdsa.PrivateKey)

	privKey.D = new(big.Int).SetBytes(privKeyBytes)
	privKey.Curve = secp256k1.S256()
	privKey.X, privKey.Y = secp256k1.S256().ScalarBaseMult(privKeyBytes)

	return
}

// *ecdsa.PublicKey 를 string type으로 Convert 해주는 Function
func ConvertToStringEcdsaPubKey(ecdsaPubKey *ecdsa.PublicKey) string {
	pubKeyBytes := secp256k1.CompressPubkey(ecdsaPubKey.X, ecdsaPubKey.Y)

	pubKeyHex := hex.EncodeToString(pubKeyBytes)

	return pubKeyHex
}

// string type을 *ecdsa.PublicKey 으로 Convert 해주는 Function
func ConvertToEcdsaPubKeyString(pubKeyHex string) (*ecdsa.PublicKey, error) {
	// Hex 문자열을 바이트 슬라이스로 디코딩
	pubKeyBytes, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		log.Fatalf("failed to decode: %v", err)
	}

	// 바이트 슬라이스를 이용해 ECDSA 공개 키 재구성
	pubKey_x, pubKey_y := secp256k1.DecompressPubkey(pubKeyBytes)

	// *ecdsa.PublicKey 생성
	ecdsaPubKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     pubKey_x,
		Y:     pubKey_y,
	}

	return &ecdsaPubKey, nil
}
