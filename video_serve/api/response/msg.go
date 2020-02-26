package response

type Resp struct {
	ResponseCode int    `json:"response_code" form:"response_code"`
	ResponseMsg  string `json:"response_msg" form:"response_msg"`
}

var (
	IntervalErr = Resp{
		ResponseCode: 500,
		ResponseMsg:  "Server Parse Failed",
	}

	UserNotFound = Resp{
		ResponseCode: 404,
		ResponseMsg:  "The User Not Found",
	}

	UserNameLimits = Resp{
		ResponseCode: 400,
		ResponseMsg:  "User Name Not be Smaller 3 or More 12 ",
	}

	Success = Resp{
		ResponseCode: 200,
		ResponseMsg:  "Successfully",
	}

	ExpireSession = Resp{
		ResponseCode: 400,
		ResponseMsg:  "Session is expire",
	}

	RequestInvalid = Resp{
		ResponseCode: 400,
		ResponseMsg:  "Request Date is invalid",
	}

	NotOpenEls = Resp{
		ResponseCode: 500,
		ResponseMsg:  "You Want use Elastic serve . Please Open Elastic and Configure it",
	}
	NotOpenOss = Resp{
		ResponseCode: 500,
		ResponseMsg:  "You Want use Oss serve . Please Open Oss and Configure it",
	}
)
