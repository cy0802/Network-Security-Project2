package processor

import (
	"net/http"

	"github.com/free5gc/openapi/models"
	"github.com/free5gc/scp/internal/logger"
)

// NOTE: The response from AMF is guaranteed to be correct
func (p *Processor) PostUeAutentications(
	authInfo models.AuthenticationInfo,
) *HandlerResponse {
	logger.ProxyLog.Debugln("[AMF->AUSF] Forward AMF UE Authentication Request")

	// TODO: Send request to target NF by setting correct uri
	var targetNfUri string
	targetNfUri = "http://10.100.200.9:8000"

	// TODO: Verify that the Information Elements (IEs) in the response body are correct
	//       Recover and handle errors if the IEs are incorrect
	response, problemDetails, err := p.Consumer().SendUeAuthPostRequest(targetNfUri, &authInfo)
	logger.ProxyLog.Debugf("response: %#v, problemDetails: %#v, err: %s", response, problemDetails, err)

	if response != nil {
		return &HandlerResponse{http.StatusCreated, nil, response}
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

func (p *Processor) PutUeAutenticationsConfirmation(
	authCtxId string,
	confirmationData models.ConfirmationData,
) *HandlerResponse {
	logger.ProxyLog.Debugln("[AMF->AUSF] Forward AMF UE Authentication Response")

	// TODO: Send request to target NF by setting correct uri
	var targetNfUri string
	targetNfUri = "http://10.100.200.9:8000"

	// TODO: Verify that the Information Elements (IEs) in the response body are correct
	//       Recover and handle errors if the IEs are incorrect
	response, problemDetails, err := p.Consumer().SendAuth5gAkaConfirmRequest(targetNfUri, authCtxId, &confirmationData)
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
