package adidas_backend

type AkamaiNeedsSolveError struct{}

func (a AkamaiNeedsSolveError) Error() string {
	return "akamai needs to be solved"
}

type WaitAfterRestock struct{}

func (w WaitAfterRestock) Error() string {
	return "waiting after restock"
}

type WaitAfterError struct{}

func (w WaitAfterError) Error() string {
	return "waiting after error"
}

type RefreshSessionError struct{}

func (r RefreshSessionError) Error() string {
	return "session needs to be refreshed"
}
