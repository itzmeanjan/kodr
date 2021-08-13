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
)

type ParallelDecoderState struct {
	field *galoisfield.GF
	// this is generation size = G
	pieceCount uint
	// #-of pieces received already
	receivedCount  uint
	coeffs, coded  Matrix
	workCount      uint
	workerQueue    []*work
	workerCount    uint
	workerChans    []chan struct{}
	supervisorChan chan uint
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
	workerChan             chan struct{}
	decoderState           *ParallelDecoderState
	currentWorkIdx         uint
	totalWorkCount         uint
	columnStart, columnEnd uint
}

func (p *ParallelDecoderState) createWork(src, dst uint, weight byte, op OP) {
	w := work{srcRow: src, dstRow: dst, weight: weight, op: op}
	p.workerQueue = append(p.workerQueue, &w)
}

func (p *ParallelDecoderState) supervise(ctx context.Context) {
	var (
		linearlyDependentPieceCount uint = 0
	)

OUT:
	for {
		select {
		case <-ctx.Done():
			break OUT

		case idx := <-p.supervisorChan:
			// useful when linearly dependent pieces are received
			idx -= linearlyDependentPieceCount

			// --- Stage A begins ---
			for j := uint(0); j < idx; j++ {
				p.createWork(j, idx, p.coeffs[idx][j], SUB_AFTER_MULT)
			}

			for j := uint(0); j < idx; j++ {
				weight := p.coeffs[idx][j]
				p.coeffs[idx][j] = 0

				for k := j; k < p.pieceCount; k++ {
					tmp := p.field.Mul(p.coeffs[j][k], weight)
					p.coeffs[idx][k] = p.field.Add(p.coeffs[idx][k], tmp)
				}
			}
			// --- Stage A ends ---

			// --- Stage B begins ---
			// first column index for row `idx`
			// which has non-zero field element
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
				linearlyDependentPieceCount += 1

				p.coeffs[idx] = nil
				copy((p.coeffs)[idx:], (p.coeffs)[idx+1:])
				p.coeffs = (p.coeffs)[:len(p.coeffs)-1]

				continue OUT
			}
			// --- Stage B ends ---

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
		}
	}
}

func (p *ParallelDecoderState) AddPiece(coding_vector kodr.CodingVector, coded_data kodr.Piece) {
	p.coeffs = append(p.coeffs, coding_vector)
	p.coded = append(p.coded, coded_data)
	p.receivedCount += 1

	// supervisor should start working only when atleast
	// 2 coded pieces are received
	if p.receivedCount < 2 {
		return
	}

	// it's blocking call, if chan is non-bufferred !
	// lets supervisor know coded piece index to work on
	//
	// -1 added due to 0 based indexing
	p.supervisorChan <- p.receivedCount - 1
}
