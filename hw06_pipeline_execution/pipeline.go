package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		out = terminate(out, done)
		out = stage(out)
	}
	return out
}

func terminate(in, done In) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		for {
			select {
			case value, opened := <-in:
				if !opened {
					return
				}
				out <- value
			case <-done:
				return
			}
		}
	}()
	return out
}
