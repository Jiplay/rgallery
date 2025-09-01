package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/robbymilo/rgallery/pkg/queries"
	"github.com/robbymilo/rgallery/pkg/types"
)

type Conf = types.Conf
type metricsCollector struct {
	totalItems     *prometheus.Desc
	totalFolders   *prometheus.Desc
	totalFavorites *prometheus.Desc
	totalTags      *prometheus.Desc
	c              Conf
}

func MetricsCollector(c Conf) *metricsCollector {
	return &metricsCollector{
		totalItems: prometheus.NewDesc("rgallery_total_items",
			"Total items in database", nil, nil,
		),
		totalFolders: prometheus.NewDesc("rgallery_total_folders",
			"Total folders in database", nil, nil,
		),
		totalFavorites: prometheus.NewDesc("rgallery_total_favorites",
			"Total favorites in database", nil, nil,
		),
		totalTags: prometheus.NewDesc("rgallery_total_tags",
			"Total tags in database", nil, nil,
		),
		c: c,
	}
}

func (collector *metricsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.totalItems
}

func (collector *metricsCollector) Collect(ch chan<- prometheus.Metric) {

	from := time.Unix(0, 0)
	to := time.Now()
	totalItems, err := queries.GetTotalMediaItems(0, from.Format(time.RFC3339), to.Format(time.RFC3339), "", "", collector.c)
	if err != nil {
		collector.c.Logger.Error("error getting total media items", "error", err)
	}
	totalFolders, err := queries.GetTotalFolders(collector.c)
	if err != nil {
		collector.c.Logger.Error("error getting total folders", "error", err)
	}
	totalFavorites, err := queries.GetTotalFavorites(5, collector.c)
	if err != nil {
		collector.c.Logger.Error("error getting total favorites", "error", err)
	}
	totalTags, err := queries.GetTotalTags(collector.c)
	if err != nil {
		collector.c.Logger.Error("error getting total favorites", "error", err)
	}

	items := prometheus.MustNewConstMetric(collector.totalItems, prometheus.GaugeValue, float64(totalItems))
	items = prometheus.NewMetricWithTimestamp(time.Now(), items)
	ch <- items

	folders := prometheus.MustNewConstMetric(collector.totalFolders, prometheus.GaugeValue, float64(totalFolders))
	folders = prometheus.NewMetricWithTimestamp(time.Now(), folders)
	ch <- folders

	favorites := prometheus.MustNewConstMetric(collector.totalFavorites, prometheus.GaugeValue, float64(totalFavorites))
	favorites = prometheus.NewMetricWithTimestamp(time.Now(), favorites)
	ch <- favorites

	tags := prometheus.MustNewConstMetric(collector.totalTags, prometheus.GaugeValue, float64(totalTags))
	tags = prometheus.NewMetricWithTimestamp(time.Now(), tags)
	ch <- tags
}
