package exception

/**
 * success
 */
const RESULT_OK="00000000"

/**
 * invaild param error
 */
const RESULT_INVAILD_PARAM_ERROR="00000001"

/**
 * analyse request error
 */
const RESULT_ANALYSE_REQUEST_ERROR="00000002"

/**
 * system is busy error
 */
const RESULT_SYSTEM_IS_BUSY_ERROR="00000003"

/**
 * query db error
 */
const RESULT_SQL_ERROR="00000004"

/**
 * query redis error
 */
const RESULT_CACHE_ERROR="00000005"

/**
 * file error
 */
const RESULT_FILE_ERROR="00000006"

/**
 * thread task error
 */
const RESULT_TASK_ERROR="00000007"

/**
 * system inter error
 */
const RESULT_SYSTEM_INTER_ERROR="00000008"

/**
 * third interface error
 */
const RESULT_THIRD_INTERFACE_ERROR="00000009"

/**
 * auth error
 */
const RESULT_ACCESS_ERROR="00000010"

/**
 * reuqest info format error
 */
const RESULT_INVALID_REQUEST_INFO_ERROR="00000011"

/**
 * invaild user token error
 */
const RESULT_INVAILD_USER_TOKEN_ERROR="00000012"

var CommonErrorInfoMap map[string]string = nil

func AddCommonErrorInfo(code string, description string) {
	if CommonErrorInfoMap == nil {
		CommonErrorInfoMap = make(map[string]string)
	}

	CommonErrorInfoMap[code]=description
}

func GetCommonErrorDecription(code string) string {
	return CommonErrorInfoMap[code]
}

