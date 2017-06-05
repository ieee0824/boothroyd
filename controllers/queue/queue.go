package queue

import (
	"net/url"
	"github.com/jobtalk/hawkeye/models/queue"
	"math/rand"
	"time"
	"github.com/pkg/errors"
	"gopkg.in/kataras/iris.v6"
	"github.com/jobtalk/hawkeye/models/work"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/aws"
)

func init(){
	rand.Seed(time.Now().UnixNano())
}

var hosts = []string{}

type queues map[string]*queue.Queue

func execLambdaFunction(param []byte)(string, error){
	sess := session.New(&aws.Config{Region: aws.String("ap-northeast-1")})
	svc := lambda.New(sess)

	result, err := svc.Invoke(&lambda.InvokeInput{
		FunctionName: aws.String("lupin_lupin"),
		Payload: param,
	})
	return result.GoString(), err
}

func shuffle(data []string) {
	n := len(data)
	for i := n - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		data[i], data[j] = data[j], data[i]
	}
}

func (q *queues) Enqueue(s string) error {
	w := &work.Work{}
	if err := json.Unmarshal([]byte(s), w); err != nil {
		return err
	}
	u, err := url.Parse(w.URL)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if _, ok := (*q)[u.Host]; !ok {
		hosts = append(hosts, u.Host)
		(*q)[u.Host] = queue.New()
	}
	(*q)[u.Host].Enqueue(s)
	return nil
}

func (q *queues) Dequeue()(string, error) {
	keys := make([]string, len(hosts))
	copy(keys, hosts)
	shuffle(keys)

	for _, key := range keys {
		ret, err := (*q)[key].Dequeue()
		if err == nil {
			return ret, nil
		}
	}
	return "", errors.New("There is no data that can be retrieved")
}

var qs = &queues{}



func Run() {
	for {
		target, err := qs.Dequeue()
		if err != nil {
			continue
		}
		fmt.Println(target)
		// lambda functionを呼び出す
		result, err := execLambdaFunction([]byte(target))
		fmt.Println(result, err)
	}
}

func Check(ctx *iris.Context) {
	ctx.JSON(
		200,
		qs,
	)
}

func Enqueue(ctx *iris.Context) {
	var w = &work.Work{}
	if err := ctx.ReadJSON(w); err != nil {
		ctx.JSON(
			500,
			map[string]interface{}{
				"status": false,
				"error": err,
			},
		)
		return
	}
	w.JobID = fmt.Sprint(time.Now().UnixNano())
	if bin, err := json.Marshal(w); err != nil {
		ctx.JSON(
			500,
			map[string]interface{}{
				"status": false,
				"error": err,
			},
		)
		return
	} else {
		qs.Enqueue(string(bin))
	}

	ctx.JSON(
		200,
		map[string]interface{}{
			"status": true,
			"error": nil,
		},
	)
}
