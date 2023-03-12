package helpers

import (
	"fmt"
	"strings"
	"time"
)

type DateConverter struct {
	UnixTimestamp int64
	Time          time.Time
}

type IDateConverter interface {
	New(UnixTimestamp int)
	ConvertIntoFirstYear() error
	ConvertIntoFirstMonth() error
	ConvertIntoFirstDay() error
}

func (d *DateConverter) New(UnixTimestamp int64) {
	d.UnixTimestamp = UnixTimestamp
}

func (d *DateConverter) ConvertIntoFirstYear() error {
	d.Time = time.Unix(d.UnixTimestamp, 0)
	year, _, _ := d.Time.Date()
	parsingString := fmt.Sprintf("%v", year)
	result, err := time.Parse("2006", parsingString)
	if err != nil {
		return err
	} else {
		d.Time = result
		return nil
	}
}
func (d *DateConverter) ConvertIntoFirstMonth() error {
	d.Time = time.Unix(d.UnixTimestamp, 0)
	year, month, _ := d.Time.Date()
	parsingString := fmt.Sprintf("%v-%v", year, month)
	fmt.Printf("parsingString : %v \n", parsingString)
	result, err := time.Parse("2006-January", parsingString)
	if err != nil {
		return err
	} else {
		d.Time = result
		return nil
	}
}
func (d *DateConverter) ConvertIntoFirstDay() error {
	d.Time = time.Unix(d.UnixTimestamp, 0)
	year, month, day := d.Time.Date()
	parsingString := fmt.Sprintf("%v-%v-%v", year, month, day)
	fmt.Printf("parsingString : %v \n", parsingString)
	result, err := time.Parse("2006-January-2", parsingString)
	if err != nil {
		return err
	} else {
		d.Time = result
		return nil
	}
}

// func (d *dateConverter) ConvertUnixToFormatType(format string) (*time.Time, error) {
// 	parsingString := strings.Replace(*d.FormatTimestamp, "/", "-", -1)
// 	//format "2006-01-02"
// 	result, errParseTime := time.Parse(format, parsingString)
// 	if errParseTime != nil {
// 		return nil, errParseTime
// 	}
// 	return &result, nil
// }

/*
   |--------------------------------------------------------------------------
   | Validate date time and convert string to time.Time
   |--------------------------------------------------------------------------
   |
   | This function convert string to time.Time,
   | Before convert to time.Time string with "/" must replace with "-".
   | This function return time.Time type end error for time validation.
*/

func ValidateDate(date string) (*time.Time, error) {

	dateFormat := strings.Replace(date, "/", "-", -1)

	result, errParseTime := time.Parse("2006-01-02", dateFormat)
	if errParseTime != nil {
		return nil, errParseTime
	}
	return &result, nil

}

/*
   |--------------------------------------------------------------------------
   | Convert Unix Timestamp to time.Time
   |--------------------------------------------------------------------------
   |
   | This function convert UnixTimestamp to time.Time
*/

func ConvertUnixToDate(date int) time.Time {
	//make timezone WIB
	loc, _ := time.LoadLocation("Asia/Jakarta")
	result := time.Unix(int64(date), 0).In(loc)

	return result
}

/*
   |--------------------------------------------------------------------------
   | Validate unixtamestamp
   |--------------------------------------------------------------------------
   |
   | This function Make Custom Binding Validation For Enum data Type
*/
// func Unix(sec int64, nsec int64) (*time.Time, error)

func ConvertUnix(date string) (*time.Time, error) {

	resultDate, err := ValidateDate(date)
	if err != nil {
		return nil, err
	}
	// if your unix timestamp is not in int64 format

	// int64 to time.Time
	myTime := time.Unix(resultDate.Unix(), 0)

	return &myTime, nil
}

/*
   |--------------------------------------------------------------------------
   | Convert Unix time.Time to Timestamp
   |--------------------------------------------------------------------------
   |
   | This function convert UnixTimestamp to time.Time
*/

func ConvertDateToUnix(date time.Time) int64 {

	result := date.Unix()

	return result

}
