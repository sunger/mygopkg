package gin_

import (
	"errors"
	"time"
	"github.com/dgrijalva/jwt-go"
)

type UserInfo struct {
	ID string
	Db string
}

const SecrectStr string = "z#@!y_q"

func CreateToken(user *UserInfo) (tokenss string, err error) {
	//自定义claim
	claim := jwt.MapClaims{
		"id": user.ID,
		// "username": user.Username,
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	// 签名后的token格式说明
	// 示例：eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1NzUyODUyNzUsImlkIjoxLCJuYmYiOjE1NzUyODUyNzUsInVzZXJuYW1lIjoicGVuZ2oifQ.bDe8UZYLxvmrK7gHcuK8TrlnoiMsIm3Jo_f0-YYle7E
	// 使用符号.，被分割成了三段
	// 第一段base64解码之后：{"alg":"HS256","typ":"JWT"}
	// 第二段base64解码之后：{"iat":1575285275,"id":1,"nbf":1575285275,"username":"pengj"}，是原始的数据。
	// 第三段是使用SigningMethodHS256加密之后的文本
	tokenss, err = token.SignedString([]byte(SecrectStr))
	return
}

func secret() jwt.Keyfunc {
	return func(token *jwt.Token) (i interface{}, e error) {
		return []byte(SecrectStr), nil
		// return i, e
	}
}

func ParseToken(tokenss string) (user UserInfo, err error) {
	user = UserInfo{}
	token, err := jwt.Parse(tokenss, secret())
	if err != nil {
		return user,err
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("cannot convert claim to mapclaim")
		return user,err
	}
	//验证token，如果token被修改过则为false
	if !token.Valid {
		err = errors.New("token is invalid")
		return user,err
	}

	user.ID = claim["id"].(string)
	// user.Username = claim["username"].(string)
	return user, nil
}
