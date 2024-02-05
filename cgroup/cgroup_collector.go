package cgroup

import (
	"bufio"
	"github.com/prometheus/client_golang/prometheus"
	"os"
	"strconv"
	"strings"
)

type cgroupMetric struct {
	Metric *prometheus.Desc
	Path   string
}

type statMetric struct {
	Metric     *prometheus.Desc
	NameInStat string
}

type cgroupMemoryCollector struct {
	metrics     []cgroupMetric
	statMetrics []statMetric
}

func NewCGroupMemoryCollector() prometheus.Collector {
	// Initialize all desired metrics with their descriptions
	metrics := []cgroupMetric{
		makeCgroupMetric("usage_in_bytes", "Memory usage of the cgroup in bytes."),
		makeCgroupMetric("max_usage_in_bytes", "Maximum memory usage of the cgroup in bytes."),
		makeCgroupMetric("limit_in_bytes", "Memory limit of the cgroup in bytes."),
		makeCgroupMetric("failcnt", "Number of memory usage limit hits."),
		makeCgroupMetric("kmem.usage_in_bytes", "Kernel memory usage of the cgroup in bytes."),
		makeCgroupMetric("kmem.limit_in_bytes", "Kernel memory limit of the cgroup in bytes."),
		makeCgroupMetric("kmem.failcnt", "Kernel number of memory usage hits limits."),
		makeCgroupMetric("kmem.tcp.usage_in_bytes", "TCP buffer memory usage in bytes"),
		makeCgroupMetric("kmem.tcp.limit_in_bytes", "TCP buffer memory limit in bytes"),
		makeCgroupMetric("kmem.tcp.max_usage_in_bytes", "TCP buffer memory max usage in bytes"),
		makeCgroupMetric("kmem.tcp.failcnt", "TCP buffer memory usage limit hits"),
	}
	statMetrics := []statMetric{
		makeStatMetric("rss"),
		makeStatMetric("cache"),
		makeStatMetric("shmem"),
		makeStatMetric("mapped_file"),
		makeStatMetric("total_rss"),
		makeStatMetric("total_cache"),
		makeStatMetric("total_shmem"),
		makeStatMetric("total_mapped_file"),
	}

	return &cgroupMemoryCollector{metrics: metrics, statMetrics: statMetrics}
}

func makeCgroupMetric(dotpath string, help string) cgroupMetric {
	return cgroupMetric{
		Path:   "/sys/fs/cgroup/memory/memory." + dotpath,
		Metric: prometheus.NewDesc("cgroup_"+strings.ReplaceAll(dotpath, ".", "_"), help, nil, nil),
	}
}

func makeStatMetric(nameInStat string) statMetric {
	return statMetric{
		NameInStat: nameInStat, // Needs to match the name in the stat file
		Metric:     prometheus.NewDesc("cgroup_memory_stat_"+nameInStat, "", nil, nil),
	}
}

func (collector *cgroupMemoryCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metricDesc := range collector.metrics {
		ch <- metricDesc.Metric
	}
}

func (collector *cgroupMemoryCollector) Collect(ch chan<- prometheus.Metric) {

	for _, metric := range collector.metrics {
		if value, err := readMetricValue(metric.Path); err == nil {
			ch <- prometheus.MustNewConstMetric(metric.Metric, prometheus.GaugeValue, value)
		}
	}

	// Handle multi-value metrics like those in memory.stat
	collectMemoryStatMetrics(collector.statMetrics, ch)
}

func collectMemoryStatMetrics(statMetrics []statMetric, ch chan<- prometheus.Metric) {
	file, err := os.Open("/sys/fs/cgroup/memory/memory.stat")
	if err != nil {
		return
	}
	defer file.Close()
	metricsMap := make(map[string]statMetric)
	for _, statMetric := range statMetrics {
		metricsMap[statMetric.NameInStat] = statMetric
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) == 2 {
			nameInStat := parts[0]
			if statMetric, exists := metricsMap[nameInStat]; exists {
				value, err := strconv.ParseFloat(parts[1], 64)
				if err == nil {
					ch <- prometheus.MustNewConstMetric(statMetric.Metric, prometheus.GaugeValue, value)
				}
			}
		}
	}
}
