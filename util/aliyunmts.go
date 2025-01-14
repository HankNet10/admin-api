package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	mts "github.com/aliyun/alibaba-cloud-sdk-go/services/mts"
)

type MtsSubmitJobsInput struct {
	Bucket   string
	Location string
	Object   string
}

func toInput(m MtsSubmitJobsInput) string {
	// m := mtsInput{
	// 	Bucket:   os.Getenv("123"),
	// 	Location: os.Getenv("ALI_REGION"),
	// 	Object:   object,
	// }
	res, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(res)

}

// "OutputObject": "{MTS_OUTPUT_OBJECT}",
// "TemplateId": "04b2436882134a9d82779ad99f63185c",
// "Encryption": {
// 	"Type": "hls-aes-128",
// 	"Key": "{ENCKEY}",
// 	"KeyType": "Base64",
// 	"KeyUri": "ZW5jLmtleQ=="
// }

const KEYURI = "{enc.key}"
const PATHNAME = "index.m3u8"

type MtsEnc struct {
	Type    string
	Key     []byte
	KeyType string
	KeyUri  string
}

type MtsSubmitJobsOutPuts struct {
	OutputObject string
	TemplateId   string
	Encryption   MtsEnc
}

func toOutputs(object string, key []byte) string {
	enc := MtsEnc{
		Type:    "hls-aes-128",
		KeyType: "Base64",
		KeyUri:  "e2VuYy5rZXl9",
		Key:     key,
	}
	m := MtsSubmitJobsOutPuts{
		OutputObject: object,
		TemplateId:   os.Getenv("ALI_MTS_TEMPLATEID"),
		Encryption:   enc,
	}
	j := []MtsSubmitJobsOutPuts{m}
	res, err := json.Marshal(j)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(res)
}

func MtsSubmitJobs(m MtsSubmitJobsInput, putObject string, key []byte) (string, error) {
	client, err := mts.NewClientWithOptions(
		os.Getenv("ALI_MTS_LOCATION"),
		sdk.NewConfig(),
		credentials.NewAccessKeyCredential(
			os.Getenv("ALI_ACCESS_KEY_ID"),
			os.Getenv("ALI_ACCESS_KEY_SECRET"),
		),
	)
	if err != nil {
		return "", err
	}
	request := mts.CreateSubmitJobsRequest()
	request.Scheme = "https"
	request.PipelineId = os.Getenv("ALI_MTS_PIPEID")
	request.OutputBucket = os.Getenv("ALI_OSS_ORIGIN")
	request.Input = toInput(m)
	request.Outputs = toOutputs(putObject, key)
	request.RegionId = os.Getenv("ALI_REGION")
	request.OutputLocation = os.Getenv("ALI_OSS_REGION")
	response, err := client.SubmitJobs(request)
	if err != nil {
		return "", err
	}
	if !response.IsSuccess() {
		return "", errors.New("请求mts网络失败")
	}
	if len(response.JobResultList.JobResult) == 1 {
		jobresult := response.JobResultList.JobResult[0]
		if !jobresult.Success {
			return "", errors.New("提交失败[" + jobresult.Code + "]" + jobresult.Message)
		}
		return jobresult.Job.JobId, nil
	} else {
		return "", errors.New("返回了异常的数据")
	}
}

func MtsQueryJob(jobid string) (*mts.QueryJobListResponse, error) {
	client, err := mts.NewClientWithOptions(
		os.Getenv("ALI_MTS_LOCATION"),
		sdk.NewConfig(),
		credentials.NewAccessKeyCredential(
			os.Getenv("ALI_ACCESS_KEY_ID"),
			os.Getenv("ALI_ACCESS_KEY_SECRET"),
		),
	)
	if err != nil {
		return nil, err
	}
	request := mts.CreateQueryJobListRequest()
	request.Scheme = "https"
	request.JobIds = jobid
	response, err := client.QueryJobList(request)
	if err != nil {
		return nil, err
	}
	if !response.IsSuccess() {
		return response, errors.New("网络请求失败")
	}
	return response, nil
}
