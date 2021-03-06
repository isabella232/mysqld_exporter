package collector

import (
	"context"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/smartystreets/goconvey/convey"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestScrapeClientStat(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error opening a stub database connection: %s", err)
	}
	defer db.Close()

	mock.ExpectQuery(sanitizeQuery(userstatCheckQuery)).WillReturnRows(sqlmock.NewRows([]string{"Variable_name", "Value"}).
		AddRow("userstat", "ON"))

	columns := []string{"CLIENT", "TOTAL_CONNECTIONS", "CONCURRENT_CONNECTIONS", "CONNECTED_TIME", "BUSY_TIME", "CPU_TIME", "BYTES_RECEIVED", "BYTES_SENT", "BINLOG_BYTES_WRITTEN", "ROWS_READ", "ROWS_SENT", "ROWS_DELETED", "ROWS_INSERTED", "ROWS_UPDATED", "SELECT_COMMANDS", "UPDATE_COMMANDS", "OTHER_COMMANDS", "COMMIT_TRANSACTIONS", "ROLLBACK_TRANSACTIONS", "DENIED_CONNECTIONS", "LOST_CONNECTIONS", "ACCESS_DENIED", "EMPTY_QUERIES"}
	rows := sqlmock.NewRows(columns).
		AddRow("localhost", 1002, 0, 127027, 286, 245, float64(2565104853), 21090856, float64(2380108042), 767691, 1764, 8778, 1210741, 0, 1764, 1214416, 293, 2430888, 0, 0, 0, 0, 0)
	mock.ExpectQuery(sanitizeQuery(clientStatQuery)).WillReturnRows(rows)

	ch := make(chan prometheus.Metric)
	go func() {
		if err = (ScrapeClientStat{}).Scrape(context.Background(), db, ch); err != nil {
			t.Errorf("error calling function on test: %s", err)
		}
		close(ch)
	}()

	expected := []MetricResult{
		{labels: labelMap{"client": "localhost"}, value: 1002, metricType: dto.MetricType_COUNTER},
		{labels: labelMap{"client": "localhost"}, value: 0, metricType: dto.MetricType_GAUGE},
		{labels: labelMap{"client": "localhost"}, value: 127027, metricType: dto.MetricType_COUNTER},
		{labels: labelMap{"client": "localhost"}, value: 286, metricType: dto.MetricType_COUNTER},
		{labels: labelMap{"client": "localhost"}, value: 245, metricType: dto.MetricType_COUNTER},
		{labels: labelMap{"client": "localhost"}, value: float64(2565104853), metricType: dto.MetricType_COUNTER},
		{labels: labelMap{"client": "localhost"}, value: 21090856, metricType: dto.MetricType_COUNTER},
		{labels: labelMap{"client": "localhost"}, value: float64(2380108042), metricType: dto.MetricType_COUNTER},
		{labels: labelMap{"client": "localhost"}, value: 767691, metricType: dto.MetricType_COUNTER},
		{labels: labelMap{"client": "localhost"}, value: 1764, metricType: dto.MetricType_COUNTER},
		{labels: labelMap{"client": "localhost"}, value: 8778, metricType: dto.MetricType_COUNTER},
		{labels: labelMap{"client": "localhost"}, value: 1210741, metricType: dto.MetricType_COUNTER},
		{labels: labelMap{"client": "localhost"}, value: 0, metricType: dto.MetricType_COUNTER},
		{labels: labelMap{"client": "localhost"}, value: 1764, metricType: dto.MetricType_COUNTER},
		{labels: labelMap{"client": "localhost"}, value: 1214416, metricType: dto.MetricType_COUNTER},
		{labels: labelMap{"client": "localhost"}, value: 293, metricType: dto.MetricType_COUNTER},
		{labels: labelMap{"client": "localhost"}, value: 2430888, metricType: dto.MetricType_COUNTER},
		{labels: labelMap{"client": "localhost"}, value: 0, metricType: dto.MetricType_COUNTER},
		{labels: labelMap{"client": "localhost"}, value: 0, metricType: dto.MetricType_COUNTER},
		{labels: labelMap{"client": "localhost"}, value: 0, metricType: dto.MetricType_COUNTER},
		{labels: labelMap{"client": "localhost"}, value: 0, metricType: dto.MetricType_COUNTER},
		{labels: labelMap{"client": "localhost"}, value: 0, metricType: dto.MetricType_COUNTER},
	}
	convey.Convey("Metrics comparison", t, func() {
		for _, expect := range expected {
			got := readMetric(<-ch)
			convey.So(expect, convey.ShouldResemble, got)
		}
	})

	// Ensure all SQL queries were executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}
