package matrix

import (
	"context"

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
	pieceCount uint
	// #-of pieces received already
	receivedCount uint
	// useful piece count i.e. linearly
	// independent pieces decoder has received
	useful         uint
	coeffs, coded  Matrix
	workerQueue    []*work
	workerChans    []chan uint64
	supervisorChan chan *kodr.CodedPiece
}

type work struct {
	// 0-based work index which monotonically increases
	idx            uint
	srcRow, dstRow uint
	// weight is element of coefficient
	// matrix i.e. field element
	weight byte
	// any of two possible row operations
	op OP
}

type workerState struct {
	workerChan             chan uint64
	decoderState           *ParallelDecoderState
	currentWorkIdx         uint
	totalWorkCount         uint
	columnStart, columnEnd uint
}

func (p *ParallelDecoderState) createWork(src, dst uint, weight byte, op OP) {
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
OUT:
	for {
		select {
		case <-ctx.Done():
			break OUT

		case codedPiece := <-p.supervisorChan:
			// done with decoding, no need to work
			// on new coded piece !
			if p.useful >= p.pieceCount {
				continue OUT
			}

			p.coeffs = append(p.coeffs, codedPiece.Vector)
			p.coded = append(p.coded, codedPiece.Piece)

			// index of current piece of interest
			idx := uint(len(p.coeffs) - 1)

			// --- Stage A begins ---
			for j := uint(0); j < idx; j++ {
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

				p.useful = uint(len(p.coeffs))
				continue OUT
			}
			// --- Stage B ends ---

			for j := uint(0); j < idx; j++ {
				weight := p.coeffs[idx][j]
				p.coeffs[idx][j] = 0
				p.createWork(j, idx, weight, SUB_AFTER_MULT)
			}
			// --- Stage A ends ---

			// --- Stage C begins ---
			p.createWork(idx, idx, p.coeffs[idx][non_zero_idx], DIVISION)

			for k := uint(non_zero_idx); k < p.pieceCount; k++ {
				p.coeffs[idx][k] = p.field.Div(p.coeffs[idx][k], p.coeffs[idx][non_zero_idx])
			}
			// --- Stage C ends ---

			// --- Stage D begins ---
			for j := uint(0); j < idx; j++ {
				p.createWork(idx, j, p.coeffs[j][non_zero_idx], SUB_AFTER_MULT)
			}

			for j := uint(0); j < idx; j++ {
				weight := p.coeffs[j][non_zero_idx]
				p.coeffs[j][non_zero_idx] = 0

				for k := uint(non_zero_idx); k < p.pieceCount; k++ {
					tmp := p.field.Mul(p.coeffs[idx][k], weight)
					p.coeffs[j][k] = p.field.Add(p.coeffs[j][k], tmp)
				}
			}
			// --- Stage D ends ---

			// these many useful pieces decoder has as of now
			p.useful = uint(len(p.coeffs))

			// because decoding is complete !
			// workers doesn't need to be alive !
			if p.useful >= p.pieceCount {
				p.createWork(0, 0, 0, STOP)
			}

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
func (p *ParallelDecoderState) AddPiece(codedPiece *kodr.CodedPiece) {
	// it's blocking call, if chan is non-bufferred !
	//
	// better to use buffered channel
	p.supervisorChan <- codedPiece
}

// If enough #-of linearly independent pieces are received
// whole data is decoded, which denotes it's good time
// to start consuming !
func (p *ParallelDecoderState) IsDecoded() bool {
	return p.useful >= p.pieceCount
}
