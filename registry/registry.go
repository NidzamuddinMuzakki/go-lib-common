package registry

import (
	"bitbucket.org/moladinTech/go-lib-common/cache"
	"bitbucket.org/moladinTech/go-lib-common/client/aws"
	"bitbucket.org/moladinTech/go-lib-common/client/moladin_evo"
	"bitbucket.org/moladinTech/go-lib-common/client/notification"
	"bitbucket.org/moladinTech/go-lib-common/client/notification/slack"
	"bitbucket.org/moladinTech/go-lib-common/encryption"
	"bitbucket.org/moladinTech/go-lib-common/middleware/gin/auth"
	"bitbucket.org/moladinTech/go-lib-common/middleware/gin/limiter"
	"bitbucket.org/moladinTech/go-lib-common/middleware/gin/panic_recovery"
	"bitbucket.org/moladinTech/go-lib-common/middleware/gin/tracer"
	"bitbucket.org/moladinTech/go-lib-common/sentry"
	"bitbucket.org/moladinTech/go-lib-common/time"
	"github.com/go-playground/validator/v10"
)

type IRegistry interface {
	GetSentry() sentry.ISentry
	GetS3() aws.S3
	GetMoladinEvo() moladin_evo.IMoladinEvo
	GetSlack() slack.ISlack
	GetNotif() notification.INotification
	GetAuthMiddleware() auth.IMiddlewareAuth
	GetPanicRecoveryMiddleware() panic_recovery.IMiddlewarePanicRecovery
	GetTraceMiddleware() tracer.IMiddlewareTracer
	GetLimiterMiddleware() limiter.IMiddlewareLimiter
	GetTime() time.TimeItf
	GetValidator() *validator.Validate
	GetCache() cache.Cacher
	GetEncryption() encryption.IEncryption
}

type registry struct {
	sentry                  sentry.ISentry
	s3                      aws.S3
	moladinEvo              moladin_evo.IMoladinEvo
	slack                   slack.ISlack
	notif                   notification.INotification
	authMiddleware          auth.IMiddlewareAuth
	panicRecoveryMiddleware panic_recovery.IMiddlewarePanicRecovery
	tracerMiddleware        tracer.IMiddlewareTracer
	limiterMiddleware       limiter.IMiddlewareLimiter
	time                    time.TimeItf
	validator               *validator.Validate
	cache                   cache.Cacher
	encryption              encryption.IEncryption
}

func WithSentry(sentry sentry.ISentry) Option {
	return func(s *registry) {
		s.sentry = sentry
	}
}

func WithS3(s3 aws.S3) Option {
	return func(s *registry) {
		s.s3 = s3
	}
}

func WithMoladinEvo(moladinEvo moladin_evo.IMoladinEvo) Option {
	return func(s *registry) {
		s.moladinEvo = moladinEvo
	}
}

func WithSlack(slack slack.ISlack) Option {
	return func(s *registry) {
		s.slack = slack
	}
}

func WithNotif(notif notification.INotification) Option {
	return func(s *registry) {
		s.notif = notif
	}
}

func WithAuthMiddleware(authMiddleware auth.IMiddlewareAuth) Option {
	return func(s *registry) {
		s.authMiddleware = authMiddleware
	}
}

func WithPanicRecoveryMiddleware(panicRecoveryMiddleware panic_recovery.IMiddlewarePanicRecovery) Option {
	return func(s *registry) {
		s.panicRecoveryMiddleware = panicRecoveryMiddleware
	}
}

func WithTracerMiddleware(tracerMiddleware tracer.IMiddlewareTracer) Option {
	return func(s *registry) {
		s.tracerMiddleware = tracerMiddleware
	}
}

func WithLimiterMiddleware(limiterMiddleware limiter.IMiddlewareLimiter) Option {
	return func(s *registry) {
		s.limiterMiddleware = limiterMiddleware
	}
}

func WithTime(time time.TimeItf) Option {
	return func(s *registry) {
		s.time = time
	}
}

func WithValidator(validator *validator.Validate) Option {
	return func(s *registry) {
		s.validator = validator
	}
}

func WithCache(cache cache.Cacher) Option {
	return func(s *registry) {
		s.cache = cache
	}
}

func WithEncryption(encryption encryption.IEncryption) Option {
	return func(s *registry) {
		s.encryption = encryption
	}
}

type Option func(r *registry)

func NewRegistry(
	options ...Option,
) IRegistry {
	registry := &registry{}

	for _, option := range options {
		option(registry)
	}

	return registry
}

func (r *registry) GetSentry() sentry.ISentry {
	return r.sentry
}

func (r *registry) GetS3() aws.S3 {
	return r.s3
}

func (r *registry) GetMoladinEvo() moladin_evo.IMoladinEvo {
	return r.moladinEvo
}

func (r *registry) GetSlack() slack.ISlack {
	return r.slack
}

func (r *registry) GetNotif() notification.INotification {
	return r.notif
}

func (r *registry) GetAuthMiddleware() auth.IMiddlewareAuth {
	return r.authMiddleware
}

func (r *registry) GetPanicRecoveryMiddleware() panic_recovery.IMiddlewarePanicRecovery {
	return r.panicRecoveryMiddleware
}

func (r *registry) GetTraceMiddleware() tracer.IMiddlewareTracer {
	return r.tracerMiddleware
}

func (r *registry) GetLimiterMiddleware() limiter.IMiddlewareLimiter {
	return r.limiterMiddleware
}

func (r *registry) GetTime() time.TimeItf {
	return r.time
}

func (r *registry) GetValidator() *validator.Validate {
	return r.validator
}

func (r *registry) GetCache() cache.Cacher {
	return r.cache
}

func (r *registry) GetEncryption() encryption.IEncryption {
	return r.encryption
}
