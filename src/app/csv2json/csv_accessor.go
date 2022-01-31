package csv2json

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-gota/gota/dataframe"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"app/utils/apiutil"
	"app/utils/logger"
)

type csvAccessor struct{}

func NewCsvAccessor() *csvAccessor {
	return &csvAccessor{}
}

// GetCSVDataFrameFromApi はAPIコールによりCSVデータを取得する
func (ca *csvAccessor) GetCSVDataFrameFromApi(apiAddress string, apiId string) (*dataframe.DataFrame, time.Time, error) {

	// get json from api
	mapBody, err := apiutil.GetJsonMapFromResponseBody(fmt.Sprintf("%s?x=%s", apiAddress, apiId))
	if err != nil {
		logger.Errors(err)
		return nil, time.Time{}, err
	}

	// get csv address from json
	csvAddress, updatedDateTime, err := getCsvAddressFromBody(mapBody, apiId)
	if err != nil {
		logger.Errors(err)
		return nil, time.Time{}, err
	}

	// キャッシュしないアドレスへ変換する
	csvAddress = strings.Replace(csvAddress, "static.hamamatsu.odpf.net", "prd-hmpf-s3-odpf-01.s3.ap-northeast-1.amazonaws.com", 1)

	logger.Infof("csv address = %v", csvAddress)
	logger.Infof("update time = %v", updatedDateTime)

	// get bytes data from csv
	bytesCsv, err := apiutil.GetBytesFromResponseBody(csvAddress)
	if err != nil {
		logger.Errors(err)
		return nil, time.Time{}, err
	}

	// convert to dataframe from csv bytes data
	ioReaderCsv := strings.NewReader(string(bytesCsv))
	prCsv := transform.NewReader(ioReaderCsv, japanese.ShiftJIS.NewDecoder())
	dfCsv := dataframe.ReadCSV(prCsv, dataframe.WithDelimiter(','), dataframe.HasHeader(true))

	return &dfCsv, updatedDateTime, nil
}

func getCsvAddressFromBody(mapBody *map[string]interface{}, apiId string) (csvAddress string, updatedDateTime time.Time, errOut error) {
	csvAddress = ""
	errOut = nil

	bodystr := (*mapBody)["odpf_body"].(string)

	datetimeIndex := strings.Index(bodystr, "<tr><th>更新日</th><td>") + len("<tr><th>更新日</th><td>")
	datetimeCount := strings.Index(bodystr[datetimeIndex:], "分") + len("分")
	datetimestr := bodystr[datetimeIndex : datetimeIndex+datetimeCount]
	logger.Infos(datetimestr)
	replacer := strings.NewReplacer(
		"午前", "AM ",
		"午後", "PM ",
		"(日)", " Sun",
		"(月)", " Mon",
		"(火)", " Tue",
		"(水)", " Wed",
		"(木)", " Thu",
		"(金)", " Fri",
		"(土)", " Sat",
	)
	datetimestr = replacer.Replace(datetimestr)
	logger.Infos(datetimestr)
	layout := "2006年1月2日 Mon PM 3時4分 (MST)"
	updatedDateTime, extime := time.Parse(layout, datetimestr+" (JST)")
	logger.Infos(updatedDateTime)

	csvpathIndexEnd := strings.Index(strings.ToLower(bodystr), ".csv") + len(".csv")
	csvpathIndexStart := strings.LastIndex(bodystr[:csvpathIndexEnd], "http")
	csvAddress = bodystr[csvpathIndexStart:csvpathIndexEnd]
	logger.Infos(csvAddress)

	if csvAddress == "" || extime != nil {
		errMsg := "not found data from body"
		logger.Errors(errMsg)
		return csvAddress, updatedDateTime, fmt.Errorf("%s", errMsg)
	}

	return csvAddress, updatedDateTime, errOut
}

// GetTimeNow は今日の日付を取得する
func (ca *csvAccessor) GetTimeNow() time.Time {
	return time.Now()
}
