package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	chTasks := make(chan Task, len(tasks))
	var wgExit sync.WaitGroup
	var mtx sync.Mutex
	wgExit.Add(n)

	// Переводим все задания в канал и закрываем его
	for _, t := range tasks {
		chTasks <- t
	}
	close(chTasks)

	// Объявим функции для учета количества ошибок
	errorCount := 0
	ignoreErrors := m <= 0
	checkErrorCount := func() bool {
		return errorCount >= m
	}
	addErrorCount := func() {
		errorCount++
		println("errorCount = ", errorCount)
	}

	// Запустим n воркеров
	for i := 0; i < n; i++ {
		go worker(chTasks, &wgExit, &mtx, i, checkErrorCount, addErrorCount, ignoreErrors)
	}

	// Ждем завершение всех n воркеров
	wgExit.Wait()

	// Обработка итогового результата
	isErrorRes := errorCount >= m
	if isErrorRes && !ignoreErrors {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func worker(
	chTasks chan Task,
	wgExit *sync.WaitGroup,
	mtx *sync.Mutex,
	i int,
	checkErrorCount func() bool,
	addErrorCount func(),
	ignoreErrors bool,
) {
	defer wgExit.Done()

	for {
		// Получаем задание
		task, ok := <-chTasks
		if !ok {
			// Если заданий больше нет, завершаем работу воркера
			println(i, ": No more tasks. Done")
			return
		}

		println(i, ": Task received")
		if !ignoreErrors {
			// Перед тем как накинуться на задание, проверим не достигнут ли лимит по ошибкам
			mtx.Lock()
			isErrorLimit := checkErrorCount()
			mtx.Unlock()

			if isErrorLimit {
				println(i, ": Error limits. Done")
				return
			}
		}

		// Собственно выполняем задание
		taskErr := task()

		// Если задание выполнилось с ошибкой - увеличиваем счетчик
		if taskErr != nil && !ignoreErrors {
			println(i, ": Task error detected")
			mtx.Lock()
			addErrorCount()
			mtx.Unlock()
		} else {
			println(i, ": Single task done")
		}
	}
}
