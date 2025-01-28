package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// Обертка канала, чтобы завершать его по сигналу done
	orDone := func(input In, done In) In {
		out := make(Bi)
		go func() {
			defer close(out)
			for {
				select {
				case v, ok := <-input:
					if !ok {
						return
					}
					select {
					case out <- v:
					case <-done:
						finishChannel(input)
						return
					}
				case <-done:
					finishChannel(input)
					return
				}
			}
		}()
		return out
	}

	// Последовательно оборачиваем входной канал через все стадии
	out := in
	for _, stage := range stages {
		out = stage(orDone(out, done))
	}

	return out
}

func finishChannel(input In) {
	for range input {
		// Дочитываем до конца канал input
		continue
	}
}
