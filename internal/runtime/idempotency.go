package runtime

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// MakeIDKey gera uma chave determinística baseada no operador, repo e params.
func MakeIDKey(op Operator, in OpInput) string {
	if in.IdempotencyKey != "" {
		return in.IdempotencyKey
	}
	h := sha256.New()
	h.Write([]byte(op.Name()))
	h.Write([]byte("|"))
	h.Write([]byte(op.Version()))
	h.Write([]byte("|"))
	h.Write([]byte(in.Repo.Owner + "/" + in.Repo.Name + "@" + in.Repo.Head))
	h.Write([]byte("|"))
	b := MarshalParamsDeterministic(in.Params)
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

// MakeCacheKey pode ser igual ao de idempotência (ou incluir versão de cache se quiser).
func MakeCacheKey(op Operator, in OpInput) string { return MakeIDKey(op, in) }

// CanonicalJSON auxilia debugging/inspeção de chaves.
func CanonicalJSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}
