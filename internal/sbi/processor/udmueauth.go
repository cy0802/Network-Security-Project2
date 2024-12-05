package processor

import (
    "fmt"
    "net/http"
    

	"github.com/free5gc/openapi/models"
	"github.com/free5gc/scp/internal/logger"
)

// NOTE: Assume Rand from UDM and ausfInstanceId from AUSF is correct
func (p *Processor) PostGenerateAuthData(
	supiOrSuci string,
	authInfo models.AuthenticationInfoRequest,
) *HandlerResponse {
	logger.ProxyLog.Debugln("[AUSF->UDM] Forward AUSF UE Authentication Request")


    // 驗證 supiOrSuci
    if supiOrSuci == "" {
        logger.DetectorLog.Errorln("models.AuthenticationInfoRequest.SupiOrSuci: Mandatory type is absent")
        problemDetails := &models.ProblemDetails{
            Status: http.StatusBadRequest,
            Cause:  "INVALID_REQUEST",
            Detail: "missing supiOrSuci",
        }
        return &HandlerResponse{http.StatusBadRequest, nil, problemDetails}
    }

    // 驗證 AuthenticationInfoRequest
    if err := validateAuthInfoRequest(authInfo); err != nil {
        problemDetails := &models.ProblemDetails{
            Status: http.StatusBadRequest,
            Cause:  "INVALID_REQUEST",
            Detail: err.Error(),
        }
        return &HandlerResponse{http.StatusBadRequest, nil, problemDetails}
    }

	// TODO: Send request to target NF by setting correct uri
	var targetNfUri string
	targetNfUri = "http://10.100.200.3:8000"

	
	// TODO: Verify that the Information Elements (IEs) in the request or response body are correct
	//       Recover and handle errors if the IEs are incorrect
	response, problemDetails, err := p.Consumer().SendGenerateAuthDataRequest(targetNfUri, supiOrSuci, &authInfo)

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

// 驗證 AuthenticationInfoRequest
func validateAuthInfoRequest(authInfo models.AuthenticationInfoRequest) error {
    // 檢查必要欄位
    if authInfo.ServingNetworkName == "" {
        logger.DetectorLog.Errorln("models.AuthenticationInfoRequest.ServingNetworkName: Mandatory type is absent")
        return fmt.Errorf("missing mandatory field: servingNetworkName")
    }

    if authInfo.AusfInstanceId == "" {
        logger.DetectorLog.Errorln("models.AuthenticationInfoRequest.AusfInstanceId: Mandatory type is absent")
        return fmt.Errorf("missing mandatory field: ausfInstanceId")
    }

    // 驗證欄位值
    /*if !strings.HasPrefix(authInfo.AusfInstanceId, "ausf-") {
        logger.DetectorLog.Errorln("models.AuthenticationInfoRequest.AusfInstanceId: Unexpected value is received")
        return fmt.Errorf("invalid ausfInstanceId format")
    }*/
    // 檢查條件性欄位
    if authInfo.ResynchronizationInfo != nil {
        if authInfo.ResynchronizationInfo.Rand == "" {
            logger.DetectorLog.Errorln("models.AuthenticationInfoRequest.ResynchronizationInfo.Rand: Miss condition")
            return fmt.Errorf("rand is required when resynchronizationInfo is present")
        }
        if authInfo.ResynchronizationInfo.Auts == "" {
            logger.DetectorLog.Errorln("models.AuthenticationInfoRequest.ResynchronizationInfo.Auts: Miss condition") 
            return fmt.Errorf("auts is required when resynchronizationInfo is present")
        }
    }

    return nil
}
