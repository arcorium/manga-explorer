package util

import (
  "errors"
  "github.com/gin-gonic/gin"
  "github.com/golang-jwt/jwt"
  "github.com/google/uuid"
  "golang.org/x/crypto/bcrypt"
  "manga-explorer/internal/util/obj"
  "math/rand"
  "strconv"
  "strings"
)

func Hash(str string) (string, error) {
  hashed, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
  return string(hashed), err
}

// DropError Used when the error is not necessary and returning only the first parameter
func DropError[T any](data T, err error, panicOnError ...bool) T {
  if len(panicOnError) > 0 && panicOnError[0] && err != nil {
    panic(err)
  }
  return data
}

// SetDefaultString Set real parameter string into def parameter if only real parameter string is empty
func SetDefaultString(real *string, def string) {
  if real == nil {
    return
  }
  if len(*real) == 0 {
    *real = def
  }
}

var NoContextValueErr = errors.New("no value with such key")
var MaltypeContextValueErr = errors.New("value has different type")

func GetContextValue[T any](ctx *gin.Context, key string) (T, error) {
  var result T
  t, ok := ctx.Get(key)
  if !ok {
    return result, NoContextValueErr
  }
  result, ok = t.(T)
  if !ok {
    return result, MaltypeContextValueErr
  }
  return result, nil
}

var EmptyStringErr = errors.New("string is empty")

func GetUintQuery(ctx *gin.Context, key string) (uint64, error) {
  str, present := ctx.GetQuery(key)
  if !present {
    return 0, EmptyStringErr
  }
  val, err := strconv.ParseUint(str, 10, 64)
  return val, err
}

func GetDefaultedUintQuery(ctx *gin.Context, key string, def uint64) uint64 {
  result, err := GetUintQuery(ctx, key)
  if err != nil {
    return def
  }
  return result
}

func Nil[T any]() *T {
  return (*T)(nil)
}

func Clone[T any](data *T) *T {
  temp := *data
  return &temp
}

func DoNothing(data ...any) {

}

func GenerateRandomString(length uint) string {
  const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

  builder := strings.Builder{}
  builder.Grow(int(length))

  var jj uint = 0
  DoNothing(jj)

  for i := 0; i < int(length); i++ {
    builder.WriteByte(letterBytes[rand.Intn(len(letterBytes))])
  }

  return builder.String()
}

func NilCount[U any, T []U | *U](datas ...T) uint {
  var counter uint = 0
  for i := 0; i < len(datas); i++ {
    if datas[i] == nil {
      counter += 1
    }
  }
  return counter
}

func GenerateJWTToken(claims jwt.Claims, method jwt.SigningMethod, secretKey []byte) (string, error) {
  if claims == nil {
    claims = jwt.MapClaims{}
  }
  return jwt.NewWithClaims(method, claims).SignedString(secretKey)
}

func IsUUID(str string) bool {
  return obj.Wrap(uuid.Parse(str)).Err() == nil
}

func IsOneOf[T comparable](original T, expecteds ...T) bool {
  for i := 0; i < len(expecteds); i++ {
    if expecteds[i] == original {
      return true
    }
  }
  return false
}

func SliceWrap[T any](datas ...T) []T {
  t := []T{}
  return append(t, datas...)
}
