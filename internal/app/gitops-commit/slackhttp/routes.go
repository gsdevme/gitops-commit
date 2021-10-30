package slackhttp

func (s *server) routes() {

	s.router.HandleFunc("/ping", s.handleHealthCheck())
	s.router.HandleFunc("/healthz", s.handleHealthCheck())

	s.router.HandleFunc(
		"/slack_command",
		s.SlackCommandMiddleware(s.handleSlackCommand(s.registry)),
	).Methods("POST")

	s.router.HandleFunc("/", s.handleNotFound())
}
