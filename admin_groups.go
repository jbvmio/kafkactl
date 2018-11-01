package kafkactl

import (
	"fmt"

	"github.com/Shopify/sarama"
)

func (kc *KClient) RemoveGroup(name string) (errMsg error) {
	var brokerErr error
	var found bool
	var code int16
	var returnCode int16
	request := sarama.DeleteGroupsRequest{
		Groups: []string{name},
	}
	for _, broker := range kc.brokers {
		response, errd := broker.DeleteGroups(&request)
		if errd == nil {
			for _, v := range response.GroupErrorCodes {
				var err error
				if found, returnCode, err = catchedErrCode(v); !found {
					return nil
				}
				if returnCode > code {
					code = returnCode
					errMsg = err
				}
			}
		} else {
			brokerErr = errd
		}
	}
	if brokerErr != nil {
		errMsg = brokerErr
	}
	return
}

func catchedErrCode(kerror sarama.KError) (bool, int16, error) {
	code := int16(kerror)
	if code > 0 {
		if code == 68 {
			return true, code, fmt.Errorf("Error Code %v Returned: NON_EMPTY_GROUP", code)
		}
		if code == 69 {
			return true, code, fmt.Errorf("Error Code %v Returned: GROUP_ID_NOT_FOUND", code)
		}
		return true, code, fmt.Errorf("%v", kerror.Error())
	}
	return false, code, nil
}
