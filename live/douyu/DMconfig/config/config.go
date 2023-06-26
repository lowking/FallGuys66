package config

type DMconfig struct {
	Rid            string
	LoginMsg       string
	LoginJoinGroup string
	Url            string
}

var SpiderConfig *DMconfig

func init() {
	SpiderConfig = &DMconfig{
		Rid:            "156277",
		LoginMsg:       "type@=loginreq/room_id@=%s/dfl@=sn@A=105@Sss@A=1/username@=%s/uid@=%s/ver@=20190610/aver@=218101901/ct@=0/",
		LoginJoinGroup: "type@=joingroup/rid@=%s/gid@=-9999/",
		Url:            "wss://danmuproxy.douyu.com:8506/",
	}
}
