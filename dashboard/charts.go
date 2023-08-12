package dashboard

import (
	db "pingo/database"
	"pingo/job"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func getLineItems(data *[]db.DBJobLog) []opts.LineData {
	items := make([]opts.LineData, len(*data))
	for i, v := range *data {
		var statusIcon string
		switch v.Status {
		case job.SUCCESS:
			statusIcon = "🟢"
		case job.FAILED:
			statusIcon = "🔴"
		}
		items[i] = opts.LineData{Name: statusIcon, Value: v.PerfTime}
	}
	return items
}

func getLineXAxisItems(data *[]db.DBJobLog) []string {
	items := make([]string, len(*data))
	for i, v := range *data {
		items[i] = v.TS
	}
	return items
}

func getPieItems(data *[]db.DBPieJob) []opts.PieData {
	items := make([]opts.PieData, len(*data))
	var idx int
	for i, v := range *data {
		items[i] = opts.PieData{Name: strings.ToTitle(v.Status), Value: v.Count}
	}
	return items
}

func getLineChart(data *[]db.DBJobLog) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithAnimation(),
		charts.WithTitleOpts(opts.Title{Title: "Ping Monitoring"}),
		charts.WithInitializationOpts(opts.Initialization{
			PageTitle: "Monitoring",
			Width:     "800px",
			Height:    "400px",
		}),
		charts.WithColorsOpts(opts.Colors{NORD10}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:  "slider",
			Start: 0,
			End:   50,
		}),
	)
	line.SetXAxis(getLineXAxisItems(data)).
		AddSeries("pings", getLineItems(data)).
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{
				ShowSymbol: true,
			}),
			charts.WithLabelOpts(opts.Label{
				Show:      true,
				Position:  "inside",
				Formatter: "{b}",
			}),
		)
	return line
}

func getPieChart(data *[]db.DBPieJob) *charts.Pie {
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Job statistics"}),
		charts.WithColorsOpts(opts.Colors{NORD14, NORD11}),
		charts.WithInitializationOpts(opts.Initialization{
			PageTitle: "Monitoring",
			Width:     "400px",
			Height:    "400px",
		}),
	)
	pie.AddSeries("pie", getPieItems(data)).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:      true,
				Formatter: "{b}: {c}",
				Color:     NORD0,
			}),
			charts.WithPieChartOpts(opts.PieChart{
				Radius: []string{"40%", "75%"},
			}),
		)
	return pie
}
