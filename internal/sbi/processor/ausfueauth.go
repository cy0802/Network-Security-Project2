package processor

import (
    "fmt"
    "net/http"
    

	"github.com/free5gc/openapi/models"
	"github.com/free5gc/scp/internal/logger"
)

// NOTE: The response from AMF is guaranteed to be correct
func (p *Processor) PostUeAutentications(
	authInfo models.AuthenticationInfo,
) *HandlerResponse {
	logger.ProxyLog.Debugln("[AMF->AUSF] Forward AMF UE Authentication Request")

    // 驗證 Information Elements
    if err := validateAuthInfo(authInfo); err != nil {
        problemDetails := &models.ProblemDetails{
            Status: http.StatusBadRequest,
            Cause:  "INVALID_REQUEST",
            Detail: err.Error(),
        }
        return &HandlerResponse{http.StatusBadRequest, nil, problemDetails}
    }

	// TODO: Send request to target NF by setting correct uri
	var targetNfUri string
    targetNfUri = "http://10.100.200.9:8000"

	// TODO: Verify that the Information Elements (IEs) in the response body are correct
	//       Recover and handle errors if the IEs are incorrect
	response, problemDetails, err := p.Consumer().SendUeAuthPostRequest(targetNfUri, &authInfo)

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

    // 驗證 authCtxId
    if authCtxId == "" {
        logger.DetectorLog.Errorln("models.ConfirmationData.AuthCtxId: Mandatory type is absent")
        problemDetails := &models.ProblemDetails{
            Status: http.StatusBadRequest,
            Cause:  "INVALID_REQUEST",
            Detail: "missing authCtxId",
        }
        return &HandlerResponse{http.StatusBadRequest, nil, problemDetails}
    }

    // 驗證 ConfirmationData
    if err := validateConfirmationData(confirmationData); err != nil {
        problemDetails := &models.ProblemDetails{
            Status: http.StatusBadRequest,
            Cause:  "INVALID_REQUEST",
            Detail: err.Error(),
        }
        return &HandlerResponse{http.StatusBadRequest, nil, problemDetails}
    }

	// TODO: Send request to target NF by setting correct uri
	var targetNfUri string
	targetNfUri = "http://10.100.200.9:8000"

	// TODO: Verify that the Information Elements (IEs) in the response body are correct
	//       Recover and handle errors if the IEs are incorrect
	
	response, problemDetails, err := p.Consumer().SendAuth5gAkaConfirmRequest(targetNfUri, authCtxId, &confirmationData)

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

// 驗證 AuthenticationInfo
func validateAuthInfo(authInfo models.AuthenticationInfo) error {
    // 檢查必要欄位
    if authInfo.SupiOrSuci == "" {
        logger.DetectorLog.Errorln("models.AuthenticationInfo.SupiOrSuci: Mandatory type is absent")
        return fmt.Errorf("missing mandatory field: supiOrSuci")
    }

    if authInfo.ServingNetworkName == "" {
        logger.DetectorLog.Errorln("models.AuthenticationInfo.ServingNetworkName: Mandatory type is absent")
        return fmt.Errorf("missing mandatory field: servingNetworkName")
    }

    // 驗證欄位值
    /*if !strings.HasPrefix(authInfo.SupiOrSuci, "5G:") {
       logger.DetectorLog.Errorln("models.AuthenticationInfo.SupiOrSuci: Unexpected value is received")
        return fmt.Errorf("invalid supiOrSuci format")
    }*/

    // 檢查條件性欄位
    if authInfo.ResynchronizationInfo != nil && authInfo.ResynchronizationInfo.Rand == "" {
        logger.DetectorLog.Errorln("models.AuthenticationInfo.ResynchronizationInfo.Rand: Miss condition")
        return fmt.Errorf("rand is required when resynchronizationInfo is present")
    }

    return nil
}

// 驗證 ConfirmationData
func validateConfirmationData(data models.ConfirmationData) error {
    // 檢查必要欄位
    if data.ResStar == "" {
        logger.DetectorLog.Errorln("models.ConfirmationData.ResStar: Mandatory type is absent")
        return fmt.Errorf("missing mandatory field: resStar")
    }

    // 驗證欄位值
    /*if len(data.ResStar) != 64 {
        logger.DetectorLog.Errorln("models.ConfirmationData.ResStar: Unexpected value is received")
        return fmt.Errorf("invalid resStar length: expected 64 characters")
    }*/

    return nil
}
