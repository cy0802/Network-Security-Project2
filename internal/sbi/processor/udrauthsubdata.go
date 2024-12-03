package processor

import (
	"net/http"

	"github.com/free5gc/openapi/models"
	"github.com/free5gc/scp/internal/logger"
)

// NOTE: The response from UDR is guaranteed to be correct
func (p *Processor) GetAuthSubsData(
	ueId string,
) *HandlerResponse {
	logger.DetectorLog.Debugln("[UDM->UDR] Forward UDM Authentication Data Query Request")

	// TODO: Send request to correct NF by setting correct uri
	var targetNfUri string
	targetNfUri = "http://10.100.200.4:8000"

	// TODO: Store UE auth subscription data
	response, problemDetails, err := p.Consumer().SendAuthSubsDataGet(targetNfUri, ueId)
	// logger.ProxyLog.Debugf("response: %#v, problemDetails: %#v, err: %s", response, problemDetails, err)

	if response != nil {
		return &HandlerResponse{http.StatusOK, nil, response}
	} else if problemDetails != nil {
		return &HandlerResponse{int(problemDetails.Status), nil, problemDetails}
	}
	logger.ProxyLog.Errorln(err)
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}

	return &HandlerResponse{http.StatusForbidden, nil, problemDetails}
}
