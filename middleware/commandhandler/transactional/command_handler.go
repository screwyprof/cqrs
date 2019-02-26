package transactional

import "github.com/screwyprof/cqrs"

func NewMiddleware(unitOfWork cqrs.UnitOfWork) cqrs.CommandHandlerMiddleware {
	return cqrs.CommandHandlerMiddleware(func(h cqrs.CommandHandler) cqrs.CommandHandler {
		return cqrs.CommandHandlerFunc(func(cmd cqrs.Command) (err error) {
			defer func() {
				// handler error occurred, don't commit
				if err != nil {
					return
				}

				if commitErr := unitOfWork.Commit(); commitErr != nil {
					err = unitOfWork.Rollback()
				}
			}()

			return h.Handle(cmd)
		})
	})
}
