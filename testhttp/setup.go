package testhttp

func Build(h httpHandler) TestHTTP {
	return TestHTTP{
		handler: h,
	}
}
