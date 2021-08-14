package matrix

import (
	"context"
	"runtime"
	"sync/atomic"

	"github.com/cloud9-tools/go-galoisfield"
	"github.com/itzmeanjan/kodr"
)

type OP uint8

const (
	SUB_AFTER_MULT OP = iota + 1
	DIVISION
	// special signal denoting workers should stop working
	STOP
)

type ParallelDecoderState struct {
	field *galoisfield.GF
	// this is generation size = G
	pieceCount uint64
	// length of coded data part of coded piece
	pieceLen uint64
	// useful piece count i.e. linearly
	// independent pieces decoder has received
	useful                 uint64
	coeffs, coded          Matrix
	workerQueue            []*work
	workerChans            []chan uint64
	supervisorAddPieceChan chan *addRequest
	supervisorGetPieceChan chan *pieceRequest
}

type addRequest struct {
	piece *kodr.CodedPiece
	err   chan error
}

// decoded piece consumption request
type pieceRequest struct {
	idx  uint64
	resp chan *kodr.Piece
	err  chan error
}

type work struct {
	// which two rows of coded data matrix are involved
	// in this row operation
	srcRow, dstRow uint64
	// weight is element of coefficient
	// matrix i.e. field element
	weight byte
	// any of two possible row operations
	op OP
}

type workerState struct {
	workerChan             chan uint64
	columnStart, columnEnd uint64
}

func (p *ParallelDecoderState) createWork(src, dst uint64, weight byte, op OP) {
	w := work{srcRow: src, dstRow: dst, weight: weight, op: op}
	p.workerQueue = append(p.workerQueue, &w)
	idx := uint(len(p.workerQueue) - 1)

	for i := 0; i < len(p.workerChans); i++ {
		// it's blocking call, better to use buffered channel,
		// then it won't probably be !
		p.workerChans[i] <- uint64(idx)
	}
}

func (p *ParallelDecoderState) supervise(ctx context.Context) {
	// how many pieces received
	receivedCount := 0

OUT:
	for {
		select {
		case <-ctx.Done():
			break OUT

		case req := <-p.supervisorAddPieceChan:
			if req.piece.Len() != uint(p.pieceCount+p.pieceLen) {
				req.err <- kodr.ErrCodedPieceSizeMismatch
				continue OUT
			}

			// done with decoding, no need to work
			// on new coded piece !
			if p.IsDecoded() {
				req.err <- kodr.ErrAllUsefulPiecesReceived
				continue OUT
			}

			// piece to be processed further, returning nil error !
			req.err <- nil
			receivedCount++
			codedPiece := req.piece
			p.coeffs = append(p.coeffs, codedPiece.Vector)
			p.coded = append(p.coded, codedPiece.Piece)

			// minimum 2 pieces are requyired to start working
			if receivedCount < 2 {
				continue OUT
			}

			// index of current piece of interest
			idx := uint64(len(p.coeffs) - 1)

			// --- Stage A begins ---
			for j := uint64(0); j < idx; j++ {
				weight := p.coeffs[idx][j]

				for k := j + 1; k < p.pieceCount; k++ {
					tmp := p.field.Mul(p.coeffs[j][k], weight)
					p.coeffs[idx][k] = p.field.Add(p.coeffs[idx][k], tmp)
				}
			}

			// --- Stage B begins ---
			// first column index for row `idx`
			// which has non-zero field element
			// after `idx-1` column
			non_zero_idx := -1
			for j := idx; j < p.pieceCount; j++ {
				if p.coeffs[idx][j] != 0 {
					non_zero_idx = int(j)
					break
				}
			}

			// if no element is found to be non-zero,
			// it's a linearly dependent piece --- not useful
			if non_zero_idx == -1 {
				p.coeffs[idx] = nil
				copy((p.coeffs)[idx:], (p.coeffs)[idx+1:])
				p.coeffs = (p.coeffs)[:len(p.coeffs)-1]

				p.coded[idx] = nil
				copy((p.coded)[idx:], (p.coded)[idx+1:])
				p.coded = (p.coded)[:len(p.coded)-1]

				atomic.StoreUint64(&p.useful, uint64(p.coeffs.Rows()))
				continue OUT
			}
			// --- Stage B ends ---

			for j := uint64(0); j < idx; j++ {
				weight := p.coeffs[idx][j]
				p.coeffs[idx][j] = 0
				p.createWork(j, idx, weight, SUB_AFTER_MULT)
			}
			// --- Stage A ends ---

			// --- Stage C begins ---
			p.createWork(idx, idx, p.coeffs[idx][non_zero_idx], DIVISION)

			for k := uint64(non_zero_idx); k < p.pieceCount; k++ {
				p.coeffs[idx][k] = p.field.Div(p.coeffs[idx][k], p.coeffs[idx][non_zero_idx])
			}
			// --- Stage C ends ---

			// --- Stage D begins ---
			for j := uint64(0); j < idx; j++ {
				p.createWork(idx, j, p.coeffs[j][non_zero_idx], SUB_AFTER_MULT)
			}

			for j := uint64(0); j < idx; j++ {
				weight := p.coeffs[j][non_zero_idx]
				p.coeffs[j][non_zero_idx] = 0

				for k := uint64(non_zero_idx); k < p.pieceCount; k++ {
					tmp := p.field.Mul(p.coeffs[idx][k], weight)
					p.coeffs[j][k] = p.field.Add(p.coeffs[j][k], tmp)
				}
			}
			// --- Stage D ends ---

			// these many useful pieces decoder has as of now
			atomic.StoreUint64(&p.useful, uint64(p.coeffs.Rows()))

			// because decoding is complete !
			// workers doesn't need to be alive !
			if p.IsDecoded() {
				p.createWork(0, 0, 0, STOP)
			}

		case req := <-p.supervisorGetPieceChan:
			if req.idx >= p.pieceCount {
				req.err <- kodr.ErrPieceOutOfBound
				continue OUT
			}

			if req.idx >= uint64(p.coeffs.Rows()) {
				req.err <- kodr.ErrPieceNotDecodedYet
				continue OUT
			}

			if p.IsDecoded() {
				req.resp <- (*kodr.Piece)(&p.coded[req.idx])
				continue OUT
			}

			cols := uint64(p.coeffs.Cols())
			decoded := true

		NESTED:
			for i := uint64(0); i < cols; i++ {
				switch i {
				case req.idx:
					if p.coeffs[req.idx][i] != 1 {
						decoded = false
						break NESTED
					}

				default:
					if p.coeffs[req.idx][i] != 0 {
						decoded = false
						break NESTED
					}

				}
			}

			if !decoded {
				req.err <- kodr.ErrPieceNotDecodedYet
				continue OUT
			}

			req.resp <- (*kodr.Piece)(&p.coded[req.idx])
			continue OUT

		}
	}
}

func (p *ParallelDecoderState) work(ctx context.Context, wState *workerState) {
OUT:
	for {
		select {
		case <-ctx.Done():
			break OUT

		case idx := <-wState.workerChan:
			w := p.workerQueue[idx]

			switch w.op {
			case SUB_AFTER_MULT:
				for i := wState.columnStart; i <= wState.columnEnd; i++ {
					tmp := p.field.Mul(p.coded[w.srcRow][i], w.weight)
					p.coded[w.dstRow][i] = p.field.Add(p.coded[w.dstRow][i], tmp)
				}

			case DIVISION:
				for i := wState.columnStart; i <= wState.columnEnd; i++ {
					p.coded[w.dstRow][i] = p.field.Add(p.coded[w.srcRow][i], w.weight)
				}

			case STOP:
				// supervisor signals decoding is complete !
				break OUT

			}
		}
	}
}

// Adds new coded piece to decoder state, so that it can process
// and progressively decoded pieces
//
// Before invoking this method, it's good idea to check
// `IsDecoded` method & refrain from invoking if already
// decoded
//
// It's concurrent safe !
func (p *ParallelDecoderState) AddPiece(codedPiece *kodr.CodedPiece) error {
	errChan := make(chan error, 1)
	req := addRequest{piece: codedPiece, err: errChan}
	p.supervisorAddPieceChan <- &req

	return <-errChan
}

// If enough #-of linearly independent pieces are received
// whole data is decoded, which denotes it's good time
// to start consuming !
//
// It's concurrent safe !
func (p *ParallelDecoderState) IsDecoded() bool {
	return atomic.LoadUint64(&p.useful) >= p.pieceCount
}

// Fetch decoded piece by index, can also return piece when not fully
// decoded, given requested piece is decoded
func (p *ParallelDecoderState) GetPiece(idx uint64) (kodr.Piece, error) {
	respChan := make(chan *kodr.Piece, 1)
	errChan := make(chan error, 1)
	req := pieceRequest{idx: idx, resp: respChan, err: errChan}

	// this may block !
	p.supervisorGetPieceChan <- &req

	// waiting for response !
	select {
	case err := <-errChan:
		return nil, err
	case piece := <-respChan:
		return *piece, nil
	}
}

// Current state of coding coefficient matrix
//
// NOTE: Don't mutate matrix, use only for writing test cases !
func (p *ParallelDecoderState) CoefficientMatrix() Matrix {
	return p.coeffs
}

// Current state of coded piece matrix, which is updated
// along side coding coefficient matrix ( during parallel rref )
//
// NOTE: Don't mutate matrix, use only for writing test cases !
func (p *ParallelDecoderState) CodedPieceMatrix() Matrix {
	return p.coded
}

// Each worker must at least take responsibility of
// 8-bytes slice of coded data & each of these
// worker slices are non-overlapping
func workerCount(pieceLen uint64) uint64 {
	wcount := pieceLen / 8
	cpus := uint64(runtime.NumCPU()) << 1
	if wcount > cpus {
		return cpus
	}
	return wcount
}

// Splitting coded data matrix mutation responsibility among workers
// Each of these slices are non-overlapping
func splitWork(pieceLen uint64) []*workerState {
	wcount := workerCount(pieceLen)
	span := pieceLen / wcount
	workers := make([]*workerState, 0, wcount)
	for i := uint64(0); i < wcount; i++ {
		start := span * i
		end := span*(i+1) - 1
		if i == wcount-1 {
			end = pieceLen - 1
		}

		ws := workerState{
			workerChan:  make(chan uint64, wcount),
			columnStart: start,
			columnEnd:   end,
		}
		workers = append(workers, &ws)
	}
	return workers
}

func NewParallelDecoderState(ctx context.Context, pieceCount, pieceLen uint64) *ParallelDecoderState {
	splitted := splitWork(pieceLen)

	dec := ParallelDecoderState{
		field:                  galoisfield.DefaultGF256,
		pieceCount:             pieceCount,
		pieceLen:               pieceLen,
		coeffs:                 make([][]byte, 0, pieceCount),
		coded:                  make([][]byte, 0, pieceCount),
		workerQueue:            make([]*work, 0),
		supervisorAddPieceChan: make(chan *addRequest, pieceCount),
		supervisorGetPieceChan: make(chan *pieceRequest, 1),
	}

	workerChans := make([]chan uint64, 0, len(splitted))
	for i := 0; i < len(splitted); i++ {
		func(idx int) {
			workerChans = append(workerChans, splitted[i].workerChan)
			// each worker runs on its own go-routine
			go dec.work(ctx, splitted[idx])
		}(i)
	}

	dec.workerChans = workerChans
	// supervisor runs on its own go-routine
	go dec.supervise(ctx)
	return &dec
}
