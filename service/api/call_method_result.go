package api

type CallMethodResult struct {
	StatusCode           `json:"StatusCode"`
	InputArgumentResults []StatusCode  `json:"InputArgumentResults"`
	OutputArguments      []interface{} `json:"OutputArguments"`
}
