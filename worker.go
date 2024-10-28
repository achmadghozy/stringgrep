package main

import (
	"fmt"
	"sync"
)

type stWorker struct {
	nId int // number
}

type stWorkerPool struct {
	findJob    chan string
	findRes    chan string
	error      chan error
	nWorkerMax int // n worker pool max
	stWorker
}

const (
	jobWalk = iota
	jobSearch
	jobWrite
	jobWalkCr
)

// Constructor worker
func newWorker(id int) *stWorker {
	return &stWorker{
		nId: id,
	}
}

func (w *stWorker) WalkJobGetter(task *stUICObjGrep, wg *sync.WaitGroup) {
	defer wg.Done()

	// Get all folder within the path
	task.getAllDirectories()
}

func (w *stWorker) WorkerJob(chTx chan string, chRx chan string, job uint, args1 int, args2 int, task *stUICObjGrep, wg *sync.WaitGroup, bJobFin *bool, chErr chan error) {
	defer w.WorkerFin(wg, bJobFin, chRx, job)

	for cJobs := range chRx {
		if cJobs == `done` {
			chTx <- `done`
			close(chRx)
			break
		} else {
			fmt.Printf(`job %d is started\n`, job)
			switch job {
			case jobWalk:
				{
					res, err := task.getFilePth2(cJobs)
					if err != nil {
						chErr <- err
					}
					for i := range res {
						chTx <- res[i]
					}
				}
			case jobSearch:
				{
					res, err := task.searchInFile(cJobs)
					if err != nil {
						chErr <- err
					}
					chTx <- res
				}
			case jobWrite:
				{
					err := task.writeResCSV(cJobs)
					if err != nil {
						chErr <- err
					}
				}
			default:
				break
			}
		}
	}

}

func (w *stWorker) WorkerFin(wg *sync.WaitGroup, bJobFin *bool, chRx chan string, nJob uint) {
	*bJobFin = true
	wg.Done()
}

func createWorkerPool(nMaxWorker int) *stWorkerPool {
	chFindJob := make(chan string)
	chFindRes := make(chan string)
	chError := make(chan error)

	return &stWorkerPool{
		findJob:    chFindJob,
		findRes:    chFindRes,
		error:      chError,
		nWorkerMax: nMaxWorker,
	}
}

func (wp *stWorkerPool) startWorkerPool(task *stUICObjGrep) {
	var wg sync.WaitGroup
	nWorker := 0
	nJob := jobWalk
	nWorkerWalkMax := 4
	nWorkerSearchMax := 3
	nWorkerWriterMax := 1
	bWalkFin := false
	bCrWalkJobFin := false
	bWriteFin := false
	bSearchFin := false
	var pFlag *bool

	for {

		if bWriteFin {
			break
		}

		if !bWalkFin {
			wp.nWorkerMax = nWorkerWalkMax
			nJob = jobWalk
			pFlag = &bWalkFin
		} else if !bCrWalkJobFin {
			nJob = jobWalkCr
			pFlag = &bCrWalkJobFin
		} else {
			wp.nWorkerMax = nWorkerSearchMax + nWorkerWriterMax
			nJob = jobSearch
			pFlag = &bSearchFin
		}

		for i := nWorker; i < wp.nWorkerMax; i++ {
			wg.Add(1)
			worker := newWorker(i)
			worker.WorkerJob(wp.findJob, wp.findRes, uint(nJob), 0, 0, task, &wg, pFlag, wp.error)
			nWorker++
			if nWorkerWriterMax < 1 && nJob == jobSearch {
				wg.Add(1)
				worker := newWorker(i)
				worker.WorkerJob(wp.findJob, wp.findRes, uint(jobWrite), 0, 0, task, &wg, &bWriteFin, wp.error)
				nWorker++
			}
		}
	}
}
