// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package table

// [START bigquery_inserting_data_types]
import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
)

// ComplexType represents a complex row item
type ComplexType struct {
	Name         string                 `bigquery:"name"`
	Age          int                    `bigquery:"age"`
	School       []byte                 `bigquery:"school"`
	Location     bigquery.NullGeography `bigquery:"location"`
	Measurements []float64              `bigquery:"measurements"`
	DatesTime    DatesTime              `bigquery:"datesTime"`
}

// DatesTime represents different date/time representation
type DatesTime struct {
	Day        civil.Date     `bigquery:"day"`
	FirstTime  civil.DateTime `bigquery:"firstTime"`
	SecondTime civil.Time     `bigquery:"secondTime"`
	ThirdTime  time.Time      `bigquery:"thirdTime"`
}

// insertingDataTypes demonstrates inserting data into a table using the streaming insert mechanism.
func insertingDataTypes(projectID, datasetID, tableID string) error {
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	// tableID := "mytable"
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("bigquery.NewClient: %w", err)
	}
	defer client.Close()

	// Manual defining schema
	schema := bigquery.Schema{
		{Name: "name", Type: bigquery.StringFieldType},
		{Name: "age", Type: bigquery.IntegerFieldType},
		{Name: "school", Type: bigquery.BytesFieldType},
		{Name: "location", Type: bigquery.GeographyFieldType},
		{Name: "measurements", Type: bigquery.FloatFieldType, Repeated: true},
		{Name: "datesTime", Type: bigquery.RecordFieldType, Schema: bigquery.Schema{
			{Name: "day", Type: bigquery.DateFieldType},
			{Name: "firstTime", Type: bigquery.DateTimeFieldType},
			{Name: "secondTime", Type: bigquery.TimeFieldType},
			{Name: "thirdTime", Type: bigquery.TimestampFieldType},
		}},
	}
	// Infer schema from struct
	// schema, err := bigquery.InferSchema(ComplexType{})

	table := client.Dataset(datasetID).Table(tableID)
	err = table.Create(ctx, &bigquery.TableMetadata{
		Schema: schema,
	})
	if err != nil {
		return fmt.Errorf("table.Create: %w", err)
	}
	day, _ := civil.ParseDate("2019-01-12")
	firstTime, _ := civil.ParseDateTime("2019-02-17T11:24:00.000")
	secondTime, _ := civil.ParseTime("14:00:00")
	thirdTime, _ := time.Parse(time.RFC3339Nano, "2020-04-27T18:07:25.356Z")
	row := &ComplexType{
		Name:         "Tom",
		Age:          30,
		School:       []byte("Test University"),
		Location:     bigquery.NullGeography{GeographyVal: "POINT(1 2)", Valid: true},
		Measurements: []float64{50.05, 100.5},
		DatesTime: DatesTime{
			Day:        day,
			FirstTime:  firstTime,
			SecondTime: secondTime,
			ThirdTime:  thirdTime,
		},
	}
	rows := []*ComplexType{row}
	inserter := table.Inserter()
	if err := inserter.Put(ctx, rows); err != nil {
		return err
	}
	return nil
}

// [END bigquery_inserting_data_types]
