package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	chart "github.com/wcharczuk/go-chart"
	datadog "gopkg.in/zorkian/go-datadog-api.v2"
)

type ConversionFamily struct {
	Name     string
	Units    []ConversionUnit
	BaseUnit ConversionUnit
}

type ConversionUnit struct {
	// Name of the unit
	Name string
	// Short name of the unit
	ShortName string
	// Multiplier is the multiplication factor from the base unit
	Multiplier float32
}

type GraphClient struct {
	DatadogClient *datadog.Client
	units         map[string]ConversionFamily
}

func NewGraphClient(apiKey, appKey string) *GraphClient {
	return &GraphClient{
		DatadogClient: datadog.NewClient(apiKey, appKey),
		units: map[string]ConversionFamily{
			"hits": ConversionFamily{
				Name: "hits",
				Units: []ConversionUnit{
					ConversionUnit{
						Name:       "thousands",
						ShortName:  "k",
						Multiplier: 1000.0,
					},
					ConversionUnit{
						Name:       "hits",
						ShortName:  "",
						Multiplier: 1.0,
					},
				},
				BaseUnit: ConversionUnit{
					Name:       "hits",
					ShortName:  "",
					Multiplier: 1.0,
				},
			},
			"time": ConversionFamily{
				Name: "time",
				Units: []ConversionUnit{
					ConversionUnit{
						Name:       "microseconds",
						ShortName:  "us",
						Multiplier: 1000000.0,
					},
					ConversionUnit{
						Name:       "milliseconds",
						ShortName:  "ms",
						Multiplier: 1000.0,
					},
					ConversionUnit{
						Name:       "seconds",
						ShortName:  "s",
						Multiplier: 1.0,
					},
				},
				BaseUnit: ConversionUnit{
					Name:       "seconds",
					ShortName:  "s",
					Multiplier: 1.0,
				},
			},
		},
	}
}

func (c *GraphClient) Graph(title string, metric string, from time.Time, to time.Time) (bytes.Buffer, error) {
	var buf bytes.Buffer

	metrics, err := c.DatadogClient.QueryMetrics(to.Unix(), from.Unix(), metric)
	if err != nil {
		return buf, err
	}

	var unitFamilyOne string
	var unitFamilyTwo string
	var hasFamilyOne bool
	var hasFamilyTwo bool
	var unitPointsOne []float64
	var unitPointsTwo []float64

	var series []chart.Series

	for _, metric := range metrics {
		unit := metric.GetUnits()[0]
		if unit == nil {
			unit = &datadog.Unit{Family: "hits", ScaleFactor: 1.0}
		}
		unitFamily, ok := c.units[unit.Family]
		if !ok {
			log.Println("Unit family not found:", unit.Family)
			return buf, errors.New("Unit not found")
		}

		var unitPoints *[]float64

		log.Println(unitFamily)

		if !hasFamilyOne || unitFamilyOne == unitFamily.Name {
			unitFamilyOne = unitFamily.Name
			hasFamilyOne = true
			unitPoints = &unitPointsOne
			goto family_chosen
		}
		if (hasFamilyOne && !hasFamilyTwo) || unitFamilyOne == unitFamily.Name {
			unitFamilyTwo = unitFamily.Name
			hasFamilyTwo = true
			unitPoints = &unitPointsTwo
			goto family_chosen
		}

		if hasFamilyOne && hasFamilyTwo && (unitFamilyOne != unitFamily.Name || unitFamilyTwo != unitFamily.Name) {
			return buf, errors.New("too many units returned by query")
		}

		return buf, errors.New("unknown error")

	family_chosen:
		baseScaleFactor := unit.ScaleFactor
		log.Println("Unit scale:", unit.ScaleFactor)
		log.Println("Base scale:", baseScaleFactor)

		for _, point := range metric.Points {
			// Converts to the base scale
			log.Println(point[1])
			val := point[1] * float64(baseScaleFactor)
			*unitPoints = append(*unitPoints, val)
		}
	}

	graph := chart.Chart{
		Title: title,
		TitleStyle: chart.Style{
			Show: true,
		},
		Width:  1280,
		Height: 720,
		Background: chart.Style{
			Padding: chart.Box{
				Top:    50,
				Bottom: 30,
			},
		},
	}

	if hasFamilyOne {
		sort.Float64s(unitPointsOne)

		unitFamily := c.units[unitFamilyOne]
		chosenUnit := unitFamily.Units[0]
		for _, unit := range unitFamily.Units {
			if float64(unit.Multiplier)*unitPointsOne[len(unitPointsOne)-1] < 1.0 {
				break
			}
			chosenUnit = unit
		}

		log.Println("Chosen unit:", chosenUnit.Name)

		for i, metric := range metrics {
			if unitFamilyOne == unitFamily.Name {
				baseScaleFactor := chosenUnit.Multiplier

				timeSeries := chart.TimeSeries{
					Style: chart.Style{
						Show:        true,
						StrokeColor: chart.GetAlternateColor(i),
						FillColor:   chart.GetAlternateColor(i).WithAlpha(100),
					},
				}

				for _, point := range metric.Points {
					// Converts to the base scale
					val := point[1] * float64(baseScaleFactor)

					timeSeries.XValues = append(timeSeries.XValues, time.Unix(0, int64(point[0])*1000000))
					timeSeries.YValues = append(timeSeries.YValues, val)
				}

				timeSeries.Name = metric.GetScope()
				series = append(series, timeSeries)
			}
		}

		graph.YAxis = chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
			ValueFormatter: func(v interface{}) string {
				typed := v.(float64)
				return fmt.Sprintf("%.1f %s", typed, chosenUnit.ShortName)
			},
			TickStyle: chart.Style{
				TextRotationDegrees: 45.0,
			},
		}
	}

	if hasFamilyTwo {
		sort.Float64s(unitPointsTwo)

		unitFamily := c.units[unitFamilyTwo]
		chosenUnit := unitFamily.Units[0]
		for _, unit := range unitFamily.Units {
			if float64(unit.Multiplier)*unitPointsTwo[len(unitPointsTwo)-1] < 1.0 {
				break
			}
			chosenUnit = unit
		}

		for i, metric := range metrics {
			if unitFamilyTwo == unitFamily.Name {
				baseScaleFactor := chosenUnit.Multiplier

				timeSeries := chart.TimeSeries{
					Style: chart.Style{
						Show:        true,
						StrokeColor: chart.GetAlternateColor(i),
						FillColor:   chart.GetAlternateColor(i).WithAlpha(100),
					},
				}

				for _, point := range metric.Points {
					// Converts to the base scale
					val := point[1] * float64(baseScaleFactor)

					timeSeries.XValues = append(timeSeries.XValues, time.Unix(0, int64(point[0])*1000000))
					timeSeries.YValues = append(timeSeries.YValues, val)
				}

				timeSeries.Name = metric.GetScope()
				series = append(series, timeSeries)
			}
		}

		graph.YAxisSecondary = chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
			ValueFormatter: func(v interface{}) string {
				typed := v.(float64)
				return fmt.Sprintf("%.1f %s", typed, chosenUnit.ShortName)
			},
			TickStyle: chart.Style{
				TextRotationDegrees: 45.0,
			},
		}
	}

	graph.Series = series
	graph.XAxis = chart.XAxis{
		Style: chart.Style{
			Show: true,
		},
		// ValueFormatter: chart.TimeValueFormatterWithFormat("Jan _2 3:04PM"),
		ValueFormatter: chart.TimeValueFormatterWithFormat("Jan _2 15:04"),
		GridMajorStyle: chart.Style{
			Show:        true,
			StrokeColor: chart.ColorAlternateGray,
			StrokeWidth: 0.5,
		},
		TickStyle: chart.Style{
			TextRotationDegrees: 50.0,
		},
	}
	graph.Elements = append(graph.Elements, chart.LegendLeft(&graph))
	err = graph.Render(chart.PNG, &buf)

	return buf, err
}

func rpcLatency(from time.Time, to time.Time) (bytes.Buffer, error) {
	var buf bytes.Buffer

	client := datadog.NewClient("dd5930a2b34093f052aea1eeb290f11b", "a138ebaed3072d4c04e063b8ee66f686980aa794")
	metrics, err := client.QueryMetrics(to.Unix(), from.Unix(), "avg:trace.rpc.request.duration{*} by {resource_name}")
	if err != nil {
		return buf, err
	}

	var series []chart.Series
	var scaleFactor float64
	var unitName string

	points := make([]float64, 0, 0)
	for _, metric := range metrics {
		for _, point := range metric.Points {
			points = append(points, point[1])
		}
	}

	if len(points) > 0 {
		sort.Float64s(points)
		kPercentIndex := int(1.0 * float64(len(points)))
		kPercentile := points[kPercentIndex-1]

		// Second units.
		scaleFactor = 1.0
		unitName = "s"
		if kPercentile < 1.0 {
			if kPercentile*1000.0 < 1.0 {
				// Convert to microseconds
				scaleFactor = 1000000.0
				unitName = "us"
			} else {
				// Convert to milliseconds
				scaleFactor = 1000.0
				unitName = "ms"
			}
		}
		log.Println("Calculated scale factor is:", scaleFactor)
	}
	scaleFactor = 1.0
	unitName = "s"

	for i, metric := range metrics {
		timeSeries := chart.TimeSeries{
			Style: chart.Style{
				Show:        true,
				StrokeColor: chart.GetAlternateColor(i),
				FillColor:   chart.GetAlternateColor(i).WithAlpha(100),
			},
		}
		for _, point := range metric.Points {
			timeSeries.XValues = append(timeSeries.XValues, time.Unix(0, int64(point[0])*1000000))
			timeSeries.YValues = append(timeSeries.YValues, point[1]*scaleFactor)
		}

		timeSeries.Name = metric.GetScope()

		series = append(series, timeSeries)
	}

	graph := chart.Chart{
		Title: "RPC Latency",
		TitleStyle: chart.Style{
			Show: true,
		},
		Width:  1280,
		Height: 720,
		Background: chart.Style{
			Padding: chart.Box{
				Top:    50,
				Bottom: 30,
			},
		},
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true,
			},
			ValueFormatter: chart.TimeValueFormatterWithFormat("3:04PM"),
			GridMajorStyle: chart.Style{
				Show:        true,
				StrokeColor: chart.ColorAlternateGray,
				StrokeWidth: 0.5,
			},
			TickStyle: chart.Style{
				TextRotationDegrees: 50.0,
			},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
			ValueFormatter: func(v interface{}) string {
				typed := v.(float64)
				return fmt.Sprintf("%.1f %s", typed, unitName)
			},
			TickStyle: chart.Style{
				TextRotationDegrees: 45.0,
			},
		},
		Series: series,
	}

	graph.Elements = append(graph.Elements, chart.LegendLeft(&graph))
	graph.Render(chart.PNG, &buf)

	return buf, nil
}

func rpcRequests(from time.Time, to time.Time) (bytes.Buffer, error) {
	var buf bytes.Buffer

	client := datadog.NewClient("dd5930a2b34093f052aea1eeb290f11b", "a138ebaed3072d4c04e063b8ee66f686980aa794")
	metrics, err := client.QueryMetrics(to.Unix(), from.Unix(), "avg:trace.rpc.request.hits{*} by {resource_name}.as_count()")
	if err != nil {
		return buf, err
	}

	var series []chart.Series
	var scaleFactor float64
	var unitName string

	// Second units.
	scaleFactor = 1.0
	unitName = "hits"

	for i, metric := range metrics {
		timeSeries := chart.TimeSeries{
			Style: chart.Style{
				Show:        true,
				StrokeColor: chart.GetAlternateColor(i),
				FillColor:   chart.GetAlternateColor(i).WithAlpha(100),
			},
		}
		for _, point := range metric.Points {
			timeSeries.XValues = append(timeSeries.XValues, time.Unix(0, int64(point[0])*1000000))
			timeSeries.YValues = append(timeSeries.YValues, point[1]*scaleFactor)
		}

		timeSeries.Name = metric.GetScope()

		series = append(series, timeSeries)
	}

	graph := chart.Chart{
		Title: "RPC Requests",
		TitleStyle: chart.Style{
			Show: true,
		},
		Width:  1280,
		Height: 720,
		Background: chart.Style{
			Padding: chart.Box{
				Top:    50,
				Bottom: 30,
			},
		},
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true,
			},
			ValueFormatter: chart.TimeValueFormatterWithFormat("3:04PM"),
			GridMajorStyle: chart.Style{
				Show:        true,
				StrokeColor: chart.ColorAlternateGray,
				StrokeWidth: 0.5,
			},
			TickStyle: chart.Style{
				TextRotationDegrees: 50.0,
			},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
			ValueFormatter: func(v interface{}) string {
				typed := v.(float64)
				return fmt.Sprintf("%.0f %s", typed, unitName)
			},
			TickStyle: chart.Style{
				TextRotationDegrees: 45.0,
			},
		},
		Series: series,
	}

	graph.Elements = append(graph.Elements, chart.LegendLeft(&graph))
	graph.Render(chart.PNG, &buf)

	return buf, nil
}
