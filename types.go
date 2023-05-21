package adidas_backend

import (
	"umbrella/internal/modules"
	globaltypes "umbrella/internal/utils/global_types"
)

const (
	PAYMENT_COD         = "CASH_ON_DELIVERY"
	PAYMENT_PAYPAL      = "PAYPAL"
	PAYMENT_CREDIT_CARD = "CREDIT_CARD"
)

const (
	ADYEN_PUBLIC = `10001|B0C2259F5CB5FECDB8F4010E526520723BDDFC6133019DF20F24E84CC199AA226663436BE449EA4FF9E058F21ED13F4F1F2BC34236AAA1171EA5989D7B486DCB147521C970575A83D9C395BCB896166A0BCE6D55C1414A13C81851306B84F513ED179F41E93C69027D83BF09DAABA9C3451AA8C8523F97A7741439B36573B3100E0FFAB08CBEA2785F3D2D2717073D6F3E243DD27BAEEB2F917502BC460B9D7F36CE6650F8CEA047C3F803C7CEC8F141B1C9194C83B8EDECB20B15614CC3A6FC773B48667C82BA726C474BCB9AB6B75D2092CE4E423A604EA3DB9E2223CE8C966FF76532F8AF37308098AE55F644DCD869EE475462E8E0B11F88B5345E200C6D`
)

var PAYMENT_CHECK = map[string]string{
	"pp":     PAYMENT_PAYPAL,
	"paypal": PAYMENT_PAYPAL,
	"cod":    PAYMENT_COD,
}

type TaskStates struct {
	Login          globaltypes.TaskState
	SessionCheck   globaltypes.TaskState
	Payment        globaltypes.TaskState
	Delivery       globaltypes.TaskState
	AddressGet     globaltypes.TaskState
	Address        globaltypes.TaskState
	CheckoutID     globaltypes.TaskState
	Discount       globaltypes.TaskState
	Monitor        globaltypes.TaskState
	Order          globaltypes.TaskState
	Akamai         globaltypes.TaskState
	Cancel         globaltypes.TaskState
	Register       globaltypes.TaskState
	SessionRefresh globaltypes.TaskState
}

type SessionTokens struct {
	AccessToken  string
	RefreshToken string
}

type RequiredRequests struct {
	Shipping bool
	Payment  bool
}

type AsyncRequestsChannels struct {
	Discount chan error
}

type Resources struct {
	SensorData            string
	SessionTokens         SessionTokens
	CheckoutID            string
	AddressId             string
	RequiredRequests      RequiredRequests
	AsyncRequestsChannels AsyncRequestsChannels
}

type Cosmetics struct {
	Price string
	Image string
	Name  string
	Size  string
}

type Config struct {
	DefaultConfig modules.DefaultConfig
	TaskStates    TaskStates
	Resources     Resources
	UsedDiscount  string
	IfDiscounted  bool
	Cosmetics     Cosmetics
}

type Variant struct {
	Pid     string
	SizePid string
	Size    string
}
