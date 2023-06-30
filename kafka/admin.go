package kafka

import (
	"github.com/Shopify/sarama"
)

// Admin contructs and returns an underlying Cluster Admin struct.
func (kc *KClient) Admin() sarama.ClusterAdmin {
	ca, err := sarama.NewClusterAdminFromClient(kc.cl)
	if err != nil {
		panic(err)
	}
	return ca
}
