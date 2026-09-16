package vine

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var (
	BASEURL = "https://apivin.cbone.uk/"
)

type VineUser struct {
	Username string
	UserId   int
	Key      string
}

func CallApi(endpoint string,
	method string,
	values url.Values,
	user *VineUser) (*json.RawMessage, error) {

	uri := BASEURL + endpoint

	if values != nil && user == nil {

		resp, _ := http.PostForm(uri, values)
		return processResp(resp)

	} else if user != nil {

		client := &http.Client{}

		req, _ := http.NewRequest(method, uri, nil)
		req.Header.Add("User-Agent", "com.vine.iphone/1.0.3 (unknown, iPhone OS 6.0.1, iPhone, Scale/2.000000)")
		req.Header.Add("Accept-Language", "en, sv, fr, de, ja, nl, it, es, pt, pt-PT, da, fi, nb, ko, zh-Hans, zh-Hant, ru, pl, tr, uk, ar, hr, cs, el, he, ro, sk, th, id, ms, en-GB, ca, hu, vi, en-us;q=0.8")
		req.Header.Add("vine-session-id", user.Key)

		resp, _ := client.Do(req)

		return processResp(resp)
	}

	log.Fatal("Wrong CallApi arguments")
	return nil, nil
}

func processResp(resp *http.Response) (*json.RawMessage, error) {
	defer resp.Body.Close()

	js, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var f map[string]*json.RawMessage
	err = json.Unmarshal(js, &f)

	if err != nil {
		return nil, err
	}

	return f["data"], nil
}

func Login(username string, password string) (*VineUser, error) {
	values := url.Values{"username": {username},
		"password": {password}}

	rsp, err := CallApi("users/authenticate", "POST", values, nil)

	if err != nil {
		return nil, err
	}

	var user *VineUser
	err = json.Unmarshal([]byte(*rsp), &user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

/*func main() {
    user, err := Login("", "")

    if err != nil {
        log.Fatal("Wrong username/password")
    }

    resp, _ := CallApi("timelines/users/"+strconv.Itoa(user.UserId), "GET", nil, user)

    if resp != nil {
        fmt.Println(string(*resp))
    }

}*/
