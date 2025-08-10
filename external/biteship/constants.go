package biteship

const (
	ErrInvalidParsingResponse = 50009001
)

const (
	ErrInvalidAuthentication = 40101001

	ErrRateInvalidPostalCode  = 40001001
	ErrRateInvalidParameter   = 40001002
	ErrRateNoCourierAvailable = 40001010
)

// ErrorMessages Map of error codes to response messages ErrRateNoCourierAvailable map to "No courier available for the given postal code or area"
var ErrorMessages = map[int]string{
	ErrInvalidParsingResponse: "Failed to parse response from Provider",
	ErrInvalidAuthentication:  "Invalid third party authentication credentials",
	ErrRateInvalidPostalCode:  "Invalid postal code provided",
	ErrRateInvalidParameter:   "Invalid parameter provided",
	ErrRateNoCourierAvailable: "No courier available for the given postal code or area",
}

// ErrorCodeToHTTPStatus Map ErrorCode to API Response HTTP Code
var ErrorCodeToHTTPStatus = map[int]int{
	ErrInvalidParsingResponse: 500, // Internal Server Error
	ErrInvalidAuthentication:  401, // Unauthorized
	ErrRateInvalidPostalCode:  400, // Bad Request
	ErrRateInvalidParameter:   400, // Bad Request
	ErrRateNoCourierAvailable: 404, // Not Found
}

type StatusTracking string

const (
	StatusTrackingConfirmed       StatusTracking = "confirmed"
	StatusTrackingAllocated       StatusTracking = "allocated"
	StatusTrackingPickingUp       StatusTracking = "pickingUp"
	StatusTrackingPicked          StatusTracking = "picked"
	StatusTrackingDroppingOff     StatusTracking = "droppingOff"
	StatusTrackingInTransit       StatusTracking = "inTransit"
	StatusTrackingOnHold          StatusTracking = "onHold"
	StatusTrackingDelivered       StatusTracking = "delivered"
	StatusTrackingRejected        StatusTracking = "rejected"
	StatusTrackingCourierNotFound StatusTracking = "courierNotFound"
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
