package cmd

// PartitionLag struct def:
type PartitionLag struct {
	Group     string
	Topic     string
	Partition int32
	Member    string
	Offset    int64
	Lag       int64
}

//tbl = table.New("GROUP", "TOPIC", "PART", "MEMBER", "LAG")
