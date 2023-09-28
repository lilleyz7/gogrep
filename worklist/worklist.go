package worklist

type FileEntry struct {
	path string
}

type Worklist struct {
	jobs chan FileEntry
}

func (f *FileEntry) GetPath() string {
	return f.path
}

func (w *Worklist) AddJob(file FileEntry) {
	w.jobs <- file
}

func (w *Worklist) NextJob() FileEntry {
	job := <-w.jobs
	return job
}

func (w *Worklist) CompleteWork(numOfWorkers int) {
	for i := 0; i < numOfWorkers; i++ {
		w.AddJob(FileEntry{""})
	}
}

func NewWorkList(bufferSize int) Worklist {
	return Worklist{
		make(chan FileEntry, bufferSize),
	}
}

func NewJob(path string) FileEntry {
	return FileEntry{path}
}
