package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// Последовательно оборачиваем входной канал через все стадии, и получаем последний выходной канал
	resultOutput := in
	for _, stage := range stages {
		resultOutput = stage(resultOutput)
	}

	// Обертка канала, чтобы завершать его по сигналу done
	out := make(Bi)
	go func() {
		defer close(out)
		for {
			select {
			case v, ok := <-resultOutput:
				if !ok {
					return
				}
				select {
				case out <- v:
				case <-done:
					return
				}
			case <-done:
				return
			}
		}
	}()

	return out
}
