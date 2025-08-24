package biteship

type BiteshipErrCode int

const (
	ErrInvalidParsingResponse BiteshipErrCode = 50009001
)

const (
	ErrInvalidAuthentication BiteshipErrCode = 40101001

	ErrRateInvalidPostalCode  BiteshipErrCode = 40001001
	ErrRateInvalidParameter   BiteshipErrCode = 40001002
	ErrRateNoCourierAvailable BiteshipErrCode = 40001010

	ErrTrackingNotFound BiteshipErrCode = 40003001
)

// ErrorMessages Map of error codes to response messages ErrRateNoCourierAvailable map to "No courier available for the given postal code or area"
var ErrorMessages = map[BiteshipErrCode]string{
	ErrInvalidParsingResponse: "Failed to parse response from Provider",
	ErrInvalidAuthentication:  "Invalid third party authentication credentials",
	ErrRateInvalidPostalCode:  "Invalid postal code provided",
	ErrRateInvalidParameter:   "Invalid parameter provided",
	ErrRateNoCourierAvailable: "No courier available for the given postal code or area",
	ErrTrackingNotFound:       "Tracking information not found for the given waybill or courier",
}

// ErrorCodeToHTTPStatus Map ErrorCode to API Response HTTP Code
var ErrorCodeToHTTPStatus = map[BiteshipErrCode]int{
	ErrInvalidParsingResponse: 500, // Internal Server Error
	ErrInvalidAuthentication:  401, // Unauthorized
	ErrRateInvalidPostalCode:  400, // Bad Request
	ErrRateInvalidParameter:   400, // Bad Request
	ErrRateNoCourierAvailable: 404, // Not Found
	ErrTrackingNotFound:       404, // Not Found
}

func (code BiteshipErrCode) GetMessage() string {
	if msg, exists := ErrorMessages[code]; exists {
		return msg
	}
	return "Unknown error"
}

func (code BiteshipErrCode) GetHTTPStatusCode() int {
	if status, exists := ErrorCodeToHTTPStatus[code]; exists {
		return status
	}
	return 500 // Default to Internal Server Error if code is not found
}

type StatusTracking string

const (
	StatusTrackingConfirmed       StatusTracking = "confirmed"
	StatusTrackingAllocated       StatusTracking = "allocated"
	StatusTrackingPickingUp       StatusTracking = "picking_up"
	StatusTrackingPicked          StatusTracking = "picked"
	StatusTrackingDroppingOff     StatusTracking = "dropping_off"
	StatusTrackingInTransit       StatusTracking = "in_transit"
	StatusTrackingOnHold          StatusTracking = "on_hold"
	StatusTrackingDelivered       StatusTracking = "delivered"
	StatusTrackingRejected        StatusTracking = "rejected"
	StatusTrackingCourierNotFound StatusTracking = "courier_not_found"
	StatusTrackingReturned        StatusTracking = "returned"
	StatusTrackingCancelled       StatusTracking = "cancelled"
	StatusTrackingDisposed        StatusTracking = "disposed"
)

var TrackingStatusMessageMap = map[StatusTracking]string{
	StatusTrackingConfirmed:       "Pesanan telah dikonfirmasi. Segera mencari kurir terdekat.",
	StatusTrackingAllocated:       "Kurir telah dialokasikan untuk mengambil pesanan Anda.",
	StatusTrackingPickingUp:       "Kurir sedang dalam perjalanan untuk mengambil pesanan Anda.",
	StatusTrackingPicked:          "Pesanan Anda telah diambil oleh kurir.",
	StatusTrackingDroppingOff:     "Kurir sedang dalam perjalanan untuk mengantarkan pesanan Anda",
	StatusTrackingInTransit:       "Pesanan Anda sedang diproses di lokasi transit.",
	StatusTrackingOnHold:          "Pesanan Anda sedang ditahan karena alasan tertentu.",
	StatusTrackingDelivered:       "Pesanan Anda telah berhasil diantarkan.",
	StatusTrackingRejected:        "Pesanan Anda telah ditolak oleh kurir.",
	StatusTrackingCourierNotFound: "Kurir tidak ditemukan untuk pesanan ini.",
	StatusTrackingReturned:        "Pesanan Anda telah dikembalikan.",
	StatusTrackingCancelled:       "Pesanan Anda telah dibatalkan.",
	StatusTrackingDisposed:        "Pesanan Anda telah dibuang.",
}

// GetMessage returns the message for a given tracking status
func (s StatusTracking) GetMessage() string {
	if msg, exists := TrackingStatusMessageMap[s]; exists {
		return msg
	}
	return "Status tidak dikenal"
}

// ToStatusTracking converts a string to StatusTracking
func ToStatusTracking(status string) StatusTracking {
	if status == "" {
		return StatusTracking("")
	}
	if st, exists := TrackingStatusMessageMap[StatusTracking(status)]; exists {
		return StatusTracking(st)
	}
	return StatusTracking(status)
}

// ToString converts StatusTracking to string
func (s StatusTracking) ToString() string {
	if s == "" {
		return ""
	}
	return string(s)
}

// CodeToBiteshipErrCode converts an error code to BiteshipErrCode
func CodeToBiteshipErrCode(code int) BiteshipErrCode {
	if code < 1000 {
		return BiteshipErrCode(code)
	}

	for errCode, errMsg := range ErrorMessages {
		if errMsg == "" {
			continue
		}
		if code == int(errCode) {
			return errCode
		}
	}

	return BiteshipErrCode(code)
}
