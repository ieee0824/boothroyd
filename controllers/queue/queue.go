package queue

import (
	"github.com/jobtalk/hawkeye/models/queue"
	"time"
	"gopkg.in/kataras/iris.v6"
	"github.com/jobtalk/hawkeye/models/work"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"net/url"
	"log"
)

var qs = queue.New()

var hosts = []string{}

func execLambdaFunction(param []byte)(string, error){
	sess := session.New(&aws.Config{Region: aws.String("ap-northeast-1")})
	svc := lambda.New(sess)

	result, err := svc.Invoke(&lambda.InvokeInput{
		FunctionName: aws.String("lupin_lupin"),
		Payload: param,
	})
	return result.GoString(), err
}

func Run() {
	for {
		select {
		case d := <- qs.C:
			var target []byte
			if _, ok := d.([]byte); ok {
				target = d.([]byte)
			} else if _, ok := d.(string); ok {
				target = []byte(d.(string))
			} else {
				continue
			}

			result, err := execLambdaFunction(target)
			fmt.Println(result, err)
		}
	}
}

func Check(ctx *iris.Context) {
	ctx.JSON(
		200,
		qs.Status(),
	)
}

func Enqueue(ctx *iris.Context) {
	var ws = []work.Work{}
	if err := ctx.ReadJSON(&ws); err != nil {
		log.Println(err)
		ctx.JSON(
			500,
			map[string]interface{}{
				"status": false,
				"error": err.Error(),
			},
		)
		return
	}
	for _, w := range ws {
		u, err := url.Parse(w.URL)
		if err != nil {
			log.Println(err)
			ctx.JSON(
				500,
				map[string]interface{}{
					"status": false,
					"error": err.Error(),
				},
			)
			return
		}

		w.JobID = fmt.Sprint(time.Now().UnixNano())
		if bin, err := json.Marshal(w); err != nil {
			log.Println(err)
			ctx.JSON(
				500,
				map[string]interface{}{
					"status": false,
					"error": err.Error(),
				},
			)
			return
		} else {
			qs.Enqueue(u.Host, bin)
		}
	}


	ctx.JSON(
		200,
		map[string]interface{}{
			"status": true,
			"error": nil,
			"data": ws,
		},
	)
}
