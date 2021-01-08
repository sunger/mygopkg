package comm

import (
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"gorm.io/gorm"
	"github.com/sunger/mygopkg/db"
	"github.com/sunger/mygopkg/log"
)

// SendRequest 发送request
func SendRequest(url string, body io.Reader, addHeaders map[string]string, method string) (resp []byte, err error) {
	// 1、创建req
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	//req.Header.Add("Content-Type", "application/json")

	// 2、设置headers
	if len(addHeaders) > 0 {
		for k, v := range addHeaders {
			req.Header.Add(k, v)
		}
	}

	// 3、发送http请求
	client := &http.Client{}
	response, err := client.Do(req)
	defer response.Body.Close()
	if err != nil {
		return nil,err
	}
	//if response.StatusCode != 200 {
	//	err = errors.New("http status err")
	//	fmt.Errorf("sendRequest failed, url=%v, response status code=%d", url, response.StatusCode)
	//	return nil
	//}

	// 4、结果读取
	resp, err = ioutil.ReadAll(response.Body)
	return resp, err
}

//加载远程数据库连接到内存
func LoadRemoteDb(url, token string, config *gorm.Config) {

	header := map[string]string{
		"token":        token,
		"Content-Type": "application/json",
	}

	bts, err := SendRequest(url, nil, header, "GET")
	if err != nil {
		log.GetLog().Error("远程请求数据库失败", zap.String("错误信息：", err.Error()))
	}

	rv := db.DbConnsResponse{}

	err = json.Unmarshal(bts, &rv)
	if err != nil {
		log.GetLog().Error("反序列化返回值失败", zap.String("错误信息：", err.Error()))
	}

	//strs := string(bts)

	if rv.Code == 0 {
		//dbstr := rv.Data.(string)
		//if dbstr == "" {
		//	log.GetLog().Error("返回的数据为空")
		//}
		//list := make([]db.DbConn, 0)
		//list := rv.Data.(string)
		//err = json.Unmarshal([]byte(dbstr), &list)
		//if err != nil {
		//	log.GetLog().Error("反序列化数据库列表失败", zap.String("错误信息：", err.Error()))
		//}

		db.MapListToDBService(rv.Data,config)
	} else {
		log.GetLog().Error("返回的Code不等于0")
	}

}
