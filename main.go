package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/go-kit/log"
	"github.com/oklog/ulid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/thanos-io/objstore/client"
	"github.com/thanos-io/thanos/pkg/block/indexheader"
	"github.com/thanos-io/thanos/pkg/store"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	objstoreConfigFile = kingpin.Flag("objstore.config-file", "Path to YAML file that contains object store configuration. See format details: https://thanos.io/tip/thanos/storage.md/#configuration").Required().File()
	datadir            = kingpin.Flag("data-dir", "Data dir for storing downloaded from storage data").Required().String()
	bucketulid         = kingpin.Flag("bucket-ulid", "Store bucket ULID as a string").Required().Short('b').String()
	labelname          = kingpin.Flag("label-name", "If provided, all values for given label name will be printed, otherwise will print all label names").Short('l').String()
	showStat           = kingpin.Flag("show-stat", "Gather stat such as labels value count and size of all values in bytes").Default("true").Short('s').Bool()
)

func main() {
	kingpin.Version("0.0.2")
	kingpin.Parse()

	ctx := context.Background()
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	bktId, err := ulid.Parse(*bucketulid)
	if err != nil {
		panic("bad bucket id")
	}
	fmt.Println("got bucket id: ", bktId.String())

	objstoreYaml, err := io.ReadAll(*objstoreConfigFile)
	if err != nil {
		panic("can't read objstore config")
	}

	bkt, err := client.NewBucket(
		logger,
		objstoreYaml,
		prometheus.NewRegistry(),
		"",
	)
	if err != nil {
		panic("can't create bucket")
	}

	reader, err := indexheader.NewBinaryReader(
		ctx,
		logger,
		bkt,
		*datadir,
		bktId,
		store.DefaultPostingOffsetInMemorySampling,
	)

	labelnames, err := reader.LabelNames()
	if err != nil {
		panic("can't read label names")
	}

	if *showStat {
		type labelStat struct {
			size      int
			valuesnum int
		}

		allStat := make(map[string]labelStat)

		for _, l := range labelnames {
			lstat := labelStat{}
			values, err := reader.LabelValues(l)
			if err != nil {
				panic("can't get label values")
			}
			for _, v := range values {
				lstat.valuesnum++
				lstat.size += len([]byte(v))
			}
			allStat[l] = lstat
		}
		for l, ls := range allStat {
			fmt.Printf("%s\t%d\t%d\n", l, ls.valuesnum, ls.size)
		}
	}

	if *labelname == "" {
		fmt.Println("No label name provided. Printing all label names:")
		for _, ln := range labelnames {
			fmt.Println(ln)
		}

		os.Exit(0)
	}

	labelValues, err := reader.LabelValues(*labelname)
	if err != nil {
		panic("can't get label values")
	}

	fmt.Printf("Printing values for label: %s\n", *labelname)
	for _, lv := range labelValues {
		fmt.Println(lv)
	}
}
