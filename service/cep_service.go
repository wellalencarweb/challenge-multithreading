package service

import (
	"context"
	"time"

	"github.com/wellalencarweb/challenge-multithreading/domain"
	"github.com/wellalencarweb/challenge-multithreading/fetchers"
)

func FindCep(ctx context.Context, cep string) (domain.Result, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 100*time.Microsecond)
	defer cancel()

	ch := make(chan domain.Result, 2)

	go fetchers.FindByBrasilAPI(ctxTimeout, cep, ch)
	go fetchers.FindByViaCEP(ctxTimeout, cep, ch)

	select {
	case res := <-ch:
		if res.Error != nil {
			return domain.Result{}, res.Error
		}
		return res, nil
	case <-ctxTimeout.Done():
		return domain.Result{}, ctxTimeout.Err()
	}
}
