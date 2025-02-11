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
	// Здесь содержатся все промежуточные выходные каналы
	middleOutputs := make([]Out, 0)
	middleOutputs = append(middleOutputs, resultOutput)
	for _, stage := range stages {
		resultOutput = stage(resultOutput)
		middleOutputs = append(middleOutputs, resultOutput)
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
					go finishChannelSlice(middleOutputs)
					return
				}
			case <-done:
				go finishChannelSlice(middleOutputs)
				return
			}
		}
	}()

	return out
}

func finishChannelSlice(middleOutputs []Out) {
	for _, middleOutput := range middleOutputs {
		finishChannel(middleOutput)
	}
}

func finishChannel(output Out) {
	for range output {
		// Дочитываем до конца канал output
		continue
	}
}
